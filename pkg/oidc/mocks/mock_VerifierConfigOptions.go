// Code generated by mockery v2.38.0. DO NOT EDIT.

package oidcmocks

import (
	oidc "github.com/coreos/go-oidc/v3/oidc"
	mock "github.com/stretchr/testify/mock"
)

// MockVerifierConfigOptions is an autogenerated mock type for the VerifierConfigOptions type
type MockVerifierConfigOptions struct {
	mock.Mock
}

type MockVerifierConfigOptions_Expecter struct {
	mock *mock.Mock
}

func (_m *MockVerifierConfigOptions) EXPECT() *MockVerifierConfigOptions_Expecter {
	return &MockVerifierConfigOptions_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: config
func (_m *MockVerifierConfigOptions) Execute(config *oidc.Config) error {
	ret := _m.Called(config)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*oidc.Config) error); ok {
		r0 = rf(config)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockVerifierConfigOptions_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockVerifierConfigOptions_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - config *oidc.Config
func (_e *MockVerifierConfigOptions_Expecter) Execute(config interface{}) *MockVerifierConfigOptions_Execute_Call {
	return &MockVerifierConfigOptions_Execute_Call{Call: _e.mock.On("Execute", config)}
}

func (_c *MockVerifierConfigOptions_Execute_Call) Run(run func(config *oidc.Config)) *MockVerifierConfigOptions_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*oidc.Config))
	})
	return _c
}

func (_c *MockVerifierConfigOptions_Execute_Call) Return(_a0 error) *MockVerifierConfigOptions_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockVerifierConfigOptions_Execute_Call) RunAndReturn(run func(*oidc.Config) error) *MockVerifierConfigOptions_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockVerifierConfigOptions creates a new instance of MockVerifierConfigOptions. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockVerifierConfigOptions(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockVerifierConfigOptions {
	mock := &MockVerifierConfigOptions{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
# (2025-03-04)