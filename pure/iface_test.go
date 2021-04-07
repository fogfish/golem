package pure_test

import (
	"runtime/debug"
	"testing"
)

type A struct {
	Value int
}

type I interface {
	Val() int
}

//go:noinline
func (a A) Val() int {
	return a.Value
}

//go:noinline
func useA(x *A) int {
	return x.Value + x.Value
}

// casting interface to struct is faster then calling the method

//go:noinline
func useI(x I) int {
	// with this 32 ms
	// v := x.Val()
	return x.(*A).Value + x.(*A).Value
	// return x.Val() + x.Val()
}

var r int

func BenchmarkA(b *testing.B) {
	for n := 0; n < b.N; n++ {
		r = useA(&A{n})
	}
}

func BenchmarkI(b *testing.B) {
	debug.SetGCPercent(-1)
	for n := 0; n < b.N; n++ {
		x := I(&A{n})
		r = useI(x)
	}
}
