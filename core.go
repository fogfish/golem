package golem

//
//
type String string

func (String) Type() {}

func (s String) Eq(x Eq) bool {
	switch v := x.(type) {
	case String:
		return s == v
	default:
		return false
	}
}

func (s String) Ne(x Eq) bool {
	return !s.Eq(x)
}

func (s String) Lt(x Ord) bool {
	switch v := x.(type) {
	case String:
		return s < v
	default:
		return false
	}
}

//
//
type Int int

func (Int) Type() {}

func (i Int) Eq(x Eq) bool {
	switch v := x.(type) {
	case Int:
		return i == v
	default:
		return false
	}
}

func (i Int) Ne(x Eq) bool {
	return !i.Eq(x)
}

func (i Int) Lt(x Ord) bool {
	switch v := x.(type) {
	case Int:
		return i < v
	default:
		return false
	}
}
