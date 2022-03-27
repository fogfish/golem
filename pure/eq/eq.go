//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package eq

import (
	"github.com/fogfish/golem/pure"
)

/*

Eq : T ⟼ T ⟼ bool
Each trait implements mapping pair of values to bool category using own
equality rules
*/
type Eq[T any] interface {
	Equal(T, T) bool
}

/*

eq generic implementation for built-in types
*/
type eq[T comparable] string

func (eq[T]) Equal(a, b T) bool { return a == b }

//
const (
	Int    = eq[int]("eq.int")
	String = eq[string]("eq.string")
)

/*

From is a combinator that lifts T ⟼ T ⟼ bool function to
an instance of Eq type trait
*/
type From[T any] func(T, T) bool

func (f From[T]) Equal(a, b T) bool { return f(a, b) }

/*

ContraMap is a combinator that build a new instance of type trait Eq[B] using
existing instance of Eq[A] and f: b ⟼ a
*/
type ContraMap[A, B any] struct {
	Eq[A]
	pure.ContraMap[A, B]
}

// Equal implementation of contra variant functor
func (f ContraMap[A, B]) Equal(a, b B) bool {
	return f.Eq.Equal(
		f.ContraMap(a),
		f.ContraMap(b),
	)
}
