# A Guide To Pure Type Combinators in Golang or How to Stop Worrying and Love the Functional Programming

Humans have developed ideas of representing things using a formal "system" since ancient history - Aristotle’s logic, Euclid’s geometry are good examples. The formalism allows anyone to proof and deduct purely within the system so that it defines a concrete solution of the current problem. System of combinators has been known for 100 years since Moses Schönfinkel developed a universal computation system that has been researched since together with mathematical logic, lambda calculus and category theory. Combinators open up an opportunity to depict computation problems in terms of fundamental elements like physics talks about the universe in terms of particles. The only definite purpose of combinators are building blocks for composition of "atomic" functions into computational structures from concrete problem "domain". So far, combinators remain as powerful symbolic expressions in computational languages.  

Combinators are simple and do not involve any advanced math in its definition. A combinator builds new "things" from previously defined "things".

```
ƒ: Thing ⟼ Thing ⟼ Thing
```

Since, "thing" can be any computational "element" including functions and other combinators. It delivers powerful **combinator patterns** for functional programming - a style of declaring a small set of primitive abstractions and collection of combinators to define advanced structures. Golang, like any other languages, supports first class functions. It allows functions to be assigned to variables, passed as arguments to other functions and returned from other functions. Therefore, combinators of pure functions is a given fact for any Golang application. 

Let's advance these patterns towards the Golang type system and define combinators over types and their instances to derive complex structures of type T. There are 7 patterns to consider and express their semantic with Golang:

