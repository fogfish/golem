//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package list_test

import (
	"testing"

	"github.com/fogfish/golem/seq/list"
	"github.com/fogfish/golem/seq/seqtest"
)

var (
	seqT = list.Trait[int]("seq.slice.int")

	defCap int           = 1000000
	defSeq list.Seq[int] = seqT.New()
)

func TestList(t *testing.T) {
	seqtest.TestSeq[list.Seq[int], int](t, seqT, seqT.New(1, 2, 3, 4, 5))
	seqtest.TestFoldable[list.Seq[int]](t, seqT)
}

func init() {
	for n := 0; n < defCap; n++ {
		defSeq = seqT.Cons(n, defSeq)
	}
}

func BenchmarkList(b *testing.B) {
	seqtest.Benchmark[list.Seq[int]](b, seqT, defSeq)
}
