//
// Copyright (C) 2022 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package optics

//
// Rank 2 Shape: Product Lens A × B
//

// Rank 2 Shape: Product Lens A × B
type Lens2[S, A, B any] interface {
	Get(*S) (A, B)
	Put(*S, A, B) *S
}

// Rank 2 Shape: Product Lens A × B
func ForShape2[T, A, B any](attr ...string) Lens2[T, A, B] {
	a, b := ForProduct2[T, A, B](attr...)
	return shape2[T, A, B]{a: a, b: b}
}

type shape2[T, A, B any] struct {
	a Lens[T, A]
	b Lens[T, B]
}

func (lens shape2[T, A, B]) Put(s *T, a A, b B) *T {
	return lens.a.Put(lens.b.Put(s, b), a)
}

func (lens shape2[T, A, B]) Get(s *T) (A, B) {
	return lens.a.Get(s), lens.b.Get(s)
}

//
// Rank 3 Shape: Product Lens A × B × C
//

// Rank 3 Shape: Product Lens A × B × C
type Lens3[S, A, B, C any] interface {
	Get(*S) (A, B, C)
	Put(*S, A, B, C) *S
}

// Rank 3 Shape: Product Lens A × B × C
func ForShape3[T, A, B, C any](attr ...string) Lens3[T, A, B, C] {
	a, b, c := ForProduct3[T, A, B, C](attr...)
	return shape3[T, A, B, C]{a: a, b: b, c: c}
}

type shape3[T, A, B, C any] struct {
	a Lens[T, A]
	b Lens[T, B]
	c Lens[T, C]
}

func (lens shape3[T, A, B, C]) Put(s *T, a A, b B, c C) *T {
	return lens.a.Put(lens.b.Put(lens.c.Put(s, c), b), a)
}

func (lens shape3[T, A, B, C]) Get(s *T) (A, B, C) {
	return lens.a.Get(s), lens.b.Get(s), lens.c.Get(s)
}

//
// Rank 4 Shape: Product Lens A × B × C × D
//

// Rank 4 Shape: Product Lens A × B × C × D
type Lens4[S, A, B, C, D any] interface {
	Get(*S) (A, B, C, D)
	Put(*S, A, B, C, D) *S
}

// Rank 4 Shape: Product Lens A × B × C × D
func ForShape4[T, A, B, C, D any](attr ...string) Lens4[T, A, B, C, D] {
	a, b, c, d := ForProduct4[T, A, B, C, D](attr...)
	return shape4[T, A, B, C, D]{a: a, b: b, c: c, d: d}
}

type shape4[T, A, B, C, D any] struct {
	a Lens[T, A]
	b Lens[T, B]
	c Lens[T, C]
	d Lens[T, D]
}

func (lens shape4[T, A, B, C, D]) Put(s *T, a A, b B, c C, d D) *T {
	return lens.a.Put(lens.b.Put(lens.c.Put(lens.d.Put(s, d), c), b), a)
}

func (lens shape4[T, A, B, C, D]) Get(s *T) (A, B, C, D) {
	return lens.a.Get(s), lens.b.Get(s), lens.c.Get(s), lens.d.Get(s)
}

//
// Rank 5 Shape: Product Lens A × B × C × D × E
//

// Rank 5 Shape: Product Lens A × B × C × D × E
type Lens5[S, A, B, C, D, E any] interface {
	Get(*S) (A, B, C, D, E)
	Put(*S, A, B, C, D, E) *S
}

// Rank 5 Shape: Product Lens A × B × C × D × E
func ForShape5[T, A, B, C, D, E any](attr ...string) Lens5[T, A, B, C, D, E] {
	a, b, c, d, e := ForProduct5[T, A, B, C, D, E](attr...)
	return shape5[T, A, B, C, D, E]{a: a, b: b, c: c, d: d, e: e}
}

