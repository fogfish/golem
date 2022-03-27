# Monoid for structural transformations

[![Documentation](https://godoc.org/github.com/fogfish/golem/generic?status.svg)](https://godoc.org/github.com/fogfish/golem/generic)


Golang uses imperative style for structural transformations. The usage of `for` loop is advertised by majority of language tutorials. Usage of `for` loops is an idiomatic replacement for well-known high-order functional constructs like `map`, `filter` and `fold`. 

```go
sum := 0
for _, x := listOfInt {
  sum = sum + x
}
```

Usage of `for` loops for anything else than primitive containers requires a lot of boilerplate, which means a lot of repeated code. For example, iteration over skip lists, trees or other complex data types requires exposure of type internals to client code, therefor making unnecessary coupling. [Monoid](https://en.wikipedia.org/wiki/Monoid) is the correct and most basic content to apply a structural transformation for any data structure:

> a monoid is an algebraic structure with a single associative binary operation and an identity element.

Monoid is just a scientific name for mostly used [computer science concept](https://en.wikipedia.org/wiki/Monoid#Examples). The example above uses "commutative monoid under addition (`+`) with identity element zero (`0`)". Any type that defines "identity element" and implements "associative binary operation" is monoid:

```go
type Monoid interface {
  // the identity value for a particular monoid.
  Mempty() Monoid

  // an associative binary operation
  Mappend(x Monoid) Monoid
}
```

The `for` loop example becomes 

```go
type MInt int
func (v MInt) Mempty() Monoid { return MInt(0) }
func (v MInt) Mappend(x Monoid) Monoid { return m + x.(MInt) }

sum := MInt(0)
for _, x := listOfInt {
  sum = sum.Mappend(x)
}
```

## Generic algorithm with monoid

Let's consider a typical structural transformation in Scala. It builds a new collection by applying a function to all elements of a container type. The data structure traversal algorithm is implemented only once, which is a huge benefit for complex data structures.

```scala
trait Seq[A] {
  def map[B](f: (A) => B): Seq[B]
} 
```

It is not possible to implement generic map in Golang due to absence of generics and type covariance. The code generation does not help, we are bloating source code due to high cardinality of `A x B` set. Fortunately, high-order structural transformation functions can be implemented in terms of monoid abstraction. Anyone is able to implement generic `map`, `filter`, `fold`, leaving concrete type `B` behind monoid interface. As an example, MapReduce programming model is a monoid with left folding. Many iterative structural transformations may be elegantly expressed by a monoid operation:
* [fold](https://en.wikipedia.org/wiki/Fold_(higher-order_function)) - analyse "recursive" data structure through use of monoid;
* [map](https://en.wikipedia.org/wiki/Map_(higher-order_function)) - immutable transformation of the structure, preserving the shape but often altering a type;
* [filter](https://en.wikipedia.org/wiki/Filter_(higher-order_function)) - produces a new data structure which contains elements accepted by predicate function;
* [comprehension](https://en.wikipedia.org/wiki/List_comprehension) - builder notation as distinct from the use of map and filter functions.

As an example of generic transformation using monoid abstraction:

```go
func (seq SeqA) Map(mB Monoid, f func (A) Monoid) Monoid {
  b := mB.Mempty()
	for _, x := range seq {
		b = b.Mappend(f(x))
	}
	return b
}
```

## Efficiency of monoid in Go

Usage of monoid interface has its own cost - it makes code "slower". The benchmark of the original example using classical `for` loop and its monoid counterpart shows significant difference. Unfortunately, an interface abstraction and its conversion to concrete types and back introduces the overhead (there are good post that shed the light about this overhead: [Research on the interface of golang](https://laptrinhx.com/research-on-the-interface-of-golang-4184713904/) and [Why Go is Not My Favorite Programming Language](http://www.cofault.com/2019/01/why-go-is-not-my-favorite-programming.html)).

```
cpu: Intel(R) Core(TM) i5-3210M CPU @ 2.50GHz
BenchmarkMonoid-4       	46114359	       136.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkForLoop-4      	425828221	        14.12 ns/op	       0 B/op	       0 allocs/op

Showing nodes accounting for 1160ms, 100% of 1160ms total
      flat  flat%   sum%        cum   cum%
     380ms 32.76% 32.76%     1140ms 98.28%  github.com/fogfish/golem/pure_test.SeqA.MMap
     360ms 31.03% 63.79%      390ms 33.62%  runtime.convT64
     190ms 16.38% 80.17%      380ms 32.76%  github.com/fogfish/golem/pure_test.SeqB.Mappend
     150ms 12.93% 93.10%      360ms 31.03%  github.com/fogfish/golem/pure_test.convAtoB
      60ms  5.17% 98.28%       60ms  5.17%  runtime.newstack
      20ms  1.72%   100%     1160ms   100%  github.com/fogfish/golem/pure_test.BenchmarkMonoid

Showing nodes accounting for 1110ms, 100% of 1110ms total
      flat  flat%   sum%        cum   cum%
     470ms 42.34% 42.34%      880ms 79.28%  github.com/fogfish/golem/pure_test.forEach
     370ms 33.33% 75.68%      400ms 36.04%  github.com/fogfish/golem/pure_test.joinAtoB (inline)
     210ms 18.92% 94.59%     1110ms   100%  github.com/fogfish/golem/pure_test.BenchmarkForLoop
```

This naive examples has shown 10x difference between "traditional" code and monoid abstraction. The overhead becomes less visible with growth of types complexity. For example, the monoid abstraction of traditional slice is 5x slower. 


## Structural transformations with clojure

The advantage of Monoid interface is the ability to apply transformation to data structure of any shape. The disadvantage is the overhead of runtime interface conversion. Usage of clojure mitigates discussed issue and leverage transformations with type safety.

> Go functions may be closures. A closure is a function value that references variables from outside its body. The function may access and assign to the referenced variables; in this sense the function is "bound" to the variables.

We can assume a generic transformation as application of type-safe closure over data structure elements:

```go
func (seq SeqA) FMap(f func(A)) {
	for _, x := range seq.value {
		f(x)
	}
}
```

The monoid definition becomes less formal in the context of clojure. The identity value is the type constructor and associative binary function is what-ever combine operation.  

```go
// associative binary function
func (seq SeqB) Append(x B) {/* ... */}

// identity value
seqB = SeqB{}
seqA.FMap(func(x A) { seqB.Append(/* A -> B */) })
```

Usage of clojure shows comparable performance with `for` loops. It is still slower for naive transformations algorithms but this difference fades away with any real life scenario.

```
cpu: Intel(R) Core(TM) i5-3210M CPU @ 2.50GHz
BenchmarkMonoid-4       	46114359	       136.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkForLoop-4      	425828221	        14.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkClojure-4      	168989125	        35.48 ns/op	       0 B/op	       0 allocs/op
```

## Afterwords

Monoid is an abstract concept in computer science that helps with validation and implementation of generic algorithms. Usage of monoid in Go makes sense mainly for structural transformations. The classical interface-based monoid abstraction introduces visible overhead which is not desired in highly loaded solutions. However, usage of clojure delivery monoid-like `A ⟼ B` type transformation while resolving issue concerning a runtime overhead.

## Related articles

1. [Monoids for Gophers](https://medium.com/@groveriffic/monoids-for-gophers-907175bb6165)
2. [Foldable Go](https://medium.com/zendesk-engineering/foldable-go-d74fb9cf2fc9)
3. [Cats: Monoid](https://typelevel.org/cats/typeclasses/monoid.html)
4. [Functors, Applicative Functors and Monoids](http://learnyouahaskell.com/functors-applicative-functors-and-monoids)
