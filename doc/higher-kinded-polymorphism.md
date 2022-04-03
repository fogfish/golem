# Why higher-kinded polymorphism is vital functional abstraction and How to implement type classes with Golang 

Most discussions about higher-kinded polymorphism are started with the functor abstraction. Indeed, the functor is a simplest and important higher-kinded abstraction in functional programming, that represents anything with a `map` function. However, the practical usage of this abstraction is seldom  explained.

There are no doubts that any complex application should be organized using divide-and-conquer approach - a style of declaring a small set of primitive abstractions and collection of combinators to define advanced structures. This style of development is not only makes easier to understand solution domain but allows to interchange elements on the go. For example, replacing data structure with improved version without changing entire application. The classical interface abstraction always requires a modification to onboard a modified structure `A'` if it depends upon data structure `A`. 

Let's advance the higher-kinded polymorphism for Golang, explore the pitfalls of the Golang type system and solve real-life examples.


## Why does anyone need higher-kinded polymorphism in Golang?

Needs of higher-kinded polymorphism is easy to illustrate using examples with few data structures to represent the sequence and algorithms over it. Golang interfaces is a logical solution to decouple data structures from algorithms and building the utility library to be used in the application.

Golang interfaces, receivers and generics provide powerful methods to abstract the sequence and any computation over it.

```go
// abstract sequence
type Seq[A any] interface {
  Head() A
  Tail() Seq[A]
}

// concrete implementation of the abstraction using built-in slices 
type SeqSlice[A any] []A

func (seq SeqSlice[A]) Head() A {/* ... */}

// an algorithm uses the abstraction to show the sequence
func show[A any](seq Seq[A]) {/* ... */}
```

