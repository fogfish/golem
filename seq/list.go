package seq

// List ...
type List struct {
	Head ISeq
	Tail *List
}

//
func (list List) Cons(hd ISeq) *List {
	if list.Head == nil && list.Tail == nil {
		return &List{Head: hd}
	}

	return &List{Head: hd, Tail: &list}
}

// FMap ...
func (list List) FMap(fmap FMap) error {
	p := &list
	for p != nil {
		if err := fmap(p.Head); err != nil {
			return err
		}
		p = p.Tail
	}
	return nil
}
