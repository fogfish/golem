package seq_test

import (
	"fmt"

	"github.com/fogfish/golem/seq"
)

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
