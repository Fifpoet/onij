package concurrent

import (
	"context"
	"errors"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	ErrLockerInvalidMutex    = errors.New("locker: invalid mutex")
	ErrLockerInvalidScope    = errors.New("locker: invalid scope")
	ErrLockerLockUnlocked    = errors.New("locker: lock unlocked")
	ErrLockerUnlockNonLocked = errors.New("locker: unlock non locked")
)

type LockOptions struct {
	lease         time.Duration
	maxLease      time.Duration
	renewInterval time.Duration

	renewFailedHook func(error)
}

func (r *LockOptions) normalize() {
	if r.lease <= 0 {
		r.lease = 5 * time.Second
	}
	if r.renewInterval <= 0 {
		r.renewInterval = time.Second
	}
}

type LockOption func(*LockOptions)

func LockWithLease(lease time.Duration) LockOption {
	return func(o *LockOptions) { o.lease = lease }
}

func LockWithMaxLease(max time.Duration) LockOption {
	return func(o *LockOptions) { o.maxLease = max }
}

func LockWithRenewInterval(interval time.Duration) LockOption {
	return func(o *LockOptions) { o.renewInterval = interval }
}

func LockWithRenewFailedHook(hook func(error)) LockOption {
	return func(o *LockOptions) { o.renewFailedHook = hook }
}

type LeaseMutex interface {
	Acquire(ctx context.Context, scope, token string, lease time.Duration) error
	Release(ctx context.Context, scope, token string) error
	Renew(ctx context.Context, scope, token string, lease time.Duration) error
}

type Locker interface {
	Lock() (bool, error)
	Unlock() (bool, error)
}

type locker struct {
	context  context.Context
	scope    string
	token    string
	mutex    LeaseMutex
	acquired atomic.Bool
	released atomic.Bool
	cancel   context.CancelFunc

	options *LockOptions
}

func NewLocker(ctx context.Context, mutex LeaseMutex, scope string, opts ...LockOption) (Locker, error) {
	if mutex == nil {
		return nil, ErrLockerInvalidMutex
	}
	if len(scope) == 0 {
		return nil, ErrLockerInvalidScope
	}

	options := new(LockOptions)
	for _, opt := range opts {
		opt(options)
	}

	options.normalize()
	context, cancel := context.WithCancel(ctx)

	r := &locker{
		context: context,
		cancel:  cancel,
		scope:   scope,
		token:   strconv.FormatInt(time.Now().UnixNano(), 16),
		mutex:   mutex,
		options: options,
	}

	return r, nil
}

func (r *locker) Lock() (bool, error) {
	if r.released.Load() {
		return false, ErrLockerLockUnlocked
	}
	if r.acquired.Load() {
		return true, nil
	}
	if !r.acquired.CompareAndSwap(false, true) {
		return false, nil
	}
	if err := r.mutex.Acquire(r.context, r.scope, r.token, r.options.lease); err != nil {
		return false, err
	}

	go r.monitor()

	return true, nil
}

func (r *locker) Unlock() (bool, error) {
	if !r.acquired.Load() {
		return false, ErrLockerUnlockNonLocked
	}
	if r.released.Load() {
		return true, nil
	}
	if !r.released.CompareAndSwap(false, true) {
		return false, nil
	}
	if r.cancel != nil {
		defer r.cancel()
	}
	if err := r.mutex.Release(r.context, r.scope, r.token); err != nil {
		return false, err
	}
	return true, nil
}

func (r *locker) monitor() {
	ticker := time.NewTicker(r.options.renewInterval)
	start := time.Now()

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if r.released.Load() {
				return
			}
			if r.options.maxLease != 0 && time.Since(start) >= r.options.maxLease {
				return
			}
			var err error
			for i := 0; i < 2; i++ {
				if i != 0 {
					time.Sleep(time.Millisecond * 100)
				}
				if err = r.mutex.Renew(r.context, r.scope, r.token, r.options.lease); err == nil {
					break
				}
			}
			if err != nil {
				if r.options.renewFailedHook != nil {
					r.options.renewFailedHook(err)
				}
				return
			}
		case <-r.context.Done():
			return
		}
	}
}
