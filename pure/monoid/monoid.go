//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package monoid

import "github.com/fogfish/golem/pure/semigroup"

// Monoid is an algebraic structure consisting of Semigroup and Empty element
type Monoid[T any] interface {
	semigroup.Semigroup[T]
	Empty() T
}

// From is a combinator that lifts Semigroup to an instance of Monoid type trait
func From[T any](empty T, combine semigroup.Semigroup[T]) Monoid[T] {
	return monoid[T]{
		Semigroup: combine,
		empty:     empty,
	}
}

// FromOp is a combinator that lifts T ⟼ T ⟼ T function (binary operator) to
// an instance of Monoid type trait
func FromOp[T any](empty T, combine func(T, T) T) Monoid[T] {
	return monoid[T]{
		Semigroup: semigroup.From[T](combine),
		empty:     empty,
	}
}

// Internal implementation of Monoid interface
type monoid[T any] struct {
	semigroup.Semigroup[T]
	empty T
}

func (m monoid[T]) Empty() T { return m.empty }
