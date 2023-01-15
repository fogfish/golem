package optics_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/fogfish/golem/optics"
	"github.com/fogfish/it"
)

func TestLensesByType(t *testing.T) {
	type Inner struct {
		E, F, G string
	}

	type T struct {
		A string
		B int
		C float64
		D Inner
	}

	a, b, c, d := optics.ForProduct4[T, string, int, float64, Inner]()

	t.Run("String", func(t *testing.T) {
		x := T{}

		it.Ok(t).
			If(a.Put(&x, "A")).Equal(nil).
			If(x.A).Equal("A").
			If(a.Get(&x)).Equal("A")
	})

	t.Run("Int", func(t *testing.T) {
		x := T{}

		it.Ok(t).
			If(b.Put(&x, 1234)).Equal(nil).
			If(x.B).Equal(1234).
			If(b.Get(&x)).Equal(1234)
	})

	t.Run("Float", func(t *testing.T) {
		x := T{}

		it.Ok(t).
			If(c.Put(&x, 1234.0)).Equal(nil).
			If(x.C).Equal(1234.0).
			If(c.Get(&x)).Equal(1234.0)
	})

	t.Run("Struct", func(t *testing.T) {
		x := T{}

		it.Ok(t).
			If(d.Put(&x, Inner{"E", "F", "G"})).Equal(nil).
			If(x.D).Equal(Inner{"E", "F", "G"}).
			If(d.Get(&x)).Equal(Inner{"E", "F", "G"})
	})

	t.Run("Struct.Codec", func(t *testing.T) {
		x := T{}
		l := newLensJSON(d)

		it.Ok(t).
			If(l.Put(&x, "{\"E\":\"E\",\"F\":\"F\",\"G\":\"G\"}")).Equal(nil).
			If(x.D).Equal(Inner{"E", "F", "G"}).
			If(l.Get(&x)).Equal("{\"E\":\"E\",\"F\":\"F\",\"G\":\"G\"}")
	})
}

func TestLensesByName(t *testing.T) {
	type Inner struct {
		E, F, G string
	}

	type T struct {
		A *string
		B *int
		C *float64
		D *Inner
	}

	a, b, c, d := optics.ForProduct4[T, string, int, float64, Inner]("A", "B", "C", "D")

	t.Run("String", func(t *testing.T) {
		x := T{}

		it.Ok(t).
			If(a.Put(&x, "A")).Equal(nil).
			If(*x.A).Equal("A").
			If(a.Get(&x)).Equal("A")
	})

	t.Run("Int", func(t *testing.T) {
		x := T{}

		it.Ok(t).
			If(b.Put(&x, 1234)).Equal(nil).
			If(*x.B).Equal(1234).
			If(b.Get(&x)).Equal(1234)
	})

	t.Run("Float", func(t *testing.T) {
		x := T{}

		it.Ok(t).
			If(c.Put(&x, 1234.0)).Equal(nil).
			If(*x.C).Equal(1234.0).
			If(c.Get(&x)).Equal(1234.0)
	})

	t.Run("Struct", func(t *testing.T) {
		x := T{}

		it.Ok(t).
			If(d.Put(&x, Inner{"E", "F", "G"})).Equal(nil).
			If(*x.D).Equal(Inner{"E", "F", "G"}).
			If(d.Get(&x)).Equal(Inner{"E", "F", "G"})
	})

	t.Run("Struct.Codec", func(t *testing.T) {
		x := T{}
		l := newLensJSON(d)

		it.Ok(t).
			If(l.Put(&x, "{\"E\":\"E\",\"F\":\"F\",\"G\":\"G\"}")).Equal(nil).
			If(*x.D).Equal(Inner{"E", "F", "G"}).
			If(l.Get(&x)).Equal("{\"E\":\"E\",\"F\":\"F\",\"G\":\"G\"}")
	})
}

