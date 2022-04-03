package main

import (
	"fmt"
)

/*

HKT[F, A] ∼ F[A]

Transforms a computation with higher-kinded type expressions into a
computation where all type expressions are of kind `*`. The abstract type
constructor `HKT` represent an idea of parametrized container type `F[A]`
*/
type HKT[F, A any] interface {
	HKT1(F)
	HKT2(A)
}

/*

SeqType is opaque type to define polymorphic context of Seq.
It makes HKT typesafe in the context of sequence trait
*/
type SeqType any

/*

Higher-Kinded Sequence type
*/
type SeqKind[A any] HKT[SeqType, A]

/*

Seq type-trait (see type-trait pattern)
*/
type Seq[F_, A any] interface {
	SeqKind[A]

	Length(F_) int
	Head(F_) A
	Tail(F_) F_
}

/*

SeqSlice is a sequence data structure implemented using slice
*/
type SeqSlice[A any] []A

// tagging as HKT[SeqType, A]
func (SeqSlice[A]) HKT1(SeqType) {}
func (SeqSlice[A]) HKT2(A)       {}

/*

SeqSliceT implements Seq type-trait for the sequence
*/
type SeqSliceT[A any] string

// tagging as HKT[SeqType, A]
func (SeqSliceT[A]) HKT1(SeqType) {}
func (SeqSliceT[A]) HKT2(A)       {}

// Seq type-trait
func (SeqSliceT[A]) Length(seq SeqSlice[A]) int       { return len(seq) }
func (SeqSliceT[A]) Head(seq SeqSlice[A]) A           { return seq[0] }
func (SeqSliceT[A]) Tail(seq SeqSlice[A]) SeqSlice[A] { return seq[1:] }

// transformer for sequence
func mapSeqSlice[A, B any](f func(A) B, a SeqSlice[A]) SeqSlice[B] {
	b := make(SeqSlice[B], len(a))
	for i, x := range a {
		b[i] = f(x)
	}
	return b
}

/*

SeqList is a sequence data structure implemented using linked list
*/
type SeqList[A any] struct {
	Head A
	Tail *SeqList[A]
}

// tagging as HKT[SeqType, A]
func (SeqList[A]) HKT1(SeqType) {}
func (SeqList[A]) HKT2(A)       {}

/*

SeqListT implements Seq type-trait for the sequence
*/
type SeqListT[A any] string

// tagging as HKT[SeqType, A]
func (SeqListT[A]) HKT1(SeqType) {}
func (SeqListT[A]) HKT2(A)       {}

// Seq type-trait
func (SeqListT[A]) Length(seq *SeqList[A]) (len int) {
	for e := seq; e != nil; e = e.Tail {
		len++
	}
	return
}
func (SeqListT[A]) Head(seq *SeqList[A]) A           { return seq.Head }
func (SeqListT[A]) Tail(seq *SeqList[A]) *SeqList[A] { return seq.Tail }

/*

Show is a generic algorithm that outputs the content of sequence
*/
type Show[F_, A any] struct{ Seq[F_, A] }

func (f Show[F_, A]) Print(fa F_) {
	fmt.Printf("==>")

	x := fa
	for f.Seq.Length(x) != 0 {
		fmt.Printf(" %v", f.Seq.Head(x))
		x = f.Seq.Tail(x)
	}

	fmt.Println()
}

/*

Unary is an abstraction of computation over `F[_]` container
*/
type Unary[A any, FA SeqKind[A]] func(FA)

func (f Unary[A, FA]) FMap(fa FA) { f(fa) }

/*

Functor is an abstraction of morphism over `F[)]` container
*/
type Functor[
	A any, FA SeqKind[A],
	B any, FB SeqKind[B],
] func(f func(A) B, a FA) FB

func (fn Functor[A, FA, B, FB]) FMap(f func(A) B, fa FA) FB { return fn(f, fa) }

/*

Generic computation to flatten convert flat sequence into sequence of sequences.
It uses a functor abstraction to make it generic over any container
*/
func Unflattening[
	A any,
	FA SeqKind[A],
	FB SeqKind[[]A],
](f Functor[A, FA, []A, FB], fa FA) FB {
	return f.FMap(func(a A) []A { return []A{a} }, fa)
}

/*

Another example of Generic computation
*/
type ShowGen[T, F_, A any] struct{ HKT[T, A] }

func (f ShowGen[T, F_, A]) Print(fa HKT[T, A]) {
	len := f.HKT.(Seq[F_, A]).Length(fa.(F_))
	fmt.Printf("==> sequence of %v\n", len)
}

//
// type instances
const (
	SliceOfInts = SeqSliceT[int]("seq.slice.int")
	ListOfInts  = SeqListT[int]("seq.list.int")
)

var (
	ShowSliceOfInts = Show[SeqSlice[int], int]{SliceOfInts}
	ShowListOfInts  = Show[*SeqList[int], int]{ListOfInts}
)

// helpers to create instances of seq
func MkSlice[A any](seq ...A) SeqSlice[A] {
	return SeqSlice[A](seq)
}

func MkList[A any](seq ...A) *SeqList[A] {
	var s *SeqList[A]
	for i := len(seq) - 1; i >= 0; i-- {
		s = &SeqList[A]{Head: seq[i], Tail: s}
	}
	return s
}

func main() {
	sInts := MkSlice(1, 2, 3, 4, 5)
	lInts := MkList(1, 2, 3, 4, 5)

	//
	ShowSliceOfInts.Print(sInts)
	ShowListOfInts.Print(lInts)

	//
	print := Unary[int, SeqSlice[int]](ShowSliceOfInts.Print)
	print.FMap(sInts)

	//
	showGen := ShowGen[SeqType, *SeqList[int], int]{ShowListOfInts}
	showGen.Print(lInts)

	//
	f := Functor[
		int, SeqSlice[int],
		[]int, SeqSlice[[]int],
	](mapSeqSlice[int, []int])

	y := Unflattening(f, sInts)
	fmt.Printf("%T %v\n", y, y)
}
