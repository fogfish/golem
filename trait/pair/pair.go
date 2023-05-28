//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package pair

import "github.com/fogfish/golem/trait/seq"

// Generic Sequence of Pairs
//
//	seq := createSeq()
//	for has := seq != nil; has; has = seq.Next() {
//		seq.Value()
//	}
type Seq[K, V any] interface {
	seq.Seq[V]
	Key() K
}

// Lift instance of T into Seq[T] trait
func From[K, V any](key K, val V) Seq[K, V] { return pair[K, V]{key, val} }

type pair[K, V any] struct {
	key K
	val V
}

func (v pair[K, V]) Key() K     { return v.key }
func (v pair[K, V]) Value() V   { return v.val }
func (v pair[K, V]) Next() bool { return false }

// Take values from iterator while predicate function true
func TakeWhile[K, V any](seq Seq[K, V], f func(K, V) bool) Seq[K, V] {
	if seq == nil || !f(seq.Key(), seq.Value()) {
		return nil
	}

	return &takeWhile[K, V]{
		Seq: seq,
		f:   f,
	}
}

type takeWhile[K, V any] struct {
	Seq[K, V]
	f func(K, V) bool
}

func (seq *takeWhile[K, V]) Next() bool {
	if seq.f == nil || seq.Seq == nil {
		return false
	}

	if !seq.Seq.Next() {
		return false
	}

	if !seq.f(seq.Key(), seq.Value()) {
		seq.f = nil
		return false
	}

	return true
}

// Drop values from iterator while predicate function true
func DropWhile[K, V any](seq Seq[K, V], f func(K, V) bool) Seq[K, V] {
	if seq == nil {
		return nil
	}

	for {
		if !f(seq.Key(), seq.Value()) {
			return seq
		}

		if !seq.Next() {
			return nil
		}
	}
}

// Filter values from iterator
func Filter[K, V any](seq Seq[K, V], f func(K, V) bool) Seq[K, V] {
	if seq == nil {
		return nil
	}

	for {
		if f(seq.Key(), seq.Value()) {
			return filter[K, V]{
				Seq: seq,
				f:   f,
			}
		}

		if !seq.Next() {
			return nil
		}
	}
}

type filter[K, V any] struct {
	Seq[K, V]
	f func(K, V) bool
}

func (seq filter[K, V]) Next() bool {
	if seq.f == nil || seq.Seq == nil {
		return false
	}

	for {
		if !seq.Seq.Next() {
			return false
		}

		if seq.f(seq.Key(), seq.Value()) {
			return true
		}
	}
}

// ForEach applies clojure on iterator
func ForEach[K, V any](seq Seq[K, V], f func(K, V) error) error {
	for has := seq != nil; has; has = seq.Next() {
		if err := f(seq.Key(), seq.Value()); err != nil {
			return err
		}
	}

	return nil
}

// FMap transform iterator type
func Map[K, A, B any](seq Seq[K, A], f func(K, A) B) Seq[K, B] {
	if seq == nil {
		return nil
	}

	return fmap[K, A, B]{Seq: seq, f: f}
}

type fmap[K, A, B any] struct {
	Seq[K, A]
	f func(K, A) B
}

func (seq fmap[K, A, B]) Value() B {
	return seq.f(seq.Seq.Key(), seq.Seq.Value())
}

// Plus operation for iterators add one after another
func Plus[K, V any](lhs, rhs Seq[K, V]) Seq[K, V] {
	if lhs == nil {
		return rhs
	}

	if rhs == nil {
		return lhs
	}

	return &plus[K, V]{Seq: lhs, rhs: rhs}
}

type plus[K, V any] struct {
	Seq[K, V]
	rhs Seq[K, V]
}

func (plus *plus[K, V]) Next() bool {
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
func Join[K1, K2, V1, V2 any](lhs Seq[K1, V1], rhs func(K1, V1) Seq[K2, V2]) Seq[K2, V2] {
	if lhs == nil {
		return nil
	}

	join := &join[K1, K2, V1, V2]{lhs: lhs, rhs: rhs}
	for {
		join.Seq = join.rhs(join.lhs.Key(), join.lhs.Value())
		if join.Seq != nil {
			return join
		}

		if !join.lhs.Next() {
			return nil
		}
	}
}

type join[K1, K2, V1, V2 any] struct {
	Seq[K2, V2]
	lhs Seq[K1, V1]
	rhs func(K1, V1) Seq[K2, V2]
}

func (join *join[K1, K2, V1, V2]) Next() bool {
	if !join.Seq.Next() {
		for {
			if !join.lhs.Next() {
				return false
			}

			join.Seq = join.rhs(join.lhs.Key(), join.lhs.Value())
			if join.Seq != nil {
				return true
			}
		}
	}

	return true
}

// Left join sequence of pairs to sequence of values
func ToSeq[K1, V1, V2 any](lhs Seq[K1, V1], rhs func(K1, V1) seq.Seq[V2]) seq.Seq[V2] {
	if lhs == nil {
		return nil
	}

	join := &toSeq[K1, V1, V2]{lhs: lhs, rhs: rhs}
	for {
		join.Seq = join.rhs(join.lhs.Key(), join.lhs.Value())
		if join.Seq != nil {
			return join
		}

		if !join.lhs.Next() {
			return nil
		}
	}
}

type toSeq[K1, V1, V2 any] struct {
	seq.Seq[V2]
	lhs Seq[K1, V1]
	rhs func(K1, V1) seq.Seq[V2]
}

func (join *toSeq[K1, V1, V2]) Next() bool {
	if !join.Seq.Next() {
		for {
			if !join.lhs.Next() {
				return false
			}

			join.Seq = join.rhs(join.lhs.Key(), join.lhs.Value())
			if join.Seq != nil {
				return true
			}
		}
	}

	return true
}

// Left join sequence of elements  to sequence of pairs
func FromSeq[K1, K2, V2 any](lhs seq.Seq[K1], rhs func(K1) Seq[K2, V2]) Seq[K2, V2] {
	if lhs == nil {
		return nil
	}

	join := &fromSeq[K1, K2, V2]{lhs: lhs, rhs: rhs}
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

type fromSeq[K1, K2, V2 any] struct {
	Seq[K2, V2]
	lhs seq.Seq[K1]
	rhs func(K1) Seq[K2, V2]
}

func (join *fromSeq[K1, K2, V2]) Next() bool {
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
