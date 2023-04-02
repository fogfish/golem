//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package ord

import "github.com/fogfish/golem/pure"

// Ordering type defines enum{ LT, EQ, GT }
type Ordering int

// enum{ LT, EQ, GT }
const (
	LT Ordering = -1
	EQ Ordering = 0
	GT Ordering = 1
)

// Ord : T ⟼ T ⟼ Ordering
// Each type implements compare rules, mapping pair of value to enum{ LT, EQ, GT }
type Ord[T any] interface {
	Compare(T, T) Ordering
}

// ord generic implementation for built-in types
type ord[T pure.AnyOrderable] string

func (ord[T]) Compare(a, b T) Ordering {
	switch {
	case a < b:
		return LT
	case a > b:
		return GT
	default:
		return EQ
	}
}

const (
	Int    = ord[int]("eq.int")
	String = ord[string]("eq.string")
)

// From is a combinator that lifts T ⟼ T ⟼ Ordering function to
// an instance of Ord type trait
type From[T any] func(T, T) Ordering

func (f From[T]) Compare(a, b T) Ordering { return f(a, b) }

// ContraMap is a combinator that build a new instance of type trait Ord[B] using
// existing instance of Ord[A] and f: b ⟼ a
type ContraMap[A, B any] struct {
	Ord[A]
	pure.ContraMap[A, B]
}

// Equal implementation of contra variant functor
func (f ContraMap[A, B]) Compare(a, b B) Ordering {
	return f.Ord.Compare(
		f.ContraMap(a),
		f.ContraMap(b),
	)
}
