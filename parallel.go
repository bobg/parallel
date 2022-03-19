package workers

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Error struct {
	N   int
	Err error
}

func (e Error) Unwrap() error {
	return e.Err
}

func Values[T any](ctx context.Context, n int, f func(context.Context, int) (T, error)) ([]T, error) {
	g, ctx := errgroup.WithContext(ctx)
	result := make([]T, n)

	for i := 0; i < n; i++ {
		i := i // Go loop var pitfall
		g.Go(func() error {
			val, err := f(ctx, i)
			result[i] = val
			if err != nil {
				return Error{N: i, Err: err}
			}
			return nil
		})
	}

	err := g.Wait()
	return result, err
}

// Producers launches n parallel workers each running the function f.
//
// Each worker receives its worker number (in the range 0 through n-1)
// and a callback to use for producing a value.
// If the callback returns an error, the worker should exit with that error.
//
// An error from any worker cancels them all.
//
// The caller gets a callback for consuming the values produced.
// Each call of the callback yields a value, a boolean, and an error.
// If the error is non-nil, it has propagated out from one of the workers and no more data may be consumed.
// Otherwise, if the boolean is false, the workers have completed and there is no more output to consume.
// Otherwise, the value is a valid next value.
//
// Example:
//
//  consume := Producers(ctx, 10, produceVals)
//  for {
//    val, ok, err := consume()
//    if err != nil { panic(err) }
//    if !ok { break }
//    ...handle val...
//  }
func Producers[T any](ctx context.Context, n int, f func(context.Context, int, func(T) error) error) func() (T, bool, error) {
	ch := make(chan T)
	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < n; i++ {
		i := i
		g.Go(func() error {
			return f(ctx, i, func(val T) error {
				select {
				case <-ctx.Done():
					return Error{N: i, Err: ctx.Err()}
				case ch <- val:
					return nil
				}
			})
		})
	}

	var err error

	go func() {
		err = g.Wait()
		if wasNil {
			close(ch)
		}
	}()

	return func() (T, bool, error) {
		var zero T

		select {
		case <-ctx.Done():
			return zero, false, ctx.Err()
		case val, ok := <-ch:
			if !ok {
				return zero, false, err
			}
			return val, true, nil
		}
	}
}

func Consumers[T any](ctx context.Context, n int, ch chan T, f func(context.Context, int, T) error) (chan<- T, func() error) {
	if ch == nil {
		ch = make(chan T)
	}

	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < n; i++ {
		i := i
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return Error{N: i, Err: ctx.Err()}
				case val, ok := <-ch:
					if !ok {
						return nil
					}
					err := f(ctx, i, val)
					if err != nil {
						return Error{N: i, Err: err}
					}
				}
			}
		})
	}

	closefn := func() error {
		close(ch)
		return g.Wait()
	}

	return ch, closefn
}
