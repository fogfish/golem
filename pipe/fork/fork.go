//
// Copyright (C) 2022 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

// Package fork simplify the creation of streaming data pipelines using
// parallel channel workers. The package is semantically compatible with pipe.
package fork

import (
	"context"
	"sync"
	"time"

	"github.com/fogfish/golem/pipe"
	"github.com/fogfish/golem/pure/monoid"
)

// Emit creates a channel and takes a function that emits data at a specified frequency.
func Emit[T any](ctx context.Context, cap int, frequency time.Duration, emit func() (T, error)) (<-chan T, <-chan error) {
	return pipe.Emit(ctx, cap, frequency, emit)
}

// Filter returns a newly-allocated channel that contains only those elements x
// of the input channel for which predicate is true.
func Filter[A any](ctx context.Context, par int, in <-chan A, f func(A) bool) <-chan A {
	var wg sync.WaitGroup
	out := make(chan A, par)

	pf := func() {
		defer wg.Done()

		var a A
		for a = range in {
			if f(a) {
				select {
				case out <- a:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	wg.Add(par)
	for i := 1; i <= par; i++ {
		go pf()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ForEach applies function for each message in the channel
func ForEach[A any](ctx context.Context, par int, in <-chan A, f func(A)) <-chan struct{} {
	var wg sync.WaitGroup
	done := make(chan struct{})

	fmap := func() {
		defer wg.Done()

		var a A
		for a = range in {
			f(a)
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}

	wg.Add(par)
	for i := 1; i <= par; i++ {
		go fmap()
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	return done
}

// FMap applies function over channel messages, flatten the output channel and
// emits it result to new channel.
func FMap[A, B any](ctx context.Context, par int, in <-chan A, fmap func(context.Context, A, chan<- B) error) (<-chan B, <-chan error) {
	var wg sync.WaitGroup
	out := make(chan B, par)
	exx := make(chan error, par)

	pmap := func() {
		defer wg.Done()

		var a A
		for a = range in {
			if err := fmap(ctx, a, out); err != nil {
				exx <- err
				return
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}

	wg.Add(par)
	for i := 1; i <= par; i++ {
		go pmap()
	}

	go func() {
		wg.Wait()
		close(out)
		close(exx)
	}()

	return out, exx
}

// Fold applies a monoid operation to the values in a channel. The final value is
// emitted though return channel when the end of the input channel is reached.
func Fold[A any](ctx context.Context, par int, in <-chan A, m monoid.Monoid[A]) <-chan A {
	var wg sync.WaitGroup
	vals := make(chan A, par)
	done := make(chan A, 1)

	pfold := func() {
		acc := m.Empty()

		defer func() {
			vals <- acc
			wg.Done()
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
	}

	wg.Add(par)
	for i := 1; i <= par; i++ {
		go pfold()
	}

	go func() {
		wg.Wait()

		var acc A
		for i := 1; i <= par; i++ {
			acc = m.Combine(acc, <-vals)
		}
		done <- acc
		close(vals)
		close(done)
	}()

	return done
}

// Map applies function over channel messages, emits result to new channel
func Map[A, B any](ctx context.Context, par int, in <-chan A, fmap func(A) (B, error)) (<-chan B, <-chan error) {
	var wg sync.WaitGroup
	out := make(chan B, par)
	exx := make(chan error, par)

	pmap := func() {
		defer wg.Done()

		var (
			a   A
			val B
			err error
		)

		for a = range in {
			val, err = fmap(a)
			if err != nil {
				exx <- err
				return
			}

			select {
			case out <- val:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(par)
	for i := 1; i <= par; i++ {
		go pmap()
	}

	go func() {
		wg.Wait()
		close(out)
		close(exx)
	}()

	return out, exx
}

// Partition channel into two channels according to predicate
func Partition[A any](ctx context.Context, par int, in <-chan A, f func(A) bool) (<-chan A, <-chan A) {
	var wg sync.WaitGroup
	lout := make(chan A, par)
	rout := make(chan A, par)

	pf := func() {
		defer wg.Done()

		sel := func(x bool) chan<- A {
			if x {
				return lout
			}
			return rout
		}

		var a A
		for a = range in {
			select {
			case sel(f(a)) <- a:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(par)
	for i := 1; i <= par; i++ {
		go pf()
	}

	go func() {
		wg.Wait()
		close(lout)
		close(rout)
	}()

	return lout, rout
}

// Unfold is the fundamental recursive constructor, it applies a function to
// each previous seed element in turn to determine the next element.
func Unfold[A any](ctx context.Context, cap int, seed A, f func(A) (A, error)) (<-chan A, <-chan error) {
	return pipe.Unfold(ctx, cap, seed, f)
}

// Join concatenate channels, returns newly-allocated channel composed of
// elements copied from input channels.
func Join[A any](ctx context.Context, in ...<-chan A) <-chan A {
	return pipe.Join(ctx, in...)
}

// returns a newly-allocated channel containing the first n elements of the input channel.
func Take[A any](ctx context.Context, in <-chan A, n int) <-chan A {
	return pipe.Take(ctx, in, n)
}

// Filter returns a newly-allocated channel that contains only those elements x
// of the input channel for which predicate is true.
func TakeWhile[A any](ctx context.Context, in <-chan A, f func(A) bool) <-chan A {
	return pipe.TakeWhile(ctx, in, f)
}

// Throttling the channel to ops per time interval
func Throttling[A any](ctx context.Context, in <-chan A, ops int, interval time.Duration) <-chan A {
	return pipe.Throttling(ctx, in, ops, interval)
}

// Lift sequence of values into channel
func Seq[T any](xs ...T) <-chan T {
	return pipe.Seq(xs...)
}

// ToSeq collects elements of channel into the sequence (slice)
func ToSeq[T any](ch <-chan T) []T {
	return pipe.ToSeq(ch)
}

// Standard error logger
func StdErr[T any](out <-chan T, exx <-chan error) <-chan T {
	return pipe.StdErr(out, exx)
}
