//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package slice_test

import (
	"testing"

	"github.com/fogfish/golem/seq/seqtest"
	"github.com/fogfish/golem/seq/slice"
)

var (
	seqT = slice.Trait[int]("seq.slice.int")

	defCap int            = 1000000
	defSeq slice.Seq[int] = seqT.New()
)

func TestSlice(t *testing.T) {
	seqtest.TestSeq[slice.Seq[int], int](t, seqT, seqT.New(1, 2, 3, 4, 5))
	seqtest.TestFoldable[slice.Seq[int]](t, seqT)
}

func init() {
	seq := make([]int, defCap)
	for n := 0; n < defCap; n++ {
		seq[n] = n
	}
	defSeq = seqT.New(seq...)
}

func BenchmarkList(b *testing.B) {
	seqtest.Benchmark[slice.Seq[int]](b, seqT, defSeq)
}
