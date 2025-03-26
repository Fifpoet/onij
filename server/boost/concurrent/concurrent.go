package concurrent

import (
	"context"
)

type RunOptions struct {
	FastFail    bool
	FastSuccess bool
	MaxLimit    int
}

type RunOption func(*RunOptions)

func WithRunOptionFastFail(fail bool) RunOption {
	return func(o *RunOptions) {
		o.FastFail = fail
	}
}

func WithRunOptionFastSuccess(success bool) RunOption {
	return func(o *RunOptions) {
		o.FastSuccess = success
	}
}

func WithRunOptionMaxLimit(maxLimit int) RunOption {
	return func(o *RunOptions) {
		o.MaxLimit = maxLimit
	}
}

type GoAction[T any] func(context.Context, T) error

func Go[T, R any](ctx context.Context, items []T, f TaskAction[T, R], opts ...RunOption) ([]R, error) {
	if len(items) == 0 {
		return nil, nil
	}

	o := new(RunOptions)
	for _, opt := range opts {
		opt(o)
	}

	schedulerOpts := []SchedulerOption{
		WithSchedulerOptionMaxWorkers(o.MaxLimit),
		WithSchedulerOptionBufferSize(len(items)),
		WithSchedulerOptionWorkerConstraint(false),
	}
	scheduler := NewScheduler[T, R](schedulerOpts...)

	defer scheduler.Shutdown()

	tasks := make([]Task[T, R], 0, len(items))
	for _, v := range items {
		task, _ := scheduler.Task(ctx, v, f)
		tasks = append(tasks, task)
	}

	return wait(o.FastFail, o.FastSuccess, tasks).Result()
}

func GoGo[T any](ctx context.Context, items []T, f GoAction[T], opts ...RunOption) error {
	if len(items) == 0 {
		return nil
	}

	o := new(RunOptions)
	for _, opt := range opts {
		opt(o)
	}

	schedulerOpts := []SchedulerOption{
		WithSchedulerOptionMaxWorkers(o.MaxLimit),
		WithSchedulerOptionBufferSize(len(items)),
		WithSchedulerOptionWorkerConstraint(false),
	}
	scheduler := NewScheduler[T, struct{}](schedulerOpts...)

	defer scheduler.Shutdown()

	tasks := make([]Task[T, struct{}], 0, len(items))
	wrapper := func(ctx context.Context, t T) (struct{}, error) { return struct{}{}, f(ctx, t) }
	for _, v := range items {
		task, _ := scheduler.Task(ctx, v, wrapper)
		tasks = append(tasks, task)
	}

	return wait(o.FastFail, o.FastSuccess, tasks).Err()
}

func Run(ctx context.Context, actions []func() error, opts ...RunOption) error {
	if len(actions) == 0 {
		return nil
	}

	o := new(RunOptions)
	for _, opt := range opts {
		opt(o)
	}

	schedulerOpts := []SchedulerOption{
		WithSchedulerOptionMaxWorkers(o.MaxLimit),
		WithSchedulerOptionBufferSize(len(actions)),
		WithSchedulerOptionWorkerConstraint(false),
	}
	scheduler := NewScheduler[struct{}, struct{}](schedulerOpts...)

	defer scheduler.Shutdown()

	tasks := make([]Task[struct{}, struct{}], 0, len(actions))
	for _, action := range actions {
		task, _ := scheduler.Task(ctx, struct{}{}, func(ctx context.Context, t struct{}) (struct{}, error) { return struct{}{}, action() })
		tasks = append(tasks, task)
	}

	return wait(o.FastFail, o.FastSuccess, tasks).Err()
}
