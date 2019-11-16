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

func (seq *MSeq) Empty() generic.Monoid {
	return &MSeq{}
}

func (seq *MSeq) Combine(x interface{}) generic.Monoid {
	switch v := x.(type) {
	case int:
		seq.value = append(seq.value, v)
	case *MSeq:
		seq.value = append(seq.value, v.value...)
	}
	return seq
}

//
// generic algorithm
type String struct {
	value []string
}

func (seq *String) Map(m generic.Monoid, f func(string) int) generic.Monoid {
	y := m.Empty()
	for _, x := range seq.value {
		y = y.Combine(f(x))
	}
	return y
}

//
// global constants
var result *MSeq
var sequence String = String{[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}}

func atoi(str string) int {
	int, _ := strconv.Atoi(str)
	return int
}

//
// unit tests
func TestIdentity(t *testing.T) {
	a := New(1)
	b := New(1)

	x := a.Combine(a.Empty())
	y := b.Empty().Combine(b)

	if !reflect.DeepEqual(x, y) {
		t.Fatalf("monoid violates identity %v != %v\n", x, y)
	}
}

func TestAssociativity(t *testing.T) {
	x := New(1).Combine(New(2)).Combine(New(3))
	y := New(1).Combine(New(2).Combine(New(3)))

	if !reflect.DeepEqual(x, y) {
		t.Fatalf("monoid violates associativity %v != %v\n", x, y)
	}
}

func TestMap(t *testing.T) {
	expect := &MSeq{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}
	value := sequence.Map(New(0), atoi)

	if !reflect.DeepEqual(expect, value) {
		t.Fatalf("failed to map sequence %v != %v\n", expect, value)
	}
}

//
// benchmarks
func BenchmarkMonoid(b *testing.B) {
	var monoid *MSeq

	for n := 0; n < b.N; n++ {
		monoid = sequence.Map(monoid, atoi).(*MSeq)
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
