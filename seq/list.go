package seq

import "github.com/fogfish/golem"

// List ...
type List struct {
	head golem.Ord
	tail *List
}

func New() *List {
	return &List{}
}

//
// https://stackoverflow.com/questions/13476349/check-for-nil-and-nil-interface-in-go
// var Empty golem.Traversable = (*tList)(nil)

/*

Cons ...
*/
func (seq *List) Cons(x golem.Ord) *List {
	if seq.head == nil && seq.tail == nil {
		return &List{head: x}
	}

	return &List{head: x, tail: seq}
}

/*

Head ...
*/
func (seq *List) Head() golem.Ord {
	return seq.head
}

/*

Tail ...
*/
func (seq *List) Tail() *List {
	return seq.tail
}
