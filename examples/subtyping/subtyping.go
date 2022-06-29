package main

import "fmt"

/*

Eq : T ⟼ T ⟼ bool
Each type implements own equality, mapping pair of value to bool category
*/
type Eq[T any] interface {
	Equal(T, T) bool
}

/*

Ord : T ⟼ T ⟼ Ordering
Each type implements compare rules, mapping pair of value to enum{ LT, EQ, GT }
*/
type Ord[T any] interface {
	Eq[T]

	Compare(T, T) Ordering
}

// Ordering type defines enum{ LT, EQ, GT }
type Ordering int

// enum{ LT, EQ, GT }
const (
	LT Ordering = -1
	EQ Ordering = 0
	GT Ordering = 1
)

type ordInt string

func (ord ordInt) Equal(a, b int) bool {
	return ord.Compare(a, b) == EQ
}

func (ordInt) Compare(a, b int) Ordering {
	switch {
	case a < b:
		return LT
	case a > b:
		return GT
	default:
		return EQ
	}
}

/*

Int create a new instance of Eq trait for int domain as immutable value so that
other functions can use this constant like `eq.Int.Eq(...)`
*/
const Int = ordInt("ord.int")

/*

Searcher is an example of algorithms that uses sub-typing
*/
type Searcher[T any] struct{ Ord[T] }

func (s Searcher[T]) Lookup(e T, seq []T) {
	l := 0
	r := len(seq) - 1

next:
	for {
		if r >= l {
			m := l + (r-l)/2

			switch s.Ord.Compare(seq[m], e) {
			case EQ:
				fmt.Printf("%v found at %v in %v\n", e, m, seq)
				return
			case LT:
				l = m + 1
				continue next
			case GT:
				r = m - 1
				continue next
			}
		}

		fmt.Printf("%v is not found in %v\n", e, seq)
		return
	}
}

func main() {
	searcher := Searcher[int]{Int}
	searcher.Lookup(23, []int{2, 5, 8, 12, 16, 23, 38, 56, 72, 91})
	searcher.Lookup(40, []int{2, 5, 8, 12, 16, 23, 38, 56, 72, 91})
}
