//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package pure

/*

String type constrain
  type MyType[T pure.String] ...
*/
type String interface {
	string
}

/*

AnyString type constrain
  type MyType[T pure.AnyString] ...
*/
type AnyString interface {
	~string
}

/*

Integer type constrain
  type MyType[T pure.Integer] ...
*/
type Integer interface {
	int | int8 | int16 | int32 | int64
}

/*

AnyInteger type constrain
  type MyType[T pure.AnyInteger] ...
*/
type AnyInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

/*

UInteger type constrain
  type MyType[T pure.UInteger] ...
*/
type UInteger interface {
	uint | uint8 | uint16 | uint32 | uint64
}

/*

AnyUInteger type constrain
  type MyType[T pure.AnyUInteger] ...
*/
type AnyUInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

/*

Float type constrain
  type MyType[T pure.Float] ...
*/
type Float interface {
	float32 | float64
}

/*

AnyFloat type constrain
  type MyType[T pure.AnyFloat] ...
*/
type AnyFloat interface {
	~float32 | ~float64
}

/*

Bool type constrain
  type MyType[T pure.Bool] ...
*/
type Bool interface {
	bool
}

/*

AnyBool type constrain
  type MyType[T pure.AnyBool] ...
*/
type AnyBool interface {
	~bool
}

/*

Number type constrain
  type MyType[T pure.Number] ...
*/
type Number interface {
	Integer | UInteger | Float
}

/*

AnyNumber type constrain
  type MyType[T pure.AnyNumber] ...
*/
type AnyNumber interface {
	AnyInteger | AnyUInteger | AnyFloat
}

/*

Orderable type constrain supports build-in compare operators
  type MyType[T pure.Orderable] ...
*/
type Orderable interface {
	Number | String
}

/*

AnyOrderable type constrain supports build-in compare operators
  type MyType[T pure.AnyOrderable] ...
*/
type AnyOrderable interface {
	AnyNumber | AnyString
}

/*

ContraMap turn morphisms around f: B ⟼ A
*/
type ContraMap[A, B any] func(B) A
