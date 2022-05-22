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

/*

mkLens instantiates a typed Lens[S, A] for hseq.Type[S]
*/
func mkLens[S, A any](t hseq.Type[S]) Lens[S, A] {
	switch hseq.AssertType[S, A](t, false) {
	case reflect.String:
		return mkLensStructString(t).(Lens[S, A])
	case reflect.Int:
		return mkLensStructInt(t).(Lens[S, A])
	case reflect.Float64:
		return mkLensStructFloat64(t).(Lens[S, A])
	case reflect.Struct:
		return mkLensStruct[S, A](t)
	default:
		panic(fmt.Errorf("Unknown lens type %v", t.Type.Name()))
	}
}

/*

lensStructString implements lens for string type
*/
type lensStructString[S any] struct{ hseq.Type[S] }

func mkLensStructString[S any](t hseq.Type[S]) Lens[S, string] {
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

func mkLensStructInt[S any](t hseq.Type[S]) Lens[S, int] {
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

func mkLensStructFloat64[S any](t hseq.Type[S]) Lens[S, float64] {
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

func mkLensStruct[S, A any](t hseq.Type[S]) Lens[S, A] {
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

/*

ForProduct1 split structure with 1 field to set of lenses
*/
func ForProduct1[T, A any]() Lens[T, A] {
	return hseq.FMap1(
		hseq.Generic[T](),
		mkLens[T, A],
	)
}

func ForProduct2[T, A, B any]() (
	Lens[T, A],
	Lens[T, B],
) {
	return hseq.FMap2(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
	)
}

func ForProduct3[T, A, B, C any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
) {
	return hseq.FMap3(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
	)
}

func ForProduct4[T, A, B, C, D any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
) {
	return hseq.FMap4(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
	)
}

func ForProduct5[T, A, B, C, D, E any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
) {
	return hseq.FMap5(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
	)
}

func ForProduct6[T, A, B, C, D, E, F any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
) {
	return hseq.FMap6(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
	)
}

func ForProduct7[T, A, B, C, D, E, F, G any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
) {
	return hseq.FMap7(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
	)
}

func ForProduct8[T, A, B, C, D, E, F, G, H any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
) {
	return hseq.FMap8(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
	)
}

func ForProduct9[T, A, B, C, D, E, F, G, H, I any]() (
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
	return hseq.FMap9(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
	)
}

func ForProduct10[T, A, B, C, D, E, F, G, H, I, J any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
	Lens[T, J],
) {
	return hseq.FMap10(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
		mkLens[T, J],
	)
}

func ForProduct11[T, A, B, C, D, E, F, G, H, I, J, K any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
	Lens[T, J],
	Lens[T, K],
) {
	return hseq.FMap11(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
		mkLens[T, J],
		mkLens[T, K],
	)
}

func ForProduct12[T, A, B, C, D, E, F, G, H, I, J, K, L any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
	Lens[T, J],
	Lens[T, K],
	Lens[T, L],
) {
	return hseq.FMap12(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
		mkLens[T, J],
		mkLens[T, K],
		mkLens[T, L],
	)
}

func ForProduct13[T, A, B, C, D, E, F, G, H, I, J, K, L, M any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
	Lens[T, J],
	Lens[T, K],
	Lens[T, L],
	Lens[T, M],
) {
	return hseq.FMap13(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
		mkLens[T, J],
		mkLens[T, K],
		mkLens[T, L],
		mkLens[T, M],
	)
}

func ForProduct14[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
	Lens[T, J],
	Lens[T, K],
	Lens[T, L],
	Lens[T, M],
	Lens[T, N],
) {
	return hseq.FMap14(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
		mkLens[T, J],
		mkLens[T, K],
		mkLens[T, L],
		mkLens[T, M],
		mkLens[T, N],
	)
}

func ForProduct15[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
	Lens[T, J],
	Lens[T, K],
	Lens[T, L],
	Lens[T, M],
	Lens[T, N],
	Lens[T, O],
) {
	return hseq.FMap15(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
		mkLens[T, J],
		mkLens[T, K],
		mkLens[T, L],
		mkLens[T, M],
		mkLens[T, N],
		mkLens[T, O],
	)
}

func ForProduct16[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
	Lens[T, J],
	Lens[T, K],
	Lens[T, L],
	Lens[T, M],
	Lens[T, N],
	Lens[T, O],
	Lens[T, P],
) {
	return hseq.FMap16(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
		mkLens[T, J],
		mkLens[T, K],
		mkLens[T, L],
		mkLens[T, M],
		mkLens[T, N],
		mkLens[T, O],
		mkLens[T, P],
	)
}

func ForProduct17[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
	Lens[T, J],
	Lens[T, K],
	Lens[T, L],
	Lens[T, M],
	Lens[T, N],
	Lens[T, O],
	Lens[T, P],
	Lens[T, Q],
) {
	return hseq.FMap17(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
		mkLens[T, J],
		mkLens[T, K],
		mkLens[T, L],
		mkLens[T, M],
		mkLens[T, N],
		mkLens[T, O],
		mkLens[T, P],
		mkLens[T, Q],
	)
}

func ForProduct18[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
	Lens[T, J],
	Lens[T, K],
	Lens[T, L],
	Lens[T, M],
	Lens[T, N],
	Lens[T, O],
	Lens[T, P],
	Lens[T, Q],
	Lens[T, R],
) {
	return hseq.FMap18(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
		mkLens[T, J],
		mkLens[T, K],
		mkLens[T, L],
		mkLens[T, M],
		mkLens[T, N],
		mkLens[T, O],
		mkLens[T, P],
		mkLens[T, Q],
		mkLens[T, R],
	)
}

func ForProduct19[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
	Lens[T, J],
	Lens[T, K],
	Lens[T, L],
	Lens[T, M],
	Lens[T, N],
	Lens[T, O],
	Lens[T, P],
	Lens[T, Q],
	Lens[T, R],
	Lens[T, S],
) {
	return hseq.FMap19(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
		mkLens[T, J],
		mkLens[T, K],
		mkLens[T, L],
		mkLens[T, M],
		mkLens[T, N],
		mkLens[T, O],
		mkLens[T, P],
		mkLens[T, Q],
		mkLens[T, R],
		mkLens[T, S],
	)
}

func ForProduct20[T, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, U any]() (
	Lens[T, A],
	Lens[T, B],
	Lens[T, C],
	Lens[T, D],
	Lens[T, E],
	Lens[T, F],
	Lens[T, G],
	Lens[T, H],
	Lens[T, I],
	Lens[T, J],
	Lens[T, K],
	Lens[T, L],
	Lens[T, M],
	Lens[T, N],
	Lens[T, O],
	Lens[T, P],
	Lens[T, Q],
	Lens[T, R],
	Lens[T, S],
	Lens[T, U],
) {
	return hseq.FMap20(
		hseq.Generic[T](),
		mkLens[T, A],
		mkLens[T, B],
		mkLens[T, C],
		mkLens[T, D],
		mkLens[T, E],
		mkLens[T, F],
		mkLens[T, G],
		mkLens[T, H],
		mkLens[T, I],
		mkLens[T, J],
		mkLens[T, K],
		mkLens[T, L],
		mkLens[T, M],
		mkLens[T, N],
		mkLens[T, O],
		mkLens[T, P],
		mkLens[T, Q],
		mkLens[T, R],
		mkLens[T, S],
		mkLens[T, U],
	)
}
