package main

import (
	"testing"
)

/*

Walk is a generic algorithm that traverse the content of sequence
*/
type Walk[F_, A any] struct{ Seq[F_, A] }

func (f Walk[F_, A]) Visit(fa F_) {
	x := fa
	for f.Seq.Length(x) != 0 {
		x = f.Seq.Tail(x)
	}
}

var (
	WalkSliceOfInts = Walk[SeqSlice[int], int]{SliceOfInts}
	WalkListOfInts  = Walk[*SeqList[int], int]{ListOfInts}

	sliceOfInts = MkSlice(1, 2, 3, 4, 5)
	listOfInts  = MkList(1, 2, 3, 4, 5)
)

func BenchmarkWalkSlice(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		WalkSliceOfInts.Visit(sliceOfInts)
	}
}

func BenchmarkWalkList(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		WalkListOfInts.Visit(listOfInts)
	}
}

func BenchmarkNativeSlice(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		mapSeqSlice(func(a int) int { return a }, sliceOfInts)
	}
}
