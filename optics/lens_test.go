package optics_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/fogfish/golem/optics"
	"github.com/fogfish/it"
)

func TestLenses(t *testing.T) {
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

//
//
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
	]()

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
	]()

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
	]()

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
	]()

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
	]()

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
	]()

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
	]()

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
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test10(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I, J string }

	a, b, c, d, e, f, g, h, i, j := optics.ForProduct10[
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
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test11(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I, J, K string }

	a, b, c, d, e, f, g, h, i, j, k := optics.ForProduct11[
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
		string,
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test12(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I, J, K, L string }

	a, b, c, d, e, f, g, h, i, j, k, l := optics.ForProduct12[
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
		string,
		string,
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test13(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I, J, K, L, M string }

	a, b, c, d, e, f, g, h, i, j, k, l, m := optics.ForProduct13[
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
		string,
		string,
		string,
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test14(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N string }

	a, b, c, d, e, f, g, h, i, j, k, l, m, n := optics.ForProduct14[
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
		string,
		string,
		string,
		string,
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test15(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O string }

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o := optics.ForProduct15[
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
		string,
		string,
		string,
		string,
		string,
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test16(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P string }

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p := optics.ForProduct16[
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
		string,
		string,
		string,
		string,
		string,
		string,
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o, "P": p,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test17(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q string }

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q := optics.ForProduct17[
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
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o, "P": p, "Q": q,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test18(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R string }

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r := optics.ForProduct18[
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
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o, "P": p, "Q": q, "R": r,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test19(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S string }

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s := optics.ForProduct19[
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
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o, "P": p, "Q": q, "R": r, "S": s,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}

func Test20(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, U string }

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s, u := optics.ForProduct20[
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
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
		string,
	]()

	for expect, f := range map[string]optics.Lens[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o, "P": p, "Q": q, "R": r, "S": s, "U": u,
	} {
		x := T{}

		it.Ok(t).
			If(f.Put(&x, expect)).Equal(nil).
			If(f.Get(&x)).Equal(expect)
	}
}
