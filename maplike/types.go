//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package maplike

import (
	"github.com/fogfish/golem/pure"
	"github.com/fogfish/golem/pure/ord"
)

/*

Type is opaque type to define polymorphic context of MapLike.
It makes HKT typesafe in the context of map trait
*/
type Type any

/*

Higher-Kinded MapLike type
*/
type Kind[K, V any] pure.HKT2[Type, K, V]

type Construction[F_, K, V any] interface {
	New(ord.Ord[K]) F_
}

type Query[F_, K, V any] interface {
	Length(F_) int
	IsEmpty(F_) bool
	// Get ?
}

//

// Traversable => Head | Tail | IsEmpty

// TODO: Traversable
// type View[F_, K, V any] interface {
// 	Head(F_) (K, V)
// 	Tail(F_) F_
// }

type KeyVal[F_, K, V any] interface {
	Put(F_, K, V) F_
	Get(F_, K) V
	Remove(F_, K) V
}

type Map[F_, K, V any] interface {
	Construction[F_, K, V]
	KeyVal[F_, K, V]
}

/*

MapLike defines a trait for container, which associate keys with values.
*/
type MapLike[K, V any] interface {
	Put(K, V) MapLike[K, V]
	Get(K) V
	Remove(K) V
}
