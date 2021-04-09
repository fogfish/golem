package eq_test

import (
	"fmt"
	"testing"

	"github.com/fogfish/golem"
	"github.com/fogfish/golem/pure/eq"
)

var bbEq bool

func BenchmarkEqInt(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		bbEq = eq.Int.Eq(n, n+1)
	}
}

type A struct {
	ID   int
	UID  int
	Name string
}

var eqA eq.Eq = eq.Struct{eq.Int, eq.Int, eq.String}.From3(
	func(x golem.T) (golem.T, golem.T, golem.T) {
		return x.(A).ID, x.(A).UID, x.(A).Name
	},
)

var eqByID = eq.ContraMap{eq.Int}.From(
	func(x golem.T) golem.T {
		return x.(A).ID
	},
)

func TestContraMap(t *testing.T) {
	a := A{1, 1, "a"}
	fmt.Println(eqByID.Eq(a, a))
}

func BenchmarkContraMap(b *testing.B) {
	b.ReportAllocs()

	a := A{1, 1, "a"}
	for n := 0; n < b.N; n++ {
		bbEq = eqByID.Eq(a, a)
	}
}

func TestStruct(t *testing.T) {
	a := A{1, 1, "a"}
	fmt.Println(eqA.Eq(a, a))
}

func BenchmarkTestStruct(b *testing.B) {
	b.ReportAllocs()

	a := A{1, 1, "a"}
	for n := 0; n < b.N; n++ {
		bbEq = eqA.Eq(a, a)
	}
}
