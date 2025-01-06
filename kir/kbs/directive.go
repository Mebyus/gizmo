package kbs

import (
	"github.com/mebyus/gizmo/ast/bop"
	"github.com/mebyus/gizmo/source"
)

type Directive interface {
	Directive()
}

type NodeDirective struct{}

func (NodeDirective) Directive() {}

type Include struct {
	NodeDirective

	Pos    source.Pos
	String string
}

var _ Directive = Include{}

type Link struct {
	NodeDirective

	Pos    source.Pos
	String string
}

var _ Directive = Link{}

type Name struct {
	NodeDirective

	Pos    source.Pos
	String string
}

var _ Directive = Name{}

type If struct {
	NodeDirective

	Clause IfClause
	ElseIf []IfClause
	Else   *Block
}

type IfClause struct {
	Pos source.Pos

	// Always not nil
	Exp  Exp
	Body Block
}

type Block struct {
	Pos source.Pos

	Directives []Directive
}

type Exp interface {
	Exp()
}

type NodeExp struct{}

func (NodeExp) Exp() {}

type String struct {
	NodeExp

	Pos source.Pos

	Value string
}

var _ Exp = String{}

// Usage of env symbol in expression.
type EnvSymbol struct {
	NodeExp

	Pos source.Pos

	Name string
}

var _ Exp = EnvSymbol{}

type BinExp struct {
	NodeExp

	Left  Exp
	Right Exp

	Op bop.Kind
}

var _ Exp = BinExp{}
