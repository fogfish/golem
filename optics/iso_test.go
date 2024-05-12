package optics_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/fogfish/golem/optics"
	"github.com/fogfish/it/v2"
)

func isomorphism[A any](t *testing.T, some A) {
	t.Helper()

	type B struct{ A A }
	type S struct{ S A }
	type T struct{ T A }
	type M = map[string]A

	getter := optics.Getter(
		optics.ForProduct1[S, A](),
		func(a A) B { return B{a} },
	)

	setter := optics.Setter(
		optics.ForProduct1[S, A](),
		func(b B) A { return b.A },
	)

	bimap := optics.BiMap(
		optics.ForProduct1[S, A](),
		func(a A) B { return B{a} },
		func(b B) A { return b.A },
	)

	iso := optics.Iso(
		optics.ForProduct1[S, A](),
		optics.ForProduct1[T, A](),
	)

	isoM := optics.Iso(
		optics.ForProduct1[S, A](),
		optics.NewLensM[string, A]("S"),
	)

	morphism := optics.Morphism(iso, iso, iso)

	t.Run("Getter", func(t *testing.T) {
		s := S{some}
		b := getter.Get(&s)
		getter.Put(&s, B{})

		it.Then(t).Should(
			it.Equiv(b.A, some),
			it.Equiv(s.S, some),
		)
	})

	t.Run("Setter", func(t *testing.T) {
		s := S{}
		setter.Put(&s, B{A: some})
		b := setter.Get(&s)

		it.Then(t).Should(
			it.Equiv(b.A, *new(A)),
			it.Equiv(s.S, some),
		)
	})

	t.Run("BiMap", func(t *testing.T) {
		s := S{}
		bimap.Put(&s, B{A: some})
		b := bimap.Get(&s)

		it.Then(t).Should(
			it.Equiv(b.A, some),
			it.Equiv(s.S, some),
		)
	})

	t.Run("Iso/Forward", func(t *testing.T) {
		s := S{some}
		y := T{}
		iso.Forward(&s, &y)

		it.Then(t).Should(
			it.Equiv(s.S, some),
			it.Equiv(y.T, some),
		)
	})

	t.Run("IsoM/Forward", func(t *testing.T) {
		s := S{some}
		y := M{}
		isoM.Forward(&s, &y)

		it.Then(t).Should(
			it.Equiv(s.S, some),
			it.Equiv(y["S"], some),
		)
	})

	t.Run("Iso/Inverse", func(t *testing.T) {
		s := S{}
		y := T{some}
		iso.Inverse(&y, &s)

		it.Then(t).Should(
			it.Equiv(s.S, some),
			it.Equiv(y.T, some),
		)
	})

	t.Run("IsoM/Inverse", func(t *testing.T) {
		s := S{}
		y := M{"S": some}
		isoM.Inverse(&y, &s)

		it.Then(t).Should(
			it.Equiv(s.S, some),
			it.Equiv(y["S"], some),
		)
	})

	t.Run("Morphism/Forward", func(t *testing.T) {
		s := S{some}
		y := T{}
		morphism.Forward(&s, &y)

		it.Then(t).Should(
			it.Equiv(s.S, some),
			it.Equiv(y.T, some),
		)
	})

	t.Run("Morphism/Inverse", func(t *testing.T) {
		s := S{}
		y := T{some}
		morphism.Inverse(&y, &s)

		it.Then(t).Should(
			it.Equiv(s.S, some),
			it.Equiv(y.T, some),
		)
	})
}

func testIsomorphism[A any](v A) func(t *testing.T) {
	return func(t *testing.T) {
		isomorphism[A](t, v)
		isomorphism[*A](t, &v)
		isomorphism[[]A](t, []A{v})
		isomorphism[*[]A](t, &[]A{v})
	}
}

func TestIsomorphism(t *testing.T) {
	type String string

	t.Run("String", testIsomorphism[string]("string"))
	t.Run("String", testIsomorphism[String]("string"))
	t.Run("Bool", testIsomorphism[bool](true))
	t.Run("Int8", testIsomorphism[int8](10))
	t.Run("UInt8", testIsomorphism[uint8](10))
	t.Run("Byte", testIsomorphism[byte](10))
	t.Run("Int16", testIsomorphism[int16](10))
	t.Run("UInt16", testIsomorphism[uint16](10))
	t.Run("Int32", testIsomorphism[int32](10))
	t.Run("Rune", testIsomorphism[rune]('\u1234'))
	t.Run("UInt32", testIsomorphism[uint32](10))
	t.Run("Int64", testIsomorphism[int64](10))
	t.Run("UInt64", testIsomorphism[uint64](10))
	type Int int
	t.Run("Int", testIsomorphism[int](10))
	t.Run("Int", testIsomorphism[Int](10))
	t.Run("UInt", testIsomorphism[uint](10))
	t.Run("UIntPtr", testIsomorphism[uintptr](10))
	t.Run("Float32", testIsomorphism[float32](10.0))
	t.Run("Float64", testIsomorphism[float64](10.0))
	t.Run("Complex64", testIsomorphism[complex64](10.0+11.0i))
	t.Run("Complex128", testIsomorphism[complex128](10.0+11.0i))

	type Struct struct{ A string }
	t.Run("StructNoName", testIsomorphism[struct{ A string }](struct{ A string }{A: "string"}))
	t.Run("Struct", testIsomorphism[Struct](Struct{A: "string"}))

	t.Run("Interface", testIsomorphism[io.Reader](&bytes.Buffer{}))
}
