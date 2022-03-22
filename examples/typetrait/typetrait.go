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

eqInt declares a new instance of Eq trait, which is a real type.
The real type "knows" everything about equality in own domain (e.g. int type).
The instance of Eq is created as type over string, it is an intentional
technique to create a namespace using Golang constants. The instance of trait is referenced as eq.Int in the code.
*/
type eqInt string

// the type "implements" equality behavior
func (eqInt) Equal(a, b int) bool { return a == b }

/*

Int is an instance of Eq trait for int domain as immutable value so that
other functions can use this constant like `eq.Int.Equal(...)`
*/
const Int = eqInt("eq.int")

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
	haystack := Haystack[int]{Int}
	haystack.Lookup(2, []int{1, 2, 3})
	haystack.Lookup(5, []int{1, 2, 3})
}
