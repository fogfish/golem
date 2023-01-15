package optics

import (
	"fmt"
	"reflect"

	"github.com/fogfish/golem/pure/hseq"
)

/*
Lens resembles concept of getters and setters, which you can compose
using functional concepts. In other words, this is combinator data
transformation for pure functional data structure.
*/
type Lens[S, A any] interface {
	Get(*S) A
	Put(*S, A) error
}

/*
Reflector is a Lens over reflect.Value
*/
type Reflector[A any] interface {
	GetValue(reflect.Value) A
	PutValue(reflect.Value, A) error
}

// NewLens instantiates a typed Lens[S, A] for hseq.Type[S]
func NewLens[S, A any](t hseq.Type[S]) Lens[S, A] {
	name, kind := hseq.Assert[S, A](t)
	switch kind {
	case reflect.String:
		if name == "string" {
			//lint:ignore SA5010 linter wrong, string is subtype of any
			return NewLensStructString(t).(Lens[S, A])
		}
		return NewLensStruct[S, A](t)
	case reflect.Int:
		if name == "int" {
			//lint:ignore SA5010 linter wrong, int is subtype of any
			return NewLensStructInt(t).(Lens[S, A])
		}
		return NewLensStruct[S, A](t)
	case reflect.Float64:
		if name == "float64" {
			//lint:ignore SA5010 linter wrong, float is subtype of any
			return NewLensStructFloat64(t).(Lens[S, A])
		}
		return NewLensStruct[S, A](t)
	case reflect.Struct:
		return NewLensStruct[S, A](t)
	default:
		panic(fmt.Errorf("unknown lens type %v", t.Type.Name()))
	}
}

/*
lensStructString implements lens for string type
*/
type lensStructString[S any] struct{ hseq.Type[S] }

func NewLensStructString[S any](t hseq.Type[S]) Lens[S, string] {
	return &lensStructString[S]{t}
}

// Put string to struct
func (lens *lensStructString[S]) Put(s *S, a string) error {
	return lens.PutValue(reflect.ValueOf(s), a)
}

func (lens *lensStructString[S]) PutValue(g reflect.Value, a string) error {
	f := g.Elem().Field(lens.ID)

	if f.Kind() == reflect.Ptr {
		p := reflect.New(lens.PureType)
		p.Elem().SetString(a)
		f.Set(p)
		return nil
	}

	f.SetString(a)
	return nil
}

// Get string from struct
func (lens *lensStructString[S]) Get(s *S) string {
	return lens.GetValue(reflect.ValueOf(s))
}

func (lens *lensStructString[S]) GetValue(g reflect.Value) string {
	f := g.Elem().Field(lens.ID)

	if f.Kind() == reflect.Ptr {
		return f.Elem().String()
	}

	return f.String()
}

/*
lensStructFloat64 implements lens for float type
*/
type lensStructInt[S any] struct{ hseq.Type[S] }

func NewLensStructInt[S any](t hseq.Type[S]) Lens[S, int] {
	return &lensStructInt[S]{t}
}

// Put int to struct
func (lens *lensStructInt[S]) Put(s *S, a int) error {
	return lens.PutValue(reflect.ValueOf(s), a)
}

func (lens *lensStructInt[S]) PutValue(g reflect.Value, a int) error {
	f := g.Elem().Field(lens.ID)

	if f.Kind() == reflect.Ptr {
		p := reflect.New(lens.PureType)
		p.Elem().SetInt(int64(a))
		f.Set(p)
		return nil
	}

	f.SetInt(int64(a))
	return nil
}

// Get float64 from struct
func (lens *lensStructInt[S]) Get(s *S) int {
	return lens.GetValue(reflect.ValueOf(s))
}

func (lens *lensStructInt[S]) GetValue(g reflect.Value) int {
	f := g.Elem().Field(lens.ID)

	if f.Kind() == reflect.Ptr {
		return int(f.Elem().Int())
	}

	return int(f.Int())
}

/*
lensStructFloat64 implements lens for float type
*/
type lensStructFloat64[S any] struct{ hseq.Type[S] }

func NewLensStructFloat64[S any](t hseq.Type[S]) Lens[S, float64] {
	return &lensStructFloat64[S]{t}
}

// Put float64 to struct
func (lens *lensStructFloat64[S]) Put(s *S, a float64) error {
	return lens.PutValue(reflect.ValueOf(s), a)
}

func (lens *lensStructFloat64[S]) PutValue(g reflect.Value, a float64) error {
	f := g.Elem().Field(lens.ID)

	if f.Kind() == reflect.Ptr {
		p := reflect.New(lens.PureType)
		p.Elem().SetFloat(a)
		f.Set(p)
		return nil
	}

	f.SetFloat(a)
	return nil
}

