package list_test

import (
	"testing"

	"github.com/fogfish/golem"
	"github.com/fogfish/golem/seq/list"
)

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
	defCap  int        = 1000000
	defList *list.Type = list.New()

	el *El
)

func init() {
	for n := 0; n < defCap; n++ {
		defList = defList.Cons(&El{ID: n})
	}
}

func BenchmarkListCons(b *testing.B) {
	b.ReportAllocs()

	l := list.New()
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
