package optics_test

import (
	"testing"

	"github.com/fogfish/golem/optics"
)

type MyT1 struct{ Name string }

type MyT5 struct{ A, B, C, D, E string }

var (
	name      = optics.ForProduct1[MyT1, string]()
	morphMyT1 = optics.Morphisms[MyT1]{
		optics.Morph(name, "hello"),
	}
	a, b, c, d, e = optics.ForProduct5[MyT5, string, string, string, string, string]()
	morphMyT5     = optics.Morphisms[MyT5]{
		optics.Morph(a, "a"),
		optics.Morph(b, "b"),
		optics.Morph(c, "c"),
		optics.Morph(d, "d"),
		optics.Morph(e, "e"),
	}
)

func BenchmarkLensPutForProduct1(mb *testing.B) {
	var val MyT1

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		_ = name.Put(&val, "name")
	}
}

func BenchmarkLensGetForProduct1(mb *testing.B) {
	val := MyT1{Name: "name"}

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		_ = name.Get(&val)
	}
}

func BenchmarkMorphismForProduct1(mb *testing.B) {
	var val MyT1

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		_ = morphMyT1.Put(&val)
	}
}

func BenchmarkMorphismForProduct5(mb *testing.B) {
	var val MyT5

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		_ = morphMyT5.Put(&val)
	}
}
