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
type Operand any

// <SubsExpression> = <ScopedIdentifier>
type SubsExpression struct {
	nodeExpression

	Identifier ScopedIdentifier
}

var _ Expression = SubsExpression{}

func (SubsExpression) Kind() exn.Kind {
	return exn.Subs
}

func (e SubsExpression) Pin() source.Pos {
	return e.Identifier.Pos()
}

// <ParenthesizedExpression> = "(" <Expression> ")"
type ParenthesizedExpression struct {
	nodeExpression

	Pos source.Pos

	Inner Expression
}

var _ Expression = ParenthesizedExpression{}

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
	Operand  Expression
}

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

// <CompoundOperand> = <Identifier> | <CallExpression> | <SelectorExpression> | <IndexExpression>
type CompoundOperand any

// <CallExpression> = <CallableExpression> <TupleLiteral>
type CallExpression struct {
	Callee    CallableExpression
	Arguments []Expression
}

// <CallableExpression> = <Identifier> | <SelectorExpression>
type CallableExpression any

// <SelectorExpression> = <SelectableExpression> "." <Identifier>
type SelectorExpression struct {
	Target   SelectableExpression
	Selected Identifier
}

// <SelectableExpression> = <Identifier> | <SelectorExpression> | <IndexExpression>
type SelectableExpression any

// <IndexExpression> = <IndexableExpression> "[" <Expression> "]"
type IndexExpression struct {
	Target IndexableExpression
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
