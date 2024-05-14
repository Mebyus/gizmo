package tt

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

type ChainOperand interface {
	Operand

	ChainOperand()

	// Depth of chain operand. Starts from zero for the first operand in chain.
	Depth() uint32
}

type nodeChainOperand struct{ nodeOperand }

func (nodeChainOperand) ChainOperand() {}

type ChainStart struct {
	nodeChainOperand

	Pos source.Pos

	// Symbol which is referenced in chain start.
	Sym *Symbol
}

// Explicit interface implementation check
var _ ChainOperand = &ChainStart{}

func (*ChainStart) Kind() exn.Kind {
	return exn.Start
}

func (s *ChainStart) Pin() source.Pos {
	return s.Pos
}

func (s *ChainStart) Type() *Type {
	return s.Sym.Type
}

func (s *ChainStart) Depth() uint32 {
	return 0
}

type IndirectExpression struct {
	nodeChainOperand

	Pos source.Pos

	Target ChainOperand

	typ *Type

	ChainDepth uint32
}

// Explicit interface implementation check
var _ ChainOperand = &IndirectExpression{}

func (*IndirectExpression) Kind() exn.Kind {
	return exn.Indirect
}

func (e *IndirectExpression) Pin() source.Pos {
	return e.Pos
}

func (e *IndirectExpression) Type() *Type {
	return e.typ
}

func (e *IndirectExpression) Depth() uint32 {
	return e.ChainDepth
}

type SymbolAddressExpression struct {
	nodeChainOperand

	Pos source.Pos

	Target *Symbol

	typ *Type
}

// Explicit interface implementation check
var _ ChainOperand = &IndirectExpression{}

func (*SymbolAddressExpression) Kind() exn.Kind {
	return exn.SymbolAddress
}

func (e *SymbolAddressExpression) Pin() source.Pos {
	return e.Pos
}

func (e *SymbolAddressExpression) Type() *Type {
	if e.typ != nil {
		return e.typ
	}

	ref := e.Target.Type
	scope := e.Target.Scope
	e.typ = scope.Types.storePointer(ref)
	return e.typ
}

func (e *SymbolAddressExpression) Depth() uint32 {
	return 0
}

type CallExpression struct {
	nodeChainOperand

	Pos source.Pos

	Callee    ChainOperand
	Arguments []Expression

	typ *Type

	ChainDepth uint32
}

// Explicit interface implementation check
var _ ChainOperand = &CallExpression{}

func (*CallExpression) Kind() exn.Kind {
	return exn.Call
}

func (e *CallExpression) Pin() source.Pos {
	return e.Pos
}

func (e *CallExpression) Depth() uint32 {
	return e.ChainDepth
}

func (e *CallExpression) Type() *Type {
	return e.typ
}

type SymbolCallExpression struct {
	nodeChainOperand

	Pos source.Pos

	Arguments []Expression
	Callee    *Symbol
}

// Explicit interface implementation check
var _ ChainOperand = &SymbolCallExpression{}

func (*SymbolCallExpression) Kind() exn.Kind {
	return exn.Call
}

func (e *SymbolCallExpression) Pin() source.Pos {
	return e.Pos
}

func (e *SymbolCallExpression) Depth() uint32 {
	return 0
}

func (e *SymbolCallExpression) Type() *Type {
	return e.Callee.Def.(*FnDef).Result
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
