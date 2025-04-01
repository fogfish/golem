//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package pipe_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/fogfish/golem/pipe/v2"
	"github.com/fogfish/it/v2"
)

func TestPipeNew(t *testing.T) {
	t.Run("Empty.Recv", func(t *testing.T) {
		ctx, close := context.WithCancel(context.Background())
		in, _ := pipe.New[int](ctx, 0)
		select {
		case <-in:
			panic("Must Not Recv Anything")
		case <-time.After(10 * time.Millisecond):
		}
		close()
	})

	t.Run("Empty.Send", func(t *testing.T) {
		ctx, close := context.WithCancel(context.Background())
		_, eg := pipe.New[int](ctx, 0)
		select {
		case eg <- 0:
			break
		case <-time.After(10 * time.Millisecond):
			panic("Must Not Blocked")
		}
		close()
	})

	t.Run("Send.Recv", func(t *testing.T) {
		ctx, close := context.WithCancel(context.Background())
		in, eg := pipe.New[int](ctx, 0)
		eg <- 100
		it.Then(t).Should(
			it.Equal(<-in, 100),
		)
		close()
	})

	t.Run("Send.Batch.Recv", func(t *testing.T) {
		ctx, close := context.WithCancel(context.Background())
		in, eg := pipe.New[int](ctx, 0)
		for i := 0; i < 1000; i++ {
			eg <- i
		}
		for i := 0; i < 1000; i++ {
			it.Then(t).Should(
				it.Equal(<-in, i),
			)
		}
		close()
	})

	t.Run("Recv.Async.Send", func(t *testing.T) {
		ctx, close := context.WithCancel(context.Background())
		in, eg := pipe.New[int](ctx, 0)
		go func() {
			for i := 0; i < 1000; i++ {
				time.Sleep(time.Duration(rand.Intn(4)) * time.Millisecond)
				it.Then(t).Should(
					it.Equal(<-in, i),
				)
			}
		}()
		for i := 0; i < 1000; i++ {
			eg <- i
			time.Sleep(time.Duration(rand.Intn(9)) * time.Millisecond)
		}
		close()
	})

	t.Run("Recv.After.Close", func(t *testing.T) {
		ctx, close := context.WithCancel(context.Background())
		in, eg := pipe.New[int](ctx, 0)
		for i := 0; i < 1000; i++ {
			eg <- i
		}
		close()

		for i := 0; i < 1000; i++ {
			it.Then(t).Should(
				it.Equal(<-in, i),
			)
		}
	})
}
