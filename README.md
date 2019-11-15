# Golem: "Scrap Your Boilerplate" for Go

> His dust was "kneaded into a shapeless husk."

> You could do this with a macro, but...
> the best macro is a macro you don't maintain

**golem** is a pure functional and generic programming for Go. It had its origins in [Purely Functional Data Structures](https://www.cs.cmu.edu/~rwh/theses/okasaki.pdf) by Chris Okasaki, on implementing a various higher rank functional abstractions and patterns, on dealing with [scrap your boilerplate](https://www.microsoft.com/en-us/research/publication/scrap-your-boilerplate-with-class/) and gaining experience from other functional languages primary Scala, Haskell and heavily inspired by Erlang twin library [datum](https://github.com/fogfish/datum). Golem is testing the limits of functional abstractions in Go.


## Inspiration

[Functional Programming](https://en.wikipedia.org/wiki/Functional_programming) is a declarative style of development that uses side-effect free functions to express solution of the problem domain. The [core concepts](http://www.se-radio.net/2007/07/episode-62-martin-odersky-on-scala/) of functional programming are elaborated by Martin Odersky - First class and high-order **functions** and **immutability**. Functional style programming can be achieved in any language, including Go. Golang's [structural type system](https://en.wikipedia.org/wiki/Structural_type_system) helps to reject invalid programs at compilation time. One of the challenge here, Go's structures, arrays, slices and maps embrace mutability rather than restricting it. Scala is a good example of the language that uses imperative runtime but provide data structure implementations that internally prevent mutation. This is a perfect approach to achieve immutability and performance through well-defined scopes. All-in-all, Go is a general purpose language with simple building blocks. This library uses these blocks to implement a functional style of development with the goal of simplicity in mind.


## Key features

* [Generic](doc/generic.md) with build time code generation.
