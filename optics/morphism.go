package optics

/*

Morphism[S] is Len[S, A] composed with the value A
*/
type Morphism[S any] interface {
	Put(*S) error
}

/*

Morphisms[S] is an ordered set of morphisms applied to S
*/
type Morphisms[S any] []Morphism[S]

func (m Morphisms[S]) Put(s *S) error {
	for _, f := range m {
		if err := f.Put(s); err != nil {
			return err
		}
	}

	return nil
}

/*

MorphismConst composes Lens[S, A] with the value A
*/
func Morph[S, A any](lens Lens[S, A], a A) Morphism[S] {
	return &morphism[S, A]{Lens: lens, Value: a}
}

type morphism[S, A any] struct {
	Lens[S, A]
	Value A
}

func (m morphism[S, A]) Put(s *S) error {
	return m.Lens.Put(s, m.Value)
}
