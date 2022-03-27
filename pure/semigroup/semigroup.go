//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package semigroup

/*

Semigroup is an algebraic structure consisting of a set together
with an associative binary operation.
*/
type Semigroup[T any] interface {
	Combine(T, T) T
}

/*

From is a combinator that lifts T ⟼ T ⟼ T function to
an instance of Semigroup type trait
*/
type From[T any] func(T, T) T

func (f From[T]) Combine(a, b T) T { return f(a, b) }
