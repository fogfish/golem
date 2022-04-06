//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package monoid

import "github.com/fogfish/golem/pure/semigroup"

/*

Monoid is an algebraic structure consisting of Semigroup and Empty element
*/
type Monoid[T any] interface {
	semigroup.Semigroup[T]
	Empty() T
}

/*

From is a combinator that lifts T ⟼ T ⟼ T function to
an instance of Monoid type trait
*/
func From[T any](empty T, combine func(T, T) T) Monoid[T] {
	return monoid[T]{
		empty:   empty,
		combine: combine,
	}
}

/*

Internal implementation of Monoid interface
*/
type monoid[T any] struct {
	empty   T
	combine func(T, T) T
}

func (m monoid[T]) Empty() T         { return m.empty }
func (m monoid[T]) Combine(a, b T) T { return m.combine(a, b) }
