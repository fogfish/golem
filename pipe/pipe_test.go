//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package pipe_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/fogfish/golem/pipe"
	"github.com/fogfish/it"
)

func TestPipeNew(t *testing.T) {
	t.Run("Empty.Recv", func(t *testing.T) {
		in, eg := pipe.New[int](0)
		select {
		case <-in:
			panic("Must Not Recv Anything")
		case <-time.After(10 * time.Millisecond):
		}
		close(eg)
	})

	t.Run("Empty.Send", func(t *testing.T) {
		_, eg := pipe.New[int](0)
		select {
		case eg <- 0:
			break
		case <-time.After(10 * time.Millisecond):
			panic("Must Not Blocked")
		}
		close(eg)
	})

	t.Run("Send.Recv", func(t *testing.T) {
		in, eg := pipe.New[int](0)
		eg <- 100
		it.Ok(t).If(<-in).Equal(100)
		close(eg)
	})

	t.Run("Send.Batch.Recv", func(t *testing.T) {
		in, eg := pipe.New[int](0)
		for i := 0; i < 1000; i++ {
			eg <- i
		}
		for i := 0; i < 1000; i++ {
			it.Ok(t).If(<-in).Equal(i)
		}
		close(eg)
	})

	t.Run("Recv.Async.Send", func(t *testing.T) {
		in, eg := pipe.New[int](0)
		go func() {
			for i := 0; i < 1000; i++ {
				time.Sleep(time.Duration(rand.Intn(4)) * time.Millisecond)
				it.Ok(t).If(<-in).Equal(i)
			}
		}()
		for i := 0; i < 1000; i++ {
			eg <- i
			time.Sleep(time.Duration(rand.Intn(9)) * time.Millisecond)
		}
		close(eg)
	})
}

func TestPipeFrom(t *testing.T) {
	t.Run("Generate", func(t *testing.T) {
		eg := pipe.From(1, 10*time.Millisecond, func() (int, error) { return 10, nil })
		it.Ok(t).If(<-eg).Equal(10)
		it.Ok(t).If(<-eg).Equal(10)
		it.Ok(t).If(<-eg).Equal(10)
		close(eg)
	})

	t.Run("NoPanic", func(t *testing.T) {
		eg := pipe.From(0, 10*time.Millisecond, func() (int, error) { return 10, nil })
		time.Sleep(30 * time.Millisecond)
		close(eg)
	})
}

func TestPipeMap(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		eg := make(chan int)
		in := pipe.Map(eg, strconv.Itoa)

		eg <- 100
		it.Ok(t).If(<-in).Equal("100")
		close(eg)
	})

	t.Run("Close", func(t *testing.T) {
		eg := make(chan int)
		in := pipe.Map(eg, strconv.Itoa)

		close(in)
		close(eg)
	})
}

func TestPipeMaybeMap(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		eg := make(chan string)
		in := pipe.MaybeMap(eg, strconv.Atoi)

		eg <- "100"
		it.Ok(t).If(<-in).Equal(100)
		close(eg)
	})

	t.Run("Map.Fail", func(t *testing.T) {
		eg := make(chan string)
		in := pipe.MaybeMap(eg, strconv.Atoi)

		eg <- "test"
		eg <- "100"
		it.Ok(t).If(<-in).Equal(100)
		close(eg)
	})

	t.Run("Close", func(t *testing.T) {
		eg := make(chan string)
		in := pipe.MaybeMap(eg, strconv.Atoi)

		close(in)
		close(eg)
	})
}

func TestPipeForEach(t *testing.T) {
	t.Run("ForEach", func(t *testing.T) {
		n := 0
		eg := make(chan int)
		pipe.ForEach(eg, func(a int) { n = n + a })

		eg <- 100
		it.Ok(t).If(n).Equal(100)
		close(eg)
	})
}

func BenchmarkPipe(b *testing.B) {
	in, eg := pipe.New[int](1)
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		eg <- n
		<-in
	}

	close(eg)
}
