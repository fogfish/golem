package generic

// Monoid is a type with an associative binary operation and empty value.
// Monoid is most basic concept to apply a structural transformation and
// variate of algorithms: map, fold, filter, comprehension and etc.
//
//   func (t *MyStruct) Empty() *Monoid {/* ... */}
//   func (t *MyStruct) Combine(x) *Monoid {/* ... */}
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
