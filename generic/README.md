# Generics in Go

[![Documentation](https://godoc.org/github.com/fogfish/golem/cmd/golem?status.svg)](https://godoc.org/github.com/fogfish/golem/cmd/golem)

> Generic programming is a style of computer programming in which algorithms are written in terms of types to-be-specified-later that are then instantiated when needed for specific types provided as parameters. -- said by [Wikipedia](https://en.wikipedia.org/wiki/Generic_programming).

The discussion about Generics in Go is [never ending story](https://github.com/golang/go/wiki/ExperienceReports#generics), may be in release 1.15. There is one article that draws a perfect bottom line on this subject: [Who needs generic? Use ... instead!](https://appliedgo.net/generics/). The primary reason for generics is the ability to [scrap your boilerplate code](https://www.microsoft.com/en-us/research/publication/scrap-your-boilerplate-with-class/). Copy-and-Paste is not an option, it only works perfectly with Stack Overflow. However, code generator is a **perfect** solution that resembles a macros in Scala or parse transforms in Erlang. The implementation of golem command-line utility is inspired by [genny](https://github.com/cheekybits/genny). 

The code generator solution allows to write a valid Go code, which is **compilable** and **testable** with standard Go tools. Generic libraries deals with valid Go code. All generic algorithms are covered with unit testing and benchmarking. You can even use non-parametrized generics algorithms directly in your applications with small performance penalty. The build time replacement of type variables with specific type mitigates all drawback.

The code generator assumes two major development workflows.

## Library

> As a generic library developer I want to define a generic type and supply its parametrized variants of standard Go type so that my generic is ready for application development.

The utility enforces a simple rule - one package defines one generic type.

```
src/
  github.com/fogfish/golem/
    seq/                     # package generic pattern (e.g. sequence)
      doc.go                 # documentation of generic pattern and go:generate
      seq.go                 # generic implementation with parametric type variable 
      seq_test.go            # testing of implementation
      int.go                 # generated code for built-in type int
      string.go              # generated code for built-in type string
      ...
```

## Application

> As a application developer I want to parametrize a generic types with my own application specific types so that the application benefits from re-use of generic implementations

The utility demands - one package defines a type and parametrization of various generic algorithms.

```
src/
  github.com/fogfish/myapp/
    main.go
    foobar/                  # a custom type definition
      foobar.go              # type definition and go:generate
      seq.go                 # parametrization of sequence pattern with foobar type
      stream.go              # parametrization of stream pattern with foobar type
      ...
```

## Package names

These workflows and source code structure refers to widely use Go [package naming](https://golang.org/doc/effective_go.html#package-names). 

```go
// When we are using a generic library then the name of 
// generic pattern is visible as package name
seq.AnyT{}
seq.Int{1, 2, 3, 4}
seq.String{"a", "b", "c", "d"}

// When we are using a custom parametrization then the name of
// generic pattern is visible as type name
foobar.Seq{FooBar{1}, FooBar{2}, FooBar{3}}
foobar.Stream{}
foobar.Stack{}
```

## How It works

Use special labels in your code to define generic types: `AnyT`, `generic.T`

```go
type AnyT struct {
  t generic.T
}

func foo(x AnyT) generic.T
func bar(x generic.T) generic.T
```

The label `generic.T` is replaced with values supplied via corresponding command line parameter `-T`. The labels `AnyT` is either replaced with a name of type `T` or generic pattern depending on the used workflow.  Insert the following comment in your source code file:

```go
// Library workflow
//go:generate golem -lib -T int -generic github.com/fogfish/golem/seq/seq.go

// Application workflow
//go:generate golem -T FooBar -generic github.com/fogfish/golem/seq/seq.go
```
