//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

/*

Package skiplist implements a probabilistic list-based data structure
that are a simple and efficient substitute for balanced trees.

http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.17.524
*/
package skiplist

import (
	"math"
	"math/rand"
	"time"

	"github.com/fogfish/golem/maplike"
	"github.com/fogfish/golem/pure/ord"
)

/*

tSkipNode ...
*/
type tSkipNode[K, V any] struct {
	key  *K
	val  *V
	next []*tSkipNode[K, V]
}

/*

tSkipList ...
*/
type tSkipList[K, V any] struct {
	ord.Ord[K]
	head   *tSkipNode[K, V]
	length int

	//
	//
	random rand.Source
	levels int
	p      []float64
}

var (
	_ maplike.MapLike[int, int] = &tSkipList[int, int]{}
	_ maplike.MapLike[int, int] = (*tSkipList[int, int])(nil)
)

/*

New creates empty skiplist
*/
func New[K, V any](c ord.Ord[K]) maplike.MapLike[K, V] {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	head := &tSkipNode[K, V]{key: nil, next: []*tSkipNode[K, V]{nil}}
	list := &tSkipList[K, V]{
		Ord:    c,
		head:   head,
		length: 0,
		random: random,
	}

	list.config(4294967296, 1/math.E)
	return list
}

/*

config generates random list of levels
*/
func (list *tSkipList[K, V]) config(n int, p float64) {
	list.levels = int(math.Log10(float64(n)) / math.Log10(1/p))
	list.p = make([]float64, list.levels+1)

	for i := 1; i <= list.levels; i++ {
		list.p[i] = math.Pow(p, float64(i-1))
	}
}

func (list *tSkipList[K, V]) mkNode(key *K, val *V) (int, *tSkipNode[K, V]) {
	// See: https://golang.org/src/math/rand/rand.go#L150
	p := float64(list.random.Int63()) / (1 << 63)

	level := 1
	for level < list.levels && p < list.p[level] {
		level++
	}

	node := &tSkipNode[K, V]{key: key, val: val, next: make([]*tSkipNode[K, V], level)}
	return level, node
}

/*

Skip ...
*/
func (list *tSkipList[K, V]) skip(key *K, node *tSkipNode[K, V]) (*tSkipNode[K, V], []*tSkipNode[K, V]) {
	rank := len(node.next)
	path := make([]*tSkipNode[K, V], rank)

	v := node
	for level := rank - 1; level >= 0; level-- {
		for v.next[level] != nil && (v.next[level].key == nil || list.Ord.Compare(*v.next[level].key, *key) == ord.LT) {
			v = v.next[level]
		}
		path[level] = v
	}

	return v.next[0], path
}

func (list *tSkipList[K, V]) skipr(key *K, node *tSkipNode[K, V]) *tSkipNode[K, V] {
	rank := len(node.next)

	v := node
	for level := rank - 1; level >= 0; level-- {
		for v.next[level] != nil && (v.next[level].key == nil || list.Ord.Compare(*v.next[level].key, *key) == ord.LT) {
			v = v.next[level]
		}
	}

	return v.next[0]
}

/*

Put
*/
func (list *tSkipList[K, V]) Put(key *K, val *V) maplike.MapLike[K, V] {
	v, path := list.skip(key, list.head)

	//
	if v != nil && list.Ord.Compare(*v.key, *key) == ord.EQ {
		v.val = val
		return list
	}

	rank, node := list.mkNode(key, val)

	if rank > len(list.head.next) {
		for i := len(list.head.next); i < rank; i++ {
			list.head.next = append(list.head.next, nil)
			path = append(path, list.head)
		}
	}

	for level := 0; level < rank; level++ {
		node.next[level] = path[level].next[level]
		path[level].next[level] = node
	}

	list.length++
	return list
}

/*

Get ...
*/
func (list *tSkipList[K, V]) Get(key *K) *V {
	v := list.skipr(key, list.head)

	if v != nil && list.Ord.Compare(*v.key, *key) == ord.EQ {
		return v.val
	}

	return nil
}
