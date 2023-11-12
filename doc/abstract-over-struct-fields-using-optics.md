# Abstract over Golang structure fields using optics

Usage of Golang `struct` is a typical approach to define type safe domain models, write correct and maintainable code. Building a complex domain and computation over it reveals a challenge - "Golang doesn't support abstraction over structure fields". Golang `struct` as type definition is too specific. Often, there is needs to explore similarities between types to avoid repetition in computation. Let's advance the polymorphic approach using Golang type system to solve this real-life problem.

## Why does any anyone need abstraction over structure fields?

**Equivalent shape** is most frequent challenge, when `struct` of different types shares subset of attributes with identical semantic: 

```go
type JPEG struct {
  Bounds image.Rectangle
  // ...
}

type PNG struct {
  Bounds image.Rectangle
  // ...
}
```

Idiomatic Go proposes `interfaces` as a solution (see [example](https://pkg.go.dev/image#Image)). However, this approach requires awkward boilerplate (interface to be implemented by every type). [Effective Go](https://go.dev/doc/effective_go#Getters) claims setters and getters are **not** that idiomatic to Go but nothing wrong to implement it. The naming challenge immediately pop-up if the domain model uses standard codec library (not possible to defined field `Bounds` and method `Bounds()`). Highly likely the final boilerplate for struct would look like:

```go
type JPEG struct {
  bounds image.Rectangle
}

func (jpeg JPEG) Bounds() image.Rectangle { /* ... */ }
func (jpeg *JPEG) SetBounds(r image.Rectangle) { /* ... */ }
func (jpeg JPEG) MarshalJSON() ([]byte, error) {/* ... */}
func (jpeg *JPEG) UnmarshalJSON(b []byte) error {/* ... */}
/* ... */
```

Idiomatic approach requires too much hassles with boilerplate code. Its maintainability over the large domain requires substation effort on implementation and maintainability. Nothing really wrong with this but complexity and number of code lines can be simplified. 

**Embedding** partially solves the challenge. Embedding does not provide the typical, type-driven notion of sub-classing. It just injects pieces of implementation within a struct. Therefore, the computation would narrow down the focus only into sub-type rather whole. It is not possible to carry the reference to whole through embedding abstraction.

```go
type GeoJSON struct { /* ... */ }

type City struct {
  GeoJSON
  Name string `json:"name,omitempty"`
  /* ... */
}

func Send(g *GeoJSON) { /* ... */ }
```

In the example above, the context about city properties is lost within the `Send` function. Boilerplate is required to solve the problem, e.g [custom interface and its implementation on each domain object](https://github.com/fogfish/geojson/tree/main#getting-started). It would be awesome if generics support it but it does not either.

```go
func Send[T struct{ GeoJSON }](g *T) { /* ... */ }

// T (type City) does not satisfy struct{GeoJSON}
Send(&City{/* ... */})
```

Embedding and companion interfaces is the most idiomatic approach to reflect the problem from domain of abstract fields of structure to domain of functions: 

```go
func (g *GeoJSON) DoSomething() { /* ... */ }

type City struct {
  GeoJSON
  /* ... */
}

func Send(g interface{ DoSomething() }) { /* ... */ }
```

This approach simple and intuitive. It only suffers from type safety, complier would not be able to distinguish a correct type fed to `Send`. Any type that implement an interface accepted by `Send`. If application posses multiple instances of `Send`, each configured for particular "final" type, then compiler would not be able to detect errors.  


**Type safeness** assumes definition of new types based on predeclared types. The purpose is explicit semantic that allows to capture errors at compile time. This is a perfect approach to define type-safeness within the domain model but it suffers from similar complexity on the implementation phase. It requires duplicate of code, casting types, etc.  

```go
type Event[T any] struct {
  Type curie.IRI `json:"@type,omitempty"`
  Object T
  // ...
}

type UserCreated Event[User]
type UserRemoved Event[User]
type NoteCreated Event[Note]
type NoteRemoved Event[Note] 
```

Making a short summary. The lack of "abstraction over structure fields" causes us to (a) writing a excessive boilerplate codes (b) giving up on type safety (c) complicating library api and requiring clients to use unnecessary semantics and (d) extensibility due to limitations on extending types outside of the module that declare it.

Typically generic abstraction equips us with operations capable to address these challenges in simple ways. It makes programs easier to understand and modify.

## How to abstract shapes of types?

>
> I know that the spades are the swords of a soldier
> I know that the clubs are weapons of war
> I know that diamonds mean money for this art
> But that's not the shape of my heart
> That's not the shape
> The shape of my heart
>
> Sting: https://www.youtube.com/watch?v=NlwIDxCjL-8
> Songwriters: Dominic James Miller / Gordon Matthew Sumner
>

Usage of Golang types (`struct`) is helpful because they are specific: it glues different pieces for code together, prevents bugs, etc. Types are too specific and Golang does have simple solution to exploit similarities between types and avoid repetition. For example, consider the following domain:

```go
type User struct {
  Name       string
  Followers  int
  Updated    time.Time
  // ...
}

type City struct {
  Name       string
  Population int
  Updated    time.Time
  // ...
}
```

These types abstracts different kinds of data but its share the same shape (contains three fields of same type). Let's assume the generic algorithm is required to update the time stamp to current value and serializing to the wire format. It would require you to write two separate implementation one per type or introduces substation amount of boilerplate. 

The purpose of Generic programming is about overcoming differences like these. New abstraction and functionality to instantiate it is required. Following [the type trait pattern](https://tech.fog.fish/2022/03/27/golang-type-combinator.html#type-trait), the shape abstraction for this example is defined as 

```go
type Shape[T any] interface {
  Put(*T, string, int, time.Time)
  Get(*T) (string, int, time.Time)
}

func Send[T any](shape Shape[T], v *T) { /* ... */ }
```

We are only missing convenient approach of converting specific types into generic ones so that common code can manipulate shape of types. Fortunately, functional programming has defined an abstraction `lens` that helps us to achieve convenient definition of shapes.


## Lenses are getters and setters over structure fields

Lenses resembles concept of getters and setters, which you can compose using functional concepts. In other words, "a lens is a first-class value that combines two operations: viewing (or getting) a subpart of a data structure, and updating (or setting) that part".  

The lens in Go is defined following approaches of Haskell library, and techniques references by [1]:

```go
type Lens[S, A any] interface {
  Get(*S) A
  Put(*S, A) *S
}
```

Well behaving lens satisfies three laws:
* **GetPut** If we get focused element `A` from `S` and immediately put `A` with no modifications back into `S`, we must get back exactly `S`.
* **PutGet** If putting `A` inside `S` yields a new `S`, then the `A` obtained from `S` is exactly `A`.
* **PutPut** A sequence of two puts is just the effect of the second, the first is completely overwritten. This law is applicable to every well behaving lenses.

The module [optics](../optics/) implements an approach to automatically derive type safe getters and setters from the product type. Absence of macros in Golang, does not allow us to make a compile type definition of lenses. It is a runtime instance but annotated with original type S, which makes it compile type safe. The lens `optics.Lens[S, A]` uniquely identify typed member of original product type. It usage in other context causes compile time error. The module provides helper function `ForProduct1` ... `ForProduct9` to automatically derive lense from the struct.

```go
// Given a product type S : A × ... × X
type S struct {
  A A
  ...
  X X
}

// build lenses
var a = optics.ForProduct1[S, A]()
...
var x = optics.ForProduct1[S, X]()
```

The shape of type or its sub-type is defined through lens product:  

```
Lens[S, A any] × Lens[S, B any] × ... × Lens[S, X any]
```

```go
var shape = optics.ForShape2[S, A, X]()
```

In the way, abstraction over structure fields problem is solved. The shape defines a type safe approach of abstracting structs. 

## Afterwords

The discussed principles make the Golang type system more powerful. The primary advantage is the ability to define a computation against higher-order types (containers) which is polymorphic on type `S`, so that computation does not concern the container type:

```go
type Writer[T any] struct {
  optics.Lens3[T, string, int, time.Time]
}

func NewWriter[T any]() Writer[T] {
	return Writer[T]{
		optics.ForShape3[T, string, int, time.Time](),
	}
}

// Writer is truly polymorphic algorithm able to read/write struct of any type
func (w Writer[T]) Write(v *T) {
  if s, i, t := w.Get(v); t.IsZero() {
		w.Put(v, s, i, time.Now())
	}
  // ...
}

var (
  user = NewWriter[User]()
  city = NewWriter[City]()
)
```

In the end, lenses is a fundamental pure functional concept, which allows anyone to solve a challenge - "Golang doesn't support abstraction over structure fields".

## References

1. [Combinators for Bi-Directional Tree Transformations: A Linguistic Approach to the View Update Problem](https://www.cis.upenn.edu/~bcpierce/papers/lenses-toplas-final.pdf)
