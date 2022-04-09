package optics_test

import (
	"testing"

	"github.com/fogfish/golem/optics"
)

type MyT1 struct {
	Name string
}

var (
	name  = optics.ForProduct1[MyT1, string]()
	nameM = optics.Morphisms[MyT1]{
		optics.Morph(name, "hello"),
	}
)

func BenchmarkLensForProduct1(mb *testing.B) {
	var val MyT1

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		name.Put(&val, "name")
	}
}

func BenchmarkMorphismForProduct1(mb *testing.B) {
	var val MyT1

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		nameM.Put(&val)
	}
}
