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
	"testing"

	"github.com/fogfish/golem/maplike"
	"github.com/fogfish/golem/maplike/skiplist"
	"github.com/fogfish/golem/pure/ord"
)

func TestMapLikeOne(t *testing.T) {
	// rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	list := skiplist.New[int, int](ord.Int)

	for i := 1; i < 10; i++ {
		key := i //rnd.Intn(1000000)
		list.Put(key, key)
	}

	fmt.Println(list)

	// key1 := 1
	// val1 := "111"

	// key2 := 2
	// val2 := "222"
	// list.Put(&key2, &val2)
	// fmt.Println(list)

	// key3 := 3
	// val3 := "333"
	// list.Put(&key3, &val3)
	// fmt.Println(list)

	key4 := 9
	o := list.Get(key4)
	// if o != nil {
	fmt.Printf("--> %v\n", o)
	// }

	// list.Remove(4)
	// fmt.Println(list)

	// 	key := golem.String("key")
	// 	val := golem.String("val")

	// 	list := skiplist.New()
	// 	it.Ok(t).
	// 		If(list.Length()).Equal(0).
	// 		IfNil(list.Get(key))

	// 	list.Put(key, val)
	// 	it.Ok(t).
	// 		If(list.Length()).Equal(1).
	// 		If(list.Get(key)).Equal(val).
	// 		IfNil(list.Get(val))

	// 	list.Remove(key)
	// 	it.Ok(t).
	// 		If(list.Length()).Equal(0).
	// 		IfNil(list.Get(key))
}

//
// go test -bench=. -benchtime=20s
//

var (
	defCap      int                       = 1000000
	defSkipList maplike.MapLike[int, int] = skiplist.New[int, int](ord.Int)
	// defShuffleList maplike.MapLike[int, int] = skiplist.New[int, int]()
	defShuffle []int = make([]int, defCap)
)

// func init() {
// 	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

// 	for i := 0; i < defCap; i++ {
// 		seqKey := i
// 		defSkipList.Put(seqKey, seqKey)

// 		rndKey := rnd.Intn(defCap)
// 		defShuffle[i] = rndKey
// 		defShuffleList.Put(rndKey, rndKey)
// 	}

// 	fmt.Println(defShuffleList)
// }

// func BenchmarkPutTail(b *testing.B) {
// 	b.ReportAllocs()

// 	list := skiplist.New[int, int](ord.Int)
// 	data := golem.String("")

// 	for i := 0; i < b.N; i++ {
// 		key := golem.Int(i)
// 		list.Put(&key, data)
// 	}
// }

// func BenchmarkSetHead(b *testing.B) {
// 	b.ReportAllocs()
// 	list := skiplist.New()
// 	data := golem.String("")

// 	for i := b.N; i > 0; i-- {
// 		key := golem.Int(i)
// 		list.Put(&key, data)
// 	}
// }

// func BenchmarkPutRand(b *testing.B) {
// 	b.ReportAllocs()
// 	list := skiplist.New[int, int](ord.Int)

// 	for i := 0; i < b.N; i++ {
// 		list.Put(defShuffle[i%defCap], defShuffle[i%defCap])
// 	}
// }

// func BenchmarkGetTail(b *testing.B) {
// 	b.ReportAllocs()

// 	for i := 0; i < b.N; i++ {
// 		key := golem.Int(i)
// 		defSkipList.Get(&key)
// 	}
// }

// func BenchmarkGetHead(b *testing.B) {
// 	b.ReportAllocs()

// 	for i := b.N; i > 0; i-- {
// 		defSkipList.Get(&i)
// 	}
// }

func BenchmarkGetRand(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		o := defSkipList.Get(defShuffle[i%defCap])
		if o != defShuffle[i%defCap] {
			panic("+++")
		}
	}
}
