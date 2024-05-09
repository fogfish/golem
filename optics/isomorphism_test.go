package optics_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/fogfish/golem/optics"
	"github.com/fogfish/it/v2"
)

func isomorphism[T any](t *testing.T, some T) {
	t.Helper()

	type M map[string]T
	type A struct{ X T }
	type B struct{ X T }
	type C struct {
		C T `optics:"X"`
	}

	isoAB := optics.Iso(optics.ForShape[A](), optics.ForShape[B]())
	isoAC := optics.Iso(optics.ForShape[A](), optics.ForShape[C]())
	isoAM := optics.Iso(optics.ForShape[A](), optics.ForShape[M]("X"))

	t.Run("FromAtoB", func(t *testing.T) {
		a := A{X: some}
		b := B{}

		isoAB.FromAtoB(&a, &b)
		it.Then(t).Should(it.Equiv(a.X, b.X))
	})

	t.Run("FromAtoB_with_optics", func(t *testing.T) {
		a := A{X: some}
		c := C{}

		isoAC.FromAtoB(&a, &c)
		it.Then(t).Should(it.Equiv(a.X, c.C))
	})

	t.Run("FromAtoB_with_map", func(t *testing.T) {
		a := A{X: some}
		m := M{}

		isoAM.FromAtoB(&a, &m)
		it.Then(t).Should(it.Equiv(a.X, m["X"]))
	})

	t.Run("FromBtoA", func(t *testing.T) {
		a := A{}
		b := B{X: some}

		isoAB.FromBtoA(&b, &a)
		it.Then(t).Should(it.Equiv(a.X, b.X))
	})

	t.Run("FromBtoA_with_optics", func(t *testing.T) {
		a := A{}
		c := C{C: some}

		isoAC.FromBtoA(&c, &a)
		it.Then(t).Should(it.Equiv(a.X, c.C))
	})

	t.Run("FromBtoA_with_map", func(t *testing.T) {
		a := A{}
		m := M{"X": some}

		isoAM.FromBtoA(&m, &a)
		it.Then(t).Should(it.Equiv(a.X, m["X"]))
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
