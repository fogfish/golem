# Monoid for structural transformations

[![Documentation](https://godoc.org/github.com/fogfish/golem/generic?status.svg)](https://godoc.org/github.com/fogfish/golem/generic)


Golang uses imperative style for structural transformations. The usage of `for` loop is advertised by majority of language tutorials. Usage of `for` loops is an idiomatic replacement for well-known functional constructs map, filter and fold. 

```go
sum := 0
for _, x := listOfInt {
  sum = sum + x
}
```

Usage of `for` loops for anything else than primitive containers requires a lot of boilerplate, which means a lot of repeated code. A functional programming techniques solves this problem with high-order functions. Monoid is most basic concept to apply a structural transformation:

> a monoid is an algebraic structure with a single associative binary operation and an identity element.

Monoid is just a scientific name for mostly used [computer science concept](https://en.wikipedia.org/wiki/Monoid#Examples). The `for` loop example above is "commutative monoid under addition with identity element zero". As an example, MapReduce programming model is a monoid with left folding. Many iterative structural transformations may be elegantly expressed by a monoid operation:
* [map](https://en.wikipedia.org/wiki/Map_(higher-order_function)) - immutable transformation of the structure, preserving the shape but often altering a type.
* [fold](https://en.wikipedia.org/wiki/Fold_(higher-order_function)) - analysis of recursive data structure through use of monoid.
* [filter](https://en.wikipedia.org/wiki/Filter_(higher-order_function)) - produces a new data structure which contains elements accepted by predicate function.
* [comprehension](https://en.wikipedia.org/wiki/List_comprehension) - builder notation as distinct from the use of map and filter functions.

## Monoid in Go

Let's consider a typical structural transformation in Scala. It builds a new collection by applying a function to all elements of a sequence. The data structure traversal algorithm is implemented only once, which is a huge benefit for complex data structures.

```scala
trait Seq[A] {
  def map[B](f: (A) => B): Seq[B]
} 
```

It is not possible to implement generic map in Golang due to absence of generics and type covariance. The code generation does not help, we are bloating source code due to high cardinality of `A x B` set. Let's define a Monoid interface and show how it can **solve** an output type parametrization problem for transformations.

```go
func (seq Seq) Map(mB Monoid, f func (A) B) Monoid {
  y := mB.Empty()
	for _, x := range seq {
		y = y.Combine(f(x))
	}
	return y
}

type Monoid interface {
  // the identity value for a particular monoid.
  Empty() Monoid

  // an associative binary function
  Combine(x interface{}) Monoid
}
```

There are a few Go specific gotchas here:
* the library uses `struct` with receiver functions to implement a monoid for concrete type. 
* an associative binary function mutates corresponding structure in place but transformation algorithm ensures immutability. 
* semantic of associative binary function is built with `interface{}` which requires dynamic casting.

The proposed solution is 64.1% slower then `for` loops if we compare a structural transformation of arrays. 

```
monoid     	30991990	       384 ns/op	     352 B/op	      15 allocs/op
for-loop   	51079478	       234 ns/op	     280 B/op	       6 allocs/op
```

