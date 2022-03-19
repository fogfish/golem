# Pure Functional Programming with type combinators in Golang

// Dr. Strangelove or: How I Learned to Stop Worrying and Love the Bomb

Humans develop ideas of representing things using formal "system" since antique history - Aristotle’s logic, Euclid’s geometry are good examples. The formalism allows anyone to proof and deduct purely within the system so that it defines a concrete solution of the current problem. System of combinators are known for 100 years since Moses Schönfinkel developed a universal computation system that has been researched since together with mathematical logic, lambda calculus and category theory. Combinators open up an opportunity to depict computation problems in terms of fundamental elements like physics talks about universe in terms of particles. The only definite purpose of combinators are building blocks for composition of "atomic" functions into computational structures from concrete problem "domain". So far, combinators remains as a powerful symbolic expressions in computational languages.  

Combinators are simple and do not involve any advanced math in its definition. A combinator builds a new "things" from previously defined "things".

```
ƒ: Thing ⟼ Thing ⟼ Thing
```

Since, "thing" can be any computational "element" including functions and other combinators. It delivers a powerful **combinator pattern** for functional programming - a style of declaring a small set of primitive abstractions and collection of combinators to defined advanced structures. Golang like any other languages supports first class functions. It allows functions to be assigned to variables, passed as arguments to other functions and returned from other functions. Therefore, combinators of pure functions is a given fact for any Golang application. 

Let's advance this patterns towards Golang type system and define combinators over types and their instances to derive complex structures of type T. There are 7 patterns to consider and express their semantic with Golang:

1.[Type laws pattern](#type-laws-pattern) (`A ⟼ Computation`) declares type class, it's laws and the intent.
2. [Sub-typing](#sub-typing) (`A <: B`) enhances existing type classes.
3. [Lifting](#lifting) (`ƒ ⟼ A`) transforms a pure function into corresponding type class.
4. `<a × b × ...> ⟼ C` homogenous composition of types to operate on containers.
5. `A ⟼ (ƒ: a ⟼ b) ⟼ B` type transformation using functor pattern.
6. `ƒ: A ⟼ B` generic computation.
7. `A × B ⟼ C` heterogenous composition of types to define generic computation. 

The implementation of each combinator is considered further in this post using simplest examples in the style of Golang.  


## Type laws pattern

Equality (==) implements an equivalence relationship where two values comparing "is equal to" if they belong to the same equivalence class in their domain (e.g. equivalence laws of boolean, numbers, strings and other types). Typical comparison operators supports only built-in types, extension of the operator over values of some abstract type T requires a definition of equality function over these types and acting differently for each type. 

```go
/*

Eq : T ⟼ T ⟼ bool
Each type implements own equality, mapping pair of value to bool category
*/
type Eq [T any] interface {
  Equal(T, T) bool
}
```

The **type law** definition looks like an interface that specifies an equality behavior for some kind of type. Anyone with OOP background might be confused. The OOP style insists on `Equal(T) bool` clamming receiver instance is compared with given value. Object-oriented languages leverage subtype polymorphism, which is not efficient to build advanced combinators. Functional programing is looking towards ad-hoc polymorphism - depict a trait of different unrelated types with type specific implementations. In this example, the interface defines "atomic" trait that knows how to compare type instances (e.g. number equality category, string equality category, etc). The equivalence relationship implementation for Golang `int` type is the following

```go
package eq

/*

eqInt declares a new instance of Eq trait, the instance is a new
concrete type, the concrete type "knows" everything about equality in
own domain (e.g. int type). The instance is created over basic string
type so that constant values enumerates all instances of Eq
*/
type eqInt string

// the type "implements" equality behavior
func (eqInt) Equal(a, b int) bool { return a == b }

/*

create a new instance of Eq trait for int domain as immutable value so that 
other functions can use this constant like `eq.Int.Eq(...)`
*/
const Int = eqInt("eq.int")
```

The type law pattern solves the problem of polymorphic algorithm implementations - using the "well-known" function names for various instances that takes different kinds of parameters. For example, equality type law helps with implementation of "haystack" algorithms: 

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

The type law facilitates the combinator pattern with "atomic" and composable element such as
1. A generic computation (an algorithm) over a some abstract type T. This generic computation is polymorphic due to usage of same function names (type law) defined by T. Although type T is polymorphic only to the computation that uses it, where there are multiple instances for each concrete type.
2. The abstract type T declares computational laws and operational intent, the type class in other words. T is declared using a common interface for an arbitrary set of individually specified types.
3. The instance of type T is declared for any concrete type by implementing all functions (laws) for the given abstract type T. Golang structural subtyping empowers different unrelated types with type specific implementations. It makes the approach flexible for pure functional combinator libraries and easy ad-hoc type extensions.

The type law pattern looks similar to type classes. The classical Golang interfaces are strictly less powerful than Haskell type classes, they are "a kind of zeroth-order type class". Only Golang generics would allow reuse of the code via higher kinded polymorphism but let's consider this subject in other post. 


## Sub-typing

**Sub-typing** (or inclusion polymorphism) creates a new type law (type trait) from existing one. Sub-typing is a classical interface embedding in Golang. For example the class of totally ordered types `Ord` is a sub-type of `Eq`. 

```go
/*

Ord : T ⟼ T ⟼ Ordering
Each type implements compare rules, mapping pair of value to
enum{ LT, EQ, GT }
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

The lifting pattern is very powerful one. It is not only leverage the gap between functional and type class domains but also facilitates the composable and re-usable definition of type classes.    


~ * * * ~

**Functor** lifts classes to unary type constructors.

```go
type Functor[A, B any] struct{ Eq[B] }

func (c Functor[A, B]) FMap(f func(A) B) Eq[A] {
	return FromEq[A](func(a, b A) bool {
		return c.Eq.Eq(f(a), f(b))
	})
}
```

https://typelevel.org/cats/typeclasses/functor.html

```go
type MyType{ ID int }

var eq2 Eq[MyType] = Functor[MyType, int]{Int}.FMap(func(a MyType) int { return a.ID })
```

~ * * * ~

**Struct**

...


~ * * * ~

We have defined a basic principles of composing "classes" into high order constructs



takes a function 

Equality "law" can be defined 

Usage of interfaces is a valid approach to declare a
High-order function is 


## Best practice of the "type class" declaration

1. define a data type. data type incorporates basic features.

```go
type MyType struct{}
```

2. type class declares a law, a set of ops about data type
e.g. Seq, Eq, New, 

```go
class Show a where
  show :: String

type Show interface {
  Show() string
}

type MyTypes interface {

}
```

3. instantiate type class for data type

```go
type showMyType string
var (
	_ Show = showMyType("")
	_ Show = (*showMyType)(nil)
)

func (showMyType) Show() string { ... }

type myTypes string

func (myTypes) ...
```

4. instantiate type class

```go
const (
	// NewAccount ...
	ShowAccount = newAccount("type.new.account")
  // Or (if data type is declared somewhere else)
  Account = newAccount("type.new.account")

  // Note: plular form here because it is mostly use
  //       package.MyTypes.Show()
  MyTypes = new showMyTypes() 
)
```

```go

type ShapeID struct{ curie.IRI }

func (id ShapeID) Child() ShapeID {
	return ShapeID{id.Join(guid.Seq.ID())}
}

const ShapeIDs = ShapeIDˆ("type.earth.ShapeID")

type ShapeIDˆ string

func (ShapeIDˆ) New(user string) ShapeID {
	return ShapeID{curie.New("geo:%s", user)}
}

```

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
