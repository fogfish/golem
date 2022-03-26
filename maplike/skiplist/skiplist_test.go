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
)

func TestMapLikeOne(t *testing.T) {
	list := skiplist.New[int, string](ord.Int)

	key1 := 1
	val1 := "abc"
	list.Put(&key1, &val1)

	key2 := 2
	val2 := "xxx"
	list.Put(&key2, &val2)

	key3 := 2
	o := list.Get(&key3)

	if o != nil {
		fmt.Printf("--> %v\n", *o)
	}

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
	defShuffle  []*int                    = make([]*int, defCap)
)

func init() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < defCap; i++ {
		seqKey := i
		defSkipList.Put(&seqKey, &seqKey)

		rndKey := rnd.Intn(defCap)
		defShuffle[i] = &rndKey
	}
}

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

// func BenchmarkSetRand(b *testing.B) {
// 	b.ReportAllocs()
// 	list := skiplist.New()
// 	data := golem.String("")

// 	for i := 0; i < b.N; i++ {
// 		list.Put(defShuffle[i%defCap], data)
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
		if *o != *defShuffle[i%defCap] {
			panic("+++")
		}
	}
}
