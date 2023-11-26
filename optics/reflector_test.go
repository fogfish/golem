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
	"io"
	"testing"

	"github.com/fogfish/golem/optics"
	"github.com/fogfish/it/v2"
)

func reflectorLaws[A any](t *testing.T, some A) {
	t.Helper()

	type T struct{ A A }
	la := optics.ForSpectrum1[T, A]()

	// Law `GetPut`:
	//  If we get focused element `A` from `S` and
	//  immediately put `A` with no modifications back into `S`,
	//  we must get back exactly `S`.
	t.Run("GettPutt", func(t *testing.T) {
		tt := T{A: some}
		la.Putt(&tt, la.Gett(&tt))

		it.Then(t).Should(
			it.Equiv(tt, T{A: some}),
		)
	})

	// Law `PutGet`:
	//  If putting `A` inside `S` yields a new `S`,
	//  then the `A` obtained from `S` is exactly `A`.
	t.Run("PuttGett", func(t *testing.T) {
		tt := T{}
		vv := la.Gett(la.Putt(&tt, some))

		it.Then(t).Should(
			it.Equiv(vv, some),
			it.Equiv(tt, T{A: some}),
		)
	})

	// Law `PutPut`:
	//  A sequence of two puts is just the effect of the second,
	//  the first is completely overwritten. This law is applicable
	//  to every well behaving lenses.
	t.Run("PuttPutt", func(t *testing.T) {
		tt := T{}
		la.Putt(la.Putt(&tt, *new(A)), some)

		it.Then(t).Should(
			it.Equiv(tt, T{A: some}),
		)
	})
}

func testReflectorLaws[A any](v A) func(t *testing.T) {
	return func(t *testing.T) {
		reflectorLaws[A](t, v)
		reflectorLaws[*A](t, &v)
		reflectorLaws[[]A](t, []A{v})
		reflectorLaws[*[]A](t, &[]A{v})
	}
}

func TestReflectorLaws(t *testing.T) {
	type String string

	t.Run("String", testReflectorLaws[string]("string"))
	t.Run("String", testReflectorLaws[String]("string"))
	t.Run("Bool", testReflectorLaws[bool](true))
	t.Run("Int8", testReflectorLaws[int8](10))
	t.Run("UInt8", testReflectorLaws[uint8](10))
	t.Run("Byte", testReflectorLaws[byte](10))
	t.Run("Int16", testReflectorLaws[int16](10))
	t.Run("UInt16", testReflectorLaws[uint16](10))
	t.Run("Int32", testReflectorLaws[int32](10))
	t.Run("Rune", testReflectorLaws[rune]('\u1234'))
	t.Run("UInt32", testReflectorLaws[uint32](10))
	t.Run("Int64", testReflectorLaws[int64](10))
	t.Run("UInt64", testReflectorLaws[uint64](10))
	type Int int
	t.Run("Int", testReflectorLaws[int](10))
	t.Run("Int", testReflectorLaws[Int](10))
	t.Run("UInt", testReflectorLaws[uint](10))
	t.Run("UIntPtr", testReflectorLaws[uintptr](10))
	t.Run("Float32", testReflectorLaws[float32](10.0))
	t.Run("Float64", testReflectorLaws[float64](10.0))
	t.Run("Complex64", testReflectorLaws[complex64](10.0+11.0i))
	t.Run("Complex128", testReflectorLaws[complex128](10.0+11.0i))

	type Struct struct{ A string }
	t.Run("StructNoName", testReflectorLaws[struct{ A string }](struct{ A string }{A: "string"}))
	t.Run("Struct", testReflectorLaws[Struct](Struct{A: "string"}))

	t.Run("Interface", testReflectorLaws[io.Reader](&bytes.Buffer{}))

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

		la := optics.ForSpectrum1[T, S]()
		lb := optics.ForSpectrum1[T, I]("I")

		t.Run("GettPutt", func(t *testing.T) {
			tt := T{A: A{"string"}, B: B{10}}
			la.Putt(&tt, la.Gett(&tt))
			lb.Putt(&tt, lb.Gett(&tt))

			it.Then(t).Should(
				it.Equiv(tt, T{A: A{"string"}, B: B{10}}),
			)
		})

		t.Run("PuttGett", func(t *testing.T) {
			tt := T{}
			va := la.Gett(la.Putt(&tt, "string"))
			vb := lb.Gett(lb.Putt(&tt, 10))

			it.Then(t).Should(
				it.Equal(va, "string"),
				it.Equal(vb, 10),
				it.Equiv(tt, T{A: A{"string"}, B: B{10}}),
			)
		})

		t.Run("PuttPutt", func(t *testing.T) {
			tt := T{}
			la.Putt(la.Putt(&tt, "foobar"), "string")
			lb.Putt(lb.Putt(&tt, 11), 10)

			it.Then(t).Should(
				it.Equiv(tt, T{A: A{"string"}, B: B{10}}),
			)
		})
	})
}

func TestReflectorIncompatibility(t *testing.T) {
	type A struct {
		S string
	}

	type B struct {
		S string
	}

	la := optics.ForSpectrum1[A, string]()
	lb := optics.ForSpectrum1[B, string]()

	aa := A{}
	bb := B{}

	it.Then(t).Should(
		it.Fail(func() { la.Putt(bb, "aa") }),
		it.Fail(func() { lb.Putt(aa, "bb") }),
	)
}
