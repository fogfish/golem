package list_test

import (
	"container/list"
	"testing"
)

var (
	defStdLib *list.List = list.New()
)

func init() {
	for n := 0; n < defCap; n++ {
		defStdLib.PushFront(&El{ID: n})
	}
}

func BenchmarkStdLibCons(b *testing.B) {
	b.ReportAllocs()

	l := list.New()
	for n := 0; n < b.N; n++ {
		l.PushFront(&El{ID: n})
	}
}

func BenchmarkStdLibTail(b *testing.B) {
	b.ReportAllocs()

	l := defStdLib.Front()

	for n := 0; n < b.N; n++ {
		switch v := l.Value.(type) {
		case *El:
			el = v
		}

		if l = l.Next(); l == nil {
			l = defStdLib.Front()
		}
	}
}
