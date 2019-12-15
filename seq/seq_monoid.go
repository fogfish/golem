package seq

import (
	"github.com/fogfish/golem/generic"
	"github.com/fogfish/golem/monoid"
)

// Mempty returns a type value that hold the identity property for
// combine operation, means the following equalities hold for any choice of x.
//   t.Combine(t.Empty()) == t.Empty().Combine(t) == t
func (seq AnyT) Mempty() monoid.AnyT {
	return AnyT{}
}

// Mappend applies a side-effect to the structure by appending a given value.
// Combine must hold associative property
//   a.Combine(b).Combine(c) == a.Combine(b.Combine(c))
func (seq AnyT) Mappend(x monoid.AnyT) monoid.AnyT {
	switch v := x.(type) {
	case AnyT:
		seq = append(seq, v...)
	}
	return seq
}

// MMap applies high-order function to all elements of sequence.
// It uses monoid interface to build a new sequence, while FMap uses clojure
func (seq AnyT) MMap(m monoid.AnyT) func(func(generic.T) AnyT) monoid.AnyT {
	return func(mapper func(generic.T) AnyT) monoid.AnyT {
		y := m.Mempty()
		for _, x := range seq {
			y = y.Mappend(mapper(x))
		}
		return y
	}
}
