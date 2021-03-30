package seq_test

import (
	"testing"

	"github.com/fogfish/golem"
	"github.com/fogfish/golem/seq"
	"github.com/fogfish/it"
)

func TestListCons(t *testing.T) {
	x := golem.String("x")
	z := seq.NewList()
	a := z.Cons(x)

	it.Ok(t).
		IfNil(z).
		IfNotNil(a).
		IfTrue(a.Head().Eq(x))
}

func TestListTail(t *testing.T) {
	x := golem.String("x")
	z := seq.NewList()
	a := z.Cons(x)
	b := a.Cons(x)
	c := b.Cons(x)
	d := c.Cons(x)

	it.Ok(t).
		IfNil(z).
		If(a.Tail()).Equal(z).
		If(b.Tail()).Equal(a).
		If(c.Tail()).Equal(b).
		If(d.Tail()).Equal(c)
}

//
//
type El struct{ ID int }

func (El) Type() {}

func (i El) Eq(x golem.Eq) bool {
	switch v := x.(type) {
	case El:
		return i.ID == v.ID
	default:
		return false
	}
}

func (i El) Lt(x golem.Ord) bool {
	switch v := x.(type) {
	case El:
		return i.ID < v.ID
	default:
		return false
	}
}

//
//
var (
	defCap  int       = 1000000
	defList *seq.List = seq.NewList()

	el *El
)

func init() {
	for n := 0; n < defCap; n++ {
		defList = defList.Cons(&El{ID: n})
	}
}

func BenchmarkListCons(b *testing.B) {
	b.ReportAllocs()

	l := seq.NewList()
	for n := 0; n < b.N; n++ {
		l = l.Cons(&El{ID: n})
	}
}

func BenchmarkListTail(b *testing.B) {
	b.ReportAllocs()

	l := defList

	for n := 0; n < b.N; n++ {
		switch v := l.Head().(type) {
		case *El:
			el = v
		}

		if l = l.Tail(); l == nil {
			l = defList
		}
	}
}
