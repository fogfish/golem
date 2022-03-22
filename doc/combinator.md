# Pure Functional Programming with type combinators in Golang

// Dr. Strangelove or: How I Learned to Stop Worrying and Love the Bomb

Humans develop ideas of representing things using formal "system" since antique history - Aristotle’s logic, Euclid’s geometry are good examples. The formalism allows anyone to proof and deduct purely within the system so that it defines a concrete solution of the current problem. System of combinators are known for 100 years since Moses Schönfinkel developed a universal computation system that has been researched since together with mathematical logic, lambda calculus and category theory. Combinators open up an opportunity to depict computation problems in terms of fundamental elements like physics talks about universe in terms of particles. The only definite purpose of combinators are building blocks for composition of "atomic" functions into computational structures from concrete problem "domain". So far, combinators remains as a powerful symbolic expressions in computational languages.  

Combinators are simple and do not involve any advanced math in its definition. A combinator builds a new "things" from previously defined "things".

```
ƒ: Thing ⟼ Thing ⟼ Thing
```

Since, "thing" can be any computational "element" including functions and other combinators. It delivers powerful **combinator patterns** for functional programming - a style of declaring a small set of primitive abstractions and collection of combinators to defined advanced structures. Golang like any other languages supports first class functions. It allows functions to be assigned to variables, passed as arguments to other functions and returned from other functions. Therefore, combinators of pure functions is a given fact for any Golang application. 

Let's advance this patterns towards Golang type system and define combinators over types and their instances to derive complex structures of type T. There are 7 patterns to consider and express their semantic with Golang:

