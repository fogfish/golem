//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package seq

import (
	"github.com/fogfish/golem/generic"
)

// AnyT is Seq data type build of generic.T elements
type AnyT []generic.T

// Contain tests if sequence contains an element
func (seq AnyT) Contain(e generic.T) bool {
	return seq.Exists(func(x generic.T) bool { return x == e })
}

// ContainSlice tests if sequence contains a sub-sequence
// func (seq AnyT) ContainSlice(subseq AnyT) bool

// Count number of elements that satisfy the predicate
func (seq AnyT) Count(p func(generic.T) bool) int {
	c := 0
	for _, x := range seq {
		if p(x) {
			c = c + 1
		}
	}
	return c
}

// Diff computes the difference between sequences: seq -- that
// func (seq AnyT) Diff(that AnyT) AnyT

// Distinct builds a new sequence without any duplicate elements
func (seq AnyT) Distinct() AnyT {
	s := AnyT{}
	for _, x := range seq {
		if !s.Contain(x) {
			s = append(s, x)
		}
	}
	return s
}

// Drop removes n elements from head of sequence
func (seq AnyT) Drop(n int) AnyT {
	return append(seq[:0:0], seq[n:]...)
}

// DropWhile removes elements from sequence head while predicate returns true
// and returns remaining sequence suffix.
func (seq AnyT) DropWhile(p func(generic.T) bool) AnyT {
	for i, x := range seq {
		if !p(x) {
			return append(seq[:0:0], seq[i:]...)
		}
	}
	return AnyT{}
}

// Exists tests if a predicate holds for at least one element
func (seq AnyT) Exists(p func(generic.T) bool) bool {
	for _, x := range seq {
		if p(x) {
			return true
		}
	}
	return false
}

// Filter selects all elements which satisfy predicate
func (seq AnyT) Filter(p func(generic.T) bool) AnyT {
	s := AnyT{}
	for _, x := range seq {
		if p(x) {
			s = append(s, x)
		}
	}
	return s
}

// Find returns the first element that satisfy predicate
func (seq AnyT) Find(p func(generic.T) bool) (e generic.T) {
	for _, x := range seq {
		if p(x) {
			return x
		}
	}
	return
}

// ForAll tests where a predicate holds for elements of sequence
func (seq AnyT) ForAll(p func(generic.T) bool) bool {
	if len(seq) == 0 {
		return false
	}
	for _, x := range seq {
		if !p(x) {
			return false
		}
	}
	return true
}

// FMap applies high-order function (clojure) to all elements of sequence
func (seq AnyT) FMap(f func(generic.T)) {
	for _, x := range seq {
		f(x)
	}
}

// Fold applies associative binary operator to sequence
func (seq AnyT) Fold(f func(generic.T, generic.T) generic.T, empty generic.T) generic.T {
	acc := empty
	for _, x := range seq {
		acc = f(x, acc)
	}
	return acc
}

// GroupBy shards sequence into map of sequences with descriminator function
func (seq AnyT) GroupBy(f func(generic.T) int) map[int]AnyT {
	s := make(map[int]AnyT)
	for _, x := range seq {
		key := f(x)
		shard, exists := s[key]
		if exists {
			s[key] = append(shard, x)
		} else {
			s[key] = AnyT{x}
		}
	}
	return s
}

// Join takes sequence of sequence, flattens and append it
func (seq AnyT) Join(subseq AnyT) AnyT {
	seq = append(seq, subseq...)
	return seq
}

// Intersect computes the intersection of sequences: seq ^ that
// func (seq AnyT) Intersect(that AnyT) AnyT

// Map applies high-order function to all element of sequence
func (seq AnyT) Map(f func(generic.T) generic.T) AnyT {
	s := AnyT{}
	for _, x := range seq {
		s = append(s, f(x))
	}
	return s
}

// Partition split sequence into two sequence accroding to predicate
// It is equivalent of consequent calls to Filter/FilterNot
func (seq AnyT) Partition(p func(generic.T) bool) (AnyT, AnyT) {
	a := AnyT{}
	b := AnyT{}
	for _, x := range seq {
		if p(x) {
			a = append(a, x)
		} else {
			b = append(b, x)
		}
	}
	return a, b
}

// Reverse returns a new sequence with elements in reserve order
// func (seq AnyT) Reverse() AnyT

// Span splits the sequences into prefix/suffix pair according to predicate
// It is equivalent of consequent calls to TakeWhile/DropWhile
// func (seq AnyT) Span(p Predicate) AnyT, AnyT

// Take accepts n elements from head of sequence
func (seq AnyT) Take(n int) AnyT {
	return append(seq[:0:0], seq[:n]...)
}

// TakeWhile accumulates elements from sequence head while predicate returns true
// and returns this prefix.
func (seq AnyT) TakeWhile(p func(generic.T) bool) AnyT {
	for i, x := range seq {
		if !p(x) {
			return append(seq[:0:0], seq[:i]...)
		}
	}
	return append(AnyT{}, seq...)
}
