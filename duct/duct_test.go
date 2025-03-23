//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package duct_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fogfish/golem/duct"
	"github.com/fogfish/it/v2"
)

type T interface{ Apply(duct.Visitor) error }
type A string
type B string
type C string
type D string

var (
	a    = duct.L1[A](nil)
	b    = duct.L1[B](nil)
	c    = duct.L1[C](nil)
	cs   = duct.L1[[]C](nil)
	fAB  = duct.L2[A, B](nil)
	fABs = duct.L2[A, []B](nil)
	fBC  = duct.L2[B, C](nil)
	fCsD = duct.L2[[]C, D](nil)
)

var tests = map[string]T{
	`
𝑚: ø ⟼ A
	ƒ: ø ⟼ A
`: duct.From(a),

	`
𝑚: ø ⟼ ø
	ƒ: ø ⟼ A
	ƒ: A ⟼ ø
`: duct.Yield(a, duct.From(a)),

	`
𝑚: ø ⟼ B
	ƒ: ø ⟼ A
	ƒ: A ⟼ B
`: duct.Join(fAB, duct.From(a)),

	`
𝑚: ø ⟼ C
	ƒ: ø ⟼ A
	ƒ: A ⟼ B
	ƒ: B ⟼ C
`: duct.Join(fBC, duct.Join(fAB, duct.From(a))),

	`
𝑚: ø ⟼ ø
	ƒ: ø ⟼ A
	ƒ: A ⟼ B
	ƒ: B ⟼ ø
`: duct.Yield(b, duct.Join(fAB, duct.From(a))),

	`
𝑚: ø ⟼ []B
	ƒ: ø ⟼ A
	ƒ: A ⟼ []B
`: duct.Join(fABs, duct.From(a)),

	`
𝑚: ø ⟼ C
	ƒ: ø ⟼ A
	ƒ: A ⟼ []B
		ƒ: B ⟼ C
`: duct.LiftF(fBC, duct.Join(fABs, duct.From(a))),

	`
𝑚: ø ⟼ ø
	ƒ: ø ⟼ A
	ƒ: A ⟼ []B
		ƒ: B ⟼ C
		ƒ: C ⟼ ø
`: duct.Yield(c, duct.LiftF(fBC, duct.Join(fABs, duct.From(a)))),

	`
𝑚: ø ⟼ []C
	ƒ: ø ⟼ A
	ƒ: A ⟼ []B
		ƒ: B ⟼ C
`: duct.Unit(duct.LiftF(fBC, duct.Join(fABs, duct.From(a)))),

	`
𝑚: ø ⟼ ø
	ƒ: ø ⟼ A
	ƒ: A ⟼ []B
		ƒ: B ⟼ C
	ƒ: []C ⟼ ø
`: duct.Yield(cs, duct.Unit(duct.LiftF(fBC, duct.Join(fABs, duct.From(a))))),

	`
𝑚: ø ⟼ D
	ƒ: ø ⟼ A
	ƒ: A ⟼ []B
		ƒ: B ⟼ C
	ƒ: []C ⟼ D
`: duct.Join(fCsD, duct.Unit(duct.LiftF(fBC, duct.Join(fABs, duct.From(a))))),

	// Note: Flat w/o Join or Yield leading to incomplete state
	`
𝑚: ø ⟼ 
	ƒ: ø ⟼ A
	ƒ: A ⟼ []B
`: duct.WrapF(duct.Join(fABs, duct.From(a))),

	`
𝑚: ø ⟼ ø
	ƒ: ø ⟼ A
	ƒ: A ⟼ []B
		ƒ: B ⟼ ø
`: duct.Yield(b, duct.WrapF(duct.Join(fABs, duct.From(a)))),
}

func TestTypeOf(t *testing.T) {
	it.Then(t).Should(
		it.Equal(duct.TypeOf[A](), "A"),
		it.Equal(duct.TypeOf[*A](), "*A"),
		it.Equal(duct.TypeOf[[]A](), "[]A"),
		it.Equal(duct.TypeOf[[]*A](), "[]*A"),
		it.Equal(duct.TypeOf[[][]A](), "[][]A"),
		it.Equal(duct.TypeOf[[][]*A](), "[][]*A"),
		it.Equal(duct.TypeOf[it.Check](), "Check"),
	)
}

func TestDuct(t *testing.T) {
	for expect, spec := range tests {
		p := &printer{}
		spec.Apply(p)

		it.Then(t).Should(
			it.Equal(p.String(), expect),
		)
	}
}

type printer struct {
	duct.AstVisitor
	strings.Builder
}

func (p *printer) OnEnterMorphism(depth int, node duct.AstSeq) error {
	nt := &typer{}
	node.Apply(0, nt)

	fmt.Fprintf(p, "\n𝑚: %s ⟼ %s\n", nt.A, nt.B)
	return nil
}

func (p *printer) OnEnterFrom(depth int, node duct.AstFrom) error {
	s := strings.Repeat("\t", depth)
	fmt.Fprintf(p, "%sƒ: ø ⟼ %s\n", s, node.Type)
	return nil
}

func (p *printer) OnEnterYield(depth int, node duct.AstYield) error {
	s := strings.Repeat("\t", depth)
	fmt.Fprintf(p, "%sƒ: %s ⟼ ø\n", s, node.Type)
	return nil
}

func (p *printer) OnEnterMap(depth int, node duct.AstMap) error {
	s := strings.Repeat("\t", depth)
	fmt.Fprintf(p, "%sƒ: %s ⟼ %s\n", s, node.TypeA, node.TypeB)
	return nil
}

// resolve morphism 𝑚: A ⟼ B types
type typer struct {
	duct.AstVisitor
	A, B string
}

func (t *typer) OnEnterSeq(depth int, node duct.AstSeq) error {
	if depth == 0 {
		return nil
	}

	nt := &typer{}
	node.Apply(0, nt)

	if len(t.A) == 0 {
		t.A = nt.A
	}
	t.B = nt.B

	return nil
}

func (t *typer) OnLeaveSeq(depth int, node duct.AstSeq) error {
	if depth == 0 {
		return nil
	}

	if !node.Deferred {
		if t.B != "ø" {
			t.B = "[]" + t.B
		}
	}
	return nil
}

func (t *typer) OnEnterMap(depth int, node duct.AstMap) error {
	if len(t.A) == 0 {
		t.A = node.TypeA
	}
	t.B = node.TypeB
	return nil
}

func (t *typer) OnEnterFrom(depth int, node duct.AstFrom) error {
	t.A = "ø"
	t.B = node.Type
	return nil
}

func (t *typer) OnEnterYield(depth int, node duct.AstYield) error {
	t.B = "ø"
	return nil
}
