package skiplist_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fogfish/golem"
	"github.com/fogfish/golem/skiplist"
	"github.com/fogfish/it"
)

func TestMapLikeOne(t *testing.T) {
	key := golem.String("key")
	val := golem.String("val")

	list := skiplist.New()
	it.Ok(t).
		If(list.Length()).Equal(0).
		IfNil(list.Get(key))

	list.Put(key, val)
	it.Ok(t).
		If(list.Length()).Equal(1).
		If(list.Get(key)).Equal(val).
		IfNil(list.Get(val))

	list.Remove(key)
	it.Ok(t).
		If(list.Length()).Equal(0).
		IfNil(list.Get(key))
}

func TestMapLikeMany(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 100000

	keys := make([]golem.Int, n)
	for i := 0; i < n; i++ {
		keys[i] = golem.Int(rnd.Intn(n))
	}

	//
	list := skiplist.New()
	for _, key := range keys {
		list.Put(key, key)
	}

	it.Ok(t).
		If(list.Length()).NotEqual(0).
		IfNil(list.Get(golem.String("x")))

	//
	for _, key := range keys {
		val := list.Get(key)
		str, ok := val.(golem.Int)

		it.Ok(t).
			IfNotNil(val).
			If(val).Equal(key).
			IfTrue(ok).
			If(str).Equal(key)
	}

	//
	for _, key := range keys {
		list.Remove(key)
		it.Ok(t).IfNil(list.Get(key))
	}

	it.Ok(t).If(list.Length()).Equal(0)
}

//
// go test -bench=. -benchtime=20s
//

var (
	defCap      int               = 1000000
	defSkipList skiplist.SkipList = skiplist.New()
	defShuffle  []golem.Ord       = make([]golem.Ord, defCap)
)

func init() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < defCap; i++ {
		seqKey := golem.Int(i)
		defSkipList.Put(&seqKey, &seqKey)

		rndKey := golem.Int(rnd.Intn(defCap))
		defShuffle[i] = &rndKey
	}
}

func BenchmarkSetTail(b *testing.B) {
	b.ReportAllocs()

	list := skiplist.New()
	data := golem.String("")

	for i := 0; i < b.N; i++ {
		key := golem.Int(i)
		list.Put(&key, data)
	}
}

func BenchmarkSetHead(b *testing.B) {
	b.ReportAllocs()
	list := skiplist.New()
	data := golem.String("")

	for i := b.N; i > 0; i-- {
		key := golem.Int(i)
		list.Put(&key, data)
	}
}

func BenchmarkSetRand(b *testing.B) {
	b.ReportAllocs()
	list := skiplist.New()
	data := golem.String("")

	for i := 0; i < b.N; i++ {
		list.Put(defShuffle[i%defCap], data)
	}
}

func BenchmarkGetTail(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := golem.Int(i)
		defSkipList.Get(&key)
	}
}

func BenchmarkGetHead(b *testing.B) {
	b.ReportAllocs()

	for i := b.N; i > 0; i-- {
		key := golem.Int(i)
		defSkipList.Get(&key)
	}
}

func BenchmarkGetRand(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		defSkipList.Get(defShuffle[i%defCap])
	}
}
