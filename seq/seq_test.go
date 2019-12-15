package seq_test

import (
	"testing"

	"github.com/fogfish/golem/generic"

	"github.com/fogfish/golem/seq"
	"github.com/fogfish/it"
)

var sequence seq.AnyT = seq.AnyT{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

func element(e generic.T) func(x generic.T) bool {
	return func(x generic.T) bool { return e == x }
}

func odd(x generic.T) bool {
	return x.(int)%2 == 0
}

func beforeThree(x generic.T) bool {
	return x.(int) < 3
}

func afterThree(x generic.T) bool {
	return x.(int) > 3
}

func beforeTen(x generic.T) bool {
	return x.(int) < 10
}

func x10(x generic.T) generic.T {
	return x.(int) * 10
}

func sum(x generic.T, y generic.T) generic.T {
	return x.(int) + y.(int)
}

func TestSeqContain(t *testing.T) {
	it.Ok(t).
		If(sequence.Contain(0)).Should().Equal(true).
		If(sequence.Contain(11)).Should().Equal(false)
}

func TestSeqCount(t *testing.T) {
	it.Ok(t).
		If(sequence.Count(odd)).Should().Equal(5)
}

func TestSeqDistinct(t *testing.T) {
	it.Ok(t).
		If(sequence.Distinct()).Should().Equal(sequence).
		If(append(sequence, sequence...).Distinct()).Should().Equal(sequence)
}

func TestSeqExists(t *testing.T) {
	it.Ok(t).
		If(sequence.Exists(element(5))).Should().Equal(true).
		If(sequence.Exists(element(15))).Should().Equal(false)
}

func TestSeqDrop(t *testing.T) {
	it.Ok(t).
		If(sequence.Drop(3)).
		Should().Equal(seq.AnyT{3, 4, 5, 6, 7, 8, 9})
}

func TestSeqDropWhile(t *testing.T) {
	it.Ok(t).
		If(sequence.DropWhile(beforeThree)).
		Should().Equal(seq.AnyT{3, 4, 5, 6, 7, 8, 9}).
		//
		If(sequence.DropWhile(beforeTen)).
		Should().Equal(seq.AnyT{})
}

func TestSeqFilter(t *testing.T) {
	it.Ok(t).
		If(sequence.Filter(odd)).Should().
		Eq(seq.AnyT{0, 2, 4, 6, 8})
}

func TestSeqFind(t *testing.T) {
	it.Ok(t).
		If(sequence.Find(afterThree)).Should().Equal(4).
		If(sequence.Find(func(generic.T) bool { return false })).Should().Equal(nil)
}

func TestSeqForAll(t *testing.T) {
	it.Ok(t).
		If(sequence.ForAll(beforeTen)).Should().Equal(true).
		If(sequence.ForAll(beforeThree)).Should().Equal(false).
		If(seq.AnyT{}.ForAll(beforeThree)).Should().Equal(false)
}

func TestSeqFMap(t *testing.T) {
	sum := 0
	sequence.FMap(func(x generic.T) { sum = sum + x.(int) })

	it.Ok(t).
		If(sum).Should().Equal(45)
}

func TestSeqFold(t *testing.T) {
	it.Ok(t).
		If(sequence.Fold(sum, 0)).Should().Equal(45)
}

func TestSeqGroupBy(t *testing.T) {
	it.Ok(t).
		If(sequence.GroupBy(func(x generic.T) int { return x.(int) % 2 })).
		Should().Equal(
		map[int]seq.AnyT{
			0: seq.AnyT{0, 2, 4, 6, 8},
			1: seq.AnyT{1, 3, 5, 7, 9},
		},
	)
}

func TestSeqJoin(t *testing.T) {
	it.Ok(t).
		If(seq.AnyT{}.Join(sequence)).Should().Equal(sequence)
}

func TestSeqMap(t *testing.T) {
	it.Ok(t).
		If(sequence.Map(x10)).
		Should().Equal(seq.AnyT{0, 10, 20, 30, 40, 50, 60, 70, 80, 90})
}

func TestSeqPartition(t *testing.T) {
	it.Ok(t).
		If(left(sequence.Partition(odd))).Should().Equal(seq.AnyT{0, 2, 4, 6, 8}).
		If(right(sequence.Partition(odd))).Should().Equal(seq.AnyT{1, 3, 5, 7, 9})
}

func left(a seq.AnyT, b seq.AnyT) seq.AnyT {
	return a
}

func right(a seq.AnyT, b seq.AnyT) seq.AnyT {
	return b
}

func TestSeqTake(t *testing.T) {
	it.Ok(t).
		If(sequence.Take(3)).Should().Equal(seq.AnyT{0, 1, 2})
}

func TestSeqTakeWhile(t *testing.T) {
	it.Ok(t).
		If(sequence.TakeWhile(beforeThree)).
		Should().Equal(seq.AnyT{0, 1, 2}).
		//
		If(sequence.TakeWhile(beforeTen)).
		Should().Equal(sequence)
}

func TestSeqMonoidLawIdentity(t *testing.T) {
	it.Ok(t).
		If(sequence.Mappend(sequence.Mempty())).Should().
		Eq(sequence.Mempty().Mappend(sequence))
}

func TestSeqMonoidLawAssociativity(t *testing.T) {
	it.Ok(t).
		If(seq.AnyT{1}.Mappend(seq.AnyT{2}).Mappend(seq.AnyT{3})).Should().
		Eq(seq.AnyT{1}.Mappend(seq.AnyT{2}.Mappend(seq.AnyT{3})))
}

func TestSeqMonoidMap(t *testing.T) {
	fmap := sequence.MMap(seq.AnyT{})
	it.Ok(t).
		If(fmap(func(x generic.T) seq.AnyT { return seq.AnyT{x.(int) * 10} })).
		Should().Equal(seq.AnyT{0, 10, 20, 30, 40, 50, 60, 70, 80, 90})
}
