package main

import (
	"fmt"
)

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

// built-in equal operator
func equal[T int | string](a, b T) bool { return a == b }

// instances of Eq type class for primitive Golang types
var (
	Int    Eq[int]    = FromEq[int](equal[int])
	String Eq[string] = FromEq[string](equal[string])
)

// UnApply2 is contra-map function for data type T that unwrap product type
type UnApply2[T, A, B any] func(T) (A, B)

/*

ProductEq2 is a shortcut due to zeroth-order type class concept in Golang and
lack of heterogenous lists. The type just a product "container" of Eq instances
*/
type ProductEq2[T, A, B any] struct {
	Eq1 Eq[A]
	Eq2 Eq[B]
	UnApply2[T, A, B]
}

// implementation of Eq type class for the product
func (eq ProductEq2[T, A, B]) Equal(a, b T) bool {
	a0, a1 := eq.UnApply2(a)
	b0, b1 := eq.UnApply2(b)
	return eq.Eq1.Equal(a0, b0) && eq.Eq2.Equal(a1, b1)
}

// ExampleType product type is product of primitive types int × string
type ExampleType struct {
	A int
	B string
}

func main() {
	eq := ProductEq2[ExampleType, int, string]{Int, String,
		func(x ExampleType) (int, string) { return x.A, x.B },
	}

	a := ExampleType{1, "a"}
	b := ExampleType{1, "a"}
	c := ExampleType{2, "a"}

	fmt.Printf("%v and %v equals %v\n", a, b, eq.Equal(a, b))
	fmt.Printf("%v and %v equals %v\n", a, c, eq.Equal(a, c))
}
