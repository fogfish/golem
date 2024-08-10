//
// Copyright (C) 2022 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package fork_test

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/fogfish/golem/pipe/fork"
	"github.com/fogfish/golem/pure/monoid"
	"github.com/fogfish/it/v2"
)

const par = 4

func TestEmit(t *testing.T) {
	seq := 0
	emit := func() (int, error) {
		seq++
		return seq, nil
	}

	ctx, close := context.WithCancel(context.Background())
	eg := fork.StdErr(fork.Emit(ctx, 0, 10*time.Microsecond, emit))

	it.Then(t).Should(
		it.Equal(<-eg, 1),
		it.Equal(<-eg, 2),
		it.Equal(<-eg, 3),
		it.Equal(<-eg, 4),
		it.Equal(<-eg, 5),
		it.Equal(<-eg, 6),
	)
	close()
}

func TestFilter(t *testing.T) {
	ctx, close := context.WithCancel(context.Background())
	seq := fork.Seq(1, 2, 3, 4, 5)
	out := fork.Filter(ctx, par, seq,
		func(x int) bool { return x%2 == 1 },
	)

	it.Then(t).Should(
		it.Seq(fork.ToSeq(out)).Contain().AllOf(1, 3, 5),
	)

	close()
}

func TestForEach(t *testing.T) {
	ctx, close := context.WithCancel(context.Background())

	seq := fork.Seq(1, 2, 3, 4, 5)

	var m sync.Mutex
	n := 0
	<-fork.ForEach(ctx, par, seq, func(a int) {
		m.Lock()
		defer m.Unlock()
		n = n + a
	})

	it.Then(t).Should(
		it.Equal(n, 15),
	)

	close()
}

func TestFMap(t *testing.T) {
	t.Run("FMap", func(t *testing.T) {
		ctx, close := context.WithCancel(context.Background())
		seq := fork.Seq(1, 2, 3, 4, 5)
		out := fork.StdErr(fork.FMap(ctx, par, seq,
			func(x int) (<-chan string, error) {
				return fork.Seq(strconv.Itoa(x)), nil
			}),
		)

		it.Then(t).Should(
			it.Seq(fork.ToSeq(out)).Contain().AllOf("1", "2", "3", "4", "5"),
		)

		close()
	})

	t.Run("Err", func(t *testing.T) {
		ctx, close := context.WithCancel(context.Background())
		seq := fork.Seq(1, 2, 3, 4, 5)
		_, exx := fork.FMap(ctx, par, seq,
			func(x int) (<-chan string, error) {
				return nil, fmt.Errorf("fail")
			},
		)

		it.Then(t).ShouldNot(
			it.Nil(<-exx),
		)

		close()
	})
}

func TestFold(t *testing.T) {
	seq := fork.Seq(1, 2, 3, 4, 5)
	n := <-fork.Fold(context.Background(), par, seq,
		monoid.FromOp(0, func(a, b int) int { return a + b }),
	)

	it.Then(t).Should(
		it.Equal(n, 15),
	)
}

func TestMap(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		ctx, close := context.WithCancel(context.Background())
		seq := fork.Seq(1, 2, 3, 4, 5)
		out := fork.StdErr(fork.Map(ctx, par, seq,
			func(x int) (string, error) { return strconv.Itoa(x), nil },
		))

		it.Then(t).Should(
			it.Seq(fork.ToSeq(out)).Contain().AllOf("1", "2", "3", "4", "5"),
		)
		close()
	})

	t.Run("Err", func(t *testing.T) {
		ctx, close := context.WithCancel(context.Background())
		seq := fork.Seq(1, 2, 3, 4, 5)
		_, exx := fork.Map(ctx, par, seq,
			func(x int) (string, error) { return "", fmt.Errorf("fail") },
		)

		it.Then(t).ShouldNot(
			it.Nil(<-exx),
		)

		close()
	})
}

func TestPartition(t *testing.T) {
	ctx, close := context.WithCancel(context.Background())
	seq := fork.Seq(1, 2, 3, 4, 5)
	lo, ro := fork.Partition(ctx, par, seq,
		func(x int) bool { return x%2 == 1 },
	)

	it.Then(t).Should(
		it.Seq(fork.ToSeq(lo)).Equal(1, 3, 5),
		it.Seq(fork.ToSeq(ro)).Equal(2, 4),
	)

	close()
}

func TestUnfold(t *testing.T) {
	ctx, close := context.WithCancel(context.Background())
	seq := fork.StdErr(fork.Unfold(ctx, 1, 0, func(x int) (int, error) { return x + 1, nil }))

	it.Then(t).Should(
		it.Equal(<-seq, 0),
		it.Equal(<-seq, 1),
		it.Equal(<-seq, 2),
		it.Equal(<-seq, 3),
		it.Equal(<-seq, 4),
		it.Equal(<-seq, 5),
	)
	close()
}

func TestJoin(t *testing.T) {
	ctx, close := context.WithCancel(context.Background())

	a := fork.Seq(1, 2, 3)
	b := fork.Seq(4, 5, 6)
	c := fork.Seq(7, 8, 9)

	out := fork.Join(ctx, a, b, c)

	it.Then(t).Should(
		it.Seq(fork.ToSeq(out)).Contain().AllOf(1, 2, 3, 4, 5, 6, 7, 8, 9),
	)

	close()
}

func TestTake(t *testing.T) {
	ctx, close := context.WithCancel(context.Background())
	seq := fork.Seq(1, 2, 3, 4, 5, 6)
	out := fork.Take(ctx, seq, 3)

	it.Then(t).Should(
		it.Seq(fork.ToSeq(out)).Equal(1, 2, 3),
	)

	close()
}

func TestTakeWhile(t *testing.T) {
	ctx, close := context.WithCancel(context.Background())
	seq := fork.Seq(1, 2, 3, 4, 5, 6)
	out := fork.TakeWhile(ctx, seq, func(x int) bool { return x < 4 })

	it.Then(t).Should(
		it.Seq(fork.ToSeq(out)).Equal(1, 2, 3),
	)

	close()
}
