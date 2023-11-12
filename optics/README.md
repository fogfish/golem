# Optics (`optics`)

Lenses resembles concept of getters and setters, which you can compose using functional concepts. In other words, "a lens is a first-class value that combines two operations: viewing (or getting) a subpart of a data structure, and updating (or setting) that part".  

Lenses solves Golang challenge **"abstraction over structure fields"**.

The lens is defined following approaches of Haskell library, and techniques references by [1]:

```go
type Lens[S, A any] interface {
  Get(*S) A
  Put(*S, A) *S
}
```

The module implements Well behaving lenses so that it satisfies three laws:
* **GetPut** If we get focused element `A` from `S` and immediately put `A` with no modifications back into `S`, we must get back exactly `S`.
* **PutGet** If putting `A` inside `S` yields a new `S`, then the `A` obtained from `S` is exactly `A`.
* **PutPut** A sequence of two puts is just the effect of the second, the first is completely overwritten. This law is applicable to every well behaving lenses.

**Ω-lenses** Lens fails if focus is not exists. Ω-lenses (Prism) are capable to recover a create a new container `Ss from nothing. The Ω-lenses are usable for practical application to construct nested data type but they are not well behaving. Ω-lenses (Prism) are not supported yet by the module.


The module unfolds product type (e.g. structs) into sequence of lenses, while preserving the original type witness. 

```
type S struct { /* ... */ } ⇒ optics.Lens[S, A] × ... × Lens[S, X any]
```

The module is an enabler for building generic algorithms over equivalent types avoiding repetition.


## Getting Started

```go
// Given a product type S : A × ... × X
type S struct {
  A A
  ...
  X X
}

// build instance of lenses
var a, x = optics.ForProduct2[S, A, X]()
```

The lenses is type safe getters and setters derived from the product type. Absence of macros in Golang, does not allow us to make a compile type definition of lenses. It is a runtime instance but annotated with original type S, which makes it compile type safe. The lens `optics.Lens[S, A]` uniquely identify typed member of original product type. It usage in other context causes compile time error. The module provides helper function `ForProduct1` ... `ForProduct9` to automatically derive lense from the struct.

```go
type T struct {
  A A
  B B
  C B
}

var (
  // build lense using only type hint A
  a = optics.ForProduct1[T, A]()
  // build lense using type and name hint "B" 
  b = optics.ForProduct1[T, B]("B")
)
```

## Quick Example

The most simplest example that shows the applicability of `optics` abstraction is building generic algorithm that abstracts structure fields

```go
// Declare type and its lenses (getters & setter)
type User struct {
  Name    string
  Updated time.Time
}

var userT = optics.ForProduct1[User, time.Time]()

type City struct {
  Name    string
  Updated time.Time
}

var cityT = optics.ForProduct1[City, time.Time]()

// Generic algorithm that modifies struct fields
func show[T any](updated optics.Lens[T, time.Time], v *T) {
  if t := updated.Get(v); t.IsZero() {
    updated.Put(v, time.Now())
  }

  b, _ := json.MarshalIndent(v, "", "  ")
  fmt.Println(string(b))
}

func main() {
  show(userT, &User{Name: "user"})
  show(cityT, &City{Name: "city"})
}
```

See runnable examples to play with the library
* [basic lense usage](./examples/lenses/main.go)
* [abstract shape of type](./examples/shapes/main.go), see problem statement [Abstract over Golang structure fields using optics](../doc/abstract-over-struct-fields-using-optics.md) for details.

## References

1. [Combinators for Bi-Directional Tree Transformations: A Linguistic Approach to the View Update Problem](https://www.cis.upenn.edu/~bcpierce/papers/lenses-toplas-final.pdf)
