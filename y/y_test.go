package y_test

import (
	"fmt"
	"testing"
)

//
type T interface{}

//
type Eq interface {
	Eq(T, T) bool
}

//
type EqInt int

func (EqInt) Eq(a, b T) bool {
	return a.(int) == b.(int)
}

func TestEqInt(t *testing.T) {
	eq := EqInt(0)
	fmt.Println("=> eq int ", eq.Eq(10, 10))
}

// naive elem
func elem(e Eq) func(T, []T) bool {
	return func(a T, seq []T) bool {
		for _, b := range seq {
			if e.Eq(a, b) {
				return true
			}
		}
		return false
	}
}

func TestElem(t *testing.T) {
	v := elem(EqInt(0))(4, []T{1, 2, 3, 4})
	fmt.Println("=> elem ", v)
}

// type elem depends on type equal
type Elem struct{ Eq }

func (e Elem) In(a T, seq []T) bool {
	for _, b := range seq {
		if e.Eq.Eq(a, b) {
			return true
		}
	}
	return false
}

func TestTypedElem(t *testing.T) {
	el := Elem{EqInt(0)}
	v := el.In(4, []T{1, 2, 3, 4})
	fmt.Println("=> elem typed ", v)
}

/*

HowTo define combinators for type

const eqPoint: Eq<Point> = getStructEq({
  x: eqNumber,
  y: eqNumber
})

*/

//
type Ord interface {
	Eq
	Compare(T, T) int
}

//
type OrdInt int

func (OrdInt) Eq(a, b T) bool {
	return a.(int) == b.(int)
}

func (OrdInt) Compare(a, b T) int {
	return a.(int) - b.(int)
}

//
// Use custom types to wrap any function to the type
type FromCompare func(T, T) int

func (f FromCompare) Eq(a, b T) bool {
	return f(a, b) == 0
}

func (f FromCompare) Compare(a, b T) int {
	return f(a, b)
}

//
func TestFromLT(t *testing.T) {
	type A struct{ V int }

	f := func(a T, b T) int { return a.(A).V - b.(A).V }
	ord := Ord(FromCompare(f))

	v := ord.Compare(A{10}, A{20})
	fmt.Println("==> ord ", v)
}

/*

export const contramap: <A, B>(f: (b: B) => A) => (fa: Ord<A>) => Ord<B> = (f) => (fa) =>
  fromCompare((first, second) => fa.compare(f(first), f(second)))

const byAge: Ord<User> = contramap((user: User) => user.age)(ordNumber)
*/
func Contramap(f func(T) T) func(Ord) Ord {
	return func(fa Ord) Ord {
		fx := func(a T, b T) int { return fa.Compare(f(a), f(b)) }
		return FromCompare(fx)
	}
}

func TestOrdStructWithContraMap(t *testing.T) {
	type A struct{ V int }
	ord := Contramap(func(x T) T { return x.(A).V })(OrdInt(0))

	v := ord.Compare(A{10}, A{20})
	fmt.Println("==> ord contra map ", v)
}

// HowTo define Contramap via type
type Contra struct{ Ord }

func (t Contra) FMap(f func(T) T) Ord {
	fx := func(a T, b T) int { return t.Ord.Compare(f(a), f(b)) }
	return FromCompare(fx)
}

func TestOrdStructWithContraTyped(t *testing.T) {
	type A struct{ V int }
	ord := Contra{OrdInt(0)}.FMap(func(x T) T { return x.(A).V })

	v := ord.Compare(A{10}, A{20})
	fmt.Println("==> ord contra map typed ", v)
}
