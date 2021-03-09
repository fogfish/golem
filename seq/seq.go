package seq

// Seq ...
type Seq []ISeq

// FMap ...
func (seq Seq) FMap(fmap FMap) error {
	for _, x := range seq {
		if err := fmap(x); err != nil {
			return err
		}
	}
	return nil
}