1. [Type trait](#type-trait-pattern) (`𝔗 ⟼ A ⟼ ƒ(𝔗[A], A)`) declares type class, it's laws and the intent.
2. [Sub-typing](#sub-typing-pattern) (`A <: B`) enhances existing type classes.
3. [Lifting](#lifting) (`ƒ ⟼ 𝔗[A]`) transforms a pure function into corresponding type class.
4. [Homogenous product](#homogenous-product) (`A × B × … ⟼ 𝔗 ⟼ 𝔗[A × B × …]`) composes type classes of same kind to operate with product type.
5. [Contra Variant Functor](#contra-variant-functor) (`(ƒ: b ⟼ a) ⟼ 𝔗[A] ⟼ 𝔗[B]`) type transformation using pure function.
6. [Compose Generic Types](#compose-generic-types) (`𝔗 ⟼ 𝕬 ⟼ 𝔗[𝕬]`) to define generic computation.
7. [Heterogenous product](#heterogenous-product) (`𝔗 × 𝕬 × … ⟼ 𝕷`) compose heterogenous type classes into a new type law pattern. 
 
The implementation of each combinator is considered further in this post using simplest examples in the style of Golang.  


## Type trait pattern

Equality (==) implements an equivalence relationship where two values comparing "is equal to" if they belong to the same equivalence class in their domain (e.g. equivalence laws of boolean, numbers, strings and other types). Typical comparison operators supports only built-in types, extension of the operator over values of some abstract type T requires a definition of equality function over these types and acting differently for each type. 

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

The **type trait** definition looks like an interface that specifies an equality behavior for some kind of type. Anyone with OOP background might be confused. The OOP style insists on `Equal(T) bool` clamming receiver instance is compared with given value. Object-oriented languages leverage subtype polymorphism, which is not efficient to build advanced combinators. Functional programing is looking towards ad-hoc polymorphism - depict a trait of different unrelated types with type specific implementations. In this example, the interface defines "atomic" trait that knows how to compare type instances (e.g. number equality category, string equality category, etc). The equivalence relationship implementation for Golang `int` type is the following

```go
package eq

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
```

This technique supports an ad-hoc polymorphism of the trait `Eq`, and detaches the data type implementation from the behavior. The application is able to implement multiple traits for same data type together with trait sub-typing, which is not achievable if `Eq` is used as standard Golang interface.     

The type trait pattern solves the problem of polymorphic algorithm implementations - using the "well-known" function names for various instances that takes different kinds of parameters. For example, equality type trait helps with implementation of "haystack" algorithms: 

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
3. The instance of type trait T is declared for any concrete type by implementing declared interface for T. Golang's [structural type system](https://en.wikipedia.org/wiki/Structural_type_system) empowers different unrelated types with type specific implementations. It makes the approach flexible for pure functional combinator libraries and easy ad-hoc type extensions.

The type trait pattern looks similar to type class. The computer science has defined "the type class is construct that supports ad hoc polymorphism. This is achieved by adding constraints to type variables in parametrically polymorphic types. Such a constraint typically involves a type class T and a type variable a, and means that a can only be instantiated to a type whose members support the overloaded operations associated with T.". However, Golang type system is  strictly less powerful than type classes, they are "a kind of zeroth-order type class", while concepts around type classes are often associated with higher-kinded polymorphism. The type trait abstraction as it is defined here provides better composability that Golang interface but less powerful in comparison with Haskell's type classes.    


## Sub-typing pattern

A **Sub-typing** pattern (or inclusion polymorphism) creates a new type trait from existing one. Sub-typing is a classical polymorphism used in object-oriented programming, sub-typing is roughly comparable with Golang embedding. The purpose of sub-typing is enhance the interface of existing trait with declaration of new behavior. For example the class of totally ordered types `Ord` is a sub-type of `Eq`. 

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

The implementation of `Ord` type instances does not differs from `Eq`, each type "implements" the behavior: 

```go
func (ordInt) Compare(a, b int) Ordering { /* ... */ }
```

TBD: subtyping of existing trait () valid sub-typing when new implementation re-uses previously defined implementation 

```go
func (ordInt) Equal(a, b, int) bool { return eq.Int.Equal(a, b) }
```


Sub-typing is a composition of "atomic" type laws into complex one. It also defines the notion of substitutability and allows the generic computation, which is written to operate elements of type T, can also operate on instances of sub-types.


## Lifting

Function and types are first class citizen in Golang. The **lifting** pattern transforms a pure function into corresponding type class - the combinator constructs a class from the function. Let's consider the example, there is equality function `equal: T ⟼ T ⟼ bool` and type class `Eq`. The combinator `FromEq` create an instance of `Eq` for the function so that

```
FromEq: (T ⟼ T ⟼ bool) -> Eq
```

Golang implementation of this combinator requires definition of type `FromEq` and implementation of corresponding type class `Eq`:

```go
/*

FromEq is a combinator that lifts T ⟼ T ⟼ bool function to Eq type class
*/
type FromEq[T any] func(T, T) bool

// implementation of Eq type class 
func (f FromEq[T]) Equal(a, b T) bool { return f(a, b)}
```

The usage of lifting pattern is straight forward:

The combinator lifts any function `T ⟼ T ⟼ bool` 

```go
// Equal is a pure function that compares two integers
func Equal(a, b int) bool { return a == b }

/*

The combinator creates a new instance of Eq data type
that "knows" how to compare int.
*/
var Int Eq[int] = eq.FromEq[int](Equal)
```

The lifting pattern is very powerful one. It is not only leverage the gap between functional and type class domains but also facilitates the composable and re-usable definition of type classes using closures and other pure functional concepts.     


## Homogenous product

Type theory and a functional programming operates with algebraic data types. They are known as a composition of other types. The theory defines two classes of compositions: product types and co-product types. Product types are strongly expressed by structs in Golang; co-products are loosely defined (let's skip them at current considerations). It is not always practical to implement type law pattern for each instance of product type. Construction of the complex type law through the composition of existing instances is an alternative solution.  

The **homogenous product** pattern allows an application to construct type laws for product types in a relatively boilerplate-free way. It composes type classes of same kind to operate on containers. Let's consider a product type `T: A × B × ...` together with set of type laws for each of elementary type `Eq[A], Eq[B], ...`. The homogenous product build `Eq[T]: Eq[A] × Eq[B] × ...`. 

```go
// ExampleType product type is product of primitive types int × string
type ExampleType struct {
	A int
	B string
}

// The type law is defined for these primitive types
var (
	Int    Eq[int]    = FromEq[int](equal[int])
	String Eq[string] = FromEq[string](equal[string])
)
```

Golang's implementation of this combinator requires definition of type `UnApplyN` to "extract" fractions on the product and implementation of corresponding product type `ProductN`for `Eq` type law:

```go
// UnApply2 is contra-map function for data type T that unwrap product type  
type UnApply2[T, A, B any] func(T) (A, B)

/*

ProductEq2 is a shortcut due to zeroth-order type class concept in Golang and lack of heterogenous lists. The type just a product "container" of Eq instances

The new instance of Eq for data type ExampleType is created using
eq := ProductEq2[ExampleType, int, string]{Int, String,
  func(x ExampleType) (int, string) { return x.A, x.B },
}
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

The homogenous product pattern is a building blocks for composition of "atomic" type classes into complex structures from concrete problem "domain".


## Contra Variant Functor

The functor pattern is one of mostly discussed patterns in functional programming. It allows "a generic type to apply a function inside without changing the structure of the generic type". In math, there are many concepts that acts as functors. The **contra variant** pattern just "turn morphisms around". 

A functional programming operates with algebraic data types. It is not always practical to implement type law pattern for each instance of data type. The type mapping is a solution to transform data types so that existing type law instances can be re-used in different context.  

Let's consider two types `A` and `B` and the instance of type class `Eq[A]`. The contra variant functor builds an instance of `Eq[B]` with help of `f: b ⟼ a` transformer. 

```go
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
```

Use the combinator to make an instance of `Eq[ExampleType]`

```go
ContraMapEq[int, ExampleType]{Int}.FMap(
  func(x ExampleType) int { /* ... */ },
)
```

## Compose generic types 

So far, the article has considered a few combinator pattern that builds a new behavior by composing generic types. The type class `𝔗` is composed with type class `𝕬`, which parametrizes `𝔗` building `𝔗 ⟼ 𝕬 ⟼ 𝔗[𝕬]` construct. Golang's type system is less powerful than than Haskell type classes, it uses "a kind of zeroth-order type class", higher kinded polymorphism is not well supported. Therefore, an alternative approach is proposed.

Let's consider a `Foldable` abstraction that represents data structures that can be reduced to a summary value one element at a time:

```go
type Foldable[T any] interface {
	Fold(a T, seq []T) (x T)
}
```

There are infinite possibilities to implement the `Foldable` type class due to unbounded definition of the "reduce to summary" function in the specification. There is an algebraic structure, called `Semigroup`, consisting of a type set together with an associative binary operation.

```go
type Semigroup[T any] interface {
	Combine(T, T) T
}
```

The composition of generic types `Foldable` over `Semigroup` allow to implement a single variant of data structure reduction algorithm. The `Semigroup` parameter allow to inject the associative binary operation into the algorithms.

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

The approach discussed by the compose generic types pattern is only the solution to parametrize one type class over another one until higher kinded polymorphism is fully supported at Golang.


## Heterogenous product

The **heterogenous product** pattern allows an application to compose type laws together building a new one, which "inherits" properties of composed types. Let's evaluate the purpose of the pattern in the context of following types `Seq` and `Eq`:

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

The definition of equality for sequence requires the definition of traversal laws. However, the application do not want to re-implement equality for various classes of the sequence. Instead, it aims for re-usability building the new equality laws from existing constructs `Seq × Eq ⟼ SeqEq`. 

```go
/*

SeqEq is heterogenous product of Seq and Eq laws.
It composes two types together that "knows" how to compare sequences.
*/
type SeqEq[S, T any] struct {
  Seq[S, T]
  Eq[T]
}

// implements equality rule for sequence using Seq & Eq type classes.
func (seq SeqEq[S, T]) Equal(a, b S) bool { /* ... */ }
``` 


## Afterwords




~ * * * ~

We have defined a basic principles of composing "classes" into high order constructs



takes a function 

Equality "law" can be defined 

Usage of interfaces is a valid approach to declare a
High-order function is 




## Absence of type-classes



* Type Classes vs Sub Typing

* Type Classes vs HoF

A pure functional design 

Use combinators for pure functional design in Golang

Functional Design in Go using combinators

https://dev.to/gcanti/functional-design-combinators-14pn

https://typelevel.org/cats/typeclasses.html
Type classes are a powerful tool used in functional programming to enable ad-hoc polymorphism, more commonly known as overloading. Where many object-oriented languages leverage subtyping for polymorphic code, functional programming tends towards a combination of parametric polymorphism (think type parameters, like Java generics) and ad-hoc polymorphism.

https://en.wikipedia.org/wiki/Type_class

Basic Logic (It was easy to represent mathematical relations) represent basically anything in terms of symbolic expressions, and transformation rules on them
* Idempotent law
* Commutative and associative law
* Distributive law
* Identity law
* Compliment law
* DeMorgan's law
* Reflective law
* Antisymmetric law
* Transitive law
^^^ testing

> In computer science, functional programming is a programming paradigm where programs are constructed by applying and composing functions. It is a declarative programming paradigm in which function definitions are trees of expressions that map values to other values, rather than a sequence of imperative statements which update the running state of the program.

> Functional Programming is a declarative style of development that uses side-effect free functions to express solution of the problem domain. The core concepts of functional programming are elaborated by Martin Odersky - First class and high-order functions and immutability. Another key feature in functional programming is the composition - a style of development to build a new things from small reusable elements. Functional code looks great only if functions clearly describe your problem. Usually lines of code per function is only a single metric that reflects quality of the code

> Composition is the essence of programming...
> The composition is a style of development to build a new things from small reusable elements.


https://writings.stephenwolfram.com/2020/12/combinators-and-the-story-of-computation/

The Abstract Representation of Things



that it’s “in the spirit of” Hilbert’s axiomatic method to build everything from as few notions as possible; then he says that what he wants to do is to “seek out those notions from which we shall best be able to construct all other notions of the branch of science in question”.


But, OK, so what had he achieved? He’d basically shown that any expression that might appear in predicate logic (with logical connectives, quantifiers, variables, etc.) could be reduced to an expression purely in terms of the combinators S, C (now K) and U.



I’m sure Schönfinkel was extremely surprised. And here I personally feel a certain commonality with him. Because in my own explorations of the computational universe, what I’ve found over and over again is that it takes only remarkably simple systems to be capable of highly complex behavior—and of universal computation. And even after exploring the computational universe for four decades, I’m still continually surprised at just how simple the systems can be.

In the end, combinators are fundamentally computational constructs

But by and large the important ideas that first arose with combinators ended up being absorbed into practical computing by quite circuitous routes, without direct reference to their origins, or to the specific structure of combinators.



Yes, it was quite often useful to think of “applying functions to things” (and SMP had its version of lambda, for example), but it was much more powerful to think about symbolic expressions as just “being there” (“x doesn’t have to have a value”)—like things in the world—with the language being able to define how things should transform.

Over time I’ve come to realize that doing this is less about what one can in principle use to construct computations, and more about making a bridge to the way humans think about things. It’s crucial that there’s an underlying structure—symbolic expressions—that can represent anything.

But that would be like having a world without nouns—a world where there’s no name for anything—and the representation of everything has to be built from scratch. But the crucial idea that’s central to human language—and now to computational language—is to be able to have layers of abstraction, where one can name things and then refer to them just by name without having to think about how they’re built up “inside”.



https://stackoverflow.com/questions/2982012/haskells-typeclasses-and-gos-interfaces


http://aquaraga.github.io/functional-programming/golang/2016/11/19/golang-interfaces-vs-functions.html

Featherweight Go
https://arxiv.org/pdf/2005.11710.pdf

Implementing, and Understanding Type Classes
http://okmij.org/ftp/Computation/typeclass.html

Inheritance vs Generics vs TypeClasses in Scala
https://dev.to/jmcclell/inheritance-vs-generics-vs-typeclasses-in-scala-20op

https://writings.stephenwolfram.com/2020/12/combinators-and-the-story-of-computation/
https://gist.github.com/Avaq/1f0636ec5c8d6aed2e45
http://okmij.org/ftp/Computation/types.html




Lecture: Parametric polymorphism, Records, and Subtyping
https://groups.seas.harvard.edu/courses/cs152/2015sp/lectures/lec14-polymorphism.pdf
