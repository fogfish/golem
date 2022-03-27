//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package eq_test

import (
	"testing"

	"github.com/fogfish/golem/pure/eq"
	"github.com/fogfish/it"
)

var bbEq bool

func TestEqual(t *testing.T) {
	it.Ok(t).
		IfTrue(eq.Int.Equal(1, 1)).
		IfFalse(eq.Int.Equal(1, 2)).
		IfTrue(eq.String.Equal("abcd", "abcd")).
		IfFalse(eq.String.Equal("abcd", "xxxx"))
}

func BenchmarkEqInt(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		bbEq = eq.Int.Equal(n, n+1)
	}
}

type A struct {
	ID   int
	Name string
}

var eqByID = eq.ContraMap[int, A]{
	eq.Int,
	func(x A) int { return x.ID },
}

func TestContraMap(t *testing.T) {
	a := A{1, "a"}
	b := A{2, "a"}

	it.Ok(t).
		IfTrue(eqByID.Equal(a, a)).
		IfFalse(eqByID.Equal(a, b))
}

func BenchmarkContraMap(b *testing.B) {
	b.ReportAllocs()

	a := A{1, "a"}
	for n := 0; n < b.N; n++ {
		bbEq = eqByID.Equal(a, a)
	}
}
