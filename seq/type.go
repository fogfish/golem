package seq

// ISeq ...
type ISeq interface {
	TSeq()
}

// Type ...
type Type struct{}

// TSeq ...
func (t Type) TSeq() {}

// FMap ...
type FMap func(ISeq) error
