package hseq

import (
	"fmt"
	"reflect"
)

/*

Type ...
*/
type Type[T any] struct {
	reflect.StructField

	ID       int
	PureType reflect.Type
}

/*

Seq ...
*/
type Seq[T any] []Type[T]

func AssertType[T, A any](t Type[T], strict bool) (string, reflect.Kind) {
	a := reflect.TypeOf(*new(A))
	k := t.Type
	if !strict && k.Kind() == reflect.Ptr {
		k = k.Elem()
	}

	if k.Kind() != a.Kind() {
		s := typeOf(*new(T))
		panic(
			fmt.Errorf("Type %s is not equal %s at %s.%s",
				t.Type.Kind(), a.Kind(), s.Name(), t.StructField.Name,
			),
		)
	}

	return a.Name(), a.Kind()
}

func AssertSeq[T any](list Seq[T], n int) {
	if len(list) != n {
		t := typeOf(*new(T))
		panic(fmt.Errorf("Unable to map type |%s| = %d to hseq of %d", t.Name(), t.NumField(), n))
	}
}

/*

Generic ...
*/
func Generic[T any]() Seq[T] {
	t := typeOf(*new(T))
	seq := make(Seq[T], t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fv := t.Field(i)
		ft := t.Field(i).Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		seq[i] = Type[T]{
			StructField: fv,
			ID:          i,
			PureType:    ft,
		}
	}
	return seq
}

func typeOf[T any](t T) reflect.Type {
	typeof := reflect.TypeOf(t)
	if typeof.Kind() == reflect.Ptr {
		typeof = typeof.Elem()
	}

	return typeof
}

/*

FMap ...
*/
func FMap[T, A any](list Seq[T], f func(Type[T]) A) []A {
	seq := make([]A, len(list))
	for i, x := range list {
		seq[i] = f(x)
	}
	return seq
}

func FMap1[T, A any](
	ts Seq[T],
	fa func(Type[T]) A,
) A {
	AssertSeq(ts, 1)
	return fa(ts[0])
}

func FMap2[T, A, B any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
) (A, B) {
	AssertSeq(ts, 2)
	return fa(ts[0]),
		fb(ts[1])
}

func FMap3[T, A, B, C any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
) (A, B, C) {
	AssertSeq(ts, 3)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2])
}

func FMap4[T, A, B, C, D any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
) (A, B, C, D) {
	AssertSeq(ts, 4)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3])
}

func FMap5[T, A, B, C, D, E any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
) (A, B, C, D, E) {
	AssertSeq(ts, 5)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4])
}

func FMap6[T, A, B, C, D, E, F any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
) (A, B, C, D, E, F) {
	AssertSeq(ts, 6)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5])
}

func FMap7[T, A, B, C, D, E, F, G any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
) (A, B, C, D, E, F, G) {
	AssertSeq(ts, 7)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6])
}

func FMap8[T, A, B, C, D, E, F, G, H any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
) (A, B, C, D, E, F, G, H) {
	AssertSeq(ts, 8)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7])
}

func FMap9[T, A, B, C, D, E, F, G, H, I any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
) (A, B, C, D, E, F, G, H, I) {
	AssertSeq(ts, 9)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8])
}

func FMap10[T, A, B, C, D, E, F, G, H, I, J any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
	fj func(Type[T]) J,
) (A, B, C, D, E, F, G, H, I, J) {
	AssertSeq(ts, 10)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8]),
		fj(ts[9])
}

func FMap11[T, A, B, C, D, E, F, G, H, I, J, K any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
	fj func(Type[T]) J,
	fk func(Type[T]) K,
) (A, B, C, D, E, F, G, H, I, J, K) {
	AssertSeq(ts, 11)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8]),
		fj(ts[9]),
		fk(ts[10])
}

func FMap12[T, A, B, C, D, E, F, G, H, I, J, K, L any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
	fj func(Type[T]) J,
	fk func(Type[T]) K,
	fl func(Type[T]) L,
) (A, B, C, D, E, F, G, H, I, J, K, L) {
	AssertSeq(ts, 12)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8]),
		fj(ts[9]),
		fk(ts[10]),
		fl(ts[11])
}

func FMap13[T, A, B, C, D, E, F, G, H, I, J, K, L, M any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
	fj func(Type[T]) J,
	fk func(Type[T]) K,
	fl func(Type[T]) L,
	fm func(Type[T]) M,
) (A, B, C, D, E, F, G, H, I, J, K, L, M) {
	AssertSeq(ts, 13)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8]),
		fj(ts[9]),
		fk(ts[10]),
		fl(ts[11]),
		fm(ts[12])
}