type shape5[T, A, B, C, D, E any] struct {
	a Lens[T, A]
	b Lens[T, B]
	c Lens[T, C]
	d Lens[T, D]
	e Lens[T, E]
}

func (lens shape5[T, A, B, C, D, E]) Put(s *T, a A, b B, c C, d D, e E) *T {
	return lens.a.Put(lens.b.Put(lens.c.Put(lens.d.Put(lens.e.Put(s, e), d), c), b), a)
}

func (lens shape5[T, A, B, C, D, E]) Get(s *T) (A, B, C, D, E) {
	return lens.a.Get(s), lens.b.Get(s), lens.c.Get(s), lens.d.Get(s), lens.e.Get(s)
}

//
// Rank 6 Shape: Product Lens A × B × C × D × E × F
//

// Rank 6 Shape: Product Lens A × B × C × D × E × F
type Lens6[S, A, B, C, D, E, F any] interface {
	Get(*S) (A, B, C, D, E, F)
	Put(*S, A, B, C, D, E, F) *S
}

// Rank 6 Shape: Product Lens A × B × C × D × E × F
func ForShape6[T, A, B, C, D, E, F any](attr ...string) Lens6[T, A, B, C, D, E, F] {
	a, b, c, d, e, f := ForProduct6[T, A, B, C, D, E, F](attr...)
	return shape6[T, A, B, C, D, E, F]{a: a, b: b, c: c, d: d, e: e, f: f}
}

type shape6[T, A, B, C, D, E, F any] struct {
	a Lens[T, A]
	b Lens[T, B]
	c Lens[T, C]
	d Lens[T, D]
	e Lens[T, E]
	f Lens[T, F]
}

func (lens shape6[T, A, B, C, D, E, F]) Put(s *T, a A, b B, c C, d D, e E, f F) *T {
	return lens.a.Put(lens.b.Put(lens.c.Put(lens.d.Put(lens.e.Put(lens.f.Put(s, f), e), d), c), b), a)
}

func (lens shape6[T, A, B, C, D, E, F]) Get(s *T) (A, B, C, D, E, F) {
	return lens.a.Get(s), lens.b.Get(s), lens.c.Get(s), lens.d.Get(s), lens.e.Get(s), lens.f.Get(s)
}

//
// Rank 7 Shape: Product Lens A × B × C × D × E × F × G
//

// Rank 7 Shape: Product Lens A × B × C × D × E × F × G
type Lens7[S, A, B, C, D, E, F, G any] interface {
	Get(*S) (A, B, C, D, E, F, G)
	Put(*S, A, B, C, D, E, F, G) *S
}

// Rank 7 Shape: Product Lens A × B × C × D × E × F × G
func ForShape7[T, A, B, C, D, E, F, G any](attr ...string) Lens7[T, A, B, C, D, E, F, G] {
	a, b, c, d, e, f, g := ForProduct7[T, A, B, C, D, E, F, G](attr...)
	return shape7[T, A, B, C, D, E, F, G]{a: a, b: b, c: c, d: d, e: e, f: f, g: g}
}

type shape7[T, A, B, C, D, E, F, G any] struct {
	a Lens[T, A]
	b Lens[T, B]
	c Lens[T, C]
	d Lens[T, D]
	e Lens[T, E]
	f Lens[T, F]
	g Lens[T, G]
}

func (lens shape7[T, A, B, C, D, E, F, G]) Put(s *T, a A, b B, c C, d D, e E, f F, g G) *T {
	return lens.a.Put(lens.b.Put(lens.c.Put(lens.d.Put(lens.e.Put(lens.f.Put(lens.g.Put(s, g), f), e), d), c), b), a)
}

func (lens shape7[T, A, B, C, D, E, F, G]) Get(s *T) (A, B, C, D, E, F, G) {
	return lens.a.Get(s), lens.b.Get(s), lens.c.Get(s), lens.d.Get(s), lens.e.Get(s), lens.f.Get(s), lens.g.Get(s)
}

