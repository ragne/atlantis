package bitbucketcloud_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/runatlantis/atlantis/server/events/models"
	"github.com/runatlantis/atlantis/server/events/vcs/bitbucketcloud"
	"github.com/runatlantis/atlantis/server/logging"
	. "github.com/runatlantis/atlantis/testing"
)

const diffstatURL = "/2.0/repositories/owner/repo/pullrequests/1/diffstat"

// Should follow pagination properly.
func TestClient_GetModifiedFilesPagination(t *testing.T) {
	logger := logging.NewNoopLogger(t)
	respTemplate := `
{
    "pagelen": 1,
    "values": [
        {
            "type": "diffstat",
            "status": "modified",
            "lines_removed": 1,
            "lines_added": 2,
            "old": {
                "path": "%s",
                "type": "commit_file",
                "links": {
                    "self": {
                        "href": "https://api.bitbucket.org/2.0/repositories/bitbucket/geordi/src/e1749643d655d7c7014001a6c0f58abaf42ad850/setup.py"
                    }
                }
            },
            "new": {
                "path": "%s",
                "type": "commit_file",
                "links": {
                    "self": {
                        "href": "https://api.bitbucket.org/2.0/repositories/bitbucket/geordi/src/d222fa235229c55dad20b190b0b571adf737d5a6/setup.py"
                    }
                }
            }
        }
    ],
    "page": 1,
    "size": 1
`
	firstResp := fmt.Sprintf(respTemplate, "file1.txt", "file2.txt")
	secondResp := fmt.Sprintf(respTemplate, "file2.txt", "file3.txt")
	var serverURL string

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		// The first request should hit this URL.
		case diffstatURL:
			resp := firstResp + fmt.Sprintf(`,"next": "%s%s?page=2"}`, serverURL, diffstatURL)
			w.Write([]byte(resp)) // nolint: errcheck
			return
			// The second should hit this URL.
		case fmt.Sprintf("%s?page=2", diffstatURL):
			w.Write([]byte(secondResp + "}")) // nolint: errcheck
		default:
			t.Errorf("got unexpected request at %q", r.RequestURI)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
	}))
	defer testServer.Close()

	serverURL = testServer.URL
	client := bitbucketcloud.NewClient(http.DefaultClient, "user", "pass", "runatlantis.io")
	client.BaseURL = testServer.URL

	files, err := client.GetModifiedFiles(
		logger,
		models.Repo{
			FullName:          "owner/repo",
			Owner:             "owner",
			Name:              "repo",
			CloneURL:          "",
			SanitizedCloneURL: "",
			VCSHost: models.VCSHost{
				Type:     models.BitbucketCloud,
				Hostname: "bitbucket.org",
			},
		}, models.PullRequest{
			Num: 1,
		})
	Ok(t, err)
	Equals(t, []string{"file1.txt", "file2.txt", "file3.txt"}, files)
}

// If the "old" key in the list of files is nil we shouldn't error.
func TestClient_GetModifiedFilesOldNil(t *testing.T) {
	logger := logging.NewNoopLogger(t)
	resp := `
{
  "pagelen": 500,
  "values": [
    {
      "status": "added",
      "old": null,
      "lines_removed": 0,
      "lines_added": 2,
      "new": {
        "path": "parent/child/file1.txt",
        "type": "commit_file",
        "links": {
          "self": {
            "href": "https://api.bitbucket.org/2.0/repositories/lkysow/atlantis-example/src/1ed8205eec00dab4f1c0a8c486a4492c98c51f8e/main.tf"
          }
        }
      },
      "type": "diffstat"
    }
  ],
  "page": 1,
  "size": 1
}`

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		// The first request should hit this URL.
		case diffstatURL:
			w.Write([]byte(resp)) // nolint: errcheck
			return
		default:
			t.Errorf("got unexpected request at %q", r.RequestURI)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
	}))
	defer testServer.Close()

	client := bitbucketcloud.NewClient(http.DefaultClient, "user", "pass", "runatlantis.io")
	client.BaseURL = testServer.URL

	files, err := client.GetModifiedFiles(
		logger,
		models.Repo{
			FullName:          "owner/repo",
			Owner:             "owner",
			Name:              "repo",
			CloneURL:          "",
			SanitizedCloneURL: "",
			VCSHost: models.VCSHost{
				Type:     models.BitbucketCloud,
				Hostname: "bitbucket.org",
			},
		}, models.PullRequest{
			Num: 1,
		})
	Ok(t, err)
	Equals(t, []string{"parent/child/file1.txt"}, files)
}

