package golem

/*

Eq interface defines equality for data types.
*/
type Eq interface {
	Eq(Eq) bool
}

/*

Ord interface defines total order of data types.
*/
type Ord interface {
	Eq
	Lt(Ord) bool
}

/*

Data is a fundamental primitive for data struct
*/
type Data interface{ Type() }

/*

Traversable ...
*/
// type Traversable interface {
// 	Cons(Ord) Traversable
// 	Head() Ord
// 	Tail() Traversable
// }

/*

MapLike defines a trait for container, which associate keys with values.
*/
type MapLike interface {
	Put(Ord, Data) MapLike
	Get(Ord) Data
	Remove(Ord) Data
	// Has(Ord) bool
	// FMap(func(Ord, Data) error) error
}

/*

Typeable is a type tag make any structure compatible the library
  type MyStruct struct{
		golem.Typeable
	}
*/
type Typeable struct{}

func (Typeable) Type() {}

var (
	_ Data = Typeable{}
	_ Data = (*Typeable)(nil)
)
