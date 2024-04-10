package tt

import (
	"github.com/mebyus/gizmo/ast/bop"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/ast/uop"
	"github.com/mebyus/gizmo/source"
)

type Expression interface {
	Node

	// dummy discriminator method
	Expression()

	Kind() exn.Kind
}

// This is dummy implementation of Expression interface.
//
// Used for embedding into other (non-dummy) expression nodes.
type nodeExpression struct{}

func (nodeExpression) Expression() {}

func (nodeExpression) Kind() exn.Kind { return 0 }

type Operand interface {
	Expression

	// dummy discriminator method
	Operand()
}

// This is dummy implementation of Operand interface.
//
// Used for embedding into other (non-dummy) operand nodes.
type nodeOperand struct{ nodeExpression }

func (nodeOperand) Operand() {}

// SymbolExpression is an operand expression which represents direct symbol usage.
// Example:
//
//	10 + a // in this expression operand "a" is SymbolExpression
type SymbolExpression struct {
	nodeOperand

	Pos source.Pos

	// Symbol to which operand refers.
	Sym *Symbol
}

// Explicit interface implementation check
var _ Operand = &SymbolExpression{}

func (*SymbolExpression) Kind() exn.Kind {
	return exn.Symbol
}

func (e *SymbolExpression) Pin() source.Pos {
	return e.Pos
}

type UnaryOperator struct {
	Pos  source.Pos
	Kind uop.Kind
}

type UnaryExpression struct {
	nodeExpression

	Operator UnaryOperator
	Inner    Expression
}

// Explicit interface implementation check
var _ Expression = &UnaryExpression{}

func (*UnaryExpression) Kind() exn.Kind {
	return exn.Unary
}

func (e *UnaryExpression) Pin() source.Pos {
	return e.Operator.Pos
}

type BinaryOperator struct {
	Pos  source.Pos
	Kind bop.Kind
}

type BinaryExpression struct {
	nodeExpression

	Operator BinaryOperator
	Left     Expression
	Right    Expression
}

// Explicit interface implementation check
var _ Expression = &BinaryExpression{}

func (*BinaryExpression) Kind() exn.Kind {
	return exn.Binary
}

func (e *BinaryExpression) Pin() source.Pos {
	return e.Left.Pin()
}