func TestClient_PullIsApproved(t *testing.T) {
	logger := logging.NewNoopLogger(t)
	cases := []struct {
		description string
		testdata    string
		exp         bool
	}{
		{
			"no approvers",
			"pull-unapproved.json",
			false,
		},
		{
			"approver is the author",
			"pull-approved-by-author.json",
			false,
		},
		{
			"single approver",
			"pull-approved.json",
			true,
		},
		{
			"two approvers one author",
			"pull-approved-multiple.json",
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			json, err := os.ReadFile(filepath.Join("testdata", c.testdata))
			Ok(t, err)
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				// The first request should hit this URL.
				case "/2.0/repositories/owner/repo/pullrequests/1":
					w.Write(json) // nolint: errcheck
					return
				default:
					t.Errorf("got unexpected request at %q", r.RequestURI)
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
			}))
			defer testServer.Close()

			client := bitbucketcloud.NewClient(http.DefaultClient, "user", "pass", "runatlantis.io")
			client.BaseURL = testServer.URL

			repo, err := models.NewRepo(models.BitbucketServer, "owner/repo", "https://bitbucket.org/owner/repo.git", "user", "token")
			Ok(t, err)
			approvalStatus, err := client.PullIsApproved(
				logger,
				repo, models.PullRequest{
					Num:        1,
					HeadBranch: "branch",
					Author:     "author",
					BaseRepo:   repo,
				})
			Ok(t, err)
			Equals(t, c.exp, approvalStatus.IsApproved)
		})
	}
}

func TestClient_PullIsMergeable(t *testing.T) {
	logger := logging.NewNoopLogger(t)
	cases := map[string]struct {
		DiffStat     string
		ExpMergeable bool
	}{
		"mergeable": {
			DiffStat: `{
				"pagelen": 500,
				"values": [
				{
					"status": "added",
					"old": null,
					"lines_removed": 0,
					"lines_added": 2,
					"new": {
						"path": "parent/child/file1.txt",
						"type": "commit_file",
						"links": {
							"self": {
								"href": "https://api.bitbucket.org/2.0/repositories/lkysow/atlantis-example/src/1ed8205eec00dab4f1c0a8c486a4492c98c51f8e/main.tf"
							}
						}
					},
					"type": "diffstat"
				}
			],
				"page": 1,
				"size": 1
			}`,
			ExpMergeable: true,
		},
		"merge conflict": {
			DiffStat: `{
			  "pagelen": 500,
			  "values": [
				{
				  "status": "merge conflict",
				  "old": {
					"path": "main.tf",
					"type": "commit_file",
					"links": {
					  "self": {
						"href": "https://api.bitbucket.org/2.0/repositories/lkysow/atlantis-example/src/6d6a8026a788621b37a9ac422a7d0ebb1500e85f/main.tf"
					  }
					}
				  },
				  "lines_removed": 1,
				  "lines_added": 0,
				  "new": {
					"path": "main.tf",
					"type": "commit_file",
					"links": {
					  "self": {
						"href": "https://api.bitbucket.org/2.0/repositories/lkysow/atlantis-example/src/742e76108714365788f5681e99e4a64f45dce147/main.tf"
					  }
					}
				  },
				  "type": "diffstat"
				}
			  ],
			  "page": 1,
			  "size": 1
			}`,
			ExpMergeable: false,
		},
		"merge conflict due to file deleted": {
			DiffStat: `{
			  "pagelen": 500,
			  "values": [
				{
				  "status": "local deleted",
				  "old": null,
				  "lines_removed": 0,
				  "lines_added": 3,
				  "new": {
					"path": "main.tf",
					"type": "commit_file",
					"links": {
					  "self": {
						"href": "https://api.bitbucket.org/2.0/repositories/lkysow/atlantis-example/src/3539b9f51c9f91e8f6280e89c62e2673ddc51144/main.tf"
					  }
					}
				  },
				  "type": "diffstat"
				}
			  ],
			  "page": 1,
			  "size": 1
			}`,
			ExpMergeable: false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case diffstatURL:
					w.Write([]byte(c.DiffStat)) // nolint: errcheck
					return
				default:
					t.Errorf("got unexpected request at %q", r.RequestURI)
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
			}))
			defer testServer.Close()

			client := bitbucketcloud.NewClient(http.DefaultClient, "user", "pass", "runatlantis.io")
			client.BaseURL = testServer.URL

			actMergeable, err := client.PullIsMergeable(
				logger,
				models.Repo{
					FullName:          "owner/repo",
					Owner:             "owner",
					Name:              "repo",
					CloneURL:          "",
					SanitizedCloneURL: "",
					VCSHost: models.VCSHost{
						Type:     models.BitbucketCloud,
						Hostname: "bitbucket.org",
					},
				}, models.PullRequest{
					Num: 1,
				}, "atlantis-test")
			Ok(t, err)
			Equals(t, c.ExpMergeable, actMergeable)
		})
	}

}

func TestClient_MarkdownPullLink(t *testing.T) {
	client := bitbucketcloud.NewClient(http.DefaultClient, "user", "pass", "runatlantis.io")
	pull := models.PullRequest{Num: 1}
	s, _ := client.MarkdownPullLink(pull)
	exp := "#1"
	Equals(t, exp, s)
}

