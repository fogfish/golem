//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package seq_test

import (
	"strconv"
	"testing"

	"github.com/fogfish/golem/trait/seq"
	"github.com/fogfish/it/v2"
)

func toSlice[T any](e seq.Seq[T]) []T {
	r := make([]T, 0)
	for has := e != nil; has; has = e.Next() {
		r = append(r, e.Value())
	}
	return r
}

func TestTakeWhile(t *testing.T) {
	for input, expect := range map[*[]int][]int{
		{1, 2, 3, 4}:     {1, 2, 3, 4},
		{1, 2, 3, 4, 10}: {1, 2, 3, 4},
		{1, 2, 10, 3, 4}: {1, 2},
		{1, 10, 2, 3, 4}: {1},
		{10}:             {},
		{}:               {},
	} {
		e := seq.TakeWhile(seq.FromSlice(*input),
			func(v int) bool { return v < 10 },
		)

		it.Then(t).Should(
			it.Seq(toSlice(e)).Equal(expect...),
		)
	}
}

func TestDropWhile(t *testing.T) {
	for input, expect := range map[*[]int][]int{
		{1, 2, 3, 4}:     {},
		{1, 2, 3, 4, 10}: {10},
		{1, 2, 10, 3, 4}: {10, 3, 4},
		{1, 10, 2, 3, 4}: {10, 2, 3, 4},
		{10, 1, 2, 3, 4}: {10, 1, 2, 3, 4},
		{}:               {},
	} {
		e := seq.DropWhile(seq.FromSlice(*input),
			func(v int) bool { return v < 10 },
		)

		it.Then(t).Should(
			it.Seq(toSlice(e)).Equal(expect...),
		)
	}
}

func TestFilter(t *testing.T) {
	for input, expect := range map[*[]int][]int{
		{1, 2, 3, 4, 5, 6}: {2, 4, 6},
		{2, 3, 4, 5, 6, 7}: {2, 4, 6},
		{2, 3, 4, 5, 6}:    {2, 4, 6},
		{1, 3, 5, 7}:       {},
		{2, 4, 6, 8}:       {2, 4, 6, 8},
		{}:                 {},
	} {
		e := seq.Filter(seq.FromSlice(*input),
			func(v int) bool { return v%2 == 0 },
		)

		it.Then(t).Should(
			it.Seq(toSlice(e)).Equal(expect...),
		)
	}
}

func TestForEach(t *testing.T) {
	for input, expect := range map[*[]int][]int{
		{1, 2, 3, 4, 5, 6}: {1, 2, 3, 4, 5, 6},
		{2, 4}:             {2, 4},
		{2}:                {2},
		{}:                 {},
	} {
		e := make([]int, 0)
		seq.ForEach(seq.FromSlice(*input),
			func(v int) error {
				e = append(e, v)
				return nil
			},
		)

		it.Then(t).Should(
			it.Seq(e).Equal(expect...),
		)
	}
}

func TestMap(t *testing.T) {
	for input, expect := range map[*[]int][]string{
		{1, 2, 3, 4, 5, 6}: {"1", "2", "3", "4", "5", "6"},
		{2, 4}:             {"2", "4"},
		{2}:                {"2"},
		{}:                 {},
	} {
		e := seq.Map(seq.FromSlice(*input), strconv.Itoa)

		it.Then(t).Should(
			it.Seq(toSlice(e)).Equal(expect...),
		)
	}
}

func TestPlus(t *testing.T) {
	for input, expect := range map[*[]int][]int{
		{1, 2, 3}: {1, 2, 3, 1, 2, 3},
		{2, 4}:    {2, 4, 2, 4},
		{2}:       {2, 2},
		{}:        {},
	} {
		e := seq.Plus(seq.FromSlice(*input), seq.FromSlice(*input))

		it.Then(t).Should(
			it.Seq(toSlice(e)).Equal(expect...),
		)

		e = seq.Plus(seq.FromSlice(*input), nil)

		it.Then(t).Should(
			it.Seq(toSlice(e)).Equal(*input...),
		)

		e = seq.Plus(nil, seq.FromSlice(*input))

		it.Then(t).Should(
			it.Seq(toSlice(e)).Equal(*input...),
		)
	}
}

func TestJoin(t *testing.T) {
	for input, expect := range map[*[]int][]string{
		{1, 2, 3}: {"1", "2", "3"},
		{2, 4}:    {"2", "4"},
		{2}:       {"2"},
		{}:        {},
	} {
		e := seq.Join(seq.FromSlice(*input),
			func(x int) seq.Seq[string] {
				return seq.From(strconv.Itoa(x))
			},
		)

		it.Then(t).Should(
			it.Seq(toSlice(e)).Equal(expect...),
		)
	}
}

func TestJoinNil(t *testing.T) {
	for input, expect := range map[*[]int][]string{
		{1, 2, 3}: {},
		{2, 4}:    {},
		{2}:       {},
		{}:        {},
	} {
		e := seq.Join(seq.FromSlice(*input),
			func(x int) seq.Seq[string] { return nil },
		)

		it.Then(t).Should(
			it.Seq(toSlice(e)).Equal(expect...),
		)
	}

	e := seq.Join(nil,
		func(x int) seq.Seq[string] {
			return seq.From(strconv.Itoa(x))
		},
	)
	it.Then(t).Should(it.Nil(e))
}
