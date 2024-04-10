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
	Expr Expression
}

// Explicit interface implementation check
var _ Statement = &VarStatement{}

func (s *VarStatement) Pin() source.Pos {
	return s.Sym.Pos
}

func (s *VarStatement) Kind() stm.Kind {
	return stm.Var
}

type ReturnStatement struct {
	nodeStatement

	Pos source.Pos

	// Equals nil if return does not have expression.
	Expr Expression
}

// Explicit interface implementation check
var _ Statement = &ReturnStatement{}

func (s *ReturnStatement) Pin() source.Pos {
	return s.Pos
}

func (s *ReturnStatement) Kind() stm.Kind {
	return stm.Return
}

// SymbolAssignStatement is a simple form of generic assign statement,
// where target is a symbol (not a complex expression). Example:
//
//	x = 10 + a;
type SymbolAssignStatement struct {
	nodeStatement

	// Position of assign statement.
	Pos source.Pos

	// Target of the assignment.
	Target *Symbol

	// Assigned expression. Always not nil.
	Expr Expression
}

// Explicit interface implementation check
var _ Statement = &SymbolAssignStatement{}

func (s *SymbolAssignStatement) Pin() source.Pos {
	return s.Pos
}

func (s *SymbolAssignStatement) Kind() stm.Kind {
	return stm.SymbolAssign
}

// IndirectAssignStatement is a simple form of generic assign statement,
// where target is an indirect on a symbol. Example:
//
//	x.@ = 10 + a;
type IndirectAssignStatement struct {
	nodeStatement

	// Position of assign statement.
	Pos source.Pos

	// Target of the assignment.
	Target *Symbol

	// Assigned expression. Always not nil.
	Expr Expression
}

// Explicit interface implementation check
var _ Statement = &IndirectAssignStatement{}

func (s *IndirectAssignStatement) Pin() source.Pos {
	return s.Pos
}

func (s *IndirectAssignStatement) Kind() stm.Kind {
	return stm.IndirectAssign
}
