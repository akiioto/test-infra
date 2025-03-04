// Code generated by mockery v2.39.1. DO NOT EDIT.

package pipelinesmocks

import (
	context "context"

	pipelines "github.com/microsoft/azure-devops-go-api/azuredevops/v7/pipelines"
	mock "github.com/stretchr/testify/mock"
)

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

type MockClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockClient) EXPECT() *MockClient_Expecter {
	return &MockClient_Expecter{mock: &_m.Mock}
}

// GetRun provides a mock function with given fields: ctx, args
func (_m *MockClient) GetRun(ctx context.Context, args pipelines.GetRunArgs) (*pipelines.Run, error) {
	ret := _m.Called(ctx, args)

	if len(ret) == 0 {
		panic("no return value specified for GetRun")
	}

	var r0 *pipelines.Run
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, pipelines.GetRunArgs) (*pipelines.Run, error)); ok {
		return rf(ctx, args)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pipelines.GetRunArgs) *pipelines.Run); ok {
		r0 = rf(ctx, args)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pipelines.Run)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, pipelines.GetRunArgs) error); ok {
		r1 = rf(ctx, args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_GetRun_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRun'
type MockClient_GetRun_Call struct {
	*mock.Call
}

// GetRun is a helper method to define mock.On call
//   - ctx context.Context
//   - args pipelines.GetRunArgs
func (_e *MockClient_Expecter) GetRun(ctx interface{}, args interface{}) *MockClient_GetRun_Call {
	return &MockClient_GetRun_Call{Call: _e.mock.On("GetRun", ctx, args)}
}

func (_c *MockClient_GetRun_Call) Run(run func(ctx context.Context, args pipelines.GetRunArgs)) *MockClient_GetRun_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(pipelines.GetRunArgs))
	})
	return _c
}

func (_c *MockClient_GetRun_Call) Return(_a0 *pipelines.Run, _a1 error) *MockClient_GetRun_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_GetRun_Call) RunAndReturn(run func(context.Context, pipelines.GetRunArgs) (*pipelines.Run, error)) *MockClient_GetRun_Call {
	_c.Call.Return(run)
	return _c
}

// RunPipeline provides a mock function with given fields: ctx, args
func (_m *MockClient) RunPipeline(ctx context.Context, args pipelines.RunPipelineArgs) (*pipelines.Run, error) {
	ret := _m.Called(ctx, args)

	if len(ret) == 0 {
		panic("no return value specified for RunPipeline")
	}

	var r0 *pipelines.Run
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, pipelines.RunPipelineArgs) (*pipelines.Run, error)); ok {
		return rf(ctx, args)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pipelines.RunPipelineArgs) *pipelines.Run); ok {
		r0 = rf(ctx, args)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pipelines.Run)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, pipelines.RunPipelineArgs) error); ok {
		r1 = rf(ctx, args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_RunPipeline_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RunPipeline'
type MockClient_RunPipeline_Call struct {
	*mock.Call
}

// RunPipeline is a helper method to define mock.On call
//   - ctx context.Context
//   - args pipelines.RunPipelineArgs
func (_e *MockClient_Expecter) RunPipeline(ctx interface{}, args interface{}) *MockClient_RunPipeline_Call {
	return &MockClient_RunPipeline_Call{Call: _e.mock.On("RunPipeline", ctx, args)}
}

func (_c *MockClient_RunPipeline_Call) Run(run func(ctx context.Context, args pipelines.RunPipelineArgs)) *MockClient_RunPipeline_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(pipelines.RunPipelineArgs))
	})
	return _c
}

func (_c *MockClient_RunPipeline_Call) Return(_a0 *pipelines.Run, _a1 error) *MockClient_RunPipeline_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_RunPipeline_Call) RunAndReturn(run func(context.Context, pipelines.RunPipelineArgs) (*pipelines.Run, error)) *MockClient_RunPipeline_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockClient creates a new instance of MockClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockClient {
	mock := &MockClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
# (2025-03-04)