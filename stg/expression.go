package stg

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/source"
)

type Exp interface {
	Node

	// dummy discriminator method
	Expression()

	Kind() exk.Kind

	Type() *Type
}

// This is dummy implementation of Expression interface.
//
// Used for embedding into other (non-dummy) expression nodes.
type nodeExpression struct{}

func (nodeExpression) Expression() {}

func (nodeExpression) Kind() exk.Kind { return 0 }

type Operand interface {
	Exp

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

func (*SymbolExpression) Kind() exk.Kind {
	return exk.Symbol
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

func (*EnumExp) Kind() exk.Kind {
	return exk.Enum
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
	Inner    Exp

	typ *Type
}

// Explicit interface implementation check
var _ Exp = &UnaryExpression{}

func (*UnaryExpression) Kind() exk.Kind {
	return exk.Unary
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

type BinExp struct {
	nodeExpression

	Operator BinaryOperator
	Left     Exp
	Right    Exp

	typ *Type
}

// Explicit interface implementation check
var _ Exp = &BinExp{}

func (*BinExp) Kind() exk.Kind {
	return exk.Binary
}

func (e *BinExp) Pin() source.Pos {
	return e.Left.Pin()
}

func (e *BinExp) Type() *Type {
	if e.typ != nil {
		return e.typ
	}
	e.typ = e.Left.Type()
	return e.typ
}

type ParenthesizedExpression struct {
	nodeOperand

	Pos source.Pos

	Inner Exp
}

// Explicit interface implementation check
var _ Operand = &ParenthesizedExpression{}

func (*ParenthesizedExpression) Kind() exk.Kind {
	return exk.Paren
}

func (e *ParenthesizedExpression) Pin() source.Pos {
	return e.Pos
}

func (e *ParenthesizedExpression) Type() *Type {
	return e.Inner.Type()
}

type CastExp struct {
	nodeOperand

	Pos source.Pos

	Target Exp

	// Type after cast is performed.
	DestType *Type
}

// Explicit interface implementation check
var _ Operand = &CastExp{}

func (*CastExp) Kind() exk.Kind {
	return exk.Cast
}

func (e *CastExp) Pin() source.Pos {
	return e.Pos
}

func (e *CastExp) Type() *Type {
	return e.DestType
}

type TintExp struct {
	nodeOperand

	Pos source.Pos

	Target Exp

	// Type after cast is performed.
	DestType *Type
}

// Explicit interface implementation check
var _ Operand = &TintExp{}

func (*TintExp) Kind() exk.Kind {
	return exk.Tint
}

func (e *TintExp) Pin() source.Pos {
	return e.Pos
}

func (e *TintExp) Type() *Type {
	return e.DestType
}

type MemCastExp struct {
	nodeOperand

	Pos source.Pos

	Target Exp

	// Type after cast is performed.
	DestType *Type
}

// Explicit interface implementation check
var _ Operand = &MemCastExp{}

func (*MemCastExp) Kind() exk.Kind {
	return exk.MemCast
}

func (e *MemCastExp) Pin() source.Pos {
	return e.Pos
}

func (e *MemCastExp) Type() *Type {
	return e.DestType
}
