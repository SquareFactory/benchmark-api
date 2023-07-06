package mocks

import (
	context "context"

	"github.com/squarefactory/benchmark-api/scheduler"
	mock "github.com/stretchr/testify/mock"
)

type Scheduler struct {
	mock.Mock
}

func (_m *Scheduler) Submit(ctx context.Context, req *scheduler.SubmitRequest) (string, error) {
	args := _m.Called(ctx, req)

	if rf, ok := args.Get(0).(func(context.Context, *scheduler.SubmitRequest) (string, error)); ok {
		return rf(ctx, req)
	}

	if rf, ok := args.Get(0).(func(context.Context, *scheduler.SubmitRequest) string); ok {
		return rf(ctx, req), nil
	}

	if rf, ok := args.Get(1).(error); ok {
		return "", rf
	}

	return "", nil
}

func (_m *Scheduler) CancelJob(ctx context.Context, req *scheduler.CancelRequest) error {
	args := _m.Called(ctx, req)

	if rf, ok := args.Get(0).(func(context.Context, *scheduler.CancelRequest) error); ok {
		return rf(ctx, req)
	}

	if rf, ok := args.Get(0).(func(context.Context, *scheduler.CancelRequest)); ok {
		rf(ctx, req)
		return nil
	}

	if rf, ok := args.Get(1).(error); ok {
		return rf
	}

	return nil
}

func (_m *Scheduler) HealthCheck(ctx context.Context) error {
	args := _m.Called(ctx)

	if rf, ok := args.Get(0).(func(context.Context) error); ok {
		return rf(ctx)
	}

	if rf, ok := args.Get(0).(func(context.Context)); ok {
		rf(ctx)
		return nil
	}

	if rf, ok := args.Get(1).(error); ok {
		return rf
	}

	return nil
}

func (_m *Scheduler) FindRunningJobByName(
	ctx context.Context,
	req *scheduler.FindRunningJobByNameRequest,
) (int, error) {
	args := _m.Called(ctx, req)

	if rf, ok := args.Get(0).(func(context.Context, *scheduler.FindRunningJobByNameRequest) (int, error)); ok {
		return rf(ctx, req)
	}

	if rf, ok := args.Get(0).(func(context.Context, *scheduler.FindRunningJobByNameRequest) int); ok {
		return rf(ctx, req), nil
	}

	if rf, ok := args.Get(1).(error); ok {
		return 0, rf
	}

	return 0, nil
}

func (_m *Scheduler) FindMemPerNode(ctx context.Context) (int, error) {
	args := _m.Called(ctx)

	if rf, ok := args.Get(0).(func(context.Context) (int, error)); ok {
		return rf(ctx)
	}

	if rf, ok := args.Get(0).(int); ok {
		return rf, nil
	}

	if rf, ok := args.Get(1).(error); ok {
		return 0, rf
	}

	return 0, nil
}

func (_m *Scheduler) FindGPUPerNode(ctx context.Context) (int, error) {
	args := _m.Called(ctx)

	if rf, ok := args.Get(0).(func(context.Context) (int, error)); ok {
		return rf(ctx)
	}

	if rf, ok := args.Get(0).(int); ok {
		return rf, nil
	}

	if rf, ok := args.Get(1).(error); ok {
		return 0, rf
	}

	return 0, nil
}

func (_m *Scheduler) FindCPUPerNode(ctx context.Context) (int, error) {
	args := _m.Called(ctx)

	if rf, ok := args.Get(0).(func(context.Context) (int, error)); ok {
		return rf(ctx)
	}

	if rf, ok := args.Get(0).(func(context.Context) int); ok {
		return rf(ctx), nil
	}

	if rf, ok := args.Get(1).(error); ok {
		return 0, rf
	}

	return 0, nil
}

type mockConstructorTestingTNewScheduler interface {
	mock.TestingT
	Cleanup(func())
}

// NewExecutor creates a new instance of Executor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewScheduler(t mockConstructorTestingTNewScheduler) *Scheduler {
	mock := &Scheduler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
