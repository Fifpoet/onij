package concurrent

import (
	"context"
	"errors"
	"onij/util/boost/collection/collext"
	"onij/util/boost/exp"
	"reflect"
	"sync"

	errext "github.com/pkg/errors"
)

type Task[T, R any] interface {
	Cancel() bool
	IsCancelled() bool
	Result() (R, error)
	Wait() <-chan struct{}
	WaitResult() <-chan *TaskResult[R]
}

type task[T, R any] struct {
	context.Context

	data   T
	action TaskAction[T, R]
	hooks  *TaskHooks[T, R]
	result *TaskResult[R]
	cancel context.CancelFunc
	donec  chan struct{}
	close  func()
}

func newTask[T, R any](ctx context.Context, t T, f TaskAction[T, R], opts []TaskHookOption[T, R]) *task[T, R] {
	ctxTask, cancel := context.WithCancel(ctx)
	donec := make(chan struct{})
	hooks := new(TaskHooks[T, R])
	for _, opt := range opts {
		opt(hooks)
	}
	return &task[T, R]{
		Context: ctxTask,
		cancel:  cancel,
		data:    t,
		action:  f,
		hooks:   hooks,
		donec:   donec,
		result:  new(TaskResult[R]),
		close:   sync.OnceFunc(func() { close(donec) }),
	}
}

func (r *task[T, R]) Cancel() bool {
	select {
	case <-r.donec:
		return false
	default:
	}
	select {
	case <-r.Context.Done():
		return true
	default:
	}

	r.cancel()
	return true
}

func (r *task[T, R]) IsCancelled() bool {
	select {
	case <-r.Context.Done():
		return true
	default:
		return false
	}
}

func (r *task[T, R]) Result() (R, error) {
	select {
	case <-r.Context.Done():
		return exp.Zero[R](), r.result.err
	default:
	}
	<-r.donec
	return r.result.r, r.result.err
}

func (r *task[T, R]) Wait() <-chan struct{} { return r.donec }

func (r *task[T, R]) WaitResult() <-chan *TaskResult[R] {
	ch := make(chan *TaskResult[R])
	go func() {
		defer close(ch)
		<-r.donec
		ch <- r.result
	}()
	return ch
}

func (r *task[T, R]) run() <-chan struct{} {
	ch := make(chan *TaskResult[R])
	go func() {
		defer close(ch)
		ch <- r.exec()
	}()
	go func() {
		defer r.close()
		select {
		case <-r.Context.Done():
			r.result.err = r.Context.Err()
		case r.result = <-ch:
		}
	}()
	return r.donec
}

func (r *task[T, R]) exec() (tr *TaskResult[R]) {
	defer func() {
		if err := recover(); err != nil {
			tr.err = errext.Errorf("panic = %v", err)
		}
	}()

	tr = new(TaskResult[R])
	tr.r, tr.err = r.action(r.Context, r.data)
	return
}

func (r *task[T, R]) discard() {
	defer r.close()

	r.result.err = ErrTaskDiscard

	if r.hooks.Discard != nil {
		r.hooks.Discard(r.Context, r, r.data)
	}

	r.Cancel()
}

type TaskResults[R any] []*TaskResult[R]

type TaskResult[R any] struct {
	err error
	r   R
}

func (r TaskResults[R]) Result() ([]R, error) {
	results := make([]R, 0, len(r))
	var firstErr error
	for _, v := range r {
		result, err := v.Result()
		results = append(results, result)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return results, firstErr
}

func (r TaskResults[R]) Err() error {
	for _, v := range r {
		if err := v.Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r TaskResults[R]) AggregateErr() error {
	var errs []error
	for _, v := range r {
		if err := v.Err(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (r *TaskResult[R]) Result() (R, error) { return r.r, r.err }

func (r *TaskResult[R]) Err() error { return r.err }

type waitOr func(chs ...<-chan struct{}) <-chan struct{}

func Wait[T, R any](tasks ...Task[T, R]) TaskResults[R] {
	results := make(TaskResults[R], 0, len(tasks))
	for _, v := range tasks {
		r, err := v.Result()
		results = append(results, &TaskResult[R]{r: r, err: err})
	}
	return results
}

func FastTouch[T, R any](tasks ...Task[T, R]) {
	var or waitOr
	or = func(chs ...<-chan struct{}) <-chan struct{} {
		if len(chs) == 0 {
			return nil
		}
		if len(chs) == 1 {
			return chs[0]
		}
		donec := make(chan struct{})
		go func() {
			defer close(donec)
			switch len(chs) {
			case 2:
				select {
				case <-chs[0]:
				case <-chs[1]:
				}
			default:
				select {
				case <-chs[0]:
				case <-chs[1]:
				case <-chs[2]:
				case <-or(append(chs[3:], donec)...):
				}
			}
		}()
		return donec
	}

	chs := collext.Pick(tasks, func(t Task[T, R]) <-chan struct{} { return t.Wait() })
	<-or(chs...)
}

func wait[T, R any](fastFail, fastSuccess bool, tasks []Task[T, R]) TaskResults[R] {
	if !fastFail && !fastSuccess {
		return Wait(tasks...)
	}

	results := make(TaskResults[R], 0, len(tasks))
	chs := collext.Pick(tasks, func(t Task[T, R]) <-chan *TaskResult[R] { return t.WaitResult() })
	cases := make([]reflect.SelectCase, len(chs))
	for i, v := range chs {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(v),
		}
	}

	for len(cases) != 0 {
		chosen, recv, _ := reflect.Select(cases)
		r, _ := recv.Interface().(*TaskResult[R])
		var ok bool
		if fastFail && fastSuccess {
			ok = true
		} else if fastFail {
			ok = r.err != nil
		} else if fastSuccess {
			ok = r.err == nil
		}
		if ok {
			return TaskResults[R]{r}
		} else {
			results = append(results, r)
		}

		cases[chosen] = cases[len(cases)-1]
		cases = cases[:len(cases)-1]
	}
	return results
}
