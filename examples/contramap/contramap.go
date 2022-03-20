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

// built-in equal operator
func equal(a, b int) bool { return a == b }

// instances of Eq type class for primitive Golang types
var (
	Int Eq[int] = FromEq[int](equal)
)

/*

ContraMapEq is a combinator that build a new instance of type class Eq[B] using
existing instance of Eq[A] and f: b ⟼ a
*/
type ContraMapEq[A, B any] struct{ Eq[A] }

// implementation of contra variant functor
func (c ContraMapEq[A, B]) FMap(f func(B) A) Eq[B] {
	return FromEq[B](func(a, b B) bool {
		return c.Eq.Equal(f(a), f(b))
	})
}

// ExampleType product type is product of primitive types int × string
type ExampleType struct {
	A int
	B string
}

func main() {
	eq := ContraMapEq[int, ExampleType]{Int}.FMap(
		func(x ExampleType) int { return x.A },
	)

	a := ExampleType{1, "a"}
	b := ExampleType{1, "a"}
	c := ExampleType{2, "a"}

	fmt.Printf("%v and %v equals %v\n", a, b, eq.Equal(a, b))
	fmt.Printf("%v and %v equals %v\n", a, c, eq.Equal(a, c))
}
