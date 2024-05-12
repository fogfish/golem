//
// Copyright (C) 2022 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package optics

// Getter composes lense with transform function, making read-only lense
func Getter[S, A, B any](lens Lens[S, A], f func(A) B) Lens[S, B] {
	return fmap[S, A, B]{lens, f}
}

type fmap[S, A, B any] struct {
	lens Lens[S, A]
	f    func(A) B
}

func (c fmap[S, A, B]) Put(s *S, b B) *S { return s }
func (c fmap[S, A, B]) Get(s *S) B       { return c.f(c.lens.Get(s)) }

// Setter composes lense with inverse transformer, making write-only lense
func Setter[S, A, B any](lens Lens[S, A], f func(B) A) Lens[S, B] {
	return cmap[S, A, B]{lens, f}
}

type cmap[S, A, B any] struct {
	lens Lens[S, A]
	f    func(B) A
}

func (c cmap[S, A, B]) Put(s *S, b B) *S { return c.lens.Put(s, c.f(b)) }
func (c cmap[S, A, B]) Get(s *S) B       { return *new(B) }

// Transformer lens is a structure-preserving mapping between two categories A, B.
func BiMap[S, A, B any](
	lens Lens[S, A],
	fmap func(A) B,
	cmap func(B) A,
) Lens[S, B] {
	return codec[S, A, B]{lens: lens, fmap: fmap, cmap: cmap}
}

type codec[S, A, B any] struct {
	lens Lens[S, A]
	fmap func(A) B
	cmap func(B) A
}

func (c codec[S, A, B]) Put(s *S, b B) *S { return c.lens.Put(s, c.cmap(b)) }
func (c codec[S, A, B]) Get(s *S) B       { return c.fmap(c.lens.Get(s)) }

// An isomorphism is a structure-preserving mapping between two structures
// of the same shape that can be reversed by an inverse mapping.
type Isomorphism[S, T any] interface {
	Forward(*S, *T)
	Inverse(*T, *S)
}

// Building Isomorphism from joining lenses
func Iso[S, T, A any](
	sa Lens[S, A],
	ta Lens[T, A],
) Isomorphism[S, T] {
	return iso[S, T, A]{sa, ta}
}

type iso[S, T, A any] struct {
	sa Lens[S, A]
	ta Lens[T, A]
}

func (iso iso[S, T, A]) Forward(s *S, t *T) { iso.ta.Put(t, iso.sa.Get(s)) }
func (iso iso[S, T, A]) Inverse(t *T, s *S) { iso.sa.Put(s, iso.ta.Get(t)) }

// Build a structure-preserving mapping between two structures
func Morphism[S, T any](seq ...Isomorphism[S, T]) Isomorphism[S, T] {
	return morphism[S, T](seq)
}

type morphism[S, T any] []Isomorphism[S, T]

func (seq morphism[S, T]) Forward(s *S, t *T) {
	for _, iso := range seq {
		if iso != nil {
			iso.Forward(s, t)
		}
	}
}

func (seq morphism[S, T]) Inverse(t *T, s *S) {
	for _, iso := range seq {
		if iso != nil {
			iso.Inverse(t, s)
		}
	}
}
