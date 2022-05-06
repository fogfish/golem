//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package pipe

import (
	"time"
)

/*

New creates unbounded channel

  in, eg := pipe.New()
  ...
  close(eg)
*/
func New[T any](n int) (<-chan T, chan<- T) {
	in := make(chan T, n)
	eg := make(chan T, n)

	mq := newq[T]()

	go func() {
		defer close(eg)

		for {
			select {
			case x, ok := <-in:
				if !ok {
					return
				}
				enq(&x, mq)

			case emit(eg, mq) <- head(mq):
				deq(mq)
			}
		}
	}()

	return eg, in
}

/*

From creates a channel periodically generates values from the function
*/
func From[T any](n int, frequency time.Duration, f func() T) chan T {
	eg := make(chan T, n)

	go func() {
		defer func() {
			// Note: recover from panic on sending to closed channel
			if recover() != nil {
			}
		}()

		for {
			select {
			case <-time.After(frequency):
				eg <- f()
			}
		}
	}()

	return eg
}

/*

Map channel type
*/
func Map[A, B any](in <-chan A, f func(A) B) chan B {
	eg := make(chan B, cap(in))

	go func() {
		defer close(eg)
		defer func() {
			// Note: recover from panic on sending to closed channel
			if recover() != nil {
			}
		}()

		var (
			x  A
			ok bool
		)

		for {
			select {
			case x, ok = <-in:
				if !ok {
					return
				}
				eg <- f(x)
			}
		}
	}()

	return eg
}

/*

ForEach applies function for each message in the channel
*/
func ForEach[A any](in <-chan A, f func(A)) {
	go func() {
		var (
			x  A
			ok bool
		)

		for {
			select {
			case x, ok = <-in:
				if !ok {
					return
				}
				f(x)
			}
		}
	}()
}
