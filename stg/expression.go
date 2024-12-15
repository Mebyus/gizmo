package stg

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/enums/smk"
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
type NodeE struct{}

func (NodeE) Expression() {}

type Operand interface {
	Exp

	// dummy discriminator method
	Operand()

	// Returns true if expression results in a stored value.
	Stored() bool
}

// This is dummy implementation of Operand interface.
//
// Used for embedding into other (non-dummy) operand nodes.
type NodeO struct{ NodeE }

func (NodeO) Operand() {}

func (NodeO) Stored() bool {
	return false
}

// SymbolExp is an operand expression which represents direct symbol usage.
//
// Example:
//
//	10 + a // in this expression operand "a" is SymbolExp
type SymbolExp struct {
	NodeO

	Pos source.Pos

	// Symbol to which operand refers.
	Sym *Symbol
}

// Explicit interface implementation check.
var _ Operand = &SymbolExp{}

func (*SymbolExp) Kind() exk.Kind {
	return exk.Symbol
}

func (e *SymbolExp) Pin() source.Pos {
	return e.Pos
}

func (e *SymbolExp) Type() *Type {
	return e.Sym.Type
}

func (e SymbolExp) Stored() bool {
	return e.Sym.Kind == smk.Var
}

// EnumExp direct usage of enum entry as value.
type EnumExp struct {
	NodeO

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

type UnaryExp struct {
	NodeE

	Operator UnaryOperator
	Inner    Exp

	typ *Type
}

// Explicit interface implementation check
var _ Exp = &UnaryExp{}

func (*UnaryExp) Kind() exk.Kind {
	return exk.Unary
}

func (e *UnaryExp) Pin() source.Pos {
	return e.Operator.Pos
}

func (e *UnaryExp) Type() *Type {
	if e.typ != nil {
		return e.typ
	}
	e.typ = e.Inner.Type()
	return e.typ
}

type BinaryOperator ast.BinaryOperator

type BinExp struct {
	NodeE

	Operator BinaryOperator
	Left     Exp
	Right    Exp

	// stored resulting type of expression
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
	if e.typ == nil {
		panic("no type")
	}
	return e.typ
}

type ParenExp struct {
	NodeO

	Pos source.Pos

	Inner Exp
}

// Explicit interface implementation check
var _ Operand = &ParenExp{}

func (*ParenExp) Kind() exk.Kind {
	return exk.Paren
}

func (e *ParenExp) Pin() source.Pos {
	return e.Pos
}

func (e *ParenExp) Type() *Type {
	return e.Inner.Type()
}

type CastExp struct {
	NodeO

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
	NodeO

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
	NodeO

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
