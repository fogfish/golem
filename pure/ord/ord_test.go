//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package ord_test

import (
	"testing"

	"github.com/fogfish/golem/pure/ord"
	"github.com/fogfish/it"
)

var bbOrd ord.Ordering

func TestCompare(t *testing.T) {
	it.Ok(t).
		IfTrue(ord.Int.Compare(0, 1) == ord.LT).
		IfTrue(ord.Int.Compare(1, 1) == ord.EQ).
		IfTrue(ord.Int.Compare(1, 0) == ord.GT)
}

func BenchmarkOrdInt(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		bbOrd = ord.Int.Compare(n, n+1)
	}
}

type A struct {
	ID   int
	Name string
}

var ordByID = ord.ContraMap[int, A]{
	ord.Int,
	func(x A) int { return x.ID },
}

func TestContraMap(t *testing.T) {
	a := A{1, "a"}
	b := A{0, "a"}

	it.Ok(t).
		IfTrue(ordByID.Compare(b, a) == ord.LT).
		IfTrue(ordByID.Compare(a, a) == ord.EQ).
		IfTrue(ordByID.Compare(a, b) == ord.GT)
}

func BenchmarkContraMap(b *testing.B) {
	b.ReportAllocs()

	a := A{1, "a"}
	for n := 0; n < b.N; n++ {
		bbOrd = ordByID.Compare(a, a)
	}
}
