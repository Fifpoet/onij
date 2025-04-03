package concurrent

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

var (
	ErrSchedulerShutdown = errors.New("scheduler shutdown")
	ErrTaskDiscard       = errors.New("task discard")
)

type Scheduler[T, R any] interface {
	Task(context.Context, T, TaskAction[T, R], ...TaskHookOption[T, R]) (Task[T, R], error)
	Shutdown() <-chan struct{}
}

type SchedulerOptions struct {
	MaxWorkers         int           // default 10
	BufferSize         int           // default 5
	TaskTimeout        time.Duration // default 0(no timeout)
	DispatchConstraint bool          // default false, true-discard, false-wait dispatch
	WorkerConstraint   bool          // default true, true-discard, false-wait worker
}

type TaskAction[T, R any] func(context.Context, T) (R, error)

type SchedulerOption func(*SchedulerOptions)

func WithSchedulerOptionMaxWorkers(maxWorkers int) SchedulerOption {
	return func(o *SchedulerOptions) { o.MaxWorkers = maxWorkers }
}

func WithSchedulerOptionBufferSize(size int) SchedulerOption {
	return func(o *SchedulerOptions) { o.BufferSize = size }
}

func WithSchedulerOptionTaskTimeout(timeout time.Duration) SchedulerOption {
	return func(o *SchedulerOptions) { o.TaskTimeout = timeout }
}

func WithSchedulerOptionDispatchConstraint(constraint bool) SchedulerOption {
	return func(o *SchedulerOptions) { o.DispatchConstraint = constraint }
}

func WithSchedulerOptionWorkerConstraint(constraint bool) SchedulerOption {
	return func(o *SchedulerOptions) { o.WorkerConstraint = constraint }
}

type TaskHook[T, R any] func(context.Context, Task[T, R], T)

type TaskHooks[T, R any] struct {
	Discard TaskHook[T, R]
}

type TaskHookOption[T, R any] func(*TaskHooks[T, R])

func WithTaskDiscardHook[T, R any](hook TaskHook[T, R]) TaskHookOption[T, R] {
	return func(h *TaskHooks[T, R]) { h.Discard = hook }
}

func (r *SchedulerOptions) normalize() {
	if r.MaxWorkers <= 0 {
		r.MaxWorkers = 10
	}
	if r.BufferSize < 0 {
		r.BufferSize = 5
	}
}

func NewScheduler[T, R any](opts ...SchedulerOption) Scheduler[T, R] {
	o := &SchedulerOptions{WorkerConstraint: true}
	for _, opt := range opts {
		opt(o)
	}
	o.normalize()

	scheduler := &scheduler[T, R]{
		options:  o,
		stopc:    make(chan struct{}),
		donec:    make(chan struct{}),
		taskPool: make(chan *task[T, R], o.BufferSize),
	}

	scheduler.workerPool = newWorkerPool[T, R](o.MaxWorkers, !o.WorkerConstraint, scheduler.stopc)

	scheduler.run()

	return scheduler
}

type scheduler[T, R any] struct {
	options *SchedulerOptions

	workerPool *workerPool[T, R]
	taskPool   chan *task[T, R]
	stopc      chan struct{}
	donec      chan struct{}
	shutdown   atomic.Bool
}

func (r *scheduler[T, R]) Shutdown() <-chan struct{} {
	if !r.shutdown.CompareAndSwap(false, true) {
		return r.donec
	}

	close(r.stopc)
	close(r.taskPool)

	return r.donec
}

func (r *scheduler[T, R]) Task(ctx context.Context, data T, action TaskAction[T, R], opts ...TaskHookOption[T, R]) (task Task[T, R], err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Task panic, msg = %v", e)
		}
	}()

	select {
	case <-r.stopc:
		return nil, ErrSchedulerShutdown
	default:
	}

	t := newTask(ctx, data, action, opts)
	if r.options.DispatchConstraint {
		select {
		case <-r.stopc:
			return nil, ErrSchedulerShutdown
		default:
			select {
			case r.taskPool <- t:
			default:
				return nil, ErrTaskDiscard
			}
		}
	} else {
		select {
		case <-r.stopc:
			return nil, ErrSchedulerShutdown
		default:
			r.taskPool <- t
		}
	}
	return t, nil
}

func (r *scheduler[T, R]) run() {
	go func() {
		defer close(r.donec)
		select {
		case <-r.stopc:
			return
		default:
		}
		for task := range r.taskPool {
			var discard bool
			select {
			case <-r.stopc:
				discard = true
			default:
				discard = !r.workerPool.Serve(task)
			}
			if discard {
				task.discard()
			}
		}
	}()
}
