//
// Copyright (C) 2022 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package optics

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/fogfish/golem/hseq"
)

// Lens resembles concept of getters and setters, which you can compose
// using functional concepts. In other words, this is combinator data
// transformation for pure functional data structure.
//
// Lens allows to abstract a "shape" of the structure rather than type itself.
type Lens[S, A any] interface {
	Get(*S) A
	Put(*S, A) *S
}

// NewLens instantiates a typed Lens[S, A] for hseq.Type[S]
func NewLens[S, A any](t hseq.Type[S]) Lens[S, A] {
	ft := t.Type
	fv := reflect.TypeOf(new(A)).Elem()

	if ft.String() == fv.String() && ft.AssignableTo(fv) {
		return &lens[S, A]{t}
	}

	cat := reflect.TypeOf(new(S)).Elem()
	panic(fmt.Errorf("invalid type: Lens[%s, %s] not compatible with %s", cat.Name(), ft.Name(), fv.Name()))
}

type lens[S, A any] struct{ hseq.Type[S] }

func (lens *lens[S, A]) Put(s *S, a A) *S {
	*(*A)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + lens.Offset + lens.RootOffs)) = a
	return s
}

func (lens *lens[S, A]) Get(s *S) A {
	return *(*A)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + lens.Offset + lens.RootOffs))
}

// NewLens instantiates a typed Lens[S, A] for map[K]A
func NewLensM[K comparable, A any](key K) Lens[map[K]A, A] {
	return &lensM[K, A]{key}
}

type lensM[K comparable, A any] struct{ key K }

func (lens *lensM[K, A]) Put(s *map[K]A, a A) *map[K]A {
	(*s)[lens.key] = a
	return s
}

func (lens *lensM[K, A]) Get(s *map[K]A) A {
	return (*s)[lens.key]
}

func Join[S, A, B any](a Lens[S, A], b Lens[A, B]) Lens[S, B] {
	return join[S, A, B]{a, b}
}

type join[S, A, B any] struct {
	a Lens[S, A]
	b Lens[A, B]
}

func (lens join[S, A, B]) Put(s *S, b B) *S {
	va := lens.a.Get(s)
	lens.b.Put(&va, b)
	lens.a.Put(s, va)
	return s
}

func (lens join[S, A, B]) Get(s *S) B {
	va := lens.a.Get(s)
	return lens.b.Get(&va)
}

// ForProduct1 unfold 1 attribute of type T
func ForProduct1[T, A any](attr ...string) Lens[T, A] {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New1[T, A]()
	} else {
		seq = hseq.New[T](attr[0])
	}

	return hseq.FMap1(seq,
		NewLens[T, A],
	)
}

// ForProduct2 unfold 2 attribute of type T
func ForProduct2[T, A, B any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New2[T, A, B]()
	} else {
		seq = hseq.New[T](attr[0:2]...)
	}

	return hseq.FMap2(seq,
		NewLens[T, A],
		NewLens[T, B],
	)
}

// ForProduct3 unfold 3 attribute of type T
func ForProduct3[T, A, B, C any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New3[T, A, B, C]()
	} else {
		seq = hseq.New[T](attr[0:3]...)
	}

	return hseq.FMap3(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
	)
}

// ForProduct4 unfold 4 attribute of type T
func ForProduct4[T, A, B, C, D any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New4[T, A, B, C, D]()
	} else {
		seq = hseq.New[T](attr[0:4]...)
	}

	return hseq.FMap4(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
	)
}

// ForProduct5 unfold 5 attribute of type T
func ForProduct5[T, A, B, C, D, E any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New5[T, A, B, C, D, E]()
	} else {
		seq = hseq.New[T](attr[0:5]...)
	}

	return hseq.FMap5(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
		NewLens[T, E],
	)
}

// ForProduct6 unfold 6 attribute of type T
func ForProduct6[T, A, B, C, D, E, F any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New6[T, A, B, C, D, E, F]()
	} else {
		seq = hseq.New[T](attr[0:6]...)
	}

	return hseq.FMap6(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
		NewLens[T, E],
		NewLens[T, F],
	)
}

// ForProduct7 unfold 7 attribute of type T
func ForProduct7[T, A, B, C, D, E, F, G any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New7[T, A, B, C, D, E, F, G]()
	} else {
		seq = hseq.New[T](attr[0:7]...)
	}

	return hseq.FMap7(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
		NewLens[T, E],
		NewLens[T, F],
		NewLens[T, G],
	)
}

// ForProduct8 unfold 8 attribute of type T
func ForProduct8[T, A, B, C, D, E, F, G, H any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New8[T, A, B, C, D, E, F, G, H]()
	} else {
		seq = hseq.New[T](attr[0:8]...)
	}

	return hseq.FMap8(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
		NewLens[T, E],
		NewLens[T, F],
		NewLens[T, G],
		NewLens[T, H],
	)
}

// ForProduct9 unfold 9 attribute of type T
func ForProduct9[T, A, B, C, D, E, F, G, H, I any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New9[T, A, B, C, D, E, F, G, H, I]()
	} else {
		seq = hseq.New[T](attr[0:9]...)
	}

	return hseq.FMap9(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
		NewLens[T, E],
		NewLens[T, F],
		NewLens[T, G],
		NewLens[T, H],
		NewLens[T, I],
	)
}
