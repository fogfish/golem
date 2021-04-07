package golem

//
//
type String string

func (String) Type() {}

func (s String) Eq(x Eq) bool {
	return s == x.(String)
}

func (s String) Lt(x Ord) bool {
	return s < x.(String)
}

//
//
type Int int

func (Int) Type() {}

func (i Int) Eq(x Eq) bool {
	return i == x.(Int)
}

func (i Int) Lt(x Ord) bool {
	return i < x.(Int)
}
