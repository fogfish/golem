//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package seqtest

import (
	"testing"

	"github.com/fogfish/golem/pure/monoid"
	"github.com/fogfish/golem/seq"
	"github.com/fogfish/it"
)

func TestSeq[F_ seq.Kind[A], A any](t *testing.T, seqT seq.Seq[F_, A], seed F_) {
	nul := seqT.New()
	one := seqT.New(*new(A))
	few := seqT.New(*new(A), *new(A), *new(A))

	t.Run("Seq.HKT", func(t *testing.T) {
		seed.HKT1(seq.Type(nil))
		seed.HKT2(*new(A))
	})

	t.Run("Seq.Length", func(t *testing.T) {
		it.Ok(t).
			If(seqT.Length(nul)).Equal(0).
			If(seqT.Length(one)).Equal(1).
			If(seqT.Length(few)).Equal(3)
	})

	t.Run("Seq.IsEmpty", func(t *testing.T) {
		it.Ok(t).
			If(seqT.IsEmpty(nul)).Equal(true).
			If(seqT.IsEmpty(one)).Equal(false).
			If(seqT.IsEmpty(few)).Equal(false)
	})

	t.Run("Seq.Head", func(t *testing.T) {
		it.Ok(t).
			If(seqT.Head(one)).Equal(*new(A)).
			If(seqT.Head(few)).Equal(*new(A))
	})

	t.Run("Seq.Tail", func(t *testing.T) {
		it.Ok(t).
			If(seqT.Length(seqT.Tail(one))).Equal(0).
			If(seqT.Length(seqT.Tail(few))).Equal(2)
	})

	t.Run("Seq.Cons", func(t *testing.T) {
		for _, s := range []F_{nul, one, few} {
			len := seqT.Length(s)
			seqCons := seqT.Cons(seqT.Head(seed), s)

			it.Ok(t).
				If(seqT.Length(seqCons)).Equal(len + 1).
				If(seqT.Head(seqCons)).Equal(seqT.Head(seed))
		}
	})
}

func TestFoldable[F_ any](t *testing.T, seqT seq.Seq[F_, int]) {
	f := seq.Foldable[F_, int]{Seq: seqT}

	t.Run("Foldable.Fold", func(t *testing.T) {
		x := f.Fold(
			monoid.From(0, func(a, b int) int { return a + b }),
			seqT.New(1, 2, 3, 4, 5),
		)

		it.Ok(t).If(x).Equal(15)
	})
}
