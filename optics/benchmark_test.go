//
// Copyright (C) 2022 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package optics_test

import (
	"testing"

	"github.com/fogfish/golem/optics"
)

type S1 string
type S2 string
type S3 string
type S4 string
type S5 string

type BT1 struct {
	S1
}

type BT5 struct {
	S1
	S2
	S3
	S4
	S5
}

var (
	bt1 = optics.ForProduct1[BT1, S1]()
	bt5 = optics.ForShape5[BT5, S1, S2, S3, S4, S5]()
)

func BenchmarkLensPutForProduct1(mb *testing.B) {
	var val BT1

	bt1.Put(&val, "string")
	if val.S1 != "string" {
		panic("lens failed")
	}

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		bt1.Put(&val, "string")
	}
}

func BenchmarkLensGetForProduct1(mb *testing.B) {
	val := BT1{S1: "string"}

	if bt1.Get(&val) != "string" {
		panic("lens failed")
	}
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		val.S1 = bt1.Get(&val)
	}
}

func BenchmarkShapePutForProduct1(mb *testing.B) {
	var val BT5

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		bt5.Put(&val, "string1", "string2", "string3", "string4", "string5")
	}
}

func BenchmarkShapeGetForProduct5(mb *testing.B) {
	var val = BT5{S1: "string1", S2: "string2", S3: "string3", S4: "string4", S5: "string5"}

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		val.S1, val.S2, val.S3, val.S4, val.S5 = bt5.Get(&val)
	}
}
