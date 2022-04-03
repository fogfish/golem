//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package pure

/*

Type `*` is seen as nullary type constructors, and also
called proper types in this context.
*/
type Type[F any] interface {
	HKT1(F)
}

/*

HKT[F, A] ∼ F[A]

Transforms a computation with higher-kinded type expressions into a
computation where all type expressions are of kind `*`. The abstract type
constructor `HKT` represent an idea of parametrized container type `F[A]`

HKT `* ⟼ *` is unary type constructor

See doc/typeclass.md for deatails about the usage of HKT.
*/
type HKT[F, A any] interface {
	Type[F]
	HKT2(A)
}
