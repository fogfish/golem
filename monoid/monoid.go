//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package monoid

// Monoid supports the generic implementation of strcutural transformations.
// It is an algebraic structure https://en.wikipedia.org/wiki/Monoid
// with identity element and associative binary function.
//
// Monoid facilitates the implementation of following transformations in
// the absence of generics and type covariance.
//
//    def map[B](f: (A) => B): Seq[B]
//
// The library uses Monoid interface to implement map, fold, filter,
// comprehension and others for complex data structures. Here is an
// example
//
//
//   type MSeq struct { value []int }
//
//   func (seq *MSeq) Mempty() monoid.AnyT {
//     return &MSeq{}
//   }
//
//   func (seq *MSeq) Mappend(x monoid.AnyT) monoid.AnyT {
//     seq.value = append(seq.value, x.(*MSeq).value...)
//     return seq
//   }
//
type AnyT interface {
	// Mempty returns a type value that hold the identity property for
	// combine operation, means the following equalities hold for any choice of x.
	//   t.Combine(t.Empty()) == t.Empty().Combine(t) == t
	Mempty() AnyT

	// Mappend applies a side-effect to the structure by appending a given value.
	// Combine must hold associative property
	//   a.Combine(b).Combine(c) == a.Combine(b.Combine(c))
	Mappend(x AnyT) AnyT
}