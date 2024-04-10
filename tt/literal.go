package tt

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

// Integer represents integer literal usage in source code.
type Integer struct {
	nodeLiteral

	Pos source.Pos

	Val uint64
}

// Explicit interface implementation check
var _ Literal = Integer{}

func (Integer) Kind() exn.Kind {
	return exn.Integer
}

func (n Integer) Pin() source.Pos {
	return n.Pos
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
