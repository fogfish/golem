package list_test

import (
	"testing"

	"github.com/fogfish/golem"
	"github.com/fogfish/golem/seq/list"
	"github.com/fogfish/it"
)

func TestListCons(t *testing.T) {
	x := golem.String("x")
	z := list.New()
	a := z.Cons(x)

	it.Ok(t).
		IfNil(z).
		IfNotNil(a).
		IfTrue(a.Head().Eq(x))
}

func TestListTail(t *testing.T) {
	x := golem.String("x")
	z := list.New()
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

func TestListFMap(t *testing.T) {
	seq := list.New().Cons(golem.Int(1)).Cons(golem.Int(2)).Cons(golem.Int(3))

	acc := golem.Int(0)
	seq.FMap(func(o golem.Ord) { acc = acc + o.(golem.Int) })

	it.Ok(t).If(acc).Equal(golem.Int(6))
}
