package seq

import (
	"github.com/fogfish/golem/generic"
	"github.com/fogfish/golem/monoid"
)

// Mempty
func (seq AnyT) Mempty() monoid.AnyT {
	return AnyT{}
}

// Mappend
func (seq AnyT) Mappend(x monoid.AnyT) monoid.AnyT {
	switch v := x.(type) {
	case AnyT:
		seq = append(seq, v...)
	}
	return seq
}

// MMap
func (seq AnyT) MMap(m monoid.AnyT) func(func(generic.T) AnyT) monoid.AnyT {
	return func(mapper func(generic.T) AnyT) monoid.AnyT {
		y := m.Mempty()
		for _, x := range seq {
			y = y.Mappend(mapper(x))
		}
		return y
	}
}
