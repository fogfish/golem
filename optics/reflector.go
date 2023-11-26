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

// Reflector is a Reflector over value of any type.
// It fails runtime if client submits invalid type.
type Reflector[A any] interface {
	Gett(any) A
	Putt(any, A) any
}

func (lens *lens[S, A]) Putt(s any, a A) any {
	switch v := s.(type) {
	case *S:
		*(*A)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + lens.Offset + lens.RootOffs)) = a
		return s
	default:
		panic(fmt.Errorf("invalid type %T passed to Reflector[%T, %T]", s, *new(S), *new(A)))
	}
}

func (lens *lens[S, A]) Gett(s any) A {
	switch v := s.(type) {
	case *S:
		return *(*A)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + lens.Offset + lens.RootOffs))
	default:
		panic(fmt.Errorf("invalid type %T passed to Reflector[%T, %T]", s, *new(S), *new(A)))
	}
}

// NewReflector instantiates a typed Reflector[S, A] for hseq.Type[S]
func NewReflector[S, A any](t hseq.Type[S]) Reflector[A] {
	ft := t.Type
	fv := reflect.TypeOf(new(A)).Elem()

	if ft.String() == fv.String() && ft.AssignableTo(fv) {
		return &lens[S, A]{t}
	}

	cat := reflect.TypeOf(new(S)).Elem()
	panic(fmt.Errorf("invalid type: Reflector[%s, %s] not compatible with %s", cat.Name(), ft.Name(), fv.Name()))
}

// ForSpectrum1 unfold 1 attribute of type T
func ForSpectrum1[T, A any](attr ...string) Reflector[A] {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New1[T, A]()
	} else {
		seq = hseq.New[T](attr[0])
	}

	return hseq.FMap1(seq,
		NewReflector[T, A],
	)
}

// ForSpectrum2 unfold 2 attribute of type T
func ForSpectrum2[T, A, B any](attr ...string) (
	Reflector[A],
	Reflector[B],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New2[T, A, B]()
	} else {
		seq = hseq.New[T](attr[0:2]...)
	}

	return hseq.FMap2(seq,
		NewReflector[T, A],
		NewReflector[T, B],
	)
}

// ForSpectrum3 unfold 3 attribute of type T
func ForSpectrum3[T, A, B, C any](attr ...string) (
	Reflector[A],
	Reflector[B],
	Reflector[C],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New3[T, A, B, C]()
	} else {
		seq = hseq.New[T](attr[0:3]...)
	}

	return hseq.FMap3(seq,
		NewReflector[T, A],
		NewReflector[T, B],
		NewReflector[T, C],
	)
}

// ForSpectrum4 unfold 4 attribute of type T
func ForSpectrum4[T, A, B, C, D any](attr ...string) (
	Reflector[A],
	Reflector[B],
	Reflector[C],
	Reflector[D],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New4[T, A, B, C, D]()
	} else {
		seq = hseq.New[T](attr[0:4]...)
	}

	return hseq.FMap4(seq,
		NewReflector[T, A],
		NewReflector[T, B],
		NewReflector[T, C],
		NewReflector[T, D],
	)
}

// ForSpectrum5 unfold 5 attribute of type T
func ForSpectrum5[T, A, B, C, D, E any](attr ...string) (
	Reflector[A],
	Reflector[B],
	Reflector[C],
	Reflector[D],
	Reflector[E],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New5[T, A, B, C, D, E]()
	} else {
		seq = hseq.New[T](attr[0:5]...)
	}

	return hseq.FMap5(seq,
		NewReflector[T, A],
		NewReflector[T, B],
		NewReflector[T, C],
		NewReflector[T, D],
		NewReflector[T, E],
	)
}

// ForSpectrum6 unfold 6 attribute of type T
func ForSpectrum6[T, A, B, C, D, E, F any](attr ...string) (
	Reflector[A],
	Reflector[B],
	Reflector[C],
	Reflector[D],
	Reflector[E],
	Reflector[F],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New6[T, A, B, C, D, E, F]()
	} else {
		seq = hseq.New[T](attr[0:6]...)
	}

	return hseq.FMap6(seq,
		NewReflector[T, A],
		NewReflector[T, B],
		NewReflector[T, C],
		NewReflector[T, D],
		NewReflector[T, E],
		NewReflector[T, F],
	)
}

// ForSpectrum7 unfold 7 attribute of type T
func ForSpectrum7[T, A, B, C, D, E, F, G any](attr ...string) (
	Reflector[A],
	Reflector[B],
	Reflector[C],
	Reflector[D],
	Reflector[E],
	Reflector[F],
	Reflector[G],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New7[T, A, B, C, D, E, F, G]()
	} else {
		seq = hseq.New[T](attr[0:7]...)
	}

	return hseq.FMap7(seq,
		NewReflector[T, A],
		NewReflector[T, B],
		NewReflector[T, C],
		NewReflector[T, D],
		NewReflector[T, E],
		NewReflector[T, F],
		NewReflector[T, G],
	)
}

// ForSpectrum8 unfold 8 attribute of type T
func ForSpectrum8[T, A, B, C, D, E, F, G, H any](attr ...string) (
	Reflector[A],
	Reflector[B],
	Reflector[C],
	Reflector[D],
	Reflector[E],
	Reflector[F],
	Reflector[G],
	Reflector[H],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New8[T, A, B, C, D, E, F, G, H]()
	} else {
		seq = hseq.New[T](attr[0:8]...)
	}

	return hseq.FMap8(seq,
		NewReflector[T, A],
		NewReflector[T, B],
		NewReflector[T, C],
		NewReflector[T, D],
		NewReflector[T, E],
		NewReflector[T, F],
		NewReflector[T, G],
		NewReflector[T, H],
	)
}

// ForSpectrum9 unfold 9 attribute of type T
func ForSpectrum9[T, A, B, C, D, E, F, G, H, I any](attr ...string) (
	Reflector[A],
	Reflector[B],
	Reflector[C],
	Reflector[D],
	Reflector[E],
	Reflector[F],
	Reflector[G],
	Reflector[H],
	Reflector[I],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New9[T, A, B, C, D, E, F, G, H, I]()
	} else {
		seq = hseq.New[T](attr[0:9]...)
	}

	return hseq.FMap9(seq,
		NewReflector[T, A],
		NewReflector[T, B],
		NewReflector[T, C],
		NewReflector[T, D],
		NewReflector[T, E],
		NewReflector[T, F],
		NewReflector[T, G],
		NewReflector[T, H],
		NewReflector[T, I],
	)
}
