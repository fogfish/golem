//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package skiplist_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/fogfish/golem/maplike"
	"github.com/fogfish/golem/maplike/skiplist"
	"github.com/fogfish/golem/pure/ord"
	"github.com/fogfish/it"
)

func TestSkipListGet(t *testing.T) {
	size := 10
	list := mkSkipList(size)

	t.Run("Get Success", func(t *testing.T) {
		for i := 1; i < size; i++ {
			it.Ok(t).If(list.Get(i)).Should().Equal(i)
		}
	})

	t.Run("Get Not Found", func(t *testing.T) {
		for i := 1; i < size; i++ {
			it.Ok(t).If(list.Get(i * size * size)).Should().Equal(0)
		}
	})
}

func TestSkipListPut(t *testing.T) {
	size := 10
	list := mkSkipList(size)

	t.Run("Put Success", func(t *testing.T) {
		for i := 1; i < size; i++ {
			list.Put(i*size*size, i)

			it.Ok(t).If(list.Get(i * size * size)).Should().Equal(i)
		}
	})
}

func TestSkipListRemove(t *testing.T) {
	size := 10
	list := mkSkipList(size)

	t.Run("Remove Success", func(t *testing.T) {
		for i := 1; i < size; i++ {
			list.Remove(i)

			it.Ok(t).If(list.Get(i)).Should().Equal(0)
		}
	})
}

func mkSkipList(n int) maplike.MapLike[int, int] {
	list := skiplist.New[int, int](ord.Int)

	for i := 1; i < n; i++ {
		key := i
		list.Put(key, key)
	}

	return list
}

// func mkSkipListRand(n int) maplike.MapLike[int, int] {
// 	prnd := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	list := skiplist.New[int, int](ord.Int)

// 	for i := 1; i < n; i++ {
// 		key := prnd.Intn(n)
// 		list.Put(key, key)
// 	}

// 	return list
// }

//
// go test -bench=. -benchtime=10s -cpu=1
//
var (
	defCap      int                       = 1000000
	defMapLike  map[int]int               = make(map[int]int)
	defSkipList maplike.MapLike[int, int] = skiplist.New[int, int](ord.Int)
	defShuffle  []int                     = make([]int, defCap)
)

func init() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < defCap; i++ {
		seqKey := i
		defSkipList.Put(seqKey, seqKey)
		defMapLike[seqKey] = seqKey

		rndKey := rnd.Intn(defCap)
		defShuffle[i] = rndKey
	}
}

func BenchmarkSkipListPutTail(b *testing.B) {
	b.ReportAllocs()

	list := skiplist.New[int, int](ord.Int)

	for i := 0; i < b.N; i++ {
		key := i
		list.Put(key, key)
	}
}

func BenchmarkSkipListPutHead(b *testing.B) {
	b.ReportAllocs()

	list := skiplist.New[int, int](ord.Int)

	for i := b.N; i > 0; i-- {
		key := i
		list.Put(key, key)
	}
}

func BenchmarkSkipListPutRand(b *testing.B) {
	b.ReportAllocs()

	list := skiplist.New[int, int](ord.Int)

	for i := 0; i < b.N; i++ {
		key := defShuffle[i%defCap]
		list.Put(key, key)
	}
}

func BenchmarkSkipListGetTail(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := i % defCap
		val := defSkipList.Get(key)
		if val != key {
			panic(fmt.Errorf("invalid state for key %v, unexpected %v", key, val))
		}
	}
}

func BenchmarkSkipListGetHead(b *testing.B) {
	b.ReportAllocs()

	for i := b.N; i > 0; i-- {
		key := i % defCap
		val := defSkipList.Get(key)
		if val != key {
			panic(fmt.Errorf("invalid state for key %v, unexpected %v", key, val))
		}
	}
}

func BenchmarkSkipListGetRand(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := defShuffle[i%defCap]
		val := defSkipList.Get(key)
		if val != key {
			panic(fmt.Errorf("invalid state for key %v, unexpected %v", key, val))
		}
	}
}

func BenchmarkMapLikeGetRand(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := defShuffle[i%defCap]
		val := defMapLike[key]
		if val != key {
			panic(fmt.Errorf("invalid state for key %v, unexpected %v", key, val))
		}
	}
}
