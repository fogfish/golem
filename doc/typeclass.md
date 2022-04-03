# A Guide To Higher-Kinded Type Classes with Golang

Golang does not support higher-kinded types `* ⟼ *`, where a type is abstracted over some type that, in turn, abstracts over another type. Golang's generic support only null-ary types `*`.

What does it mean in practice? Let's try to draft a functor using [Scala's style](https://typelevel.org/cats/typeclasses/functor.html) as inspiration.

```go
package main

type Functor[F[_] any] interface {
  Map[A, B any](fa F[A], f func(A) B) F[B]
}
```

Golang fails to compile with the following errors

```
3:16: undeclared name _ for array length
4:6: interface method must have no type parameters
4:20: F is not a generic type
4:33: undefined: A
4:36: undefined: B
4:39: F is not a generic type
```

It is not possible to define [unary type constructor](https://en.wikipedia.org/wiki/Kind_(type_theory)) `F[_]` using Golang standard syntax.

The problem has been solved long time ago by John Reynolds, who has introduced defunctionalization as a technique for translating higher-order programs into a first-order language. His approach has been adopted by Jeremy Yallop and Leo White to OCaml type system at [Lightweight higher-kinded polymorphism](https://www.cl.cam.ac.uk/~jdy22/papers/lightweight-higher-kinded-polymorphism.pdf).

## Type classes

The defunctionalization transforms a computation with higher-kinded type expressions into a computation where all type expressions are of kind `*`. The solution is the abstract unary type constructor with new type `HKT` that represents an idea of parametrized container type `F[A]`.

```go
type HKT[F, A any] interface {
	HKT1(F)
	HKT2(A)
}
```

Now, the functor looks like

```go
type Functor[F, A, B any] interface {
  Map(func(A) B, HKT[F, A]) HKT[F, B]
}
```

`Maybe` is the simplest functor anyone can implement. The usage of `HKT` abstraction requires the definition of polymorphic context, any opaque type `MaybeType`, which is also called brand. It restricts the type instance to the kind of container. 

The opaque type helps to define two containers of `HKT[MaybeType, _]` type. The first one `Some` contains a successful value, `None` contains nothing.

```go
type MaybeType any

type MaybeKind[A any] HKT[MaybeType, A]

type Some[A any] struct{ Value A }

func (Some[A]) HKT1(MaybeType) {}
func (Some[A]) HKT2(A)         {}

type None[A any] int

func (None[A]) HKT1(MaybeType) {}
func (None[A]) HKT2(A)         {}
```

Using `Some` and `None` types can be used to define a morphism for `Maybe`. 

```go
func fmap[A, B any](f func(A) B, fa MaybeKind[A]) MaybeKind[B] {
  switch x := fa.(type) {
  case Some[A]:
    return Some[B]{Value: f(x.Value)}
  default:
    return None[B](0)
  }
}
```

Final touch, let's lift the morphism to the `Functor` interface. 

```go
type Maybe[A, B any] func(func(A) B, MaybeKind[A]) MaybeKind[B]

func (eff Maybe[A, B]) Map(
  f func(A) B,
  fa HKT[MaybeType, A],
) HKT[MaybeType, B] {
  return eff(f, fa)
}

var maybe Functor[MaybeType, int, string] = Maybe[int, string](
  fmap[int, string],
)
```

Higher-kinded type class in action.

```go
double := func(x int) string { return strings.Repeat("x", x) }

// outputs: {xxxxx}
fmt.Println(
  maybe.Map(double, Some[int]{5}),
)

// outputs: 0
fmt.Println(
  maybe.Map(double, None[int](0)),
)
```

[Why higher-kinded polymorphism is vital functional abstraction and How to implement type classes with Golang](./higher-kinded-polymorphism.md) discusses this pattern in depth.

A full [source code of Maybe functor](https://gist.github.com/fogfish/6df9d9e0b09c88efed27f05c0c84cf18) is available in [the gist](https://gist.github.com/fogfish/6df9d9e0b09c88efed27f05c0c84cf18).
