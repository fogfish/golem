package seq_test

import (
	"testing"

	"github.com/fogfish/golem/seq"
)

func BenchmarkSeqAppend(b *testing.B) {
	mkSeq(b.N)
}

func BenchmarkSeqFMap(b *testing.B) {
	b.StopTimer()
	s := mkSeq(100)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		s.FMap(func(seq.ISeq) error { return nil })
	}
}

func BenchmarkSeqFMapTyped(b *testing.B) {
	b.StopTimer()
	s := mkSeq(100)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		s.FMap(isValue)
	}
}

//
//
// utility functions
//
//
func mkSeq(size int) seq.Seq {
	s := seq.Seq{}
	for n := 0; n < size; n++ {
		s = append(s, &Value{ID: n})
	}
	return s
}

// func pure(n int) []*Value {
// 	seq := []*Value{}
// 	for i := 0; i < n; i++ {
// 		seq = append(seq, &Value{ID: i})
// 	}
// 	return seq
// }

// func unit(n int) seq.Seq {
// 	seq := seq.Seq{}
// 	for i := 0; i < n; i++ {
// 		seq.Append(&Value{ID: i})
// 	}
// 	return seq
// }

// var valPure *Value
// var seqPure []*Value = pure(100)

// var valType seq.ISeq
// var seqType seq.Seq = unit(100)

// func BenchmarkAppendPure(b *testing.B) {
// 	pure(b.N)
// }

// func BenchmarkAppendType(b *testing.B) {
// 	unit(b.N)
// }

// func BenchmarkFMapPure(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		for _, x := range seqPure {
// 			if err := usePure(x); err != nil {
// 				break
// 			}
// 		}
// 	}
// }

// func BenchmarkFMapType(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		seqType.FMap(useType)
// 	}
// }

// func useType(x seq.ISeq) error {
// 	switch t := x.(type) {
// 	case *Value:
// 		valPure = t
// 	}
// 	return nil
// }

// func usePure(x *Value) error {
// 	// valPure = x
// 	return nil
// }

// func BenchmarkX(b *testing.B) {

// 	for n := 0; n < b.N; n++ {
// 		v := Value{}
// 		r := v.X()
// 		runtime.KeepAlive(r)
// 	}
// }
