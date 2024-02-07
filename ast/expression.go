package ast

import (
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/ast/oper"
	"github.com/mebyus/gizmo/source"
)

// <Expression> = <PrimaryOperand> | <BinaryExpression>
type Expression interface {
	Node

	// dummy discriminator method
	Expression()

	Kind() exn.Kind
}

type nodeExpression struct{ uidHolder }

func (nodeExpression) Expression() {}

// <PrimaryOperand> = <Operand> | <UnaryExpression>
type PrimaryOperand any

// <Operand> = <Literal> | <SubsExpression> | <ParenthesizedExpression> | <SelectorExpression> |
// <IndexExpression> | <CallExpression>
type Operand interface {
	Expression

	// dummy discriminator method
	Operand()
}

// Dummy operand node provides quick, easy to use implementation of discriminator Operand() method
//
// Used for embedding into other (non-dummy) operand nodes
type nodeOperand struct{ nodeExpression }

func (nodeOperand) Operand() {}

// <SubsExpression> = <ScopedIdentifier>
type SubsExpression struct {
	nodeOperand

	Identifier ScopedIdentifier
}

var _ Operand = SubsExpression{}

func (SubsExpression) Kind() exn.Kind {
	return exn.Subs
}

func (e SubsExpression) Pin() source.Pos {
	return e.Identifier.Pos()
}

// <ParenthesizedExpression> = "(" <Expression> ")"
type ParenthesizedExpression struct {
	nodeOperand

	Pos source.Pos

	Inner Expression
}

var _ Operand = ParenthesizedExpression{}

func (ParenthesizedExpression) Kind() exn.Kind {
	return exn.Paren
}

func (e ParenthesizedExpression) Pin() source.Pos {
	return e.Pos
}

// <UnaryExpression> = <UnaryOperator> <UnaryOperand>
//
// <UnaryOperand> = <Operand> | <UnaryExpression>
type UnaryExpression struct {
	nodeExpression

	Operator oper.Unary
	Inner    Expression
}

var _ Expression = UnaryExpression{}

func (UnaryExpression) Kind() exn.Kind {
	return exn.Unary
}

func (e UnaryExpression) Pin() source.Pos {
	return e.Operator.Pos
}

// <BinaryExpression> = <Expression> <BinaryOperator> <Expression>
type BinaryExpression struct {
	nodeExpression

	Operator oper.Binary
	Left     Expression
	Right    Expression
}

var _ Expression = BinaryExpression{}

func (BinaryExpression) Kind() exn.Kind {
	return exn.Binary
}

func (e BinaryExpression) Pin() source.Pos {
	return e.Left.Pin()
}

// <ChainOperand> = <ChainStart> | <CallExpression> | <SelectorExpression> | <IndexExpression>
type ChainOperand interface {
	Operand

	ChainOperand()
}

type nodeChainOperand struct{ nodeOperand }

func (nodeChainOperand) ChainOperand() {}

// <ChainStart> = <ScopedIdentifier>
type ChainStart struct {
	nodeChainOperand

	Identifier ScopedIdentifier
}

var _ Operand = SubsExpression{}

func (ChainStart) Kind() exn.Kind {
	return exn.Start
}

func (s ChainStart) Pin() source.Pos {
	return s.Identifier.Pos()
}

// <CallExpression> = <CallableExpression> <TupleLiteral>
type CallExpression struct {
	Callee    ChainOperand
	Arguments []Expression
}

// <SelectorExpression> = <SelectableExpression> "." <Identifier>
type SelectorExpression struct {
	Target   ChainOperand
	Selected Identifier
}

// <IndexExpression> = <IndexableExpression> "[" <Expression> "]"
type IndexExpression struct {
	Target ChainOperand
	Index  Expression
}

// <SliceExpression> = <Target> "[" [ <Start> ] ":" [ <End> ] "]"
type SliceExpression struct {
	Target IndexableExpression

	// Part before ":". Can be nil if expression is omitted
	Start Expression

	// Part after ":". Can be nil if expression is omitted
	End Expression
}

// <IndexableExpression> = <Identifier> | <SelectorExpression> | <IndexExpression>
type IndexableExpression any
