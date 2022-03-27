package main

import "fmt"

/*

Eq : T ⟼ T ⟼ bool
Each type implements own equality, mapping pair of value to bool category
*/
type Eq[T any] interface {
	Equal(T, T) bool
}

/*

Seq defines fundamental general purpose sequence
*/
type Seq[S, T any] interface {
	Head(S) *T
	Tail(S) S
	IsVoid(S) bool
}

/*

type "implements" equality for integers
*/
type eqInt string

func (eqInt) Equal(a, b int) bool { return a == b }

// EqInt instance of Eq trait for int
const EqInt = eqInt("eq.int")

/*

type "implements" sequence for slice of integers
*/
type seqSlice[T any] string

func (seqSlice[T]) Head(seq []T) *T {
	if len(seq) == 0 {
		return nil
	}

	return &seq[0]
}

func (seqSlice[T]) Tail(seq []T) []T {
	if len(seq) == 0 {
		return seq
	}

	return seq[1:]
}

func (seqSlice[T]) IsVoid(seq []T) bool {
	return len(seq) == 0
}

// SliceInt instance of Seq trait for slice of ints
const SliceInt = seqSlice[int]("seq.slice")

/*

SeqEq is heterogenous product of Seq and Eq laws.
It composes two types together that "knows" how to compare sequences.
*/
type SeqEq[S, T any] struct {
	Seq[S, T]
	Eq[T]
}

// implements equality rule for sequence using Seq & Eq type classes.
func (seq SeqEq[S, T]) Equal(a, b S) bool {
	seqA := a
	seqB := b
	for !seq.Seq.IsVoid(seqA) && !seq.Seq.IsVoid(seqB) {
		headA := seq.Seq.Head(seqA)
		headB := seq.Seq.Head(seqB)
		if headA == nil || headB == nil || !seq.Eq.Equal(*headA, *headB) {
			return false
		}

		seqA = seq.Seq.Tail(seqA)
		seqB = seq.Seq.Tail(seqB)
	}

	return seq.Seq.IsVoid(seqA) && seq.Seq.IsVoid(seqB)
}

func main() {
	seq := SeqEq[[]int, int]{SliceInt, EqInt}

	a := []int{1, 2, 3, 4, 5}
	b := []int{1, 2, 3, 4, 6}
	c := []int{1, 2, 3, 4}

	fmt.Printf("Seq %v and %v are equal %v\n", a, a, seq.Equal(a, a))
	fmt.Printf("Seq %v and %v are equal %v\n", a, b, seq.Equal(a, b))
	fmt.Printf("Seq %v and %v are equal %v\n", a, c, seq.Equal(a, c))
}
