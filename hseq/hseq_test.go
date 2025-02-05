//
// Copyright (C) 2022 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package hseq_test

import (
	"reflect"
	"testing"

	"github.com/fogfish/golem/hseq"
	"github.com/fogfish/it/v2"
)

//
// ForType[A, T] selector
//

func forType[A any]() func(t *testing.T) {
	type T struct {
		F1 A
		F2 *A
		F3 []A
		F4 *[]A
	}
	seq := hseq.New[T]()

	return func(t *testing.T) {
		for expect, tt := range map[string]hseq.Type[T]{
			"F1": hseq.ForType[A, T](seq),
			"F2": hseq.ForType[*A, T](seq),
			"F3": hseq.ForType[[]A, T](seq),
			"F4": hseq.ForType[*[]A, T](seq),
		} {
			it.Then(t).Should(
				it.Equal(tt.Name, expect),
			)
		}
	}
}

func TestForType(t *testing.T) {
	type S string
	t.Run("String", forType[string]())
	t.Run("StringAlias", forType[S]())
	t.Run("Bool", forType[bool]())
	t.Run("Int8", forType[int8]())
	t.Run("UInt8", forType[uint8]())
	t.Run("Byte", forType[byte]())
	t.Run("Int16", forType[int16]())
	t.Run("UInt16", forType[uint16]())
	t.Run("Int32", forType[int32]())
	t.Run("Rune", forType[rune]())
	t.Run("UInt32", forType[uint32]())
	t.Run("Int64", forType[int64]())
	t.Run("UInt64", forType[uint64]())
	t.Run("Int", forType[int]())
	t.Run("UInt", forType[uint]())
	t.Run("UIntPtr", forType[uintptr]())
	t.Run("Float32", forType[float32]())
	t.Run("Float64", forType[float64]())
	t.Run("Complex64", forType[complex64]())
	t.Run("Complex128", forType[complex128]())

	type Struct struct{ A string }
	t.Run("StructNoName", forType[struct{ A string }]())
	t.Run("Struct", forType[Struct]())

	type Interface interface{ A() string }
	t.Run("InterfaceNoName", forType[interface{ A() string }]())
	t.Run("Interface", forType[Interface]())

	t.Run("Embedded", func(t *testing.T) {
		type S string
		type I int
		type A struct{ S }
		type B struct{ I }
		type T struct {
			A
			*B
		}

		seq := hseq.New[T]()

		for expect, tt := range map[string]hseq.Type[T]{
			"S": hseq.ForType[S, T](seq),
			"I": hseq.ForType[I, T](seq),
			"A": hseq.ForType[A, T](seq),
			"B": hseq.ForType[*B, T](seq),
		} {
			it.Then(t).Should(
				it.Equal(tt.Name, expect),
			)
		}
	})

	t.Run("Unknown", func(t *testing.T) {
		type A string
		type T struct{ A }

		seq := hseq.New[T]()
		it.Then(t).Should(
			it.Fail(
				func() { hseq.ForType[int](seq) },
			).Contain("Critical Error"),
		)
	})
}

//
// ForName[T] selector
//

func forName[A any]() func(t *testing.T) {
	type T struct {
		F1 A
		F2 *A
		F3 []A
		F4 *[]A
	}
	seq := hseq.New[T]()

	return func(t *testing.T) {
		for field, kind := range map[string]reflect.Type{
			"F1": reflect.TypeOf(new(A)).Elem(),
			"F2": reflect.TypeOf(new(*A)).Elem(),
			"F3": reflect.TypeOf(new([]A)).Elem(),
			"F4": reflect.TypeOf(new(*[]A)).Elem(),
		} {
			tt := hseq.ForName[T](seq, field)

			it.Then(t).Should(
				it.Equal(tt.Name, field),
				it.Equal(tt.Type, kind),
			)
		}
	}
}

func TestForName(t *testing.T) {
	type S string
	t.Run("String", forName[string]())
	t.Run("StringAlias", forName[S]())
	t.Run("Bool", forName[bool]())
	t.Run("Int8", forName[int8]())
	t.Run("UInt8", forName[uint8]())
	t.Run("Byte", forName[byte]())
	t.Run("Int16", forName[int16]())
	t.Run("UInt16", forName[uint16]())
	t.Run("Int32", forName[int32]())
	t.Run("Rune", forName[rune]())
	t.Run("UInt32", forName[uint32]())
	t.Run("Int64", forName[int64]())
	t.Run("UInt64", forName[uint64]())
	t.Run("Int", forName[int]())
	t.Run("UInt", forName[uint]())
	t.Run("UIntPtr", forName[uintptr]())
	t.Run("Float32", forName[float32]())
	t.Run("Float64", forName[float64]())
	t.Run("Complex64", forName[complex64]())
	t.Run("Complex128", forName[complex128]())

	type Struct struct{ A string }
	t.Run("StructNoName", forName[struct{ A string }]())
	t.Run("Struct", forName[Struct]())

	type Interface interface{ A() string }
	t.Run("InterfaceNoName", forName[interface{ A() string }]())
	t.Run("Interface", forName[Interface]())

	t.Run("Embedded", func(t *testing.T) {
		type S string
		type I int
		type A struct{ S }
		type B struct{ I }
		type T struct {
			A
			*B
		}

		seq := hseq.New[T]()
		it.Then(t).Should(
			it.Equal(hseq.ForName[T](seq, "S").Name, "S"),
			it.Equal(hseq.ForName[T](seq, "I").Name, "I"),
			it.Equal(hseq.ForName[T](seq, "A").Name, "A"),
			it.Equal(hseq.ForName[T](seq, "B").Name, "B"),
		)
	})

	t.Run("Unknown", func(t *testing.T) {
		type A string
		type T struct{ A }

		seq := hseq.New[T]()
		it.Then(t).Should(
			it.Fail(
				func() { hseq.ForName(seq, "xxx") },
			).Contain("Critical Error"),
		)
	})

}

