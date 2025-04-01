//
// Copyright (C) 2022 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

// Package pipe simplify the creation of streaming data pipelines using
// sequential channel workers. The package is semantically compatible with fork.
package pipe

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/fogfish/golem/pure/monoid"
)

// Emit creates a channel and takes a function that emits data at a specified frequency.
func Emit[T any](ctx context.Context, cap int, frequency time.Duration, f F[int, T]) (<-chan T, <-chan error) {
	out := make(chan T, cap)
	exx := f.errch(cap)

	go func() {
		defer close(out)
		defer close(exx)

		var (
			val T
			err error
		)

		for i := 0; true; i++ {
			time.Sleep(frequency)

			val, err = f.apply(i)
			if err != nil {
				if !f.catch(ctx, err, exx) {
					return
				}
				continue
			}

			select {
			case out <- val:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, exx
}

// Filter returns a newly-allocated channel that contains only those elements x
// of the input channel for which predicate is true.
func Filter[A any](ctx context.Context, in <-chan A, f F[A, bool]) <-chan A {
	out := make(chan A, cap(in))

	go func() {
		defer close(out)

		var a A
		for a = range in {
			if take, err := f.apply(a); take && err == nil {
				select {
				case out <- a:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}

// ForEach applies function for each message in the channel
func ForEach[A any](ctx context.Context, in <-chan A, f F[A, A]) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		var x A
		for x = range in {
			f.apply(x)
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return done
}

// FMap applies function over channel messages, flatten the output channel and
// emits it result to new channel.
func FMap[A, B any](ctx context.Context, in <-chan A, fmap FF[A, B]) (<-chan B, <-chan error) {
	out := make(chan B, cap(in))
	exx := fmap.errch(cap(in))

	go func() {
		defer close(out)
		defer close(exx)

		var a A
		for a = range in {
			if err := fmap.apply(ctx, a, out); err != nil {
				if !fmap.catch(ctx, err, exx) {
					return
				}
				continue
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return out, exx
}

// Fold applies a monoid operation to the values in a channel. The final value is
// emitted though return channel when the end of the input channel is reached.
func Fold[A any](ctx context.Context, in <-chan A, m monoid.Monoid[A]) <-chan A {
	done := make(chan A, 1)

	go func() {
		acc := m.Empty()

		defer func() {
			done <- acc
			close(done)
		}()

		var x A
		for x = range in {
			acc = m.Combine(acc, x)
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return done
}

// Map applies function over channel messages, emits result to new channel
func Map[A, B any](ctx context.Context, in <-chan A, f F[A, B]) (<-chan B, <-chan error) {
	out := make(chan B, cap(in))
	exx := f.errch(cap(in))

	go func() {
		defer close(out)
		defer close(exx)

		var (
			a   A
			val B
			err error
		)

		for a = range in {
			val, err = f.apply(a)
			if err != nil {
				if !f.catch(ctx, err, exx) {
					return
				}
				continue
			}

			select {
			case out <- val:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, exx
}

// Partition channel into two channels according to predicate
func Partition[A any](ctx context.Context, in <-chan A, f F[A, bool]) (<-chan A, <-chan A) {
	lout := make(chan A, cap(in))
	rout := make(chan A, cap(in))

	go func() {
		defer close(rout)
		defer close(lout)

		sel := func(x bool, err error) chan<- A {
			if x && err == nil {
				return lout
			}
			return rout
		}

		var a A
		for a = range in {
			select {
			case sel(f.apply(a)) <- a:
			case <-ctx.Done():
				return
			}
		}
	}()

	return lout, rout
}

// Unfold is the fundamental recursive constructor, it applies a function to
// each previous seed element in turn to determine the next element.
func Unfold[A any](ctx context.Context, cap int, seed A, f F[A, A]) (<-chan A, <-chan error) {
	out := make(chan A, cap)
	exx := f.errch(cap)

	go func() {
		defer close(out)
		defer close(exx)

		var err error
		for {
			select {
			case out <- seed:
			case <-ctx.Done():
				return
			}

			seed, err = f.apply(seed)
			if err != nil {
				if !f.catch(ctx, err, exx) {
					return
				}
				continue
			}
		}
	}()

	return out, exx
}

// Join concatenate channels, returns newly-allocated channel composed of
// elements copied from input channels.
func Join[A any](ctx context.Context, in ...<-chan A) <-chan A {
	var wg sync.WaitGroup
	out := make(chan A, len(in))

	join := func(c <-chan A) {
		defer wg.Done()

		for x := range c {
			select {
			case out <- x:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(len(in))
	for _, c := range in {
		go join(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// returns a newly-allocated channel containing the first n elements of the input channel.
func Take[A any](ctx context.Context, in <-chan A, n int) <-chan A {
	out := make(chan A, cap(in))

	go func() {
		defer close(out)

		var a A
		for a = range in {

			select {
			case out <- a:
			case <-ctx.Done():
				return
			}

			n--
			if n == 0 {
				return
			}
		}
	}()

	return out
}

// Filter returns a newly-allocated channel that contains only those elements x
// of the input channel for which predicate is true.
func TakeWhile[A any](ctx context.Context, in <-chan A, f F[A, bool]) <-chan A {
	out := make(chan A, cap(in))

	go func() {
		defer close(out)

		var a A
		for a = range in {
			if take, err := f.apply(a); !take || err != nil {
				return
			}

			select {
			case out <- a:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// Throttling the channel to ops per time interval.
func Throttling[A any](ctx context.Context, in <-chan A, ops int, interval time.Duration) <-chan A {
	out := make(chan A, cap(in))
	ctl := make(chan struct{}, ops)

	go func() {
		defer close(ctl)

		for {
			for i := 0; i < ops; i++ {
				select {
				case ctl <- struct{}{}:
				case <-ctx.Done():
					return
				}
			}
			select {
			case <-time.After(interval):
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer close(out)

		var a A
		for a = range in {
			select {
			case <-ctl:
			case <-ctx.Done():
				return
			}

			select {
			case out <- a:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// Lift sequence of values into channel
func Seq[T any](xs ...T) <-chan T {
	out := make(chan T, len(xs))
	for _, x := range xs {
		out <- x
	}
	close(out)
	return out
}

// ToSeq collects elements of channel into the sequence (slice)
func ToSeq[T any](ch <-chan T) []T {
	seq := make([]T, 0)
	for x := range ch {
		seq = append(seq, x)
	}
	return seq
}

// Standard error logger
func StdErr[T any](out <-chan T, exx <-chan error) <-chan T {
	go func() {
		var err error
		for err = range exx {
			if err != nil {
				slog.Error("pipe stage failed.", "error", err)
			}
		}
	}()

	return out
}
