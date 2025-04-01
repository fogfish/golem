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
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/fogfish/golem/pipe/v2"
	"github.com/fogfish/golem/pure/monoid"
	"github.com/fogfish/it/v2"
)

func TestEmit(t *testing.T) {
	t.Run("Emit", func(t *testing.T) {
		emit := pipe.Pure(func(x int) int { return x })

		ctx, close := context.WithCancel(context.Background())
		eg := pipe.StdErr(pipe.Emit(ctx, 0, 10*time.Microsecond, emit))

		it.Then(t).Should(
			it.Equal(<-eg, 0),
			it.Equal(<-eg, 1),
			it.Equal(<-eg, 2),
			it.Equal(<-eg, 3),
			it.Equal(<-eg, 4),
			it.Equal(<-eg, 5),
		)
		close()
	})

	t.Run("Err", func(t *testing.T) {
		fail := pipe.Lift(func(int) (int, error) { return 0, fmt.Errorf("fail") })

		ctx, close := context.WithCancel(context.Background())
		_, exx := pipe.Emit(ctx, 0, 10*time.Millisecond, fail)

		it.Then(t).ShouldNot(
			it.Nil(<-exx),
		)

		close()
	})
}

func TestFilter(t *testing.T) {
	fun := pipe.Pure(func(x int) bool { return x%2 == 1 })

	ctx, close := context.WithCancel(context.Background())
	seq := pipe.Seq(1, 2, 3, 4, 5)
	out := pipe.Filter(ctx, seq, fun)

	it.Then(t).Should(
		it.Seq(pipe.ToSeq(out)).Equal(1, 3, 5),
	)

	close()
}

func TestForEach(t *testing.T) {
	n := 0
	fun := pipe.Pure(func(x int) int {
		n = n + x
		return n
	})

	ctx, close := context.WithCancel(context.Background())
	seq := pipe.Seq(1, 2, 3, 4, 5)
	<-pipe.ForEach(ctx, seq, fun)

	it.Then(t).Should(
		it.Equal(n, 15),
	)

	close()
}

func TestFMap(t *testing.T) {
	t.Run("FMap", func(t *testing.T) {
		fun := pipe.LiftF(
			func(ctx context.Context, x int, ch chan<- string) error {
				ch <- strconv.Itoa(x)
				return nil
			},
		)

		ctx, close := context.WithCancel(context.Background())
		seq := pipe.Seq(1, 2, 3, 4, 5)
		out := pipe.StdErr(pipe.FMap(ctx, seq, fun))

		it.Then(t).Should(
			it.Seq(pipe.ToSeq(out)).Equal("1", "2", "3", "4", "5"),
		)

		close()
	})

	t.Run("Err", func(t *testing.T) {
		fun := pipe.LiftF(
			func(ctx context.Context, x int, ch chan<- string) error {
				return fmt.Errorf("fail")
			},
		)

		ctx, close := context.WithCancel(context.Background())
		seq := pipe.Seq(1, 2, 3, 4, 5)
		_, exx := pipe.FMap(ctx, seq, fun)

		it.Then(t).ShouldNot(
			it.Nil(<-exx),
		)

		close()
	})

	t.Run("Cancel", func(t *testing.T) {
		emit := pipe.Pure(func(x int) int { return x })
		fun := pipe.LiftF(
			func(ctx context.Context, x int, ch chan<- int) error {
				ch <- x
				return nil
			},
		)

		ctx, close := context.WithCancel(context.Background())
		seq := pipe.StdErr(pipe.Emit(ctx, 1000, 10*time.Microsecond, emit))
		out := pipe.StdErr(pipe.FMap(ctx, seq, fun))

		vals := pipe.ToSeq(pipe.Take(ctx, out, 10))
		close()

		it.Then(t).Should(
			it.Seq(vals).Contain().AllOf(0, 1, 2, 3, 4, 5, 6, 7, 8, 9),
		)
	})
}

func TestFold(t *testing.T) {
	acc := monoid.FromOp(0, func(a, b int) int { return a + b })
	seq := pipe.Seq(1, 2, 3, 4, 5)
	n := <-pipe.Fold(context.Background(), seq, acc)

	it.Then(t).Should(
		it.Equal(n, 15),
	)
}

