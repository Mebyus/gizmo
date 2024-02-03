package ast

import (
	"github.com/mebyus/gizmo/ast/oper"
	"github.com/mebyus/gizmo/token"
)

// <Expression> = <PrimaryOperand> | <BinaryExpression>
type Expression any

// <PrimaryOperand> = <Operand> | <UnaryExpression>
type PrimaryOperand any

// <Operand> = <Literal> | <Identifier> | <ParenthesizedExpression> | <SelectorExpression> |
// <IndexExpression> | <CallExpression>
type Operand any

// token.Kind is DecimalInteger, DecimalFloat, Character, String, Nil, True, False
type BasicLiteral token.Token

// <ParenthesizedExpression> = "(" <Expression> ")"
type ParenthesizedExpression struct {
	Inner Expression
}

// <UnaryExpression> = <UnaryOperator> <UnaryOperand>
type UnaryExpression struct {
	Operator     oper.Unary
	UnaryOperand UnaryOperand
}

// <UnaryOperand> = <Operand> | <UnaryExpression>
type UnaryOperand any

// <BinaryExpression> = <Expression> <BinaryOperator> <Expression>
type BinaryExpression struct {
	Operator  oper.Binary
	LeftSide  Expression
	RightSide Expression
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
