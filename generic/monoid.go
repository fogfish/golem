package generic

// Monoid supports the generic implementation of strcutural transformations.
// It is an algebraic structure https://en.wikipedia.org/wiki/Monoid
// with identity element and associative binary function.
//
// Monoid facilitates the implementation of following transformations in
// the absence of generics and type covariance.
//
//    def map[B](f: (A) => B): Seq[B]
//
// The library uses Monoid interface to implement map, fold, filter,
// comprehension and others for complex data structures. Here is an
// example
//
//
//   type MSeq struct { value []int }
//
//   func (seq *MSeq) Empty() generic.Monoid {
//     return &MSeq{}
//   }
//
//   func (seq *MSeq) Combine(x interface{}) generic.Monoid {
//     seq.value = append(seq.value, x.(int))
//     return seq
//   }
//
type Monoid interface {
	// Empty returns a type value that hold the identity property for
	// combine operation, means the following equalities hold for any choice of x.
	//   t.Combine(t.Empty()) == t.Empty().Combine(t) == t
	Empty() Monoid

	// Combine applies a side-effect to the structure by appending a given value.
	// Combine must hold associative property
	//   a.Combine(b).Combine(c) == a.Combine(b.Combine(c))
	Combine(x interface{}) Monoid
}
