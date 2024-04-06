package tt

import (
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
)

type Node interface {
	source.Pin
}

type Statement interface {
	Node

	// dummy discriminator method
	Statement()

	Kind() stm.Kind
}

// Dummy provides quick, easy to use implementation of discriminator Statement() method
//
// Used for embedding into other (non-dummy) statement nodes
type nodeStatement struct{}

func (nodeStatement) Statement() {}

type VarStatement struct {
	nodeStatement

	// Symbol created by this statement.
	Sym *Symbol

	// Equals nil if init expression is dirty.
	Expr any
}

// Explicit interface implementation check
var _ Statement = &VarStatement{}

func (s *VarStatement) Pin() source.Pos {
	return s.Sym.Pos
}

func (s *VarStatement) Kind() stm.Kind {
	return stm.Var
}