It looks as an OOP style that insists on `Seq` behavior to be implemented by the type `SeqSlice`. Object-oriented languages leverage subtype polymorphism, which is not efficient and is not properly supported by Golang. The complication also arises when building ad-hoc polymorphism while extending or tagging existed type (e.g. `type UserIDs SeqSlice[int]`). Often, a proxy boilerplate implementation is required to ensure the compatibility of the new type with the interface so that any existing algorithms are reusable. Anyone might argue that type assertion and casting can be used, which is expensive and do not guarantee type safety at compile time. [The previous post](https://tech.fog.fish/2022/03/27/golang-type-combinator.html#type-trait) has discussed an alternative technique, so called "type trait" pattern to decouple behavior from data type.

The "type trait" pattern supports an ad-hoc polymorphism of `Seq`, and detaches the structural implementation of data type from its behavior. However, the application is not able to implement multiple data types for the same trait if `Seq` is used as a standard Golang interface.

```go
type Seq[A any] interface {
  Head([]A) A
  Tail([]A) A[]
}
```

The interface `Seq` is bound with the implementation of data structure `[]A`, changing the implementation of sequence to something else (e.g. linked list) would require refactoring of this interface and all dependencies. An unary type constructor is required to properly abstract container data structure from computations.  

```go
type Seq[F[_] any, A any] interface {
  Head(F[A]) A
  Tail(F[A]) F[A]
}
```

A simple `Seq` interface can be utilized with several concrete container types like slice, list, linked-list, etc. The aim to build a solution where `Seq` is instantiated without any restrictions to a specific type. The application just defines `Seq` as a parameterized interface that takes a type F as a parameter; F takes another type as a parameter so that F[_] implies a type F of type _ (anything). Unfortunately, Golang does not support higher-kinded types `* ⟼ *`, where a type is abstracted over some type that, in turn, abstracts over another type. Golang's generic support only null-ary types `*`, a solution is required.

The purpose of higher-kinded types is the abstraction that helps to implement
* polymorphic containers is a container that holds any type of items;
* shared libraries uses higher-kinded types to provide the ability to customize exposed interfaces while reducing boilerplate code;
* data morphism involves reading, transforming and writing varieties of data;
* compatibility and harness testing makes a proof that specified protocol is implemented using polymorphic tests;

A higher-kinded type abstraction equips us with a set of operations, which are the only operations applicable to given type of types. The underlying representation can be changed without affecting the rest of computations. It makes programs easier to understand and modify.


## Higher-kinded types first pursuit 

Let's try to resolve the `F[_]` type construct issue with the following shortcut.

```go
type Seq[F_, A any] interface {
  Head(F_) A
  Tail(F_) F_
}
```

This abstraction is enough to implement highly customizable computation using earlier discussed pattern [a composed generic type](https://tech.fog.fish/2022/03/27/golang-type-combinator.html#compose-generic-types). The computation `Show` knows how to iterate sequence and print each element to the console.    

```go
type Show[F_, A any] struct{ Seq[F_, A] }

func (f Show[F_, A]) Print(fa F_) { /* ... f.Seq.Head(x) ... */ }

//
// Then, any one implements the sequence data type
type SeqSlice[A any] []A

type SeqSliceT[A any] string

func (SeqSliceT[A]) Length(seq SeqSlice[A]) int       { return len(seq) }
func (SeqSliceT[A]) Head(seq SeqSlice[A]) A           { return seq[0] }
func (SeqSliceT[A]) Tail(seq SeqSlice[A]) SeqSlice[A] { return seq[1:] }
```

This pattern provides good enough type safety to protect the `Print` computation at compile time.

```go
// compiles successfully
show.Print(SeqSlice[int]{1, 2, 3, 4, 5})

// cannot use SeqSlice[int64]{…} as type SeqSlice[int] in argument to show.
show.Print(SeqSlice[int64]{1, 2, 3, 4, 5})
```

The abstraction is suitable for the abstraction purpose but the nature of interface at Golang makes it possible to bypass compile type checks, which makes this pattern weak to implement pure functional abstractions like Functor, Applicative, Monads and others. The better and type safe solution for higher-kinded polymorphism is required.


## Higher-kinded types with defunctionalization

Since Golang does not support higher-kinded type variables to represent type constructors, the application suffers from abstracting over type expressions of higher kind. Fortunately, the problem has been solved a long time ago by John Reynolds, who has introduced defunctionalization as a technique for translating higher-order programs into a first-order language. His approach has been adopted by Jeremy Yallop and Leo White to OCaml type system. Please read a detailed explanation about the solution from the "Lightweight higher-kinded polymorphism" paper.

The defunctionalization transforms a computation with higher-kinded type expressions into a computation where all type expressions are of kind `*`. The solution is the abstract type constructor `HKT` that represents an idea of parametrized container type `F[A]`.

```go
// HKT[F, A] ∼ F[A]
type HKT[F, A any] interface{}
```

The `HKT` type eliminates higher-kinded type expressions `F[A]`, each container type is bound with a distinct instance of the kind `type SeqKind[A any] HKT[_, A]`. The usage of this approach requires the definition of polymorphic context, any opaque type (called brand) that restricts the type instance to the kind of container. The `HKT` allows abstract computation over generic container types and also annotates a type trait. The following definition implies that trait `Seq` operates with higher-kinded type `F[A]`.

```go
// opaque type to define polymorphic context of Seq
type SeqType interface{}

// `* ⟼ *` type constructor  
type SeqKind[A any] HKT[SeqType, A]

type Seq[F_, A any] interface {
  SeqKind[A]
  // ...
}
```

With `HKT`, it is possible to define polymorphic and composable computation over the kind of container `F[A]` types.

```go
type Unary[T_, F_, A any] func(F_)

func (f Unary[T_, F_, A]) FMap(fa HKT[T_, A]) { f(fa.(F_)) }
```


## Hardening higher-kinded types  

The defined `HKT` type suffers from type safety on the edge cases because the defined abstract type is equivalent to `interface{}`, which stands for any type in Golang. As a result, computation can leak types. Let's consider a Functor abstraction   

```go
// just an abstraction that map f over container F[A], producing container F[B]
type Functor[
  A any, FA SeqKind[A],
  B any, FB SeqKind[B],
] func(f func(A) B, a FA) FB

func (fn Functor[A, FA, B, FB]) FMap(f func(A) B, fa FA) FB { return fn(f, fa) }

/*

With functor type, anyone can build a generic algorithm to convert flat sequence
into a sequence of sequences.
*/
func Unflattening[
	A any,
	FA SeqKind[A],
	FB SeqKind[[]A],
](f Functor[A, FA, []A, FB], fa FA) FB {
	return f.FMap(func(a A) []A { return []A{a} }, fa)
}

/*

The following snippet illustrates, where compiler fails to type check
the results of the algorithm to correct type because SeqKind[_] is
nothing more than interface.
*/
f := Unflattening[int, SeqSlice[int], SeqSlice[[]int]]
var a SeqKind[[]int] = f(SeqSlice[int]{1, 2, 3, 4, 5})
var b SeqKind[int] = f(SeqSlice[int]{1, 2, 3, 4, 5})
```

The original solution proposed by "Lightweight higher-kinded polymorphism" addresses this issue with distinct internal representation of data structure and its external interface. The external interface is built upon `HKT` abstraction and ensures a type safety at compile type for the computation outside of the container. The internal representation is a type safe implementation of the container type. There are two functions injection and projection that connect internal and external types together by doing runtime casting of types.

Fortunately, injection and projection can be avoided in Golang due to interfaces. It is possible to harden `HKT` type into safe construct and benefit from compile type checks. Essentially, the approach follows injections/projections but it is automatically implemented by Golang, while dealing with interfaces.

```go
/*

changing the phantom types F, A into type tags makes it possible for compiler
to distinct HKT instances
*/
type HKT[F, A any] interface {
	HKT1(F)
	HKT2(A)
}

// type instance just tags itself with corresponding parameters
type SeqSlice[A any] []A

func (SeqSlice[A]) HKT1(SeqType) {}
func (SeqSlice[A]) HKT2(A)       {}
```

In contrast with the implementation of lightweight higher-kinded polymorphism in TypScript or OCaml, hardening of `HKT` by avoiding the usage of phantom types is required due to differences in the Golang type system.


## Afterwords

The discussed principles make the Golang type system more powerful. The solution just followed Hilbert’s axiomatic method "to build everything from as few notions as possible". It explains how to build a lightweight higher-kinded polymorphism just using standard Golang notation. 

The primary advantage it the ability to define a computation against higher-order types (containers) which is polymorphic on type A, so that computation is only concern the container type, like it is shown in the following example:

```go
type Show[F_, A any] struct{ Seq[F_, A] }

func (f Show[F_, A]) Print(fa F_) {
	fmt.Printf("==>")

	x := fa
	for f.Seq.Length(x) != 0 {
		fmt.Printf(" %v", f.Seq.Head(x))
		x = f.Seq.Tail(x)
	}

	fmt.Println()
}
```

In the end, higher-kinded polymorphism is a fundamental pure functional concept, which allows anyone to represent advanced abstractions such as Functor, Applicative, Monad and many others.

## References

* [Why doesn't Go have variance in its type system?](https://blog.merovius.de/2018/06/03/why-doesnt-go-have-variance-in.html)
* [Kind (type theory)](https://en.wikipedia.org/wiki/Kind_(type_theory))
* [Lightweight higher-kinded polymorphism](https://www.cl.cam.ac.uk/~jdy22/papers/lightweight-higher-kinded-polymorphism.pdf)
* [Typed functional programming in TypeScript](https://github.com/gcanti/fp-ts)
