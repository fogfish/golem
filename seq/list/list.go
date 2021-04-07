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

Iterator of major list properties
*/
type Iterator interface {
	// First (head) element of sequence
	Head() golem.Ord
	// All elements except the first (tail) of sequence
	Tail() *Type
	// The number of elements in the sequence
	Length() int
	// Is this the empty sequence
	IsEmpty() bool
}

//
type Predicate func(golem.Ord) bool

// TODO: rename to "bucket-(ing)" function
type GroupBy func(golem.Ord) bool

type Ord func(golem.Ord, golem.Ord) bool

type Zip func(golem.Ord, golem.Ord) golem.Ord
type UnZip func(golem.Ord) (golem.Ord, golem.Ord)

/*

Seacher ...
*/
type Searcher interface {
	// TakeWhile accumulates elements from sequence head while predicate returns true and returns this prefix.
	TakeWhile(Predicate) *Type
	// DropWhile removes elements from sequence head while predicate returns true and returns remaining sequence suffix.
	DropWhile(Predicate) *Type
	// Span splits the sequences into prefix/suffix pair. The prefix contains elements that
	// satisfy predicates. The suffix is the remainder of the sequence. It is equivalent of consequent calls to TakeWhile/DropWhile
	SpanWith(Predicate) (*Type, *Type)
	// Partition split sequence into two sequence according to predicate. The first sequence contains
	// It is equivalent of consequent calls to Filter/FilterNot
	Partition(Predicate) (*Type, *Type)
	GroupBy(GroupBy) map[int]*Type
	Filter(Predicate) *Type
	FilterNot(Predicate) *Type
}

/*

Sorter ...
*/
type Sorter interface {
	// Sorts the specified Seq by the natural ordering of its elements.
	Sort() *Type
	// Sorts the specified Seq according to the specified comparator.
	SortBy(Ord) *Type
	// Distinct builds a new sequence without any duplicate elements
	Distinct() *Type
}

type X interface {
	// Finds the first element that satisfy predicate
	Find(Predicate) golem.Ord
	Take(int) *Type
	Drop(int) *Type
	Span(int) (*Type, *Type)
}

type Zipper interface {
	Zip(*Type) *Type
	ZipWith(Zip, *Type)
	UnZip(*Type) (*Type, *Type)
	UnZipWith(UnZip) (*Type, *Type)
}

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

/*

FMap applies closure to each element of the list
*/
func (list *Type) FMap(f func(golem.Ord)) {
	for l := list; l != nil; l = l.Tail() {
		f(l.head)
	}
}
