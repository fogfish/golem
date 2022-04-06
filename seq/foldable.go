//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package seq

import (
	"github.com/fogfish/golem/pure/monoid"
)

/*

Foldable Sequence
*/
type Foldable[F_, A any] struct{ Seq[F_, A] }

/*

Fold sequence with Monoid
*/
func (f Foldable[F_, A]) Fold(m monoid.Monoid[A], seq F_) A {
	x := m.Empty()
	s := seq

	for !f.Seq.IsEmpty(s) {
		x = m.Combine(x, f.Seq.Head(s))
		s = f.Seq.Tail(s)
	}

	return x
}
