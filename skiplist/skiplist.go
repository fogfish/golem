package skiplist

import (
	"math"
	"math/rand"
	"time"

	"github.com/fogfish/golem"
)

/*

New creates empty skip list
*/
func New() SkipList {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	head := &tSkipNode{key: tVoid(0), next: []*tSkipNode{nil}}
	list := &tSkipList{
		head:   head,
		length: 0,
		random: random,
	}
	list.Config(4294967296, 1/math.E)
	return list
}

/*

Config ...
*/
func (list *tSkipList) Config(n int, p float64) {
	list.levels = int(math.Log10(float64(n)) / math.Log10(1/p))
	list.p = make([]float64, list.levels+1)

	for i := 1; i <= list.levels; i++ {
		list.p[i] = math.Pow(p, float64(i-1))
	}
}

func (list *tSkipList) mkNode(key golem.Ord, val golem.Data) (int, *tSkipNode) {
	// See: https://golang.org/src/math/rand/rand.go#L150
	p := float64(list.random.Int63()) / (1 << 63)

	level := 1
	for level < list.levels && p < list.p[level] {
		level++
	}

	node := &tSkipNode{key: key, val: val, next: make([]*tSkipNode, level)}
	return level, node
}

func (list *tSkipList) Length() int {
	return list.length
}

/*
func (list *SkipList) Head() *SkipNode {
	return list.head
}
*/

/*

Put
*/
func (list *tSkipList) Put(key golem.Ord, val golem.Data) golem.MapLike {
	v, path := list.head.Skip(key)

	//
	if v != nil && v.key.Eq(key) {
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
func (list *tSkipList) Get(key golem.Ord) golem.Data {
	v, _ := list.head.Skip(key)

	if v != nil && v.key.Eq(key) {
		return v.val
	}

	return nil
}

/*

Remove ...
*/
func (list *tSkipList) Remove(key golem.Ord) golem.Data {
	rank := len(list.head.next)
	v, path := list.head.Skip(key)

	if v != nil && v.key.Eq(key) {
		for level := 0; level < rank; level++ {
			if path[level].next[level] == v {
				if len(v.next) > level {
					path[level].next[level] = v.next[level]
				} else {
					path[level].next[level] = nil
				}
			}
		}
		list.length--
		return v.val
	}

	return nil
}

/*

Skip ...
*/
func (node *tSkipNode) Skip(key golem.Ord) (*tSkipNode, []*tSkipNode) {
	rank := len(node.next)
	path := make([]*tSkipNode, rank)

	v := node
	for level := rank - 1; level >= 0; level-- {
		for v.next[level] != nil && v.next[level].key.Lt(key) {
			v = v.next[level]
		}
		path[level] = v
	}

	return v.next[0], path
}
