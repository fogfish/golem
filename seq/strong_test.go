package seq_test

import "testing"

// List ...
type ListEl struct {
	head *El
	tail *ListEl
}

func (seq *ListEl) Cons(x *El) *ListEl {
	if seq.head == nil && seq.tail == nil {
		return &ListEl{head: x}
	}

	return &ListEl{head: x, tail: seq}
}

func (seq *ListEl) Head() *El {
	return seq.head
}

func (seq *ListEl) Tail() *ListEl {
	return seq.tail
}

var (
	defStrong *ListEl = &ListEl{}
)

func init() {
	for n := 0; n < defCap; n++ {
		defStrong = defStrong.Cons(&El{ID: n})
	}
}

func BenchmarkStrongCons(b *testing.B) {
	b.ReportAllocs()

	l := &ListEl{}
	for n := 0; n < b.N; n++ {
		l = l.Cons(&El{ID: n})
	}
}

func BenchmarkStrongTail(b *testing.B) {
	b.ReportAllocs()

	l := defStrong

	for n := 0; n < b.N; n++ {
		el = l.Head()

		if l = l.Tail(); l == nil {
			l = defStrong
		}
	}
}
