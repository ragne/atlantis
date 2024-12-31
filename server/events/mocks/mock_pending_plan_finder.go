// Code generated by pegomock. DO NOT EDIT.
// Source: github.com/runatlantis/atlantis/server/events (interfaces: PendingPlanFinder)

package mocks

import (
	pegomock "github.com/petergtz/pegomock/v4"
	events "github.com/runatlantis/atlantis/server/events"
	"reflect"
	"time"
)

type MockPendingPlanFinder struct {
	fail func(message string, callerSkip ...int)
}

func NewMockPendingPlanFinder(options ...pegomock.Option) *MockPendingPlanFinder {
	mock := &MockPendingPlanFinder{}
	for _, option := range options {
		option.Apply(mock)
	}
	return mock
}

func (mock *MockPendingPlanFinder) SetFailHandler(fh pegomock.FailHandler) { mock.fail = fh }
func (mock *MockPendingPlanFinder) FailHandler() pegomock.FailHandler      { return mock.fail }

func (mock *MockPendingPlanFinder) DeletePlans(pullDir string) error {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockPendingPlanFinder().")
	}
	_params := []pegomock.Param{pullDir}
	_result := pegomock.GetGenericMockFrom(mock).Invoke("DeletePlans", _params, []reflect.Type{reflect.TypeOf((*error)(nil)).Elem()})
	var _ret0 error
	if len(_result) != 0 {
		if _result[0] != nil {
			_ret0 = _result[0].(error)
		}
	}
	return _ret0
}

func (mock *MockPendingPlanFinder) Find(pullDir string) ([]events.PendingPlan, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockPendingPlanFinder().")
	}
	_params := []pegomock.Param{pullDir}
	_result := pegomock.GetGenericMockFrom(mock).Invoke("Find", _params, []reflect.Type{reflect.TypeOf((*[]events.PendingPlan)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var _ret0 []events.PendingPlan
	var _ret1 error
	if len(_result) != 0 {
		if _result[0] != nil {
			_ret0 = _result[0].([]events.PendingPlan)
		}
		if _result[1] != nil {
			_ret1 = _result[1].(error)
		}
	}
	return _ret0, _ret1
}

func (mock *MockPendingPlanFinder) VerifyWasCalledOnce() *VerifierMockPendingPlanFinder {
	return &VerifierMockPendingPlanFinder{
		mock:                   mock,
		invocationCountMatcher: pegomock.Times(1),
	}
}

func (mock *MockPendingPlanFinder) VerifyWasCalled(invocationCountMatcher pegomock.InvocationCountMatcher) *VerifierMockPendingPlanFinder {
	return &VerifierMockPendingPlanFinder{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
	}
}

func (mock *MockPendingPlanFinder) VerifyWasCalledInOrder(invocationCountMatcher pegomock.InvocationCountMatcher, inOrderContext *pegomock.InOrderContext) *VerifierMockPendingPlanFinder {
	return &VerifierMockPendingPlanFinder{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		inOrderContext:         inOrderContext,
	}
}

func (mock *MockPendingPlanFinder) VerifyWasCalledEventually(invocationCountMatcher pegomock.InvocationCountMatcher, timeout time.Duration) *VerifierMockPendingPlanFinder {
	return &VerifierMockPendingPlanFinder{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		timeout:                timeout,
	}
}

type VerifierMockPendingPlanFinder struct {
	mock                   *MockPendingPlanFinder
	invocationCountMatcher pegomock.InvocationCountMatcher
	inOrderContext         *pegomock.InOrderContext
	timeout                time.Duration
}

func (verifier *VerifierMockPendingPlanFinder) DeletePlans(pullDir string) *MockPendingPlanFinder_DeletePlans_OngoingVerification {
	_params := []pegomock.Param{pullDir}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "DeletePlans", _params, verifier.timeout)
	return &MockPendingPlanFinder_DeletePlans_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockPendingPlanFinder_DeletePlans_OngoingVerification struct {
	mock              *MockPendingPlanFinder
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockPendingPlanFinder_DeletePlans_OngoingVerification) GetCapturedArguments() string {
	pullDir := c.GetAllCapturedArguments()
	return pullDir[len(pullDir)-1]
}

func (c *MockPendingPlanFinder_DeletePlans_OngoingVerification) GetAllCapturedArguments() (_param0 []string) {
	_params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(_params) > 0 {
		if len(_params) > 0 {
			_param0 = make([]string, len(c.methodInvocations))
			for u, param := range _params[0] {
				_param0[u] = param.(string)
			}
		}
	}
	return
}

func (verifier *VerifierMockPendingPlanFinder) Find(pullDir string) *MockPendingPlanFinder_Find_OngoingVerification {
	_params := []pegomock.Param{pullDir}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "Find", _params, verifier.timeout)
	return &MockPendingPlanFinder_Find_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockPendingPlanFinder_Find_OngoingVerification struct {
	mock              *MockPendingPlanFinder
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockPendingPlanFinder_Find_OngoingVerification) GetCapturedArguments() string {
	pullDir := c.GetAllCapturedArguments()
	return pullDir[len(pullDir)-1]
}

func (c *MockPendingPlanFinder_Find_OngoingVerification) GetAllCapturedArguments() (_param0 []string) {
	_params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(_params) > 0 {
		if len(_params) > 0 {
			_param0 = make([]string, len(c.methodInvocations))
			for u, param := range _params[0] {
				_param0[u] = param.(string)
			}
		}
	}
	return
}
