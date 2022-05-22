package hseq_test

import (
	"reflect"
	"testing"

	"github.com/fogfish/golem/pure/hseq"
	"github.com/fogfish/it"
)

type Identity[S, A any] struct{ hseq.Type[S] }

func (id Identity[S, A]) Value(s S) A {
	f := reflect.ValueOf(s).FieldByIndex(id.Index)
	return f.Interface().(A)
}

func mkIdentity[S, A any](s hseq.Type[S]) Identity[S, A] {
	hseq.AssertType[S, A](s, true)
	return Identity[S, A]{s}
}

func TestSafeSeq(t *testing.T) {
	type T struct {
		A string
	}

	it.Ok(t).If(
		func() {
			hseq.FMap2(
				hseq.Generic[T](),
				mkIdentity[T, string],
				mkIdentity[T, string],
			)
		},
	).Should().Fail()
}

func TestSafeType(t *testing.T) {
	type T struct {
		A string
	}

	it.Ok(t).If(
		func() {
			hseq.FMap1(
				hseq.Generic[T](),
				mkIdentity[T, int],
			)
		},
	).Should().Fail()
}

func Test1(t *testing.T) {
	type T struct {
		A string
	}
	v := T{"A"}

	a := hseq.FMap1(
		hseq.Generic[T](),
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test2(t *testing.T) {
	type T struct {
		A, B string
	}
	v := T{"A", "B"}

	a, b := hseq.FMap2(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test3(t *testing.T) {
	type T struct {
		A, B, C string
	}
	v := T{"A", "B", "C"}

	a, b, c := hseq.FMap3(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test4(t *testing.T) {
	type T struct {
		A, B, C, D string
	}
	v := T{"A", "B", "C", "D"}

	a, b, c, d := hseq.FMap4(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test5(t *testing.T) {
	type T struct {
		A, B, C, D, E string
	}
	v := T{"A", "B", "C", "D", "E"}

	a, b, c, d, e := hseq.FMap5(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test6(t *testing.T) {
	type T struct {
		A, B, C, D, E, F string
	}
	v := T{"A", "B", "C", "D", "E", "F"}

	a, b, c, d, e, f := hseq.FMap6(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test7(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G"}

	a, b, c, d, e, f, g := hseq.FMap7(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test8(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H"}

	a, b, c, d, e, f, g, h := hseq.FMap8(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test9(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

	a, b, c, d, e, f, g, h, i := hseq.FMap9(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test10(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I, J string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	a, b, c, d, e, f, g, h, i, j := hseq.FMap10(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test11(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I, J, K string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}

	a, b, c, d, e, f, g, h, i, j, k := hseq.FMap11(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test12(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I, J, K, L string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"}

	a, b, c, d, e, f, g, h, i, j, k, l := hseq.FMap12(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test13(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I, J, K, L, M string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M"}

	a, b, c, d, e, f, g, h, i, j, k, l, m := hseq.FMap13(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test14(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I, J, K, L, M, N string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N"}

	a, b, c, d, e, f, g, h, i, j, k, l, m, n := hseq.FMap14(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test15(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I, J, K, L, M, N, O string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O"}

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o := hseq.FMap15(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test16(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P"}

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p := hseq.FMap16(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o, "P": p,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test17(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q"}

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q := hseq.FMap17(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o, "P": p, "Q": q,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test18(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R"}

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r := hseq.FMap18(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o, "P": p, "Q": q, "R": r,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test19(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S"}

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s := hseq.FMap19(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o, "P": p, "Q": q, "R": r, "S": s,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}

func Test20(t *testing.T) {
	type T struct {
		A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, U string
	}
	v := T{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "U"}

	a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s, u := hseq.FMap20(
		hseq.Generic[T](),
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
		mkIdentity[T, string],
	)

	for expect, f := range map[string]Identity[T, string]{
		"A": a, "B": b, "C": c, "D": d, "E": e, "F": f, "G": g, "H": h, "I": i, "J": j, "K": k, "L": l, "M": m, "N": n, "O": o, "P": p, "Q": q, "R": r, "S": s, "U": u,
	} {
		it.Ok(t).If(f.Value(v)).Equal(expect)
	}
}
