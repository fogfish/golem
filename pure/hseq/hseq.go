package hseq

import (
	"fmt"
	"reflect"
)

// Element of product type
type Type[T any] struct {
	reflect.StructField
	PureType reflect.Type

	ID int // TODO: rename to Index
}

// Heterogenous projection of product type
type Seq[T any] []Type[T]

// Unfold type T to heterogenous sequence
func New[T any](names ...string) Seq[T] {
	cat := typeOf(*new(T))
	seq := make(Seq[T], 0)
	seq = unfold(cat, seq)

	if len(names) == 0 {
		return seq
	}

	nseq := make(Seq[T], len(names))
	for i, name := range names {
		nseq[i] = ForName(seq, name)
	}

	return nseq
}

func unfold[T any](cat reflect.Type, seq Seq[T]) Seq[T] {
	for i := 0; i < cat.NumField(); i++ {
		fv := cat.Field(i)
		ft := cat.Field(i).Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if fv.Anonymous && ft.Kind() == reflect.Struct {
			seq = unfold(ft, seq)
		} else {
			seq = append(seq, Type[T]{
				StructField: fv,
				PureType:    ft,
				ID:          len(seq),
			})
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

// Unfold type T to heterogenous sequence
func New1[T, A any]() Seq[T] {
	seq := New[T]()
	return Seq[T]{
		ForType[A](seq),
	}
}

// Unfold type T to heterogenous sequence
func New2[T, A, B any]() Seq[T] {
	seq := New[T]()
	return Seq[T]{
		ForType[A](seq),
		ForType[B](seq),
	}
}

// Unfold type T to heterogenous sequence
func New3[T, A, B, C any]() Seq[T] {
	seq := New[T]()
	return Seq[T]{
		ForType[A](seq),
		ForType[B](seq),
		ForType[C](seq),
	}
}

// Unfold type T to heterogenous sequence
func New4[T, A, B, C, D any]() Seq[T] {
	seq := New[T]()
	return Seq[T]{
		ForType[A](seq),
		ForType[B](seq),
		ForType[C](seq),
		ForType[D](seq),
	}
}

// Unfold type T to heterogenous sequence
func New5[T, A, B, C, D, E any]() Seq[T] {
	seq := New[T]()
	return Seq[T]{
		ForType[A](seq),
		ForType[B](seq),
		ForType[C](seq),
		ForType[D](seq),
		ForType[E](seq),
	}
}

// Unfold type T to heterogenous sequence
func New6[T, A, B, C, D, E, F any]() Seq[T] {
	seq := New[T]()
	return Seq[T]{
		ForType[A](seq),
		ForType[B](seq),
		ForType[C](seq),
		ForType[D](seq),
		ForType[E](seq),
		ForType[F](seq),
	}
}

// Unfold type T to heterogenous sequence
func New7[T, A, B, C, D, E, F, G any]() Seq[T] {
	seq := New[T]()
	return Seq[T]{
		ForType[A](seq),
		ForType[B](seq),
		ForType[C](seq),
		ForType[D](seq),
		ForType[E](seq),
		ForType[F](seq),
		ForType[G](seq),
	}
}

// Unfold type T to heterogenous sequence
func New8[T, A, B, C, D, E, F, G, H any]() Seq[T] {
	seq := New[T]()
	return Seq[T]{
		ForType[A](seq),
		ForType[B](seq),
		ForType[C](seq),
		ForType[D](seq),
		ForType[E](seq),
		ForType[F](seq),
		ForType[G](seq),
		ForType[H](seq),
	}
}

// Unfold type T to heterogenous sequence
func New9[T, A, B, C, D, E, F, G, H, I any]() Seq[T] {
	seq := New[T]()
	return Seq[T]{
		ForType[A](seq),
		ForType[B](seq),
		ForType[C](seq),
		ForType[D](seq),
		ForType[E](seq),
		ForType[F](seq),
		ForType[G](seq),
		ForType[H](seq),
		ForType[I](seq),
	}
}

// Lookup type heterogenous sequence by "witness" type
func ForType[A, T any](seq Seq[T]) Type[T] {
	val := typeOf(*new(A))

	for _, f := range seq {
		ft := f.Type
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}

		if ft.String() == val.String() && ft.AssignableTo(val) {
			return f
		}
	}

	cat := typeOf(*new(T))
	panic(fmt.Errorf("%s is not member of %s type", val.Name(), cat.Name()))
}

// Lookup type in heterogenous sequence by name of member
func ForName[T any](seq Seq[T], field string) Type[T] {
	for _, f := range seq {
		if f.Name == field {
			return f
		}
	}

	cat := typeOf(*new(T))
	panic(fmt.Errorf("%s is not member of %s type", field, cat.Name()))
}

// Transform heterogenous sequence to something else
func FMap[T, A any](seq Seq[T], f func(Type[T]) A) []A {
	val := make([]A, len(seq))
	for i, x := range seq {
		val[i] = f(x)
	}
	return val
}

// Assert equality of type
func Assert[T, A any](t Type[T]) (string, reflect.Kind) {
	return assertType[T, A](t, false)
}

// Assert strict equality of type
func AssertStrict[T, A any](t Type[T]) (string, reflect.Kind) {
	return assertType[T, A](t, true)
}

func assertType[T, A any](t Type[T], strict bool) (string, reflect.Kind) {
	a := reflect.TypeOf(*new(A))
	k := t.Type
	if !strict && k.Kind() == reflect.Ptr {
		k = k.Elem()
	}

	if k.Kind() != a.Kind() {
		s := typeOf(*new(T))
		panic(
			fmt.Errorf("type %s is not equal %s at %s.%s",
				t.Type.Kind(), a.Kind(), s.Name(), t.StructField.Name,
			),
		)
	}

	return a.Name(), a.Kind()
}

func FMap1[T, A any](
	ts Seq[T],
	fa func(Type[T]) A,
) A {
	return fa(ts[0])
}

func FMap2[T, A, B any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
) (A, B) {
	return fa(ts[0]),
		fb(ts[1])
}

func FMap3[T, A, B, C any](
	ts Seq[T],
	fa func(Type[T]) A,
	fb func(Type[T]) B,
	fc func(Type[T]) C,
) (A, B, C) {
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
