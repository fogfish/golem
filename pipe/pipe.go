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
		defer func() {
			// Note: recover from panic on sending to closed channel
			recover()
		}()
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
func From[T any](n int, frequency time.Duration, f func() (T, error)) chan T {
	eg := make(chan T, n)

	go func() {
		defer func() {
			// Note: recover from panic on sending to closed channel
			recover()
		}()

		for {
			time.Sleep(frequency)
			if x, err := f(); err == nil {
				eg <- x
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
		defer func() {
			// Note: recover from panic on sending to closed channel
			recover()
		}()
		defer close(eg)

		var x A
		for x = range in {
			eg <- f(x)
		}
	}()

	return eg
}

/*

MaybeMap channel type
*/
func MaybeMap[A, B any](in <-chan A, f func(A) (B, error)) chan B {
	eg := make(chan B, cap(in))

	go func() {
		defer func() {
			// Note: recover from panic on sending to closed channel
			recover()
		}()
		defer close(eg)

		var x A
		for x = range in {
			if y, err := f(x); err == nil {
				eg <- y
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
		var x A
		for x = range in {
			f(x)
		}
	}()
}
