# Monoids for structural transformations

Golang uses imperative style for structural transformations. The usage of `for` loop is advertised by majority of language tutorials. Usage of `for` loops is an idiomatic replacement for well-known functional constructs map, filter and fold. 

```go
sum := 0
for _, x := listOfInt {
  sum = sum + x
}
```

Usage of `for` loops for anything else than primitive containers requires a lot of boilerplate, which means a lot of repeated code. A functional programming techniques solves this problem with high-order functions. Monoid is most basic concept to apply a structural transformation:

> a monoid is an algebraic structure with a single associative binary operation and an identity element.

Monoid is just a scientific name for mostly used [computer science concept](https://en.wikipedia.org/wiki/Monoid#Examples). The `for` loop example above is "commutative monoid under addition with identity element zero". As an example, MapReduce programming model is a monoid with left folding. Many iterative structural transformations may be elegantly expressed by a monoid operation:
* [map](https://en.wikipedia.org/wiki/Map_(higher-order_function)) - immutable transformation of the structure, preserving the shape but often altering a type.
* [fold](https://en.wikipedia.org/wiki/Fold_(higher-order_function)) - analysis of recursive data structure through use of monoid.
* [filter](https://en.wikipedia.org/wiki/Filter_(higher-order_function)) - produces a new data structure which contains elements accepted by predicate function.
* [comprehension](https://en.wikipedia.org/wiki/List_comprehension) - builder notation as distinct from the use of map and filter functions.

