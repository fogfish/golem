//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package slice

import "github.com/fogfish/golem/seq"

/*

Seq type build over slice data structure
*/
type Seq[A any] []A

func (Seq[A]) HKT1(seq.Type) {}
func (Seq[A]) HKT2(A)        {}

/*

Trait implements seq.Seq type law for Seq type
*/
type Trait[A any] string

var _ seq.Seq[Seq[any], any] = Trait[any]("seq.any")

func (Trait[A]) New(seq ...A) Seq[A] { return seq }

func (Trait[A]) Cons(x A, seq Seq[A]) Seq[A] {
	return append([]A{x}, seq...)
}

func (Trait[A]) Head(seq Seq[A]) A      { return seq[0] }
func (Trait[A]) Tail(seq Seq[A]) Seq[A] { return seq[1:] }

func (Trait[A]) Length(seq Seq[A]) int { return len(seq) }

func (Trait[A]) IsEmpty(seq Seq[A]) bool { return len(seq) == 0 }
