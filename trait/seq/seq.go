//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/skiplist
//

package seq

// Generic Sequence
//
//	seq := createSeq()
//	for has := seq != nil; has; has = seq.Next() {
//		seq.Value()
//	}
type Seq[T any] interface {
	Value() T
	Next() bool
}

// Lift instance of T into Seq[T] trait
func From[T any](xs T) Seq[T] { return element[T]{xs} }

type element[T any] struct{ v T }

func (v element[T]) Value() T   { return v.v }
func (v element[T]) Next() bool { return false }

// List instance of []T into Seq[T] trait
func FromSlice[T any](xs []T) Seq[T] {
	if len(xs) == 0 {
		return nil
	}

	return &seqOf[T]{xs}
}

type seqOf[T any] struct{ el []T }

func (s *seqOf[T]) Value() T { return s.el[0] }
func (s *seqOf[T]) Next() bool {
	if len(s.el) == 1 {
		return false
	}

	s.el = s.el[1:]
	return true
}

// Take values from iterator while predicate function true
func TakeWhile[T any](seq Seq[T], f func(T) bool) Seq[T] {
	if seq == nil || !f(seq.Value()) {
		return nil
	}

	return &takeWhile[T]{
		Seq: seq,
		f:   f,
	}
}

type takeWhile[T any] struct {
	Seq[T]
	f func(T) bool
}

func (seq *takeWhile[T]) Next() bool {
	if seq.f == nil || seq.Seq == nil {
		return false
	}

	if !seq.Seq.Next() {
		return false
	}

	if !seq.f(seq.Value()) {
		seq.f = nil
		return false
	}

	return true
}

// Drop values from iterator while predicate function true
func DropWhile[T any](seq Seq[T], f func(T) bool) Seq[T] {
	if seq == nil {
		return nil
	}

	for {
		if !f(seq.Value()) {
			return seq
		}

		if !seq.Next() {
			return nil
		}
	}
}

// Filter values from iterator
func Filter[T any](seq Seq[T], f func(T) bool) Seq[T] {
	if seq == nil {
		return nil
	}

	for {
		if f(seq.Value()) {
			return filter[T]{
				Seq: seq,
				f:   f,
			}
		}

		if !seq.Next() {
			return nil
		}
	}
}

type filter[T any] struct {
	Seq[T]
	f func(T) bool
}

func (seq filter[T]) Next() bool {
	if seq.f == nil || seq.Seq == nil {
		return false
	}

	for {
		if !seq.Seq.Next() {
			return false
		}

		if seq.f(seq.Value()) {
			return true
		}
	}
}

// ForEach applies clojure on iterator
func ForEach[T any](seq Seq[T], f func(T) error) error {
	for has := seq != nil; has; has = seq.Next() {
		if err := f(seq.Value()); err != nil {
			return err
		}
	}

	return nil
}

// Map transform iterator type
func Map[A, B any](seq Seq[A], f func(A) B) Seq[B] {
	if seq == nil {
		return nil
	}

	return fmap[A, B]{Seq: seq, f: f}
}

type fmap[A, B any] struct {
	Seq[A]
	f func(A) B
}

func (seq fmap[A, B]) Value() B {
	return seq.f(seq.Seq.Value())
}

// Plus operation for iterators add one after another
func Plus[T any](lhs, rhs Seq[T]) Seq[T] {
	if lhs == nil {
		return rhs
	}

	if rhs == nil {
		return lhs
	}

	return &plus[T]{Seq: lhs, rhs: rhs}
}

type plus[T any] struct {
	Seq[T]
	rhs Seq[T]
}

func (plus *plus[T]) Next() bool {
	hasNext := plus.Seq.Next()

	if !hasNext && plus.rhs != nil {
		plus.Seq, plus.rhs = plus.rhs, nil
		return true
	}

	if !hasNext && plus.rhs == nil {
		return false
	}

	return true
}

// Left join
func Join[A, B any](lhs Seq[A], rhs func(A) Seq[B]) Seq[B] {
	if lhs == nil {
		return nil
	}

	join := &join[A, B]{lhs: lhs, rhs: rhs}
	for {
		join.Seq = join.rhs(join.lhs.Value())
		if join.Seq != nil {
			return join
		}

		if !join.lhs.Next() {
			return nil
		}
	}
}

type join[A, B any] struct {
	Seq[B]
	lhs Seq[A]
	rhs func(A) Seq[B]
}

func (join *join[A, B]) Next() bool {
	if !join.Seq.Next() {
		for {
			if !join.lhs.Next() {
				return false
			}

			join.Seq = join.rhs(join.lhs.Value())
			if join.Seq != nil {
				return true
			}
		}
	}

	return true
}
