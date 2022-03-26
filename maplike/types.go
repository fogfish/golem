//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package maplike

/*

MapLike defines a trait for container, which associate keys with values.
*/
type MapLike[K, V any] interface {
	Put(*K, *V) MapLike[K, V]
	Get(*K) *V
	// Remove(K) V
}