func TestClient_GetMyUUID(t *testing.T) {
	json, err := os.ReadFile(filepath.Join("testdata", "user.json"))
	Ok(t, err)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/2.0/user":
			w.Write(json) // nolint: errcheck
			return
		default:
			t.Errorf("got unexpected request at %q", r.RequestURI)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
	}))
	defer testServer.Close()

	client := bitbucketcloud.NewClient(http.DefaultClient, "user", "pass", "runatlantis.io")
	client.BaseURL = testServer.URL
	v, _ := client.GetMyUUID()
	Equals(t, v, "{00000000-0000-0000-0000-000000000001}")
}

func TestClient_GetComment(t *testing.T) {
	json, err := os.ReadFile(filepath.Join("testdata", "comments.json"))
	Ok(t, err)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/2.0/repositories/myorg/myrepo/pullrequests/5/comments":
			w.Write(json) // nolint: errcheck
			return
		default:
			t.Errorf("got unexpected request at %q", r.RequestURI)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
	}))
	defer testServer.Close()

	client := bitbucketcloud.NewClient(http.DefaultClient, "user", "pass", "runatlantis.io")
	client.BaseURL = testServer.URL
	v, _ := client.GetPullRequestComments(
		models.Repo{
			FullName:          "myorg/myrepo",
			Owner:             "owner",
			Name:              "myrepo",
			CloneURL:          "",
			SanitizedCloneURL: "",
			VCSHost: models.VCSHost{
				Type:     models.BitbucketCloud,
				Hostname: "bitbucket.org",
			},
		}, 5)

	Equals(t, len(v), 5)
	exp := "Plan"
	Assert(t, strings.Contains(v[1].Content.Raw, exp), "Comment should contain word \"%s\", has \"%s\"", exp, v[1].Content.Raw)
}

func TestClient_DeleteComment(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/2.0/repositories/myorg/myrepo/pullrequests/5/comments/1":
			if r.Method == "DELETE" {
				w.WriteHeader(http.StatusNoContent)
			}
			return
		default:
			t.Errorf("got unexpected request at %q", r.RequestURI)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
	}))
	defer testServer.Close()

	client := bitbucketcloud.NewClient(http.DefaultClient, "user", "pass", "runatlantis.io")
	client.BaseURL = testServer.URL
	err := client.DeletePullRequestComment(
		models.Repo{
			FullName:          "myorg/myrepo",
			Owner:             "owner",
			Name:              "myrepo",
			CloneURL:          "",
			SanitizedCloneURL: "",
			VCSHost: models.VCSHost{
				Type:     models.BitbucketCloud,
				Hostname: "bitbucket.org",
			},
		}, 5, 1)
	Ok(t, err)
}

func TestClient_HidePRComments(t *testing.T) {
	logger := logging.NewNoopLogger(t)
	comments, err := os.ReadFile(filepath.Join("testdata", "comments.json"))
	Ok(t, err)
	json, err := os.ReadFile(filepath.Join("testdata", "user.json"))
	Ok(t, err)

	called := 0

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		// we have two comments in the test file
		// The code is going to delete them all and then create a new one
		case "/2.0/repositories/myorg/myrepo/pullrequests/5/comments/498931882":
			if r.Method == "DELETE" {
				w.WriteHeader(http.StatusNoContent)
			}
			w.Write([]byte("")) // nolint: errcheck
			called += 1
			return
			// This is the second one
		case "/2.0/repositories/myorg/myrepo/pullrequests/5/comments/498931784":
			if r.Method == "DELETE" {
				http.Error(w, "", http.StatusNoContent)
			}
			w.Write([]byte("")) // nolint: errcheck
			called += 1
			return
		case "/2.0/repositories/myorg/myrepo/pullrequests/5/comments/49893111":
			Assert(t, r.Method != "DELETE", "Shouldn't delete this one")
			return
		case "/2.0/repositories/myorg/myrepo/pullrequests/5/comments":
			w.Write(comments) // nolint: errcheck
			return
		case "/2.0/user":
			w.Write(json) // nolint: errcheck
			return
		default:
			t.Errorf("got unexpected request at %q", r.RequestURI)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
	}))
	defer testServer.Close()

	client := bitbucketcloud.NewClient(http.DefaultClient, "user", "pass", "runatlantis.io")
	client.BaseURL = testServer.URL
	err = client.HidePrevCommandComments(logger,
		models.Repo{
			FullName:          "myorg/myrepo",
			Owner:             "owner",
			Name:              "myrepo",
			CloneURL:          "",
			SanitizedCloneURL: "",
			VCSHost: models.VCSHost{
				Type:     models.BitbucketCloud,
				Hostname: "bitbucket.org",
			},
		}, 5, "plan", "")
	Ok(t, err)
	Equals(t, 2, called)
}
