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

type T1 string
type T2 string
type Foo struct{ T1 }
type Poo struct{ T2 }
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
}

var FIELDS = []string{"T1", "T2", "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9"}

func TestNew(t *testing.T) {

	t.Run("New()", func(t *testing.T) {
		seq := hseq.FMap(
			hseq.New[Bar](),
			func(t hseq.Type[Bar]) string { return t.Name },
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
