//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package seq

import "github.com/fogfish/golem/pure"

/*

Type is opaque type to define polymorphic context of Seq.
It makes HKT typesafe in the context of sequence trait
*/
type Type any

/*

Higher-Kinded Sequence type
*/
type Kind[A any] pure.HKT[Type, A]

type Construction[F_, A any] interface {
	New(...A) F_
	Cons(A, F_) F_
}

type View[F_, A any] interface {
	Head(F_) A
	Tail(F_) F_
}

type Query[F_, A any] interface {
	Length(F_) int
	IsEmpty(F_) bool
}

type Seq[F_, A any] interface {
	Construction[F_, A]
	View[F_, A]
	Query[F_, A]
}
