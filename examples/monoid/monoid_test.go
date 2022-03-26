//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package monoid

import (
	"testing"
)

//
// container type
type SeqA []int

type FMap func(int) Monoid
type Mapper func(FMap) Monoid

// Map with Monoid, instantiate mapper
func (seq SeqA) Map(m Monoid) Mapper {
	return func(f FMap) Monoid {
		y := m.Mempty()
		for _, x := range seq {
			y = y.Mappend(f(x))
		}
		return y
	}
}

// Map with Monoid
func (seq SeqA) MMap(m Monoid, f FMap) Monoid {
	y := m.Mempty()
	for _, x := range seq {
		y = y.Mappend(f(x))
	}
	return y
}

func (seq SeqA) MM(m M, f func(int) interface{}) interface{} {
	y := m.Empty()
	for _, x := range seq {
		y = m.Append(y, f(x))
	}
	return y
}

// Map with Closure
func (seq SeqA) FMap(f func(int)) {
	for _, x := range seq {
		f(x)
	}
}

//
// accumulator type implements monoid abstraction
type SeqB int

func (seq SeqB) Mempty() Monoid {
	return SeqB(0)
}

func (seq SeqB) Mappend(x Monoid) Monoid {
	return seq + x.(SeqB)
}

// usage of interface instead of pure type causes per penalties

type M interface {
	Empty() interface{}
	Append(a, b interface{}) interface{}
}

type SeqC int

func (SeqC) Empty() interface{} {
	return 0 //SeqC(0)
}

func (SeqC) Append(a, b interface{}) interface{} {
	return a.(int) + b.(int)
}

func Mc() M {
	return SeqC(0)
}

//
// global constants
var (
	sequence SeqA = SeqA{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	result   Monoid
	xxx      M
)

// convAtoB is mapper function, converts type A to B
func convAtoB(s int) Monoid {
	return SeqB(s)
}

//
// forEach make benchmark of forLoop comparable to FMap
func forEach() SeqB {
	seq := SeqB(0)
	for _, x := range sequence {
		seq = joinAtoB(seq, x)
	}
	return seq
}

func joinAtoB(b SeqB, a int) SeqB {
	return b + SeqB(a)
}

//
// benchmarks
func BenchmarkMonoid(b *testing.B) {
	b.ReportAllocs()

	var monoid Monoid = SeqB(0)

	for n := 0; n < b.N; n++ {
		monoid = sequence.MMap(monoid, convAtoB)
	}
	result = monoid
}

func BenchmarkMonoidC(b *testing.B) {
	b.ReportAllocs()

	// var monoid SeqC = SeqC(0)
	mC := Mc()
	for n := 0; n < b.N; n++ {
		sequence.MM(mC, func(x int) interface{} { return x })
	}
	// xxx = monoid
}

func BenchmarkMonoidHoF(b *testing.B) {
	b.ReportAllocs()

	var monoid Monoid = SeqB(0)
	mapper := sequence.Map(monoid)

	for n := 0; n < b.N; n++ {
		monoid = mapper(convAtoB)
	}
	result = monoid
}

func BenchmarkForLoop(b *testing.B) {
	b.ReportAllocs()

	var seq SeqB

	for n := 0; n < b.N; n++ {
		forEach()
	}
	result = seq
}

func BenchmarkClojure(b *testing.B) {
	b.ReportAllocs()

	var seq SeqB

	for n := 0; n < b.N; n++ {
		seq = SeqB(0)
		sequence.FMap(func(x int) { seq = joinAtoB(seq, x) })
	}
	result = seq
}
