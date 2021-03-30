//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package seq

import "github.com/fogfish/golem"

func NewList() *List {
	return nil
}

/*

IsEmpty return true if list does not contain any elements
*/
func (list *List) IsEmpty() bool {
	return list == nil
}

/*

Cons creates a new list
*/
func (list *List) Cons(head golem.Ord) *List {
	return &List{head: head, tail: list}
}

/*

Head of the list
*/
func (seq *List) Head() golem.Ord {
	return seq.head
}

/*

Tail of the list
*/
func (seq *List) Tail() *List {
	return seq.tail
}
