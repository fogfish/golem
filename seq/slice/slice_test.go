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

func TestSlice(t *testing.T) {
	seqT := slice.Trait[int]("seq.slice.int")

	seqtest.TestSeq[slice.Seq[int], int](t, seqT, seqT.New(1, 2, 3, 4, 5))
	seqtest.TestFoldable[slice.Seq[int]](t, seqT)
}
