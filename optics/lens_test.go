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

func testLensLaws[A any](v A) func(t *testing.T) {
	return func(t *testing.T) {
		lensLaws[A](t, v)
		lensLaws[*A](t, &v)
		lensLaws[[]A](t, []A{v})
		lensLaws[*[]A](t, &[]A{v})
	}
}

func TestLensLaws(t *testing.T) {
	type String string

	t.Run("String", testLensLaws[string]("string"))
	t.Run("String", testLensLaws[String]("string"))
	t.Run("Bool", testLensLaws[bool](true))
	t.Run("Int8", testLensLaws[int8](10))
	t.Run("UInt8", testLensLaws[uint8](10))
	t.Run("Byte", testLensLaws[byte](10))
	t.Run("Int16", testLensLaws[int16](10))
	t.Run("UInt16", testLensLaws[uint16](10))
	t.Run("Int32", testLensLaws[int32](10))
	t.Run("Rune", testLensLaws[rune]('\u1234'))
	t.Run("UInt32", testLensLaws[uint32](10))
	t.Run("Int64", testLensLaws[int64](10))
	t.Run("UInt64", testLensLaws[uint64](10))
	type Int int
	t.Run("Int", testLensLaws[int](10))
	t.Run("Int", testLensLaws[Int](10))
	t.Run("UInt", testLensLaws[uint](10))
	t.Run("UIntPtr", testLensLaws[uintptr](10))
	t.Run("Float32", testLensLaws[float32](10.0))
	t.Run("Float64", testLensLaws[float64](10.0))
	t.Run("Complex64", testLensLaws[complex64](10.0+11.0i))
	t.Run("Complex128", testLensLaws[complex128](10.0+11.0i))

	type Struct struct{ A string }
	t.Run("StructNoName", testLensLaws[struct{ A string }](struct{ A string }{A: "string"}))
	t.Run("Struct", testLensLaws[Struct](Struct{A: "string"}))

	t.Run("Interface", testLensLaws[io.Reader](&bytes.Buffer{}))

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