func TestMap(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		fun := pipe.Pure(strconv.Itoa)

		ctx, close := context.WithCancel(context.Background())
		seq := pipe.Seq(1, 2, 3, 4, 5)
		out := pipe.StdErr(pipe.Map(ctx, seq, fun))

		it.Then(t).Should(
			it.Seq(pipe.ToSeq(out)).Equal("1", "2", "3", "4", "5"),
		)
		close()
	})

	t.Run("Err", func(t *testing.T) {
		fun := pipe.Lift(func(x int) (string, error) { return "", fmt.Errorf("fail") })

		ctx, close := context.WithCancel(context.Background())
		seq := pipe.Seq(1, 2, 3, 4, 5)
		_, exx := pipe.Map(ctx, seq, fun)

		it.Then(t).ShouldNot(
			it.Nil(<-exx),
		)

		close()
	})
}

func TestPartition(t *testing.T) {
	fun := pipe.Pure(func(x int) bool { return x%2 == 1 })

	ctx, close := context.WithCancel(context.Background())
	seq := pipe.Seq(1, 2, 3, 4, 5)
	lo, ro := pipe.Partition(ctx, seq, fun)

	it.Then(t).Should(
		it.Seq(pipe.ToSeq(lo)).Equal(1, 3, 5),
		it.Seq(pipe.ToSeq(ro)).Equal(2, 4),
	)

	close()
}

func TestUnfold(t *testing.T) {
	t.Run("Unfold", func(t *testing.T) {
		fun := pipe.Pure(func(x int) int { return x + 1 })

		ctx, close := context.WithCancel(context.Background())
		seq := pipe.StdErr(pipe.Unfold(ctx, 1, 0, fun))

		it.Then(t).Should(
			it.Equal(<-seq, 0),
			it.Equal(<-seq, 1),
			it.Equal(<-seq, 2),
			it.Equal(<-seq, 3),
			it.Equal(<-seq, 4),
			it.Equal(<-seq, 5),
		)
		close()
	})

	t.Run("Err", func(t *testing.T) {
		fun := pipe.Lift(func(x int) (int, error) { return x, fmt.Errorf("fail") })

		ctx, close := context.WithCancel(context.Background())
		out, exx := pipe.Unfold(ctx, 1, 0, fun)

		it.Then(t).Should(
			it.Equal(<-out, 0),
		).
			ShouldNot(
				it.Nil(<-exx),
			)

		close()
	})
}

func TestJoin(t *testing.T) {
	ctx, close := context.WithCancel(context.Background())

	a := pipe.Seq(1, 2, 3)
	b := pipe.Seq(4, 5, 6)
	c := pipe.Seq(7, 8, 9)

	out := pipe.Join(ctx, a, b, c)

	it.Then(t).Should(
		it.Seq(pipe.ToSeq(out)).Contain().AllOf(1, 2, 3, 4, 5, 6, 7, 8, 9),
	)

	close()
}

func TestTake(t *testing.T) {
	ctx, close := context.WithCancel(context.Background())
	seq := pipe.Seq(1, 2, 3, 4, 5, 6)
	out := pipe.Take(ctx, seq, 3)

	it.Then(t).Should(
		it.Seq(pipe.ToSeq(out)).Equal(1, 2, 3),
	)

	close()
}

func TestTakeWhile(t *testing.T) {
	fun := pipe.Pure(func(x int) bool { return x < 4 })

	ctx, close := context.WithCancel(context.Background())
	seq := pipe.Seq(1, 2, 3, 4, 5, 6)
	out := pipe.TakeWhile(ctx, seq, fun)

	it.Then(t).Should(
		it.Seq(pipe.ToSeq(out)).Equal(1, 2, 3),
	)

	close()
}

func TestThrottling(t *testing.T) {
	ctx, close := context.WithCancel(context.Background())
	seq := pipe.Seq(1, 2, 3, 4, 5, 6, 7, 8, 9, 0)
	slowSeq := pipe.Throttling(ctx, seq, 1, 100*time.Millisecond)
	out := pipe.StdErr(pipe.Map(ctx, slowSeq,
		pipe.Lift(func(_ int) (time.Time, error) {
			return time.Now(), nil
		}),
	))
	wt := pipe.ToSeq(out)
	for i := 1; i < len(wt); i++ {
		diff := wt[i].Sub(wt[i-1])

		it.Then(t).Should(
			it.Less(diff, 110*time.Millisecond),
			it.Greater(diff, 99*time.Millisecond),
		)
	}

	close()
}

func BenchmarkPipe(b *testing.B) {
	ctx, close := context.WithCancel(context.Background())
	in, eg := pipe.New[int](ctx, 0)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		eg <- n
		<-in
	}

	close()
}