func FMap14[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
	fj func(Type[T]) J,
	fk func(Type[T]) K,
	fl func(Type[T]) L,
	fm func(Type[T]) M,
	fn func(Type[T]) N,
) (A, B, C, D, E, F, G, H, I, J, K, L, M, N) {
	AssertSeq(ts, 14)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8]),
		fj(ts[9]),
		fk(ts[10]),
		fl(ts[11]),
		fm(ts[12]),
		fn(ts[13])
}

func FMap15[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
	fj func(Type[T]) J,
	fk func(Type[T]) K,
	fl func(Type[T]) L,
	fm func(Type[T]) M,
	fn func(Type[T]) N,
	fo func(Type[T]) O,
) (A, B, C, D, E, F, G, H, I, J, K, L, M, N, O) {
	AssertSeq(ts, 15)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8]),
		fj(ts[9]),
		fk(ts[10]),
		fl(ts[11]),
		fm(ts[12]),
		fn(ts[13]),
		fo(ts[14])
}

func FMap16[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
	fj func(Type[T]) J,
	fk func(Type[T]) K,
	fl func(Type[T]) L,
	fm func(Type[T]) M,
	fn func(Type[T]) N,
	fo func(Type[T]) O,
	fp func(Type[T]) P,
) (A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P) {
	AssertSeq(ts, 16)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8]),
		fj(ts[9]),
		fk(ts[10]),
		fl(ts[11]),
		fm(ts[12]),
		fn(ts[13]),
		fo(ts[14]),
		fp(ts[15])
}

func FMap17[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
	fj func(Type[T]) J,
	fk func(Type[T]) K,
	fl func(Type[T]) L,
	fm func(Type[T]) M,
	fn func(Type[T]) N,
	fo func(Type[T]) O,
	fp func(Type[T]) P,
	fq func(Type[T]) Q,
) (A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q) {
	AssertSeq(ts, 17)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8]),
		fj(ts[9]),
		fk(ts[10]),
		fl(ts[11]),
		fm(ts[12]),
		fn(ts[13]),
		fo(ts[14]),
		fp(ts[15]),
		fq(ts[16])
}

func FMap18[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
	fj func(Type[T]) J,
	fk func(Type[T]) K,
	fl func(Type[T]) L,
	fm func(Type[T]) M,
	fn func(Type[T]) N,
	fo func(Type[T]) O,
	fp func(Type[T]) P,
	fq func(Type[T]) Q,
	fr func(Type[T]) R,
) (A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R) {
	AssertSeq(ts, 18)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8]),
		fj(ts[9]),
		fk(ts[10]),
		fl(ts[11]),
		fm(ts[12]),
		fn(ts[13]),
		fo(ts[14]),
		fp(ts[15]),
		fq(ts[16]),
		fr(ts[17])
}

func FMap19[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
	fj func(Type[T]) J,
	fk func(Type[T]) K,
	fl func(Type[T]) L,
	fm func(Type[T]) M,
	fn func(Type[T]) N,
	fo func(Type[T]) O,
	fp func(Type[T]) P,
	fq func(Type[T]) Q,
	fr func(Type[T]) R,
	fs func(Type[T]) S,
) (A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S) {
	AssertSeq(ts, 19)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8]),
		fj(ts[9]),
		fk(ts[10]),
		fl(ts[11]),
		fm(ts[12]),
		fn(ts[13]),
		fo(ts[14]),
		fp(ts[15]),
		fq(ts[16]),
		fr(ts[17]),
		fs(ts[18])
}

func FMap20[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, U any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
	fd func(Type[T]) D,
	fe func(Type[T]) E,
	ff func(Type[T]) F,
	fg func(Type[T]) G,
	fh func(Type[T]) H,
	fi func(Type[T]) I,
	fj func(Type[T]) J,
	fk func(Type[T]) K,
	fl func(Type[T]) L,
	fm func(Type[T]) M,
	fn func(Type[T]) N,
	fo func(Type[T]) O,
	fp func(Type[T]) P,
	fq func(Type[T]) Q,
	fr func(Type[T]) R,
	fs func(Type[T]) S,
	fu func(Type[T]) U,
) (A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, U) {
	AssertSeq(ts, 20)
	return fa(ts[0]),
		fb(ts[1]),
		fc(ts[2]),
		fd(ts[3]),
		fe(ts[4]),
		ff(ts[5]),
		fg(ts[6]),
		fh(ts[7]),
		fi(ts[8]),
		fj(ts[9]),
		fk(ts[10]),
		fl(ts[11]),
		fm(ts[12]),
		fn(ts[13]),
		fo(ts[14]),
		fp(ts[15]),
		fq(ts[16]),
		fr(ts[17]),
		fs(ts[18]),
		fu(ts[19])
}