//
// Rank 8 Shape: Product Lens A × B × C × D × E × F × G × H
//

// Rank 8 Shape: Product Lens A × B × C × D × E × F × G × H
type Lens8[S, A, B, C, D, E, F, G, H any] interface {
	Get(*S) (A, B, C, D, E, F, G, H)
	Put(*S, A, B, C, D, E, F, G, H) *S
}

// Rank 8 Shape: Product Lens A × B × C × D × E × F × G × H
func ForShape8[T, A, B, C, D, E, F, G, H any](attr ...string) Lens8[T, A, B, C, D, E, F, G, H] {
	a, b, c, d, e, f, g, h := ForProduct8[T, A, B, C, D, E, F, G, H](attr...)
	return shape8[T, A, B, C, D, E, F, G, H]{a: a, b: b, c: c, d: d, e: e, f: f, g: g, h: h}
}

type shape8[T, A, B, C, D, E, F, G, H any] struct {
	a Lens[T, A]
	b Lens[T, B]
	c Lens[T, C]
	d Lens[T, D]
	e Lens[T, E]
	f Lens[T, F]
	g Lens[T, G]
	h Lens[T, H]
}

func (lens shape8[T, A, B, C, D, E, F, G, H]) Put(s *T, a A, b B, c C, d D, e E, f F, g G, h H) *T {
	return lens.a.Put(lens.b.Put(lens.c.Put(lens.d.Put(lens.e.Put(lens.f.Put(lens.g.Put(lens.h.Put(s, h), g), f), e), d), c), b), a)
}

func (lens shape8[T, A, B, C, D, E, F, G, H]) Get(s *T) (A, B, C, D, E, F, G, H) {
	return lens.a.Get(s), lens.b.Get(s), lens.c.Get(s), lens.d.Get(s), lens.e.Get(s), lens.f.Get(s), lens.g.Get(s), lens.h.Get(s)
}

//
// Rank 9 Shape: Product Lens A × B × C × D × E × F × G × H × I
//

// Rank 9 Shape: Product Lens A × B × C × D × E × F × G × H × I
type Lens9[S, A, B, C, D, E, F, G, H, I any] interface {
	Get(*S) (A, B, C, D, E, F, G, H, I)
	Put(*S, A, B, C, D, E, F, G, H, I) *S
}

// Rank 9 Shape: Product Lens A × B × C × D × E × F × G × H × I
func ForShape9[T, A, B, C, D, E, F, G, H, I any](attr ...string) Lens9[T, A, B, C, D, E, F, G, H, I] {
	a, b, c, d, e, f, g, h, i := ForProduct9[T, A, B, C, D, E, F, G, H, I](attr...)
	return shape9[T, A, B, C, D, E, F, G, H, I]{a: a, b: b, c: c, d: d, e: e, f: f, g: g, h: h, i: i}
}

type shape9[T, A, B, C, D, E, F, G, H, I any] struct {
	a Lens[T, A]
	b Lens[T, B]
	c Lens[T, C]
	d Lens[T, D]
	e Lens[T, E]
	f Lens[T, F]
	g Lens[T, G]
	h Lens[T, H]
	i Lens[T, I]
}

func (lens shape9[T, A, B, C, D, E, F, G, H, I]) Put(s *T, a A, b B, c C, d D, e E, f F, g G, h H, i I) *T {
	return lens.a.Put(lens.b.Put(lens.c.Put(lens.d.Put(lens.e.Put(lens.f.Put(lens.g.Put(lens.h.Put(lens.i.Put(s, i), h), g), f), e), d), c), b), a)
}

func (lens shape9[T, A, B, C, D, E, F, G, H, I]) Get(s *T) (A, B, C, D, E, F, G, H, I) {
	return lens.a.Get(s), lens.b.Get(s), lens.c.Get(s), lens.d.Get(s), lens.e.Get(s), lens.f.Get(s), lens.g.Get(s), lens.h.Get(s), lens.i.Get(s)
}
