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

Please see the article that depicts the data structure
http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.17.524
*/
package skiplist

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/fogfish/golem/maplike"
	"github.com/fogfish/golem/pure/ord"
)

/*

tSkipList implements MapLike with skip data structure
*/
type tSkipList[K, V any] struct {
	ord.Ord[K]

	//
	// head of the list, the node is a lowerst element
	head *tSkipNode[K, V]

	//
	// max levels of the nodes
	levels int

	//
	// number of elements in the list, O(1)
	length int

	//
	// random generator
	random rand.Source

	//
	// probability table to determine node level
	p []float64

	//
	// buffer to estimate the skip path during insert / remove
	// the buffer implements optimization of memory allocations
	path []*tSkipNode[K, V]
}

var (
	_ maplike.MapLike[int, int] = &tSkipList[int, int]{Ord: ord.Int}
	_ maplike.MapLike[int, int] = (*tSkipList[int, int])(nil)
)

/*

New creates empty skiplist
*/
func New[K, V any](compare ord.Ord[K]) maplike.MapLike[K, V] {
	levels, ptable := probability(4294967296, 1/math.E)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &tSkipList[K, V]{
		Ord:    compare,
		head:   &tSkipNode[K, V]{fingers: make([]*tSkipNode[K, V], levels)},
		levels: levels,
		length: 0,
		random: random,
		p:      ptable,
		path:   make([]*tSkipNode[K, V], levels),
	}
}

/*

calculates probability table
*/
func probability(n int, p float64) (int, []float64) {
	level := int(math.Log10(float64(n)) / math.Log10(1/p))
	table := make([]float64, level+1)

	for i := 1; i <= level; i++ {
		table[i-1] = math.Pow(p, float64(i-1))
	}

	return level, table
}

// String converts table to string
func (list *tSkipList[K, V]) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("--- SkipList %p ---\n", &list))

	v := list.head
	for v != nil {
		buffer.WriteString(v.String())
		buffer.WriteString("\n")
		v = v.fingers[0]
	}

	return buffer.String()
}

/*

Put insters the element into the list
*/
func (list *tSkipList[K, V]) Put(key K, val V) maplike.MapLike[K, V] {
	v, path := list.skip(key)

	if v != nil && list.Ord.Compare(v.key, key) == ord.EQ {
		v.val = val
		return list
	}

	rank, node := list.mkNode(key, val)

	// re-bind fingers to new node
	for level := 0; level < rank; level++ {
		node.fingers[level] = path[level].fingers[level]
		path[level].fingers[level] = node
	}

	list.length++
	return list
}

/*

skip algorithm is similar to search algorithm that traversing forward pointers.
skip maintain the vector path that contains a pointer to the rightmost node
of level i or higher that is to the left of the location of the
insertion/deletion.
*/
func (list *tSkipList[K, V]) skip(key K) (*tSkipNode[K, V], []*tSkipNode[K, V]) {
	path := list.path

	node := list.head
	next := node.fingers
	for level := list.levels - 1; level >= 0; level-- {
		for next[level] != nil && list.Ord.Compare(next[level].key, key) == ord.LT {
			node = node.fingers[level]
			next = node.fingers
		}
		path[level] = node
	}

	return next[0], path
}

/*

mkNode creates a new node, randomly defines empty fingers (level of the node)
*/
func (list *tSkipList[K, V]) mkNode(key K, val V) (int, *tSkipNode[K, V]) {
	// See: https://golang.org/src/math/rand/rand.go#L150
	p := float64(list.random.Int63()) / (1 << 63)

	level := 0
	for level < list.levels && p < list.p[level] {
		level++
	}

	node := &tSkipNode[K, V]{
		key:     key,
		val:     val,
		fingers: make([]*tSkipNode[K, V], level),
	}

	return level, node
}

/*

Get looks up the element in the list
*/
func (list *tSkipList[K, V]) Get(key K) V {
	node := list.search(key)

	if node != nil && list.Ord.Compare(node.key, key) == ord.EQ {
		return node.val
	}

	return *new(V)
}

/*

search algorithm traversing forward pointers that do not jumps over the node
containing the element (for each level the finger shall be less than key).
When no more progress can be made at the current level of forward pointers,
the search moves down to the next level. When we can make no more progress at
level 0, we must be immediately in front of the node that contains
the desired element (if it is in the list).
*/
func (list *tSkipList[K, V]) search(key K) *tSkipNode[K, V] {
	node := list.head
	next := list.head.fingers
	for level := list.levels - 1; level >= 0; level-- {
		for next[level] != nil && list.Ord.Compare(next[level].key, key) == ord.LT {
			node = node.fingers[level]
			next = node.fingers
		}
	}

	return next[0]
}

/*

Remove element from the list
*/
func (list *tSkipList[K, V]) Remove(key K) V {
	rank := len(list.head.fingers)
	v, path := list.skip(key)

	if v != nil && list.Ord.Compare(v.key, key) == ord.EQ {
		for level := 0; level < rank; level++ {
			if path[level].fingers[level] == v {
				if len(v.fingers) > level {
					path[level].fingers[level] = v.fingers[level]
				} else {
					path[level].fingers[level] = nil
				}
			}
		}
		list.length--
		return v.val
	}

	return *new(V)
}
