//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package list

import "github.com/fogfish/golem"

/*

List is a linked list
*/
type Type struct {
	head golem.Ord
	tail *Type
}

/*

New empty list
*/
func New() *Type {
	return nil
}

/*

IsEmpty return true if list does not contain any elements
*/
func (list *Type) IsEmpty() bool {
	return list == nil
}

/*

Cons creates a new list
*/
func (list *Type) Cons(head golem.Ord) *Type {
	return &Type{head: head, tail: list}
}

/*

Head of the list
*/
func (list *Type) Head() golem.Ord {
	return list.head
}

/*

Tail of the list
*/
func (list *Type) Tail() *Type {
	return list.tail
}
