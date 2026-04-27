//
// Copyright (C) 2022 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package optics

// Read-only lense, get value from structure, but do not modify it
func View[S, A any](lens Lens[S, A]) Lens[S, A] { return view[S, A]{lens} }

type view[S, A any] struct{ Lens[S, A] }

func (view[S, A]) Put(s *S, a A) *S { return s }

// Write-only lense, modify structure, but do not get value from it
func Update[S, A any](lens Lens[S, A]) Lens[S, A] { return update[S, A]{lens} }

type update[S, A any] struct{ Lens[S, A] }

func (update[S, A]) Get(s *S) A { return *new(A) }

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

// Automatic transformer for string superset, preserve mapping between
// two categories A, B, where both rooted to strings
func BiMapS[S any, A, B String](attr ...string) Lens[S, B] {
	return BiMap(
		ForProduct1[S, A](attr...),
		func(a A) B { return B(a) },
		func(b B) A { return A(b) },
	)
}

type String interface {
	~string
}

// Automatic transformer for []byte superset, preserve mapping between
// two categories A, B, where both rooted to []byte
func BiMapB[S any, A, B Byte](attr ...string) Lens[S, B] {
	return BiMap(
		ForProduct1[S, A](attr...),
		func(a A) B { return B(a) },
		func(b B) A { return A(b) },
	)
}

type Byte interface {
	~[]byte
}

// Automatic transformer for int superset, preserve mapping between
// two categories A, B, where both rooted to int
func BiMapI[S any, A, B Int](attr ...string) Lens[S, B] {
	return BiMap(
		ForProduct1[S, A](attr...),
		func(a A) B { return B(a) },
		func(b B) A { return A(b) },
	)
}

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Automatic transformer for float32 superset, preserve mapping between
// two categories A, B, where both rooted to float
func BiMapF[S any, A, B Float](attr ...string) Lens[S, B] {
	return BiMap(
		ForProduct1[S, A](attr...),
		func(a A) B { return B(a) },
		func(b B) A { return A(b) },
	)
}

type Float interface {
	~float32 | ~float64
}

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

func CodecS1T1[S, T, A, B any]() Isomorphism[S, T] {
	sa, sb := ForProduct2[S, A, B]()
	ta, tb := ForProduct2[T, A, B]()
	return Morphism(
		Iso(View(sa), Update(ta)),
		Iso(Update(sb), View(tb)),
	)
}

func CodecS2T1[S, T, A, B, C any]() Isomorphism[S, T] {
	sa, sb, sc := ForProduct3[S, A, B, C]()
	ta, tb, tc := ForProduct3[T, A, B, C]()
	return Morphism(
		Iso(View(sa), Update(ta)),
		Iso(View(sb), Update(tb)),
		Iso(Update(sc), View(tc)),
	)
}

func CodecS2T2[S, T, A, B, C, D any]() Isomorphism[S, T] {
	sa, sb, sc, sd := ForProduct4[S, A, B, C, D]()
	ta, tb, tc, td := ForProduct4[T, A, B, C, D]()
	return Morphism(
		Iso(View(sa), Update(ta)),
		Iso(View(sb), Update(tb)),
		Iso(Update(sc), View(tc)),
		Iso(Update(sd), View(td)),
	)
}
