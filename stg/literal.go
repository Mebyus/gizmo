package stg

import (
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/source"
)

type Literal interface {
	Operand

	// dummy discriminator method
	Literal()
}

// This is dummy implementation of Literal interface.
//
// Used for embedding into other (non-dummy) operand nodes.
type nodeLiteral struct{ nodeOperand }

func (nodeLiteral) Literal() {}

// True represents "true" literal usage in source code.
type True struct {
	nodeLiteral

	Pos source.Pos
}

// Explicit interface implementation check
var _ Literal = True{}

func (True) Kind() exn.Kind {
	return exn.True
}

func (t True) Pin() source.Pos {
	return t.Pos
}

func (True) Type() *Type {
	return StaticBoolean
}

// False represents "false" literal usage in source code.
type False struct {
	nodeLiteral

	Pos source.Pos
}

// Explicit interface implementation check
var _ Literal = False{}

func (False) Kind() exn.Kind {
	return exn.False
}

func (f False) Pin() source.Pos {
	return f.Pos
}

func (False) Type() *Type {
	return StaticBoolean
}

// Integer represents integer literal usage
// or statically evaluated integer.
type Integer struct {
	nodeLiteral

	Pos source.Pos

	Val uint64

	typ *Type

	// True if integer if negative.
	// Only evaluated integer can be negative.
	Neg bool
}

// Explicit interface implementation check
var _ Literal = Integer{}

func (Integer) Kind() exn.Kind {
	return exn.Integer
}

func (n Integer) Pin() source.Pos {
	return n.Pos
}

func (n Integer) Type() *Type {
	return n.typ
}

// String represents string literal usage in source code.
type String struct {
	nodeLiteral

	Pos source.Pos

	Val string
}

// Explicit interface implementation check
var _ Literal = String{}

func (String) Kind() exn.Kind {
	return exn.String
}

func (s String) Pin() source.Pos {
	return s.Pos
}

func (String) Type() *Type {
	return StaticString
}
