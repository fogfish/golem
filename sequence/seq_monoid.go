package seq

import (
	"github.com/fogfish/golem/generic"
	"github.com/fogfish/golem/pure"
)

// Mempty returns a type value that hold the identity property for
// combine operation, means the following equalities hold for any choice of x.
//   t.Combine(t.Empty()) == t.Empty().Combine(t) == t
func (seq AnyT) Mempty() pure.Monoid {
	return AnyT{}
}

// Mappend applies a side-effect to the structure by appending a given value.
// Combine must hold associative property
//   a.Combine(b).Combine(c) == a.Combine(b.Combine(c))
func (seq AnyT) Mappend(x pure.Monoid) pure.Monoid {
	switch v := x.(type) {
	case AnyT:
		seq = append(seq, v...)
	}
	return seq
}

// MMap applies high-order function to all elements of sequence.
// It uses monoid interface to build a new sequence, while FMap uses clojure
func (seq AnyT) MMap(m pure.Monoid) func(func(generic.T) pure.Monoid) pure.Monoid {
	return func(mapper func(generic.T) pure.Monoid) pure.Monoid {
		y := m.Mempty()
		for _, x := range seq {
			y = y.Mappend(mapper(x))
		}
		return y
	}
}
