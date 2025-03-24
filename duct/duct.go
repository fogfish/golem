//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package duct

import (
	"reflect"
)

// Unary type A
type T[A any] struct{ v any }

// Lifts a value into a unary type A
func L1[A any](f any) T[A] { return T[A]{v: f} }

// Binary type A, B
type F[A, B any] struct{ f any }

// Lifts a value into a binary type A, B
func L2[A, B any](f any) F[A, B] { return F[A, B]{f: f} }

// Morphism 𝑚: A ⟼ B is an abstract transformer of category `A` to `B`.
type Morphism[A, B any] struct {
	code *AstSeq
}

// Applies visitor over computation elements within morphism 𝑚: A ⟼ B.
func (seq Morphism[A, B]) Apply(v Visitor) error {
	return seq.code.Apply(0, v)
}

// From create new morphism 𝑓: ø ⟼ A, binding it with source of category `A`.
func From[A any](source T[A]) Morphism[A, A] {
	code := &AstSeq{
		Root:     true,
		Deferred: true,
		Seq:      make([]Ast, 0),
	}

	in := &AstFrom{
		Type:   TypeOf[A](),
		Source: source.v,
	}
	code.append(in)

	return Morphism[A, A]{code: code}
}

// Compose transformer 𝑓: B ⟼ C with morphism 𝑚: A ⟼ B producing a new morphism 𝑚: A ⟼ C.
func Join[A, B, C any](f F[B, C], m Morphism[A, B]) Morphism[A, C] {
	code := m.code

	join := &AstMap{
		TypeA: TypeOf[B](),
		TypeB: TypeOf[C](),
		F:     f.f,
	}
	code.append(join)

	return Morphism[A, C]{code: code}
}

// Compose the transformer 𝑓: B ⟼ C with the morphism 𝑚: A ⟼ []B to produce
// a free monad 𝑚⁺: A ⟼ []C that enables transformation within a functorial
// context while preserving the computational structure without immediate
// collapsing. Unlike traditional flatMap, LiftF preserves the transformer 𝑓's
// context, enabling further composition inside the transformation. Specifically,
// Join(g, LiftF(f)) is equivalent to flatMap(f ∘ g), where the function g is
// applied after f, and the results are flattened into a single structure.
//
// It is a responsibility of creator of such a free monadic value to do something
// with those nested contexts either yielding individual elements or uniting into
// the monad (e.g. use Unit(Join(g, LiftF(f))) to leave nested context into
// the morphism 𝑚: A ⟼ []C).
func LiftF[A, B, C any](f F[B, C], m Morphism[A, []B]) Morphism[A, C] {
	inner := &AstSeq{
		Root:     false,
		Deferred: true,
		Seq:      make([]Ast, 0),
	}

	join := &AstMap{
		TypeA: TypeOf[B](),
		TypeB: TypeOf[C](),
		F:     f.f,
	}
	inner.append(join)

	code := m.code
	code.append(inner)

	return Morphism[A, C]{code: code}
}

// WrapF is equivalent to LiftF but operates directly on the inner structure of
// the morphism 𝑚: A ⟼ []B, extracting individual elements of B while
// preserving the transformation context, enabling further composition.
//
// Usable to Yield elements of []B without transformation
func WrapF[A, B any](m Morphism[A, []B]) Morphism[A, B] {
	inner := &AstSeq{
		Root:     false,
		Deferred: true,
		Seq:      make([]Ast, 0),
	}

	code := m.code
	code.append(inner)

	return Morphism[A, B]{code: code}
}

// Unit finalizes a transformation context by collapsing the free monadic value
// of the morphism 𝑚: A ⟼ []B. It acts as the terminal operation, ensuring
// that all staged compositions, such as those built with LiftF and WrapF,
// are fully resolved into a single, consumable form.
func Unit[A, B any](m Morphism[A, B]) Morphism[A, []B] {
	code := m.code
	code.unit()

	return Morphism[A, []B]{code: code}
}

// Yield results of 𝑚: A ⟼ B binding it with target of category `B`.
func Yield[A, B any](target T[B], m Morphism[A, B]) Morphism[A, Void] {
	code := m.code

	eg := &AstYield{
		Type:   TypeOf[B](),
		Target: target.v,
	}
	code.append(eg)

	return Morphism[A, Void]{code: code}
}

// TypeOf returns normalized name of the type T.
func TypeOf[T any]() string {
	return typeName(reflect.TypeOf(new(T)).Elem())
}

func typeName(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr:
		return "*" + typeName(t.Elem())
	case reflect.Slice:
		return "[]" + typeName(t.Elem())
	default:
		return t.Name()
	}
}
