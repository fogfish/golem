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

FromEq is a combinator that lifts T ⟼ T ⟼ bool function to Eq type class
*/
type FromEq[T any] func(T, T) bool

// implementation of Eq type class
func (f FromEq[T]) Equal(a, b T) bool { return f(a, b) }

// Equal is a pure function that compares two integers
func Equal(a, b int) bool { return a == b }

/*

Haystack is an example of algorithms that uses type law
*/
type Haystack[T any] struct{ Eq[T] }

func (h Haystack[T]) Lookup(e T, seq []T) {
	for _, x := range seq {
		if h.Eq.Equal(e, x) {
			fmt.Printf("needle %v found at %v\n", e, seq)
			return
		}
	}

	fmt.Printf("needle %v is not found at %v\n", e, seq)
}

func main() {
	haystack := Haystack[int]{FromEq[int](Equal)}
	haystack.Lookup(2, []int{1, 2, 3})
	haystack.Lookup(5, []int{1, 2, 3})
}