//
// New
//

type T1 string
type T2 string
type Foo struct{ T1 }
type Poo struct{ T2 }
type Boo interface{ Boo() }
type Bar struct {
	Foo
	*Poo
	F1 string
	F2 []byte
	F3 bool
	F4 uint16
	F5 uint
	F6 int32
	F7 int
	F8 float32
	F9 float64
	FA Boo
}

var FIELDS = []string{"Foo", "T1", "Poo", "T2", "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "FA"}

func TestNew(t *testing.T) {

	t.Run("New[T]()", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New[Bar](),
			func(t hseq.Type[Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal(FIELDS...),
		)
	})

	t.Run("New[*T]()", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New[*Bar](),
			func(t hseq.Type[*Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal(FIELDS...),
		)
	})

	t.Run("New(...)", func(t *testing.T) {
		for _, f := range FIELDS {
			seq := hseq.New[Bar](f)
			it.Then(t).Should(
				it.Equal(len(seq), 1),
				it.Equal(seq[0].Name, f),
			)
		}

		for _, f := range FIELDS[:9] {
			seq := hseq.FMap(
				hseq.New[Bar]("F9", f),
				func(t hseq.Type[Bar]) string { return t.Name },
			)
			it.Then(t).Should(
				it.Seq(seq).Equal("F9", f),
			)
		}
	})

	t.Run("New1", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New1[Bar, string](),
			func(t hseq.Type[Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal("F1"),
		)
	})

	t.Run("New2", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New2[Bar, string, []byte](),
			func(t hseq.Type[Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal("F1", "F2"),
		)
	})

	t.Run("New3", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New3[Bar, string, []byte, bool](),
			func(t hseq.Type[Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal("F1", "F2", "F3"),
		)
	})

	t.Run("New4", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New4[Bar, string, []byte, bool, uint16](),
			func(t hseq.Type[Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal("F1", "F2", "F3", "F4"),
		)
	})

	t.Run("New5", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New5[Bar, string, []byte, bool, uint16, uint](),
			func(t hseq.Type[Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal("F1", "F2", "F3", "F4", "F5"),
		)
	})

	t.Run("New6", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New6[Bar, string, []byte, bool, uint16, uint, int32](),
			func(t hseq.Type[Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal("F1", "F2", "F3", "F4", "F5", "F6"),
		)
	})

	t.Run("New7", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New7[Bar, string, []byte, bool, uint16, uint, int32, int](),
			func(t hseq.Type[Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal("F1", "F2", "F3", "F4", "F5", "F6", "F7"),
		)
	})

	t.Run("New8", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New8[Bar, string, []byte, bool, uint16, uint, int32, int, float32](),
			func(t hseq.Type[Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal("F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8"),
		)
	})

	t.Run("New9", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New9[Bar, string, []byte, bool, uint16, uint, int32, int, float32, float64](),
			func(t hseq.Type[Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal("F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9"),
		)
	})

	t.Run("New1/Interface", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New1[Bar, Boo](),
			func(t hseq.Type[Bar]) string { return t.Name },
		)
		it.Then(t).Should(
			it.Seq(seq).Equal("FA"),
		)
	})
}

type Getter[S, A any] struct{ hseq.Type[S] }

func (id Getter[S, A]) Value(s S) A {
	f := reflect.ValueOf(s).FieldByName(id.Name)
	return f.Interface().(A)
}

func newGetter[S, A any](s hseq.Type[S]) Getter[S, A] {
	hseq.AssertStrict[S, A](s)
	return Getter[S, A]{s}
}

func TestFMap(t *testing.T) {
	t.Run("FMap1", func(t *testing.T) {
		hseq.FMap1(
			hseq.New[Bar]("F1"),
			newGetter[Bar, string],
		)
	})

	t.Run("FMap1/Interface", func(t *testing.T) {
		hseq.FMap1(
			hseq.New[Bar]("FA"),
			newGetter[Bar, Boo],
		)
	})

	t.Run("FMap2", func(t *testing.T) {
		hseq.FMap2(
			hseq.New[Bar]("F1", "F2"),
			newGetter[Bar, string],
			newGetter[Bar, []byte],
		)
	})

	t.Run("FMap2", func(t *testing.T) {
		hseq.FMap2(
			hseq.New[Bar]("F1", "F2"),
			newGetter[Bar, string],
			newGetter[Bar, []byte],
		)
	})

	t.Run("FMap3", func(t *testing.T) {
		hseq.FMap3(
			hseq.New[Bar]("F1", "F2", "F3"),
			newGetter[Bar, string],
			newGetter[Bar, []byte],
			newGetter[Bar, bool],
		)
	})

	t.Run("FMap9", func(t *testing.T) {
		hseq.FMap9(
			hseq.New[Bar]("F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9"),
			newGetter[Bar, string],
			newGetter[Bar, []byte],
			newGetter[Bar, bool],
			newGetter[Bar, uint16],
			newGetter[Bar, uint],
			newGetter[Bar, int32],
			newGetter[Bar, int],
			newGetter[Bar, float32],
			newGetter[Bar, float64],
		)
	})
}
