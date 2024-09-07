package stg

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/source"
)

type Expression interface {
	Node

	// dummy discriminator method
	Expression()

	Kind() exn.Kind

	Type() *Type
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

func (e *SymbolExpression) Type() *Type {
	return e.Sym.Type
}

// EnumExp direct usage of enum entry as value.
type EnumExp struct {
	nodeOperand

	Pos source.Pos

	Enum *Type

	Entry *EnumEntry
}

// Explicit interface implementation check
var _ Operand = &EnumExp{}

func (*EnumExp) Kind() exn.Kind {
	return exn.Enum
}

func (e *EnumExp) Pin() source.Pos {
	return e.Pos
}

func (e *EnumExp) Type() *Type {
	return e.Enum
}

type UnaryOperator ast.UnaryOperator

type UnaryExpression struct {
	nodeExpression

	Operator UnaryOperator
	Inner    Expression

	typ *Type
}

// Explicit interface implementation check
var _ Expression = &UnaryExpression{}

func (*UnaryExpression) Kind() exn.Kind {
	return exn.Unary
}

func (e *UnaryExpression) Pin() source.Pos {
	return e.Operator.Pos
}

func (e *UnaryExpression) Type() *Type {
	if e.typ != nil {
		return e.typ
	}
	e.typ = e.Inner.Type()
	return e.typ
}

type BinaryOperator ast.BinaryOperator

type BinaryExpression struct {
	nodeExpression

	Operator BinaryOperator
	Left     Expression
	Right    Expression

	typ *Type
}

// Explicit interface implementation check
var _ Expression = &BinaryExpression{}

func (*BinaryExpression) Kind() exn.Kind {
	return exn.Binary
}

func (e *BinaryExpression) Pin() source.Pos {
	return e.Left.Pin()
}

func (e *BinaryExpression) Type() *Type {
	if e.typ != nil {
		return e.typ
	}
	e.typ = e.Left.Type()
	return e.typ
}

type ParenthesizedExpression struct {
	nodeOperand

	Pos source.Pos

	Inner Expression
}

// Explicit interface implementation check
var _ Operand = &ParenthesizedExpression{}

func (*ParenthesizedExpression) Kind() exn.Kind {
	return exn.Paren
}

func (e *ParenthesizedExpression) Pin() source.Pos {
	return e.Pos
}

func (e *ParenthesizedExpression) Type() *Type {
	return e.Inner.Type()
}

type CastExpression struct {
	nodeOperand

	Pos source.Pos

	Target Expression

	// Type after cast is performed.
	DestType *Type
}

// Explicit interface implementation check
var _ Operand = &CastExpression{}

func (*CastExpression) Kind() exn.Kind {
	return exn.Cast
}

func (e *CastExpression) Pin() source.Pos {
	return e.Pos
}

func (e *CastExpression) Type() *Type {
	return e.DestType
}

type TintExp struct {
	nodeOperand

	Pos source.Pos

	Target Expression

	// Type after cast is performed.
	DestType *Type
}

// Explicit interface implementation check
var _ Operand = &TintExp{}

func (*TintExp) Kind() exn.Kind {
	return exn.Tint
}

func (e *TintExp) Pin() source.Pos {
	return e.Pos
}

func (e *TintExp) Type() *Type {
	return e.DestType
}
