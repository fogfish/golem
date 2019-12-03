//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package generic_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/fogfish/golem/generic"
)

//
// example of monoid implementation
type MSeq struct {
	value []int
}

func New(x int) *MSeq {
	return &MSeq{[]int{x}}
}

func (seq *MSeq) Mempty() generic.Monoid {
	return &MSeq{}
}

func (seq *MSeq) Mappend(x interface{}) generic.Monoid {
	switch v := x.(type) {
	case *MSeq:
		seq.value = append(seq.value, v.value...)
	case int:
		seq.value = append(seq.value, v)
	}
	return seq
}

// type safe associative binary function
func (seq *MSeq) Append(x int) {
	seq.value = append(seq.value, x)
}

//
// generic algorithm
type String struct {
	value []string
}

// Map with Monoid
func (seq *String) Map(m generic.Monoid) func(func(string) interface{}) generic.Monoid {
	return func(f func(string) interface{}) generic.Monoid {
		y := m.Mempty()
		for _, x := range seq.value {
			y = y.Mappend(f(x))
		}
		return y
	}
}

func (seq *String) MMap(m generic.Monoid, f func(string) interface{}) generic.Monoid {
	y := m.Mempty()
	for _, x := range seq.value {
		y = y.Mappend(f(x))
	}
	return y
}

func (seq *String) FMap(f func(string)) {
	for _, x := range seq.value {
		f(x)
	}
}

//
// String x Monoid
type Product struct {
	String
	M generic.Monoid
}

func (p *Product) Map(f func(string) interface{}) generic.Monoid {
	y := p.M.Mempty()
	for _, x := range p.String.value {
		y = y.Mappend(f(x))
	}
	return y
}

//
// global constants
var result *MSeq
var sequence String = String{[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}}

func atog(str string) interface{} {
	int, _ := strconv.Atoi(str)
	return int
}

func atoi(str string) int {
	int, _ := strconv.Atoi(str)
	return int
}

//
// unit tests
func TestIdentity(t *testing.T) {
	a := New(1)
	b := New(1)

	x := a.Mappend(a.Mempty())
	y := b.Mempty().Mappend(b)

	if !reflect.DeepEqual(x, y) {
		t.Fatalf("monoid violates identity %v != %v\n", x, y)
	}
}

func TestAssociativity(t *testing.T) {
	x := New(1).Mappend(New(2)).Mappend(New(3))
	y := New(1).Mappend(New(2).Mappend(New(3)))

	if !reflect.DeepEqual(x, y) {
		t.Fatalf("monoid violates associativity %v != %v\n", x, y)
	}
}

func TestMap(t *testing.T) {
	expect := &MSeq{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}
	value := sequence.MMap(New(0), atog)

	if !reflect.DeepEqual(expect, value) {
		t.Fatalf("failed to map sequence %v != %v\n", expect, value)
	}
}

//
// benchmarks
func BenchmarkMonoid(b *testing.B) {
	var monoid *MSeq

	for n := 0; n < b.N; n++ {
		monoid = sequence.MMap(monoid, atog).(*MSeq)
	}
	result = monoid
}

func BenchmarkMonoidHoF(b *testing.B) {
	var monoid *MSeq
	mapper := sequence.Map(monoid)

	for n := 0; n < b.N; n++ {
		monoid = mapper(atog).(*MSeq)
	}
	result = monoid
}

func BenchmarkMonoidProduct(b *testing.B) {
	var monoid *MSeq
	product := &Product{sequence, monoid}

	for n := 0; n < b.N; n++ {
		monoid = product.Map(atog).(*MSeq)
	}
	result = monoid
}

func BenchmarkForLoop(b *testing.B) {
	var seq *MSeq

	for n := 0; n < b.N; n++ {
		seq = &MSeq{}
		for _, x := range sequence.value {
			seq.value = append(seq.value, atoi(x))
		}
	}
	result = seq
}

func BenchmarkClojure(b *testing.B) {
	var seq *MSeq

	for n := 0; n < b.N; n++ {
		seq = &MSeq{}
		sequence.FMap(func(x string) { seq.Append(atoi(x)) })
	}
	result = seq
}
