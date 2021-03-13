package golem

/*

Eq interface defines equality and inequality for data types.
*/
type Eq interface {
	Eq(Eq) bool
	Ne(Eq) bool
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
