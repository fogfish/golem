//
// Copyright (C) 2022 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package optics_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/fogfish/golem/optics"
	"github.com/fogfish/it/v2"
)

func lensLaws[A any](t *testing.T, some A) {
	t.Helper()

	type T struct{ A A }
	la := optics.ForProduct1[T, A]()

	// Law `GetPut`:
	//  If we get focused element `A` from `S` and
	//  immediately put `A` with no modifications back into `S`,
	//  we must get back exactly `S`.
	t.Run("GetPut", func(t *testing.T) {
		tt := T{A: some}
		la.Put(&tt, la.Get(&tt))

		it.Then(t).Should(
			it.Equiv(tt, T{A: some}),
		)
	})

	// Law `PutGet`:
	//  If putting `A` inside `S` yields a new `S`,
	//  then the `A` obtained from `S` is exactly `A`.
	t.Run("PutGet", func(t *testing.T) {
		tt := T{}
		vv := la.Get(la.Put(&tt, some))

		it.Then(t).Should(
			it.Equiv(vv, some),
			it.Equiv(tt, T{A: some}),
		)
	})

	// Law `PutPut`:
	//  A sequence of two puts is just the effect of the second,
	//  the first is completely overwritten. This law is applicable
	//  to every well behaving lenses.
	t.Run("PutPut", func(t *testing.T) {
		tt := T{}
		la.Put(la.Put(&tt, *new(A)), some)

		it.Then(t).Should(
			it.Equiv(tt, T{A: some}),
		)
	})
}

func testLaws[A any](v A) func(t *testing.T) {
	return func(t *testing.T) {
		lensLaws[A](t, v)
		lensLaws[*A](t, &v)
		lensLaws[[]A](t, []A{v})
		lensLaws[*[]A](t, &[]A{v})
	}
}

func TestLaws(t *testing.T) {
	type String string

	t.Run("String", testLaws[string]("string"))
	t.Run("String", testLaws[String]("string"))
	t.Run("Bool", testLaws[bool](true))
	t.Run("Int8", testLaws[int8](10))
	t.Run("UInt8", testLaws[uint8](10))
	t.Run("Byte", testLaws[byte](10))
	t.Run("Int16", testLaws[int16](10))
	t.Run("UInt16", testLaws[uint16](10))
	t.Run("Int32", testLaws[int32](10))
	t.Run("Rune", testLaws[rune]('\u1234'))
	t.Run("UInt32", testLaws[uint32](10))
	t.Run("Int64", testLaws[int64](10))
	t.Run("UInt64", testLaws[uint64](10))
	type Int int
	t.Run("Int", testLaws[int](10))
	t.Run("Int", testLaws[Int](10))
	t.Run("UInt", testLaws[uint](10))
	t.Run("UIntPtr", testLaws[uintptr](10))
	t.Run("Float32", testLaws[float32](10.0))
	t.Run("Float64", testLaws[float64](10.0))
	t.Run("Complex64", testLaws[complex64](10.0+11.0i))
	t.Run("Complex128", testLaws[complex128](10.0+11.0i))

	type Struct struct{ A string }
	t.Run("StructNoName", testLaws[struct{ A string }](struct{ A string }{A: "string"}))
	t.Run("Struct", testLaws[Struct](Struct{A: "string"}))

	t.Run("Interface", testLaws[io.Reader](&bytes.Buffer{}))

	t.Run("Embedded", func(t *testing.T) {
		type S string
		type I int
		type A struct{ S }
		type B struct{ I }
		type T struct {
			A
			B
			// Note: lens package does not support embedded pointers
		}

		la := optics.ForProduct1[T, S]()
		lb := optics.ForProduct1[T, I]("I")

		t.Run("GetPut", func(t *testing.T) {
			tt := T{A: A{"string"}, B: B{10}}
			la.Put(&tt, la.Get(&tt))
			lb.Put(&tt, lb.Get(&tt))

			it.Then(t).Should(
				it.Equiv(tt, T{A: A{"string"}, B: B{10}}),
			)
		})

		t.Run("PutGet", func(t *testing.T) {
			tt := T{}
			va := la.Get(la.Put(&tt, "string"))
			vb := lb.Get(lb.Put(&tt, 10))

			it.Then(t).Should(
				it.Equal(va, "string"),
				it.Equal(vb, 10),
				it.Equiv(tt, T{A: A{"string"}, B: B{10}}),
			)
		})

		t.Run("PutPut", func(t *testing.T) {
			tt := T{}
			la.Put(la.Put(&tt, "foobar"), "string")
			lb.Put(lb.Put(&tt, 11), 10)

			it.Then(t).Should(
				it.Equiv(tt, T{A: A{"string"}, B: B{10}}),
			)
		})
	})
}

//
//
//

type lensStructJSON[S, A any] struct{ optics.Lens[S, A] }

func newLensJSON[S, A any](l optics.Lens[S, A]) optics.Lens[S, string] {
	return &lensStructJSON[S, A]{l}
}

func (l lensStructJSON[S, A]) Put(s *S, a string) *S {
	var o A

	if err := json.Unmarshal([]byte(a), &o); err != nil {
		panic(err)
	}
	return l.Lens.Put(s, o)
}

func (l lensStructJSON[S, A]) Get(s *S) string {
	v, err := json.Marshal(l.Lens.Get(s))
	if err != nil {
		panic(err)
	}

	return string(v)
}

func TestLensCompose(t *testing.T) {
	type Inner struct {
		E, F, G string
	}

	type T struct {
		D *Inner
	}

	ld := newLensJSON(optics.ForProduct1[T, *Inner]())

	t.Run("GetPut", func(t *testing.T) {
		tt := T{D: &Inner{E: "E", F: "F", G: "G"}}
		ld.Put(&tt, ld.Get(&tt))

		it.Then(t).Should(
			it.Equiv(tt, T{D: &Inner{E: "E", F: "F", G: "G"}}),
		)
	})

	t.Run("PutGet", func(t *testing.T) {
		tt := T{}
		vv := ld.Get(ld.Put(&tt, "{\"E\":\"E\",\"F\":\"F\",\"G\":\"G\"}"))

		it.Then(t).Should(
			it.Equal(vv, "{\"E\":\"E\",\"F\":\"F\",\"G\":\"G\"}"),
			it.Equiv(tt, T{D: &Inner{E: "E", F: "F", G: "G"}}),
		)
	})

	t.Run("PutPut", func(t *testing.T) {
		tt := T{}
		ld.Put(ld.Put(&tt, "{\"E\":\"e\",\"F\":\"f\",\"G\":\"g\"}"), "{\"E\":\"E\",\"F\":\"F\",\"G\":\"G\"}")

		it.Then(t).Should(
			it.Equiv(tt, T{D: &Inner{E: "E", F: "F", G: "G"}}),
		)
	})
}

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
