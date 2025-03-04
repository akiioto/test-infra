// Code generated by mockery v2.46.3. DO NOT EDIT.

package oidcmocks

import (
	mock "github.com/stretchr/testify/mock"
	context "golang.org/x/net/context"

	oidc "github.com/kyma-project/test-infra/pkg/oidc"
)

// MockTokenVerifierInterface is an autogenerated mock type for the TokenVerifierInterface type
type MockTokenVerifierInterface struct {
	mock.Mock
}

type MockTokenVerifierInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTokenVerifierInterface) EXPECT() *MockTokenVerifierInterface_Expecter {
	return &MockTokenVerifierInterface_Expecter{mock: &_m.Mock}
}

// Verify provides a mock function with given fields: _a0, _a1
func (_m *MockTokenVerifierInterface) Verify(_a0 context.Context, _a1 string) (oidc.Token, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Verify")
	}

	var r0 oidc.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (oidc.Token, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) oidc.Token); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(oidc.Token)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTokenVerifierInterface_Verify_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Verify'
type MockTokenVerifierInterface_Verify_Call struct {
	*mock.Call
}

// Verify is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *MockTokenVerifierInterface_Expecter) Verify(_a0 interface{}, _a1 interface{}) *MockTokenVerifierInterface_Verify_Call {
	return &MockTokenVerifierInterface_Verify_Call{Call: _e.mock.On("Verify", _a0, _a1)}
}

func (_c *MockTokenVerifierInterface_Verify_Call) Run(run func(_a0 context.Context, _a1 string)) *MockTokenVerifierInterface_Verify_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockTokenVerifierInterface_Verify_Call) Return(_a0 oidc.Token, _a1 error) *MockTokenVerifierInterface_Verify_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTokenVerifierInterface_Verify_Call) RunAndReturn(run func(context.Context, string) (oidc.Token, error)) *MockTokenVerifierInterface_Verify_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTokenVerifierInterface creates a new instance of MockTokenVerifierInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTokenVerifierInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTokenVerifierInterface {
	mock := &MockTokenVerifierInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
# (2025-03-04)