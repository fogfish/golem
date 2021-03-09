package seq_test

import (
	"testing"

	"github.com/fogfish/golem/seq"
)

func BenchmarkListCons(b *testing.B) {
	mkList(b.N)
}

func BenchmarkListFMap(b *testing.B) {
	b.StopTimer()
	s := mkList(100)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		s.FMap(func(seq.ISeq) error { return nil })
	}
}

func BenchmarkListFMapTyped(b *testing.B) {
	b.StopTimer()
	s := mkList(100)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		s.FMap(isValue)
	}
}

//
//
func BenchmarkListTypedCons(b *testing.B) {
	mkListValue(b.N)
}

func BenchmarkListTypedFMap(b *testing.B) {
	b.StopTimer()
	s := mkListValue(100)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		s.FMap(func(Value) error { return nil })
	}
}

//
//
// utility functions
//
//
func mkList(size int) *seq.List {
	s := &seq.List{}
	for n := 0; n < size; n++ {
		s = s.Cons(&Value{ID: n})
	}
	return s
}

func mkListValue(size int) *ListValue {
	s := NewListValue()
	for n := 0; n < size; n++ {
		s = s.Cons(&Value{ID: n})
	}
	return s
}
