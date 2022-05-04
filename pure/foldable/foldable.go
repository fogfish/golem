//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package foldable

import "github.com/fogfish/golem/pure/monoid"

/*

Foldable type trait define rules of folding data structures to a summary value.
*/
type Foldable[F_, A any] interface {
	Fold(monoid.Monoid[A], F_) A
}
