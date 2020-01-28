package stream

import "github.com/fogfish/golem/generic"

//
type AnyT struct {
	Head     generic.T
	tail     func() AnyT
	nonempty bool
}

//
func NewAnyT(head generic.T, tail func() AnyT) AnyT {
	return AnyT{head, tail, true}
}

//
func (s AnyT) Tail() AnyT {
	if s.tail == nil {
		return AnyT{}
	}

	return s.tail()
}

//
func (s AnyT) Empty() bool {
	return !s.nonempty
}

// FMap applies high-order function (clojure) to all elements of sequence.
func (s AnyT) FMap(f func(generic.T)) {
	stream := s
	for !stream.Empty() {
		f(stream.Head)
		stream = stream.Tail()
	}
}