// Get float64 from struct
func (lens *lensStructFloat64[S]) Get(s *S) float64 {
	return lens.GetValue(reflect.ValueOf(s))
}

func (lens *lensStructFloat64[S]) GetValue(g reflect.Value) float64 {
	f := g.Elem().Field(lens.ID)

	if f.Kind() == reflect.Ptr {
		return f.Elem().Float()
	}

	return f.Float()
}

/*
lensStructFloat64 implements lens for float type
*/
type lensStruct[S, A any] struct{ hseq.Type[S] }

func NewLensStruct[S, A any](t hseq.Type[S]) Lens[S, A] {
	return &lensStruct[S, A]{t}
}

// Put reflect.Value to struct
func (lens *lensStruct[S, A]) Put(s *S, a A) error {
	return lens.PutValue(reflect.ValueOf(s), a)
}
func (lens *lensStruct[S, A]) PutValue(g reflect.Value, a A) error {
	f := g.Elem().Field(lens.ID)

	if f.Kind() == reflect.Ptr {
		if lens.Type.Type.Kind() == reflect.Ptr {
			f.Set(reflect.ValueOf(&a))
		} else {
			f.Set(reflect.ValueOf(a))
		}
		return nil
	}

	f.Set(reflect.ValueOf(a))
	return nil
}

// Get reflect.Value from struct
func (lens *lensStruct[S, A]) Get(s *S) A {
	return lens.GetValue(reflect.ValueOf(s))
}

func (lens *lensStruct[S, A]) GetValue(g reflect.Value) A {
	f := g.Elem().Field(lens.ID)

	if f.Kind() == reflect.Ptr {
		return f.Elem().Interface().(A)
	}

	return f.Interface().(A)
}

// ForProduct1 unfold 1 attribute of type
func ForProduct1[T, A any](attr ...string) Lens[T, A] {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New1[T, A]()
	} else {
		seq = hseq.New[T](attr[0])
	}

	return hseq.FMap1(seq,
		NewLens[T, A],
	)
}

// ForProduct2 unfold 2 attribute of type
func ForProduct2[T, A, B any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New2[T, A, B]()
	} else {
		seq = hseq.New[T](attr[0:2]...)
	}

	return hseq.FMap2(seq,
		NewLens[T, A],
		NewLens[T, B],
	)
}

func ForProduct3[T, A, B, C any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New3[T, A, B, C]()
	} else {
		seq = hseq.New[T](attr[0:3]...)
	}

	return hseq.FMap3(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
	)
}

func ForProduct4[T, A, B, C, D any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New4[T, A, B, C, D]()
	} else {
		seq = hseq.New[T](attr[0:4]...)
	}

	return hseq.FMap4(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
	)
}

func ForProduct5[T, A, B, C, D, E any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New5[T, A, B, C, D, E]()
	} else {
		seq = hseq.New[T](attr[0:5]...)
	}

	return hseq.FMap5(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
		NewLens[T, E],
	)
}

func ForProduct6[T, A, B, C, D, E, F any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New6[T, A, B, C, D, E, F]()
	} else {
		seq = hseq.New[T](attr[0:6]...)
	}

	return hseq.FMap6(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
		NewLens[T, E],
		NewLens[T, F],
	)
}

func ForProduct7[T, A, B, C, D, E, F, G any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New7[T, A, B, C, D, E, F, G]()
	} else {
		seq = hseq.New[T](attr[0:7]...)
	}

	return hseq.FMap7(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
		NewLens[T, E],
		NewLens[T, F],
		NewLens[T, G],
	)
}

func ForProduct8[T, A, B, C, D, E, F, G, H any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New8[T, A, B, C, D, E, F, G, H]()
	} else {
		seq = hseq.New[T](attr[0:8]...)
	}

	return hseq.FMap8(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
		NewLens[T, E],
		NewLens[T, F],
		NewLens[T, G],
		NewLens[T, H],
	)
}

func ForProduct9[T, A, B, C, D, E, F, G, H, I any](attr ...string) (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New9[T, A, B, C, D, E, F, G, H, I]()
	} else {
		seq = hseq.New[T](attr[0:9]...)
	}

	return hseq.FMap9(seq,
		NewLens[T, A],
		NewLens[T, B],
		NewLens[T, C],
		NewLens[T, D],
		NewLens[T, E],
		NewLens[T, F],
		NewLens[T, G],
		NewLens[T, H],
		NewLens[T, I],
	)
}
