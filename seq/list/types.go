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

Builder declares list type constructors
*/
type Builder interface {
	// Creates a new list by adding a new head element
	Cons(golem.Ord) *Type
	// Create a new list by adding a another sequence to the head
	Join(*Type) *Type
	// Builds a sequence of given length with a function representing sequence
	// element.
	From(int, func(int) golem.Ord) *Type
	// Unfold is generic sequence constructor. It builds a sequence from a seed
	// value and generator function. It continues to construct the sequence until
	// generator return nil value
	Unfold(golem.Data, func(golem.Data) golem.Ord) *Type
}
