//
// Copyright (C) 2022 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package optics_test

import (
	"testing"

	"github.com/fogfish/golem/optics"
	"github.com/fogfish/it/v2"
)

func TestShape(t *testing.T) {
	type T struct{ F1, F2, F3, F4, F5, F6, F7, F8, F9 int }

	t.Run("Shape2", func(t *testing.T) {
		ln := optics.ForShape2[T, int, int]("F1", "F2")

		tt := T{}
		a, b := ln.Get(
			ln.Put(&tt, 1, 2),
		)

		it.Then(t).Should(
			it.Equiv(tt, T{1, 2, 0, 0, 0, 0, 0, 0, 0}),
			it.Seq([]int{a, b}).Equal(1, 2),
		)
	})

	t.Run("Shape3", func(t *testing.T) {
		ln := optics.ForShape3[T, int, int, int]("F1", "F2", "F3")

		tt := T{}
		a, b, c := ln.Get(
			ln.Put(&tt, 1, 2, 3),
		)

		it.Then(t).Should(
			it.Equiv(tt, T{1, 2, 3, 0, 0, 0, 0, 0, 0}),
			it.Seq([]int{a, b, c}).Equal(1, 2, 3),
		)
	})

	t.Run("Shape4", func(t *testing.T) {
		ln := optics.ForShape4[T, int, int, int, int]("F1", "F2", "F3", "F4")

		tt := T{}
		a, b, c, d := ln.Get(
			ln.Put(&tt, 1, 2, 3, 4),
		)

		it.Then(t).Should(
			it.Equiv(tt, T{1, 2, 3, 4, 0, 0, 0, 0, 0}),
			it.Seq([]int{a, b, c, d}).Equal(1, 2, 3, 4),
		)
	})

	t.Run("Shape5", func(t *testing.T) {
		ln := optics.ForShape5[T, int, int, int, int, int]("F1", "F2", "F3", "F4", "F5")

		tt := T{}
		a, b, c, d, e := ln.Get(
			ln.Put(&tt, 1, 2, 3, 4, 5),
		)

		it.Then(t).Should(
			it.Equiv(tt, T{1, 2, 3, 4, 5, 0, 0, 0, 0}),
			it.Seq([]int{a, b, c, d, e}).Equal(1, 2, 3, 4, 5),
		)
	})

	t.Run("Shape6", func(t *testing.T) {
		ln := optics.ForShape6[T, int, int, int, int, int, int]("F1", "F2", "F3", "F4", "F5", "F6")

		tt := T{}
		a, b, c, d, e, f := ln.Get(
			ln.Put(&tt, 1, 2, 3, 4, 5, 6),
		)

		it.Then(t).Should(
			it.Equiv(tt, T{1, 2, 3, 4, 5, 6, 0, 0, 0}),
			it.Seq([]int{a, b, c, d, e, f}).Equal(1, 2, 3, 4, 5, 6),
		)
	})

	t.Run("Shape7", func(t *testing.T) {
		ln := optics.ForShape7[T, int, int, int, int, int, int, int]("F1", "F2", "F3", "F4", "F5", "F6", "F7")

		tt := T{}
		a, b, c, d, e, f, g := ln.Get(
			ln.Put(&tt, 1, 2, 3, 4, 5, 6, 7),
		)

		it.Then(t).Should(
			it.Equiv(tt, T{1, 2, 3, 4, 5, 6, 7, 0, 0}),
			it.Seq([]int{a, b, c, d, e, f, g}).Equal(1, 2, 3, 4, 5, 6, 7),
		)
	})

	t.Run("Shape8", func(t *testing.T) {
		ln := optics.ForShape8[T, int, int, int, int, int, int, int, int]("F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8")

		tt := T{}
		a, b, c, d, e, f, g, h := ln.Get(
			ln.Put(&tt, 1, 2, 3, 4, 5, 6, 7, 8),
		)

		it.Then(t).Should(
			it.Equiv(tt, T{1, 2, 3, 4, 5, 6, 7, 8, 0}),
			it.Seq([]int{a, b, c, d, e, f, g, h}).Equal(1, 2, 3, 4, 5, 6, 7, 8),
		)
	})

	t.Run("Shape9", func(t *testing.T) {
		ln := optics.ForShape9[T, int, int, int, int, int, int, int, int, int]("F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9")

		tt := T{}
		a, b, c, d, e, f, g, h, i := ln.Get(
			ln.Put(&tt, 1, 2, 3, 4, 5, 6, 7, 8, 9),
		)

		it.Then(t).Should(
			it.Equiv(tt, T{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			it.Seq([]int{a, b, c, d, e, f, g, h, i}).Equal(1, 2, 3, 4, 5, 6, 7, 8, 9),
		)
	})
}
