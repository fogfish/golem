package seq

import "github.com/fogfish/golem/monoid"

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

//
// Functor
//

/*
func (seq AnyT) FMap(m pure.Monoid) func(func(generic.T) AnyT) pure.Monoid {
	return func(mapper func(generic.T) AnyT) pure.Monoid {
		y := m.Mempty()
		for _, x := range seq {
			y = y.Mappend(mapper(x))
		}
		return y
	}
}
*/
