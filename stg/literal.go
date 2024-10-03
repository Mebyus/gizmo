package stg

import (
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/source"
)

type Dirty struct {
	NodeE

	Pos source.Pos
}

// Explicit interface implementation check
var _ Exp = Dirty{}

func (Dirty) Kind() exk.Kind {
	return exk.Dirty
}

func (d Dirty) Pin() source.Pos {
	return d.Pos
}

func (Dirty) Type() *Type {
	return nil
}

type Literal interface {
	Operand

	// dummy discriminator method
	Literal()
}

// This is dummy implementation of Literal interface.
//
// Used for embedding into other (non-dummy) operand nodes.
type nodeLiteral struct{ NodeO }

func (nodeLiteral) Literal() {}

// True represents "true" literal usage in source code.
type True struct {
	nodeLiteral

	Pos source.Pos
}

// Explicit interface implementation check
var _ Literal = True{}

func (True) Kind() exk.Kind {
	return exk.True
}

func (t True) Pin() source.Pos {
	return t.Pos
}

func (True) Type() *Type {
	return StaticBooleanType
}

// False represents "false" literal usage in source code.
type False struct {
	nodeLiteral

	Pos source.Pos
}

// Explicit interface implementation check
var _ Literal = False{}

func (False) Kind() exk.Kind {
	return exk.False
}

func (f False) Pin() source.Pos {
	return f.Pos
}

func (False) Type() *Type {
	return StaticBooleanType
}

// Integer represents integer literal usage
// or statically evaluated integer.
type Integer struct {
	nodeLiteral

	Pos source.Pos

	Val uint64

	typ *Type

	// True if integer is negative.
	// Only evaluated integer can be negative.
	Neg bool
}

// Explicit interface implementation check
var _ Literal = Integer{}

func (Integer) Kind() exk.Kind {
	return exk.Integer
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

	// Text inside string literal in escaped form.
	Val string

	// Binary size of the string in unescaped form.
	Size uint64
}

// Explicit interface implementation check
var _ Literal = String{}

func (String) Kind() exk.Kind {
	return exk.String
}

func (s String) Pin() source.Pos {
	return s.Pos
}

func (String) Type() *Type {
	return StaticStringType
}
