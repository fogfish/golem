# Pure Functional Abstractions (`pure`)

The module provides generic purely functional abstractions:

## Type traits

* [hkt](./hkt.go) Higher-Kinded Type
* [eq](./eq/eq.go) is `Eq` (equality) type trait
* [foldable](./foldable/foldable.go) is `Foldable` type trait define rules of folding data structures to a summary value.
* [monoid](./monoid/monoid.go) is `Monoid` type trait defined an algebraic structure consisting of Semigroup and Empty element. 
* [ord](./ord/ord.go) is `Ord` (ordering) type trait
* [semigroup](./semigroup/semigroup.go) is `Semigroup` type trait defined an associative binary operation for a set.
