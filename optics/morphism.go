package optics

import "reflect"

/*

Morphism[S] is Len[S, A] composed with the value A
*/
type Morphism[S any] interface {
	Put(*S) error
	PutValue(reflect.Value) error
}

/*

Morphisms[S] is an ordered set of morphisms applied to S
*/
type Morphisms[S any] []Morphism[S]

func (m Morphisms[S]) Put(s *S) error {
	g := reflect.ValueOf(s)

	for _, f := range m {
		if err := f.PutValue(g); err != nil {
			return err
		}
	}

	return nil
}

/*

Morph composes Lens[S, A] with the value A
*/
func Morph[S, A any](lens Lens[S, A], a A) Morphism[S] {
	return &morphism[S, A]{
		Lens:      lens,
		Reflector: lens.(Reflector[A]),
		Value:     a,
	}
}

type morphism[S, A any] struct {
	Lens[S, A]
	Reflector[A]
	Value A
}

func (m morphism[S, A]) Put(s *S) error {
	return m.Lens.Put(s, m.Value)
}

func (m morphism[S, A]) PutValue(s reflect.Value) error {
	return m.Reflector.PutValue(s, m.Value)
}
