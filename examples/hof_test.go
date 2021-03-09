package main_test

import (
	"testing"
)

type Foo struct {
	a int
}

type Handler interface {
	Serve(x int)
}

type HoF func(x int)

func (foo *Foo) Serve(x int) {
	foo.a = foo.a + x
}

func FooHoF(foo *Foo) HoF {
	return func(x int) {
		foo.a = foo.a + x
	}
}

var result *Foo

func BenchmarkInterface(b *testing.B) {
	f := &Foo{}
	for n := 0; n < b.N; n++ {
		f.Serve(n)
	}
	// fmt.Println(f)
	result = f
}

func BenchmarkHoF(b *testing.B) {
	f := &Foo{}
	hof := FooHoF(f)
	for n := 0; n < b.N; n++ {
		hof(n)
	}
	result = f
}
