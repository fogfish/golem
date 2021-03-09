package main

import (
	"fmt"
)

type Thing interface {
	UID() string
}

//
type Key struct {
	ID string
}

var (
	_ Thing = Key{}
	_ Thing = (*Key)(nil)
)

func (x Key) UID() string {
	return x.ID
}

type Yek struct {
	ID string
}

//
type Val struct {
	Key
	Yek
	Name string
}

//
type Foo struct {
	Bar string
}

var _ Thing = Foo{}

//
func (x Foo) UID() string {
	return "foobar"
}

func use(x Thing) {
	switch t := x.(type) {
	case *Val:
		fmt.Printf(">>> %T %v\n", t, t)
	default:
		fmt.Printf("+++ %T %v\n", t, t)
	}
}

//
//
type Dict map[string]Thing

func (d Dict) Store(x Thing) {
	d[x.UID()] = x
}

func (d Dict) Fetch(x Key) Thing {
	v := d[x.ID]
	return v
}

func (d Dict) FMap(f func(Thing) error) error {
	for _, v := range d {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}

type ValDict struct{ Dict }

func (v ValDict) Store(x *Val) {
	v.Dict.Store(x)
}

func (v ValDict) Fetch(x Key) *Val {
	switch t := v.Dict.Fetch(x).(type) {
	case *Val:
		return t
	default:
		return nil
	}
}

func (v ValDict) FMap(f func(*Val) error) error {
	return v.Dict.FMap(func(x Thing) error {
		switch t := x.(type) {
		case *Val:
			return f(t)
		default:
			return nil
		}
	})
}

//
//
func main() {
	dict := ValDict{Dict: Dict{}} // map[string]Thing{}

	a := &Val{Key: Key{ID: "a"}, Name: ">> a"}
	use(a)

	dict.Store(a)
	// dict[a.ID] = a
	// fmt.Println(dict["a"])

	// b := dict["a"]
	b := dict.Fetch(a.Key)
	use(b)

	// >>> foo := &Foo{Bar: "xxx"}
	// >>> dict.Store(foo)
	// dict["x"] = foo
	// fmt.Println(dict["x"])

	c := dict.Fetch(Key{ID: "foobar"})
	use(c)

	dict.FMap(func(x *Val) error { fmt.Println(*x); return nil })
}
