package eq

import "github.com/fogfish/golem"

//
type Eq interface {
	Eq(golem.T, golem.T) bool
}

//
const (
	Int    = eqInt("eq.int")
	String = eqString("eq.string")
)

//
type eqInt string

func (eqInt) Eq(a, b golem.T) bool {
	return a.(int) == b.(int)
}

//
type eqString string

func (eqString) Eq(a, b golem.T) bool {
	return a.(string) == b.(string)
}

//
type FromEq func(golem.T, golem.T) bool

func (f FromEq) Eq(a, b golem.T) bool {
	return f(a, b)
}

var (
	_ Eq = FromEq(func(golem.T, golem.T) bool { return false })
)

//
type ContraMap struct{ Eq }

func (c ContraMap) From(f func(golem.T) golem.T) Eq {
	return FromEq(func(a, b golem.T) bool {
		return c.Eq.Eq(f(a), f(b))
	})
}

//
type Struct []Eq

func (seq Struct) From2(f func(golem.T) (golem.T, golem.T)) Eq {
	return FromEq(func(a, b golem.T) bool {
		a0, a1 := f(a)
		b0, b1 := f(b)
		return seq[0].Eq(a0, b0) && seq[1].Eq(a1, b1)
	})
}

func (seq Struct) From3(f func(golem.T) (golem.T, golem.T, golem.T)) Eq {
	return FromEq(func(a, b golem.T) bool {
		a0, a1, a2 := f(a)
		b0, b1, b2 := f(b)
		return seq[0].Eq(a0, b0) && seq[1].Eq(a1, b1) && seq[2].Eq(a2, b2)
	})
}
