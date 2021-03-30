package seq_test

import (
	"fmt"
	"testing"

	"github.com/fogfish/golem/seq"
)

type Els seq.List

func NewEls() *Els {
	return (*Els)(seq.NewList())
}

func (els *Els) Cons(x *El) *Els {
	return (*Els)((*seq.List)(els).Cons(x))
}

func (els *Els) Head() *El {
	switch v := (*seq.List)(els).Head().(type) {
	case *El:
		return v
	default:
		panic(fmt.Errorf("Invalid element type %T %v", v, v))
	}
}

func (els *Els) Tail() *Els {
	return (*Els)((*seq.List)(els).Tail())
}

//
//
var (
	defTyped *Els = NewEls()
)

func init() {
	for n := 0; n < defCap; n++ {
		defTyped = defTyped.Cons(&El{ID: n})
	}
}

func BenchmarkTypedCons(b *testing.B) {
	b.ReportAllocs()

	l := NewEls()
	for n := 0; n < b.N; n++ {
		l = l.Cons(&El{ID: n})
	}
}

func BenchmarkTypedTail(b *testing.B) {
	b.ReportAllocs()

	l := defTyped

	for n := 0; n < b.N; n++ {
		el = l.Head()

		if l = l.Tail(); l == nil {
			l = defTyped
		}
	}
}

/*
type Value struct {
	seq.Type
	ID int
}

func isValue(x seq.ISeq) error {
	switch x.(type) {
	case *Value:
		return nil
	}
	return fmt.Errorf("Invalid type")
}

//
//
type ListValue struct{ *seq.List }

func NewListValue() *ListValue {
	return &ListValue{&seq.List{}}
}

func (list ListValue) Cons(x *Value) *ListValue {
	return &ListValue{list.List.Cons(x)}
}

func (list ListValue) FMap(fmap func(Value) error) {
	list.List.FMap(func(x seq.ISeq) error {
		switch t := x.(type) {
		case *Value:
			return fmap(*t)
		}
		return nil
	})
}
*/
