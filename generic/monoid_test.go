package generic_test

import (
	"reflect"
	"testing"

	"github.com/fogfish/golem/generic"
)

//
// example of monoid implementation
type seq struct {
	x []int
}

func (s *seq) Empty() generic.Monoid {
	return &seq{}
}

func (s *seq) Combine(x interface{}) generic.Monoid {
	switch v := x.(type) {
	case int:
		s.x = append(s.x, v)
	case *seq:
		s.x = append(s.x, v.x...)
	}
	return s
}

//
// generic algorithm that uses monoid
func (s *seq) fmap(f func(int) int) *seq {
	y := s.Empty()
	for _, x := range s.x {
		y.Combine(f(x))
	}
	return y.(*seq)
}

func TestIdentity(t *testing.T) {
	a := &seq{[]int{1}}
	b := &seq{[]int{1}}

	x := a.Combine(a.Empty())
	y := b.Empty().Combine(b)

	if !reflect.DeepEqual(x, y) {
		t.Fatalf("monoid violates identity %v != %v\n", x, y)
	}
}

func TestAssociativity(t *testing.T) {
	x := (&seq{[]int{1}}).Combine(&seq{[]int{2}}).Combine(&seq{[]int{3}})
	y := (&seq{[]int{1}}).Combine((&seq{[]int{2}}).Combine(&seq{[]int{3}}))

	if !reflect.DeepEqual(x, y) {
		t.Fatalf("monodi violates associativity %v != %v\n", x, y)
	}
}

var bSeq *seq

func BenchmarkMonoid(b *testing.B) {
	var m *seq
	x := &seq{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}}

	for n := 0; n < b.N; n++ {
		m = x.fmap(func(x int) int { return x * 10 })
	}
	bSeq = m
}

func BenchmarkForLoop(b *testing.B) {
	var m *seq
	x := &seq{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}}
	for n := 0; n < b.N; n++ {
		m = &seq{}
		for _, v := range x.x {
			m.x = append(m.x, v*10)
		}
	}
	bSeq = m
}
