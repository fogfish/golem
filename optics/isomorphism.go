//
// Copyright (C) 2022 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package optics

import (
	"reflect"
	"strings"

	"github.com/fogfish/golem/hseq"
)

// An isomorphism is a structure-preserving mapping between two structures
// of the same shape that can be reversed by an inverse mapping.
type Isomorphism[A, B any] interface {
	FromAtoB(*A, *B)
	FromBtoA(*B, *A)
}

// Shape of the category S within isomorphism
type Shape[S any] map[string]Value[S]

// Value is an element of the category S
type Value[S any] interface {
	get(reflect.Value) reflect.Value
	set(reflect.Value, reflect.Value)
}

type valueMap[S any] struct{ key reflect.Value }

func (p valueMap[S]) get(kv reflect.Value) reflect.Value {
	return kv.MapIndex(p.key)
}

func (p valueMap[S]) set(kv reflect.Value, v reflect.Value) {
	kv.SetMapIndex(p.key, v)
}

type valueStruct[S any] struct{ key hseq.Type[S] }

func (p valueStruct[S]) get(kv reflect.Value) reflect.Value {
	return kv.FieldByIndex(p.key.Index)
}

func (p valueStruct[S]) set(kv reflect.Value, v reflect.Value) {
	kv.FieldByIndex(p.key.Index).Set(v)
}

// Unfolds category S to shape
func ForShape[S any](attr ...string) Shape[S] {
	s := *new(S)
	v := reflect.ValueOf(s)

	shape := Shape[S]{}

	if v.Kind() == reflect.Struct {
		hseq.FMap(
			hseq.New[S](attr...),
			func(t hseq.Type[S]) any {
				name := strings.Split(t.StructField.Tag.Get("optics"), ",")[0]
				if name == "" {
					name = t.Name
				}
				shape[name] = valueStruct[S]{t}
				return nil
			},
		)
	}

	if v.Kind() == reflect.Map {
		for _, name := range attr {
			shape[name] = valueMap[S]{reflect.ValueOf(name)}
		}
	}

	return shape
}

// Isomorphism from category A to B
func Iso[A, B any](a Shape[A], b Shape[B]) Isomorphism[A, B] {
	iso := isomorphism[A, B]{}
	for key, va := range a {
		if vb, has := b[key]; has {
			iso = append(iso, arrow[A, B]{va, vb})
		}
	}
	return iso
}

type arrow[A, B any] struct {
	a Value[A]
	b Value[B]
}

func (arr arrow[A, B]) fromAtoB(a, b reflect.Value) { arr.b.set(b, arr.a.get(a)) }
func (arr arrow[A, B]) fromBtoA(a, b reflect.Value) { arr.a.set(a, arr.b.get(b)) }

type isomorphism[A, B any] []arrow[A, B]

func (iso isomorphism[A, B]) FromAtoB(a *A, b *B) {
	va := reflect.ValueOf(a).Elem()
	vb := reflect.ValueOf(b).Elem()
	for _, m := range iso {
		m.fromAtoB(va, vb)
	}
}

func (iso isomorphism[A, B]) FromBtoA(b *B, a *A) {
	va := reflect.ValueOf(a).Elem()
	vb := reflect.ValueOf(b).Elem()
	for _, m := range iso {
		m.fromBtoA(va, vb)
	}
}
