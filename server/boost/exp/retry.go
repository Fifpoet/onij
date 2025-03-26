package exp

import (
	"context"
	"math"
	"time"

	"github.com/pkg/errors"
)

type OnRetryFunc func(attempted int, err error)
type RetryAction func(context.Context) error
type RetryFunc[R any] func(context.Context) (R, error)
type RetryFunc2[T, R any] func(context.Context, T) (R, error)

const (
	defaultRetryDelay    = 500 * time.Millisecond
	defaultRetryAttempts = 2
)

type RetryOptions struct {
	attempts int           // 重试次数(<=0为不重试)
	delay    time.Duration // 重试间隔
	backoff  bool          // 是否开启指数退避
	onRetry  OnRetryFunc   // 重试回调, 重试前触发
}

type RetryOption func(*RetryOptions)

func WithRetryOptionsAttempts(attempts int) RetryOption {
	return func(o *RetryOptions) { o.attempts = attempts }
}

func WithRetryOptionsDelay(delay time.Duration) RetryOption {
	return func(o *RetryOptions) { o.delay = delay }
}

func WithRetryOptionsBackoff() RetryOption {
	return func(o *RetryOptions) { o.backoff = true }
}

func WithRetryOptionsOnRetry(fn OnRetryFunc) RetryOption {
	return func(o *RetryOptions) { o.onRetry = fn }
}

// Retry: 默认固定500ms间隔重试2次; 调用方通过ctx控制取消或超时;
func Retry(ctx context.Context, action RetryAction, opts ...RetryOption) error {
	if action == nil {
		return nil
	}

	wrapper := func(ctx context.Context, _ struct{}) (struct{}, error) {
		return struct{}{}, action(ctx)
	}
	if _, err := RetryTR(ctx, struct{}{}, wrapper, opts...); err != nil {
		return err
	}
	return nil
}

// RetryR: 默认固定500ms间隔重试2次; 调用方通过ctx控制取消或超时;
func RetryR[R any](ctx context.Context, fn RetryFunc[R], opts ...RetryOption) (r R, err error) {
	if fn == nil {
		return
	}

	wrapper := func(ctx context.Context, _ struct{}) (R, error) {
		return fn(ctx)
	}
	if _, err = RetryTR(ctx, struct{}{}, wrapper, opts...); err != nil {
		return
	}
	return
}

// RetryTR: 默认固定500ms间隔重试2次; 调用方通过ctx控制取消或超时;
func RetryTR[T, R any](ctx context.Context, t T, fn RetryFunc2[T, R], opts ...RetryOption) (result R, err error) {
	if fn == nil {
		return
	}

	options := &RetryOptions{
		attempts: defaultRetryAttempts,
		delay:    defaultRetryDelay,
	}
	for _, opt := range opts {
		opt(options)
	}
	if options.attempts < 0 {
		options.attempts = 0
	}
	if options.delay < 0 {
		options.delay = 0
	}

	for attempted := 0; ; attempted++ {
		if err = ctx.Err(); err != nil {
			return
		}

		attempt := options.attempts > 0 && attempted < options.attempts

		var panicErr error
		func() {
			defer func() {
				if msg := recover(); msg != nil {
					panicErr = errors.Errorf("panic = %v", msg)
				}
			}()
			if result, err = fn(ctx, t); err != nil && attempt {
				if options.onRetry != nil {
					options.onRetry(attempted, err)
				}
			}
		}()

		if panicErr != nil {
			err = panicErr
			return
		}
		if err == nil {
			return
		}
		if !attempt {
			return
		}

		delay := options.delay
		if options.backoff {
			delay = options.delay * time.Duration(math.Pow(2, float64(attempted)))
		}
		select {
		case <-ctx.Done():
			return Zero[R](), ctx.Err()
		case <-time.After(delay):
		}
	}
}
