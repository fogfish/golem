# Heterogenous Sequence of Types (`hseq`)

The module unfolds product type (e.g. structs) into sequence of types, while preserving the original type witness. 

```
type T struct { /* ... */ } ⇒ hseq.Seq[T]
```

It is extension over reflection `reflect.StructField` but with compile and runtime type safeness even product types have equivalent set of members

```go
type A struct { string }
type B struct { string }

// hseq.Seq[A] != hseq.Seq[B]
```

The module is an enabler for building generic algorithms over equivalent types avoiding repetition. For example,
* Building an optics abstraction requires definition of `Lens[T, A]`. The `hseq` module simplifies the type safe implementation allowing to lift any type `A` from the product `T`.
* Building Domain Specific Language might require exposure of type safe utilities over subset of product attributes.
* Implement genetic and type safe witness of type `T` (e.g. pair T, A). 

## Getting Started

```go
// hseq.Seq[T] is equivalent presentation of type T
// both represents product type T
type T struct { /* ... */ }

// make instance of heterogenous sequence
var seq = hseq.New[T]()
```

The heterogenous sequence is an ordered list of types derived from the product type. Absence of macros in Golang, does not allow us to make a compile type definition of heterogenous sequence. It is a runtime reflection but annotated with original type T, which makes it compile type safe. The sequence consists of `hseq.Type[T]`, which uniquely identify typed member of original product type. The module provides two helper functions `ForType/ForName`, which lookup heterogenous sequence. 

```go
type T struct {
  A A
  B B
}

var (
  // lookup sequence using type hint A
  a = hseq.ForType[T, A](seq)
  // lookup sequence using name hint "B"
  b = hseq.ForName[T, B](seq, "B")
)
```

So far, examples above made nothing but unfold (decompose) the product type on individual components. Then, the module implements `hseq.FMap1 ... hseq.FMap9` functions to construct type safe representation of the product type `A × B × ... × I` where the type safe witness is casted per element forming pair `T, A`. 

```go
var a, b = hseq.Map2[T, A, B](seq, /* ... */)
```

## Quick Example

The most simplest example that shows the applicability of `hseq` abstraction is building generic getter - the function that focuses on struct attribute.

```go
// Getter type just defines a generic function: T ⇒ A  
type Getter[T, A any] struct{ hseq.Type[T] }

func (field Getter[T, A]) Value(s T) A {
	f := reflect.ValueOf(s).FieldByName(field.Name)
	return f.Interface().(A)
}

// Creates type safe instance of Getter 
func NewGetter[T, A any](s hseq.Type[T]) Getter[T, A] {
	hseq.AssertStrict[T, A](s)
	return Getter[T, A]{s}
}

// For given product type Someone
type Someone struct {
  Name string
  Age  int
}

// Getters are defined. 
var name, age = hseq.FMap2(
	hseq.New[Someone](),
	newGetter[Someone, string],
  newGetter[Someone, int]
)

// While Getter is generic but instances are type safe.
// Compiler fails if name is used with other type than Someone and string
name.Value(/* ... */)
```
