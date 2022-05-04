//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package list

import "github.com/fogfish/golem/seq"

/*

Seq a type represents a finite sequence of values of type a, which is built
using linked-list
*/
type Seq[A any] struct {
	len  int
	list *list[A]
}

func (Seq[A]) HKT1(seq.Type) {}
func (Seq[A]) HKT2(A)        {}

type list[A any] struct {
	head A
	tail *list[A]
}

/*

Trait implements seq.Seq type law for Seq type
*/
type Trait[A any] string

var _ seq.Seq[Seq[any], any] = Trait[any]("seq.any")

func (Trait[A]) New(seq ...A) Seq[A] {
	var tail *list[A]

	for i := len(seq) - 1; i >= 0; i-- {
		tail = &list[A]{head: seq[i], tail: tail}
	}

	return Seq[A]{len: len(seq), list: tail}
}

func (Trait[A]) Cons(x A, seq Seq[A]) Seq[A] {
	return Seq[A]{
		len:  seq.len + 1,
		list: &list[A]{head: x, tail: seq.list},
	}
}

func (Trait[A]) Head(seq Seq[A]) A      { return seq.list.head }
func (Trait[A]) Tail(seq Seq[A]) Seq[A] { return Seq[A]{len: seq.len - 1, list: seq.list.tail} }

func (Trait[A]) Length(seq Seq[A]) int { return seq.len }

func (Trait[A]) IsEmpty(seq Seq[A]) bool { return seq.len == 0 }
