package generic_test

import (
	"reflect"
	"testing"

	"github.com/fogfish/golem/generic"
)

//
// example of monoid implementation
type Int int

func (v Int) Empty() generic.Monoid {
	return Int(0)
}

func (v Int) Combine(x interface{}) generic.Monoid {
	return v + x.(Int)
}

//
// generic algorithm
type seq struct {
	x []Int
}

func (s *seq) fold(f func(Int) Int, m generic.Monoid) generic.Monoid {
	y := m.Empty()
	for _, x := range s.x {
		y = y.Combine(f(x))
	}
	return y
}

func TestIdentity(t *testing.T) {
	x := Int(1).Combine(Int(1).Empty())
	y := Int(1).Empty().Combine(Int(1))

	if !reflect.DeepEqual(x, y) {
		t.Fatalf("monoid violates identity %v != %v\n", x, y)
	}
}

func TestAssociativity(t *testing.T) {
	x := Int(1).Combine(Int(2)).Combine(Int(3))
	y := Int(1).Combine(Int(2).Combine(Int(3)))

	if !reflect.DeepEqual(x, y) {
		t.Fatalf("monoid violates associativity %v != %v\n", x, y)
	}
}

func TestFold(t *testing.T) {
	s := seq{[]Int{0, 1, 2, 3, 4, 5}}
	y := s.fold(func(x Int) Int { return x }, Int(0))
	if !reflect.DeepEqual(Int(15), y) {
		t.Fatalf("failed to fold with monoid sum( %v ) != %v\n", s, y)
	}
}

var result Int

func BenchmarkMonoid(b *testing.B) {
	var m Int
	x := seq{[]Int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}

	for n := 0; n < b.N; n++ {
		m = x.fold(func(x Int) Int { return x * 10 }, Int(0)).(Int)
	}
	result = m
}

func BenchmarkForLoop(b *testing.B) {
	var m Int
	x := seq{[]Int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}

	for n := 0; n < b.N; n++ {
		m = 0
		for _, v := range x.x {
			m = m + v*10
		}
	}
	result = m
}
