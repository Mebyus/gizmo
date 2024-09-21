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
type NodeS struct{}

func (NodeS) Statement() {}

type VarStatement struct {
	NodeS

	// Symbol created by this statement.
	Sym *Symbol

	// Equals nil if init expression is dirty.
	Exp Exp
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
	NodeS

	// Symbol created by this statement.
	Sym *Symbol

	// Init expression for created symbol. Always not nil.
	Expr Exp
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
	NodeS

	Pos source.Pos

	// Equals nil if return does not have expression.
	Expr Exp
}

// Explicit interface implementation check
var _ Statement = &ReturnStatement{}

func (s *ReturnStatement) Pin() source.Pos {
	return s.Pos
}

func (s *ReturnStatement) Kind() stm.Kind {
	return stm.Return
}

type NeverStatement struct {
	NodeS

	Pos source.Pos
}

// Explicit interface implementation check
var _ Statement = &NeverStatement{}

func (s *NeverStatement) Pin() source.Pos {
	return s.Pos
}

func (s *NeverStatement) Kind() stm.Kind {
	return stm.Never
}

type AssignStatement struct {
	NodeS

	// Target of the assignment.
	Target ChainOperand

	// Assigned expression. Always not nil.
	Expr Exp

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
	NodeS

	Pos source.Pos

	// Always not nil.
	Condition Exp

	Body Block
}

var _ Statement = &SimpleIfStatement{}

func (*SimpleIfStatement) Kind() stm.Kind {
	return stm.SimpleIf
}

func (s *SimpleIfStatement) Pin() source.Pos {
	return s.Pos
}

type MatchStatement struct {
	NodeS

	Pos source.Pos

	Exp Exp

	Cases []MatchCase

	Else *Block
}

type MatchCase struct {
	Pos source.Pos

	ExpList []Exp

	Body Block
}

var _ Statement = &MatchStatement{}

func (*MatchStatement) Kind() stm.Kind {
	return stm.Match
}

func (s *MatchStatement) Pin() source.Pos {
	return s.Pos
}

type CallStatement struct {
	NodeS

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

type DeferStatement struct {
	NodeS

	Pos source.Pos

	Args []Exp

	// defer index (among all defers inside the function)
	Index uint32

	Uncertain bool
}

var _ Statement = &DeferStatement{}

func (*DeferStatement) Kind() stm.Kind {
	return stm.Defer
}

func (s *DeferStatement) Pin() source.Pos {
	return s.Pos
}
