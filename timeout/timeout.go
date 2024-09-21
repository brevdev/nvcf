package timeout

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

// https://github.com/zeromicro/go-zero/blob/master/core/fx/timeout.go

var (
	// ErrCanceled is the error returned when the context is canceled.
	ErrCanceled = context.Canceled
	// ErrTimeout is the error returned when the context's deadline passes.
	ErrTimeout = context.DeadlineExceeded
)

// DoOption defines the method to customize a DoWithTimeout call.
type DoOption func() context.Context

type Options struct {
	ParentContext context.Context
	CatchPanic    bool
}

// if you loop forever, make sure you have a way to break the loop
// see Test_DoWithTimeoutTimeoutLoop
func DoWithTimeout(fn func(ctx context.Context) error, timeout time.Duration, opts ...Options) error {
	_, err := DoWithTimeoutData(func(ctx context.Context) (interface{}, error) {
		return nil, fn(ctx)
	}, timeout, opts...)
	return err
}

// if you loop forever, make sure you have a way to break the loop
// see Test_DoWithTimeoutTimeoutLoop
func DoWithTimeoutData[T any](fn func(ctx context.Context) (T, error), timeout time.Duration, opts ...Options) (T, error) {
	type result struct {
		res T
		err error
	}
	options := Options{
		ParentContext: context.Background(),
	}
	for _, opt := range opts {
		options = opt
	}
	if options.ParentContext == nil {
		options.ParentContext = context.Background()
	}
	ctx, cancel := context.WithTimeout(options.ParentContext, timeout)
	defer cancel()

	// create channel with buffer size 1 to avoid goroutine leak
	resChan := make(chan result, 1)
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// attach call stack to avoid missing in different goroutine
				panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
			}
		}()
		res, err := fn(ctx)
		resChan <- result{res, err}
	}()

	var emptyT T

	select {
	case p := <-panicChan:
		if options.CatchPanic {
			return emptyT, fmt.Errorf("panic: %v", p)
		} else {
			panic(p)
		}
	case result := <-resChan:
		return result.res, result.err
	case <-ctx.Done():
		return emptyT, ctx.Err() //nolint:wrapcheck // no need to wrap
	}
}

// WithContext customizes a DoWithTimeout call with given ctx.
func WithContext(ctx context.Context) DoOption {
	return func() context.Context {
		return ctx
	}
}
