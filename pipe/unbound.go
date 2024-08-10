//
// Copyright (C) 2022 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package pipe

import "context"

// New creates an unbounded channel.
// By default, Go channels have a fixed capacity, which can cause the producer
// or consumer to block if the capacity is exceeded. This function returns a
// pair of channels backed by an in-memory queue, simulating an "unbounded" channel.
//
// However, it is important to note that while unbounded channels avoid the
// blocking issues of bounded channels, they come with their own trade-offs.
// Specifically, they can lead to uncontrolled memory usage if the consumer is
// too slow relative to the producer, as the in-memory queue backing the channel
// can grow indefinitely. Therefore, unbounded channels should be used carefully,
// considering the memory implications.
//
//	ctx, close := context.WithCancel(context.Background())
//	rcv, snd := pipe.New(ctx, 0)
//	...
//	close()
func New[T any](ctx context.Context, cap int) (<-chan T, chan<- T) {
	eg := make(chan T, cap)
	in := make(chan T, cap)

	mq := newq[T]()

	go func() {
		defer close(eg)
		defer close(in)

		for {
			select {
			case <-ctx.Done():
				for mq.head != nil {
					eg <- head(mq)
					deq(mq)
				}
				return

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
