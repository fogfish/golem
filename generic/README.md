# Generics in Go 

> Generic programming is a style of computer programming in which algorithms are written in terms of types to-be-specified-later that are then instantiated when needed for specific types provided as parameters. -- said by [Wikipedia](https://en.wikipedia.org/wiki/Generic_programming).

The discussion about Generics in Go is [never ending story](https://github.com/golang/go/wiki/ExperienceReports#generics), fortunately there is [a light in the end of the tunnel](https://www.youtube.com/watch?v=TborQFPY2IM). Often engineers associate generics as an important feature of functional programming. The computing history plays a joke with us. ML used generics as a foundation of its type system. All its descendants, including Haskell, inherited this property. Generics simplify function composition and makes combinator functions like map or filter usable with application specific types.

There is one article that draws a perfect bottom line on this subject in context of Golang: [Who needs generic? Use ... instead!](https://appliedgo.net/generics/). The primary reason for generics is the ability to [scrap your boilerplate code](https://www.microsoft.com/en-us/research/publication/scrap-your-boilerplate-with-class/). Copy-and-Paste is not an option, it only works perfectly with Stack Overflow. A group of Gophers proposes a code generator as a "perfect" solution. It resembles concept of macros in Scala or parse transforms in Erlang. [Genny](https://github.com/cheekybits/genny) is one of the famous solution. The code generator solution allows to write a valid Go code, which is compilable and testable with standard Go tools. All generic algorithms are covered with unit testing and benchmarking. Unfortunately, the code generators suffers from `O(mn)` complexity, e.g. implementation of simple transformer `map[B](f: (A) ⇒ B): Seq[B]` is a challenge. Therefore **interfaces** is only viable approach to address generic algorithms now and then in the future. In Golang interfaces are behavioral traits. Let's study this subject in depth.


## Generic list

A simple linked list is a common data type in imperative and functional programming. It holds a sequence of generic elements:

```scala
trait List[T] {
  def cons(x: T): List[T]
  def head(): T
  def tail(): List[T]
}
```

Golang has built-in [List container](https://golang.org/pkg/container/list/). It uses empty `interface{}`, which allows any data type to be used with this container. Usage of `interface{}` is a valid technique however usage of "specialized" traits brings a flavour of type safeness to your program. For example a totally ordered data types is good abstraction for list elements, anyone can imagine sorting, filtering, searching, etc.

```go
type Ord interface {
	Eq
	Lt(Ord) bool
}
```

The declaration of generic list type becomes a trivial: 

```go
type List struct {
	head golem.Ord
	tail *List
}

func (list *List) Cons(x golem.Ord) *List
func (list *List) Head() golem.Ord
func (list *List) Tail() *List
```

Compiler would not allow to use the list container unless the element type implements `Ord` interface. This technique gives sufficient protection but requires annotation of your types with the interface implementation:

```go
type String string

func (s String) Lt(x Ord) bool {
	switch v := x.(type) {
	case String:
		return s < v
	default:
		return false
	}
}
```


## "Generic" overhead

Usage of interface requires a type assertion (see `Lt` implementation above) to cast the type behind the interface into "concrete" type. Golang is efficient to make it fast. The micro benchmark of generic list (`List`), the strongly typed list of pointers (`Strong`) and standard Golang list (`StdLib`) does not show significant difference.

```
cpu: Intel(R) Core(TM) i5-3210M CPU @ 2.50GHz
BenchmarkListCons-4     	100000000	       163.9 ns/op	      32 B/op	       2 allocs/op
BenchmarkStrongCons-4   	100000000	       116.5 ns/op	      24 B/op	       2 allocs/op
BenchmarkStdLibCons-4   	100000000	       215.0 ns/op	      56 B/op	       2 allocs/op

BenchmarkListTail-4     	1000000000	     2.524 ns/op	       0 B/op	       0 allocs/op
BenchmarkStrongTail-4   	1000000000	     2.325 ns/op	       0 B/op	       0 allocs/op
BenchmarkStdLibTail-4   	1000000000	     4.305 ns/op	       0 B/op	       0 allocs/op
```

Note: the list container from standard library is slower just because it implements double-linked list. 


## Strictly typed implementation

Usage of interfaces has one disadvantage: a type T (whether its a concrete type, or itself an interface) implements an interface I, then T can be viewed as a subtype of I. The example list container fits any type that implements `Ord` interface. It does not really work as real generic, where the container type is fixed to deal with concrete type. Type casting is only the opportunity to make a type safe implementation. Unfortunately, it requires a bit of boilerplate but it is wrappers only. Just declare a type alias as set of functions to "magically" cast pointers here-there. 

```go
type ListT List

func (seq *ListT) Cons(x *T) *ListT {
	return (*ListT)((*List)(seq).Cons(x))
}

func (seq *ListT) Head() *T {
	switch v := (*List)(els).Head().(type) {
	case *T:
		return v
	default:
		panic(fmt.Errorf("Invalid element type %T %v", v, v))
	}
}

func (seq *ListT) Tail() *ListT {
	return (*ListT)((*List)(seq).Tail())
}
```

The given implementation inherits a generic list behavior. The compiler fails if you try to use other type than `T`. The back-and-forth type casting brings extra overhead, the list becomes slower:

```
cpu: Intel(R) Core(TM) i5-3210M CPU @ 2.50GHz
BenchmarkTypedCons-4    	100000000	       167.8 ns/op	      32 B/op	       2 allocs/op
BenchmarkTypedTail-4    	1000000000	     5.136 ns/op	       0 B/op	       0 allocs/op
```

## Afterwords

There are two perspective on generics

> **As a** generic library developer **I want to** define a generic type and supply its parametrized variants of standard Go type **so that** my generic is ready for application development.

> **As a** application developer **I want to** parametrize a generic types with my own application specific types **so that** the application benefits from re-use of generic implementations.

Usage of interfaces to declare a generic type trait servers both perspective. Interfaces does not bring significant overhead, its a common technique in Golang. The interface based approach still suffers a strong type safeness to compare with generics in Scala, Haskell or other especially if developers uses `interface{}` type, which is literally matches anythings. Type aliases and casting helps to wrap a "loosely" typed implementation to its strongly typed counter part.  
