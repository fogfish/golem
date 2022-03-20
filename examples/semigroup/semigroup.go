package main

import "fmt"

/*

The Foldable class represents data structures that can be reduced to
a summary value one element at a time.
*/
type Foldable[T any] interface {
	Fold(a T, seq []T) (x T)
}

/*

Semigroup is an algebraic structure consisting of a set together
with an associative binary operation.
*/
type Semigroup[T any] interface {
	Combine(T, T) T
}

//
type Folder[T any] struct{ Semigroup[T] }

func (f Folder[T]) Fold(a T, seq []T) (x T) {
	x = a
	for _, y := range seq {
		x = f.Semigroup.Combine(x, y)
	}
	return
}

type semigroupInt string

func (semigroupInt) Combine(a, b int) int { return a + b }

const Int = semigroupInt("semigroup.int")

func main() {
	var folder Foldable[int] = Folder[int]{Int}

	seq := []int{1, 2, 3, 4, 5}
	sum := folder.Fold(0, seq)

	fmt.Printf("sum of %v is %v\n", seq, sum)
}
