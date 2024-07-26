package stg

import (
	"github.com/mebyus/gizmo/ast/aop"
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

// This is dummy implementation of Statement interface.
//
// Used for embedding into other (non-dummy) statement nodes.
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

type LetStatement struct {
	nodeStatement

	// Symbol created by this statement.
	Sym *Symbol

	// Init expression for created symbol. Always not nil.
	Expr Expression
}

// Explicit interface implementation check
var _ Statement = &LetStatement{}

func (s *LetStatement) Pin() source.Pos {
	return s.Sym.Pos
}

func (s *LetStatement) Kind() stm.Kind {
	return stm.Let
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

type AssignStatement struct {
	nodeStatement

	// Target of the assignment.
	Target ChainOperand

	// Assigned expression. Always not nil.
	Expr Expression

	Operation aop.Kind
}

// Explicit interface implementation check
var _ Statement = &AssignStatement{}

func (s *AssignStatement) Pin() source.Pos {
	return s.Target.Pin()
}

func (s *AssignStatement) Kind() stm.Kind {
	return stm.Assign
}

type SimpleIfStatement struct {
	nodeStatement

	Pos source.Pos

	// Always not nil.
	Condition Expression

	Body Block
}

var _ Statement = &SimpleIfStatement{}

func (*SimpleIfStatement) Kind() stm.Kind {
	return stm.SimpleIf
}

func (s *SimpleIfStatement) Pin() source.Pos {
	return s.Pos
}

type WhileStatement struct {
	nodeStatement

	Pos source.Pos

	// Always not nil.
	Condition Expression

	Body Block
}

var _ Statement = &WhileStatement{}

func (*WhileStatement) Kind() stm.Kind {
	return stm.ForCond
}

func (s *WhileStatement) Pin() source.Pos {
	return s.Pos
}

type LoopStatement struct {
	nodeStatement

	Pos source.Pos

	Body Block
}

var _ Statement = &LoopStatement{}

func (*LoopStatement) Kind() stm.Kind {
	return stm.For
}

func (s *LoopStatement) Pin() source.Pos {
	return s.Pos
}

type CallStatement struct {
	nodeStatement

	Pos source.Pos

	Call *CallExpression
}

var _ Statement = &CallStatement{}

func (*CallStatement) Kind() stm.Kind {
	return stm.Call
}

func (s *CallStatement) Pin() source.Pos {
	return s.Pos
}
