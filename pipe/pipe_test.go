package pipe_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fogfish/golem/pipe"
	"github.com/fogfish/it"
)

func TestPipe(t *testing.T) {
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