1. [Type Trait](#type-trait) (`𝔗 ⟼ A ⟼ ƒ(𝔗[A], A)`) declares type class, it's laws and the intent.
2. [Sub-typing](#sub-typing) (`A <: B`) enhances existing type classes.
3. [Lifting](#lifting) (`ƒ ⟼ 𝔗[A]`) transforms a pure function into corresponding type trait.
4. [Homogenous product](#homogenous-product) (`A × B × … ⟼ 𝔗 ⟼ 𝔗[A × B × …]`) composes type classes of same kind to operate with product type.
5. [Contra Variant Functor](#contra-variant-functor) (`(ƒ: b ⟼ a) ⟼ 𝔗[A] ⟼ 𝔗[B]`) applies type transformation using pure function.
6. [Compose Generic Types](#compose-generic-types) (`𝔗 ⟼ 𝕬 ⟼ 𝔗[𝕬]`) to define generic computation.
7. [Heterogeneous product](#heterogeneous-product) (`𝔗 × 𝕬 × … ⟼ 𝕷`) compose heterogenous type classes into a new type law pattern. 
 
The implementation of each combinator is considered further in this post using simplest examples in the style of Golang.  


## Type trait

Equality (==) implements an equivalence relationship where two values comparing "is equal to" if they belong to the same equivalence class in their domain (e.g. equivalence laws of boolean, numbers, strings and other types). Typical comparison operators support only built-in types, extension of the operator over values of some abstract type T requires a definition of equality function over these types and acting differently for each type. 

```go
/*

Eq : T ⟼ T ⟼ bool
Each trait implements mapping pair of value to bool category using own
equality rules 
*/
type Eq [T any] interface {
  Equal(T, T) bool
}
```

The **type trait** definition looks like an interface that specifies an equality behavior for some kind of type. Anyone with an OOP background might be confused. The OOP style insists on `Equal(T) bool` clamming receiver instance being compared with given value. Object-oriented languages leverage subtype polymorphism, which is not efficient to build advanced combinators. Functional programing is looking towards ad-hoc polymorphism - depicting a trait of different unrelated types with type specific implementations. In this example, the interface defines an "atomic" trait that knows how to compare type instances (e.g. number equality category, string equality category, etc). The equivalence relationship implementation for Golang `int` type is the following

```go
package eq

/*

eqInt declares a new instance of Eq trait, which is a real type.
The real type "knows" everything about equality in its own domain.
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
```

This technique supports an ad-hoc polymorphism of the trait `Eq`, and detaches the data type implementation from the behavior. The application is able to implement multiple traits for the same data type together with trait sub-typing, which is not achievable if `Eq` is used as standard Golang interface.     

The type trait pattern solves the problem of polymorphic algorithm implementations - using the "well-known" function names for various instances that take different kinds of parameters. For example, equality type trait helps with implementation of "haystack" algorithms: 

```go
type Haystack[T any] struct{ Eq[T] }

func (h Haystack[T]) Lookup(a T, b []T) bool {
  for _, x := range b {
    if h.Eq.Equal(a, x) {
      return true
    }
  }
  return false
}
```

The type trait is the combinator pattern with "atomic" and composable element such as
1. A generic computation (an algorithm) over a some type trait T. This generic computation is polymorphic due to usage of same function names defined by T. Although type T is polymorphic only to the computation that uses it, where there are multiple instances for each concrete type.
2. The type trait T declares computational laws and operational intent. T is declared using a common interface for an arbitrary set of individually specified types.
3. The instance of type trait T is declared for any concrete type by implementing a declared interface for T. Golang's [structural type system](https://en.wikipedia.org/wiki/Structural_type_system) empowers different unrelated types with type specific implementations. It makes the approach flexible for pure functional combinator libraries and easy ad-hoc type extensions.

The type trait pattern looks similar to type class. Computer Science has defined the type class as a construct "that supports ad hoc polymorphism. This is achieved by adding constraints to type variables in parametrically polymorphic types. Such a constraint typically involves a type class T and a type variable a, and means that a can only be instantiated to a type whose members support the overloaded operations associated with T.". However, the Golang type system is strictly less powerful and does not support type classes, it only supports "a kind of zeroth-order type class", while concepts around type classes are often associated with higher-kinded polymorphism. The type trait abstraction as it is defined here provides better composability than the Golang interface but less powerful in comparison with Haskell's type classes.    


## Sub-typing

A **Sub-typing** pattern (or inclusion polymorphism) creates a new type trait from existing one. Sub-typing is a classical polymorphism used in object-oriented programming, sub-typing is roughly comparable with Golang embedding. The purpose of sub-typing is to enhance the interface of existing traits with declaration of new behavior. For example the class of totally ordered types `Ord` is a sub-type of `Eq`. 

```go
package ord

/*

Ord : T ⟼ T ⟼ Ordering
Each type implements compare rules, mapping pair of value to enum{ LT, EQ, GT }
*/
type Ord [T any] interface {
  Eq[T]

  Compare(T, T) Ordering
}
```

Instances of `Ord` trait do not differs from `Eq`, each type "implements" the declared specification: 

```go
func (ordInt) Compare(a, b int) Ordering { /* ... */ }
```

Instances of type trait can aggregate other traits that allow re-use previously defined implementations.

```go
func (ordInt) Equal(a, b, int) bool { return eq.Int.Equal(a, b) }
```

Sub-typing combinator is simple but yet powerful to built complex structures from "atomic" type traits. It also defines the notion of substitutability in the generic computation, which is written to operate on elements of type T, and can also operate on instances of sub-types.


## Lifting

Function and types are first class objects in Golang. The **lifting** pattern transforms a pure function into the instance of corresponding type trait - the combinator constructs a new instance from the function. Let's consider the example, there is equality function `equal: T ⟼ T ⟼ bool` and type trait `Eq`. The combinator `FromEq` create an instance of `Eq` for the function so that

```
FromEq: (T ⟼ T ⟼ bool) -> Eq
```

Golang implementation of this combinator requires definition of type `FromEq` and implementation of corresponding type class `Eq`:

```go
/*

FromEq is a combinator that lifts T ⟼ T ⟼ bool function to
an instance of Eq type trait
*/
type FromEq[T any] func(T, T) bool

// implementation of Eq type class 
func (f FromEq[T]) Equal(a, b T) bool { return f(a, b)}
```

The constructor `FromEq` takes a function `T ⟼ T ⟼ bool` as parameter, the output returns an instance of `Eq`: 

```go
// Equal is a pure function that compares two integers
func Equal(a, b int) bool { return a == b }

/*

The combinator creates a new instance of Eq trait
that "knows" how to compare int.
*/
var Int Eq[int] = eq.FromEq[int](Equal)
```

The lifting pattern is a very powerful one. It not only leverages the gap between functional and type traits domains but also facilitates the composable and re-usable definition of type traits using closures and other pure functional concepts.


## Homogenous product

Type theory and functional programming operates with algebraic data types. They are known as a composition of other types. The theory defines two classes of compositions: product types and co-product types. Product types are strongly expressed by structs in Golang; co-products are loosely defined (let's skip them at current considerations). It is not always practical to implement instances of type trait for each product type. Construction of new instances of type trait through the composition of existing instances is an alternative solution.  

The **homogenous product** pattern allows an application to create an instance of type trait for product types in a relatively boilerplate-free way. It composes type traits of the same kind to operate on algebraic data types. Let's consider a product type `T: A × B × ...` together with set of type traits for each of "elementary" type `𝔗: Eq[A], Eq[B], ...`. The homogenous product build `Eq[T]: Eq[A] × Eq[B] × ...`. 

```go
// ExampleType product type is product of primitive types int × string
type ExampleType struct {
  A int
  B string
}

// Instances of Eq type trait for primitive types
var (
  Int    Eq[int]    = FromEq[int](equal[int])
  String Eq[string] = FromEq[string](equal[string])
)
```

Golang's implementation of homogenous product requires definition of type `UnApplyN` to "extract" fractions on the product and implementation of corresponding product type `ProductN`for `Eq` type law:

```go
/*

UnApply2 is like contra-map function for data type T that unwrap product type
*/  
type UnApply2[T, A, B any] func(T) (A, B)

/*

ProductEq2 is a container, product of type trait instances.
Here, the implementation is a shortcut due to the absence of heterogeneous lists in Golang. 
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
```

The homogenous product pattern is a building block for composition of "atomic" type traits into complex structures to deal with algebraic data types. The new instance of `Eq` trait for `ExampleType` is created with

```go
ProductEq2[ExampleType, int, string]{Int, String,
  func(x ExampleType) (int, string) { return x.A, x.B },
}
```

## Contra Variant Functor

The functor pattern is one of the most discussed patterns in functional programming. It allows "a generic type to apply a function inside without changing the structure of the generic type". In math, there are many concepts that act as functors. The **contra variant** pattern just "turn morphisms around". 

A functional programming operates with algebraic data types. It is not always practical to implement type trait instances for each data type. The type mapping is a solution to transform data types so that existing type trait instances can be reused in different contexts.  

Let's consider two types `A` and `B` and the instance of type class `Eq[A]`. The contra variant functor builds an instance of `Eq[B]` with help of `f: b ⟼ a` transformer. 

```go
/*

ContraMapEq is a combinator that build a new instance of type trait Eq[B] using
existing instance of Eq[A] and f: b ⟼ a
*/
type ContraMapEq[A, B any] struct{ Eq[A] }

// implementation of contra variant functor
func (c ContraMapEq[A, B]) FMap(f func(B) A) Eq[B] {
  return FromEq[B](func(a, b B) bool {
    return c.Eq.Equal(f(a), f(b))
  })
}
```

Use the combinator to make an instance of `Eq[ExampleType]`

```go
ContraMapEq[int, ExampleType]{Int}.FMap(
  func(x ExampleType) int { /* ... */ },
)
```

Often, functional programming literature explains the purpose of contra variant functors on the example of type safe sorting or filtering algorithms where applications specific algebraic data types are processed with transformation into the domain of primitive built-in types.


## Compose generic types

Building complex type traits from "atomic" traits is the ultimate goal that allows anyone to declare complex concepts. The **compose generic types** pattern operates with two or more distinct generic types. The operation takes type trait `𝔗` and composes it with another type trait `𝕬` as the result it produces `𝔗[𝕬]` (`𝔗 ⟼ 𝕬 ⟼ 𝔗[𝕬]`).

Let's consider a `Foldable` abstraction that represents data structures that can be reduced to a summary value one element at a time:

```go
type Foldable[T any] interface {
  Fold(a T, seq []T) (x T)
}
```

There are infinite possibilities to implement the `Foldable` type trait due to the unbounded definition of the "reduce to summary" function in the specification. There is an algebraic structure, called `Semigroup`, consisting of a type set together with an associative binary operation.

```go
type Semigroup[T any] interface {
  Combine(T, T) T
}
```

The composition of generic types `Foldable` over `Semigroup` allows the implementation of a single variant of data structure reduction algorithm. The `Semigroup` parameter allow injecting the associative binary operation into the algorithms.

The composition of generic types is defined through a new type `Folder`. This type embeds `Semigroup` and implements the `Foldable` type definition:

```go
type Folder[T any] struct{ Semigroup[T] }

func (f Folder[T]) Fold(a T, seq []T) (x T) {
  x = a
  for _, y := range seq {
    x = f.Semigroup.Combine(x, y)
  }
  return
}
```

The compose generic types pattern follows Hilbert’s axiomatic method "to build everything from as few notions as possible". It uses standard Golang notations from which other combinator notations are constructed. The approach discussed by the compose generic types pattern is only the solution to parametrize one type class over another one until higher kinded polymorphism is fully supported at Golang.



## Heterogeneous product

The **heterogeneous product** pattern allows an application to compose type traits together building a new one, which "inherits" properties of composed types. Let's evaluate the purpose of the pattern in the context of following types `Seq` and `Eq`:

```go
/*

Seq defines fundamental general purpose sequence
*/
type Seq[S any, T any] interface {
  Head(S) *T
  Tail(S) S
  IsVoid(S) bool
}

/*

Eq : T ⟼ T ⟼ bool
Each type implements own equality, mapping pair of value to bool category
*/
type Eq[T any] interface {
  Equal(T, T) bool
}
```

The definition of equality for sequence requires the definition of traversal trait. However, the application does not want to re-implement equality for various classes of the sequence. Instead, it aims for re-usability building the new equality laws from existing constructs `Seq × Eq ⟼ SeqEq`. 

```go
/*

SeqEq is a heterogeneous product of Seq and Eq laws.
It composes two types together that "knows" how to compare sequences.
*/
type SeqEq[S, T any] struct {
  Seq[S, T]
  Eq[T]
}

// implements equality rule for sequence using Seq & Eq type classes.
func (seq SeqEq[S, T]) Equal(a, b S) bool { /* ... */ }
``` 

Newly composed trait is a product of two types.


## Afterwords

The post have defined a basic principles of composing types into high order constructs. It follows Hilbert’s axiomatic method "to build everything from as few notions as possible". All defined combinators use standard Golang notations from which other combinator notations are constructed. Think about combinator expressions in the terms of the language being able to define how types and instances should transform. These expressions shown that any universal computational could be reduced to an expression purely in terms of types and combinators. The crucial idea of these expressions is the computational language, which delivers abstractions, where anyone can declare things and then reuse them without having to think about how they're built inside.  

In the end, combinators are fundamentally computational constructs. It is surprising, just how simple the combinator systems can be. Combinators make a bridge to the way humans think, allowing anyone to represent anything using structured symbolic expressions. 


## References

The post is based on the following articles.  

1. [Featherweight Go](https://arxiv.org/pdf/2005.11710.pdf)
2. [Functional design: combinators](https://dev.to/gcanti/functional-design-combinators-14pn)
3. [Combinators and the Story of Computation](https://writings.stephenwolfram.com/2020/12/combinators-and-the-story-of-computation/)
4. [Parametric polymorphism, Records, and Subtyping](https://groups.seas.harvard.edu/courses/cs152/2015sp/lectures/lec14-polymorphism.pdf)
5. [Type classes](https://typelevel.org/cats/typeclasses.html)
6. [Implementing, and Understanding Type Classes](https://okmij.org/ftp/Computation/typeclass.html)
6. [Higher-order functions vs interfaces in golang](http://aquaraga.github.io/functional-programming/golang/2016/11/19/golang-interfaces-vs-functions.html)
7. [Inheritance vs Generics vs TypeClasses in Scala](https://dev.to/jmcclell/inheritance-vs-generics-vs-typeclasses-in-scala-20op)
8. [Common combinators in JavaScript](https://gist.github.com/Avaq/1f0636ec5c8d6aed2e45)
9. [Haskell's TypeClasses and Go's Interfaces](https://stackoverflow.com/questions/2982012/haskells-typeclasses-and-gos-interfaces)