func TestLensesCustomTypes(t *testing.T) {
	type MyString string
	type MyInt int
	type MyFloat float64

	type T struct {
		A MyString
		B MyInt
		C MyFloat
	}

	a, b, c := optics.ForProduct3[T, MyString, MyInt, MyFloat]()

	t.Run("String", func(t *testing.T) {
		x := T{}

		it.Ok(t).
			If(a.Put(&x, "A")).Equal(nil).
			If(x.A).Equal(MyString("A")).
			If(a.Get(&x)).Equal(MyString("A"))
	})

	t.Run("Int", func(t *testing.T) {
		x := T{}

		it.Ok(t).
			If(b.Put(&x, 1234)).Equal(nil).
			If(x.B).Equal(MyInt(1234)).
			If(b.Get(&x)).Equal(MyInt(1234))
	})

	t.Run("Float", func(t *testing.T) {
		x := T{}

		it.Ok(t).
			If(c.Put(&x, 1234.0)).Equal(nil).
			If(x.C).Equal(MyFloat(1234.0)).
			If(c.Get(&x)).Equal(MyFloat(1234.0))
	})
}

func TestMorphism(t *testing.T) {
	type T struct{ A string }
	a := optics.ForProduct1[T, string]()

	m := optics.Morph(a, "hello")

	x := T{}
	y := T{}

	it.Ok(t).
		If(m.Put(&x)).Equal(nil).
		If(x.A).Equal("hello").
		If(m.PutValue(reflect.ValueOf(&y))).Equal(nil).
		If(y.A).Equal("hello")
}

func TestMorphisms(t *testing.T) {
	type T struct {
		A string
		B int
		C float64
	}
	a, b, c := optics.ForProduct3[T, string, int, float64]()

	m := optics.Morphisms[T]{
		optics.Morph(a, "hello"),
		optics.Morph(b, 1234),
		optics.Morph(c, 1234.0),
	}

	x := T{}

	it.Ok(t).
		If(m.Put(&x)).Equal(nil).
		If(x.A).Equal("hello").
		If(x.B).Equal(1234).
		If(x.C).Equal(1234.0)
}

/*
Custom Lens to parse JSON
*/
type lensStructJSON[S, A any] struct{ optics.Lens[S, A] }

func newLensJSON[S, A any](l optics.Lens[S, A]) optics.Lens[S, string] {
	return &lensStructJSON[S, A]{l}
}

func (l lensStructJSON[S, A]) Put(s *S, a string) error {
	var o A

	if err := json.Unmarshal([]byte(a), &o); err != nil {
		return err
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

func Test1(t *testing.T) {
	type T struct{ A string }

	a := optics.ForProduct1[
		T,
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test2(t *testing.T) {
	type T struct{ A, B string }

	a, b := optics.ForProduct2[
		T,
		string,
		string,
	]("A", "B")

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test3(t *testing.T) {
	type T struct{ A, B, C string }

	a, b, c := optics.ForProduct3[
		T,
		string,
		string,
		string,
	]("A", "B", "C")

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test4(t *testing.T) {
	type T struct{ A, B, C, D string }

	a, b, c, d := optics.ForProduct4[
		T,
		string,
		string,
		string,
		string,
	]("A", "B", "C", "D")

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test5(t *testing.T) {
	type T struct{ A, B, C, D, E string }

	a, b, c, d, e := optics.ForProduct5[
		T,
		string,
		string,
		string,
		string,
		string,
	]("A", "B", "C", "D", "E")

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test6(t *testing.T) {
	type T struct{ A, B, C, D, E, F string }

	a, b, c, d, e, f := optics.ForProduct6[
		T,
		string,
		string,
		string,
		string,
		string,
		string,
	]("A", "B", "C", "D", "E", "F")

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test7(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G string }

	a, b, c, d, e, f, g := optics.ForProduct7[
		T,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
	]("A", "B", "C", "D", "E", "F", "G")

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test8(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H string }

	a, b, c, d, e, f, g, h := optics.ForProduct8[
		T,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
	]("A", "B", "C", "D", "E", "F", "G", "H")

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test9(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I string }

	a, b, c, d, e, f, g, h, i := optics.ForProduct9[
		T,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
	]("A", "B", "C", "D", "E", "F", "G", "H", "I")

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}
