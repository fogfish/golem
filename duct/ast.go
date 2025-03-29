//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package duct

// Abstract syntax tree (AST) of the computation defined by morphism
type Ast interface {
	Apply(depth int, v Visitor) error
}

// Visitor over abstract syntax tree
type Visitor interface {
	OnMorphism
	OnSeq
	OnMap
	OnFrom
	OnYield
}

// Visits root morphism
type OnMorphism interface {
	OnEnterMorphism(depth int, node AstSeq) error
	OnLeaveMorphism(depth int, node AstSeq) error
}

// Visist nested transformers
type OnSeq interface {
	OnEnterSeq(depth int, node AstSeq) error
	OnLeaveSeq(depth int, node AstSeq) error
}

type OnMap interface {
	OnEnterMap(depth int, node AstMap) error
	OnLeaveMap(depth int, node AstMap) error
}

type OnFrom interface {
	OnEnterFrom(depth int, node AstFrom) error
	OnLeaveFrom(depth int, node AstFrom) error
}

type OnYield interface {
	OnEnterYield(depth int, node AstYield) error
	OnLeaveYield(depth int, node AstYield) error
}

// Empty abstract syntax tree visitor
type AstVisitor struct{}

func (AstVisitor) OnEnterMorphism(depth int, node AstSeq) error { return nil }
func (AstVisitor) OnLeaveMorphism(depth int, node AstSeq) error { return nil }
func (AstVisitor) OnEnterSeq(depth int, node AstSeq) error      { return nil }
func (AstVisitor) OnLeaveSeq(depth int, node AstSeq) error      { return nil }
func (AstVisitor) OnEnterMap(depth int, node AstMap) error      { return nil }
func (AstVisitor) OnLeaveMap(depth int, node AstMap) error      { return nil }
func (AstVisitor) OnEnterFrom(depth int, node AstFrom) error    { return nil }
func (AstVisitor) OnLeaveFrom(depth int, node AstFrom) error    { return nil }
func (AstVisitor) OnEnterYield(depth int, node AstYield) error  { return nil }
func (AstVisitor) OnLeaveYield(depth int, node AstYield) error  { return nil }

//------------------------------------------------------------------------------

// AST element for the input to morphism ƒ: ø ⟼ A
type AstFrom struct {
	Type   string
	Source any
}

func (node AstFrom) Apply(depth int, v Visitor) error {
	if err := v.OnEnterFrom(depth, node); err != nil {
		return err
	}
	if err := v.OnLeaveFrom(depth, node); err != nil {
		return err
	}
	return nil
}

//------------------------------------------------------------------------------

// AST element for the output of morphism ƒ: A ⟼ ø
type AstYield struct {
	Type   string
	Target any
}

// Void type to emulate empty type ø
type Void any

func (node AstYield) Apply(depth int, v Visitor) error {
	if err := v.OnEnterYield(depth, node); err != nil {
		return err
	}
	if err := v.OnLeaveYield(depth, node); err != nil {
		return err
	}
	return nil
}

//------------------------------------------------------------------------------

// AST element representing transformation ƒ: A ⟼ B
type AstMap struct {
	TypeA, TypeB string
	F            any
}

func (node AstMap) Apply(depth int, v Visitor) error {
	if err := v.OnEnterMap(depth, node); err != nil {
		return err
	}
	if err := v.OnLeaveMap(depth, node); err != nil {
		return err
	}
	return nil
}

//------------------------------------------------------------------------------

// AST element representing morphism, sequence of transformations 𝑚: A ⟼ B ⟼ ... ⟼ F
type AstSeq struct {
	Root     bool
	Deferred bool
	Seq      []Ast
}

func (n AstSeq) Apply(depth int, v Visitor) error {
	if n.Root {
		if err := v.OnEnterMorphism(depth, n); err != nil {
			return err
		}
	} else {
		if err := v.OnEnterSeq(depth, n); err != nil {
			return err
		}
	}

	for _, x := range n.Seq {
		if err := x.Apply(depth+1, v); err != nil {
			return err
		}
	}

	if n.Root {
		if err := v.OnLeaveMorphism(depth, n); err != nil {
			return err
		}
	} else {
		if err := v.OnLeaveSeq(depth, n); err != nil {
			return err
		}
	}

	return nil
}

func (f *AstSeq) unit() bool {
	if !f.Deferred {
		return false
	}

	if len(f.Seq) == 0 {
		if !f.Root {
			f.Deferred = false
		}
		return true
	}

	switch v := f.Seq[len(f.Seq)-1].(type) {
	case *AstSeq:
		if ok := v.unit(); ok {
			return true
		}
	}

	if !f.Root {
		f.Deferred = false
	}
	return true
}

func (f *AstSeq) append(n Ast) bool {
	if !f.Deferred {
		return false
	}

	if len(f.Seq) == 0 {
		f.Seq = append(f.Seq, n)
		return true
	}

	switch v := f.Seq[len(f.Seq)-1].(type) {
	case *AstSeq:
		if ok := v.append(n); ok {
			return true
		}
	}

	f.Seq = append(f.Seq, n)
	return true
}
