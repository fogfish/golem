package stream_test

import (
	"testing"

	"github.com/fogfish/golem/stream"
	"github.com/fogfish/it"
)

func ints(x int) stream.AnyT {
	return stream.NewAnyT(x,
		func() stream.AnyT { return ints(x + 1) },
	)
}

func id(a stream.AnyT) stream.AnyT {
	return stream.NewAnyT(
		a.Head,
		func() stream.AnyT { return id(a.Tail()) },
	)
}

func x10(a stream.AnyT) stream.AnyT {
	return stream.NewAnyT(
		a.Head.(int)*10,
		func() stream.AnyT { return x10(a.Tail()) },
	)
}

func TestStreamNew(t *testing.T) {
	a := stream.NewAnyT(1, nil)
	it.Ok(t).
		If(a.Head).Should().Equal(1).
		If(a.Tail()).Should().Equal(stream.AnyT{})

	b := stream.NewAnyT(2, func() stream.AnyT { return a })
	it.Ok(t).
		If(b.Head).Should().Equal(2).
		If(b.Tail().Head).Should().Equal(1).
		If(b.Tail().Tail()).Should().Equal(stream.AnyT{})
}

func BenchmarkStreamNew(b *testing.B) {
	s := ints(0)
	for n := 0; n < b.N; n++ {
		s = s.Tail()
	}
}

func TestStreamMap(t *testing.T) {
	a := x10(ints(1))

	it.Ok(t).
		If(a.Head).Should().Equal(10).
		If(a.Tail().Head).Should().Equal(20).
		If(a.Tail().Tail().Head).Should().Equal(30).
		If(a.Tail().Tail().Tail().Head).Should().Equal(40)
}

func BenchmarkStreamMap1(b *testing.B) {
	s := id(ints(1))
	for n := 0; n < b.N; n++ {
		s = s.Tail()
	}
}

func BenchmarkStreamMap2(b *testing.B) {
	s := id(id(ints(1)))
	for n := 0; n < b.N; n++ {
		s = s.Tail()
	}
}

func BenchmarkStreamMap3(b *testing.B) {
	s := id(id(id(ints(1))))
	for n := 0; n < b.N; n++ {
		s = s.Tail()
	}
}

func BenchmarkStreamMap4(b *testing.B) {
	s := id(id(id(id(ints(1)))))
	for n := 0; n < b.N; n++ {
		s = s.Tail()
	}
}

func BenchmarkStreamMap5(b *testing.B) {
	s := id(id(id(id(id(ints(1))))))
	for n := 0; n < b.N; n++ {
		s = s.Tail()
	}
}

func BenchmarkStreamMap6(b *testing.B) {
	s := id(id(id(id(id(id(ints(1)))))))
	for n := 0; n < b.N; n++ {
		s = s.Tail()
	}
}

func BenchmarkStreamMap7(b *testing.B) {
	s := id(id(id(id(id(id(id(ints(1))))))))
	for n := 0; n < b.N; n++ {
		s = s.Tail()
	}
}

func BenchmarkStreamMap8(b *testing.B) {
	s := id(id(id(id(id(id(id(id(ints(1)))))))))
	for n := 0; n < b.N; n++ {
		s = s.Tail()
	}
}

func BenchmarkStreamMap9(b *testing.B) {
	s := id(id(id(id(id(id(id(id(id(ints(1))))))))))
	for n := 0; n < b.N; n++ {
		s = s.Tail()
	}
}

func BenchmarkStreamMapA(b *testing.B) {
	s := id(id(id(id(id(id(id(id(id(id(ints(1)))))))))))
	for n := 0; n < b.N; n++ {
		s = s.Tail()
	}
}
