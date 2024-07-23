package ast

import (
	"github.com/mebyus/gizmo/ast/bop"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/ast/uop"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

// <Expression> = <PrimaryOperand> | <BinaryExpression>
type Expression interface {
	Node

	// dummy discriminator method
	Expression()

	Kind() exn.Kind
}

type nodeExpression struct{}

func (nodeExpression) Expression() {}

// <PrimaryOperand> = <Operand> | <UnaryExpression>
type PrimaryOperand any

// <Operand> = <Literal> | <SubsExpression> | <ParenthesizedExpression> | <SelectorExpression> |
// <IndexExpression> | <CallExpression> | <AddressExpression> | <CastExpression>
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

// <SymbolExpression> = <Identifier>
type SymbolExpression struct {
	nodeOperand

	Identifier Identifier
}

var _ Operand = SymbolExpression{}

func (SymbolExpression) Kind() exn.Kind {
	return exn.Symbol
}

func (e SymbolExpression) Pin() source.Pos {
	return e.Identifier.Pos
}

// <Receiver> = "rv"
type Receiver struct {
	nodeOperand

	Pos source.Pos
}

var _ Operand = Receiver{}

func (Receiver) Kind() exn.Kind {
	return exn.Receiver
}

func (r Receiver) Pin() source.Pos {
	return r.Pos
}

func (r Receiver) AsIdentifier() Identifier {
	return Identifier{Pos: r.Pos}
}

func (r Receiver) AsChainOperand() ChainOperand {
	return ChainOperand{Identifier: r.AsIdentifier()}
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

	Operator UnaryOperator
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

	Operator BinaryOperator
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

// <CastExpression> = "cast" "[" <Expression> ":" <TypeSpecifier> "]"
type CastExpression struct {
	nodeOperand

	Target Expression
	Type   TypeSpec
}

var _ Expression = CastExpression{}

func (CastExpression) Kind() exn.Kind {
	return exn.Cast
}

func (e CastExpression) Pin() source.Pos {
	return e.Target.Pin()
}

// <BitCastExpression> = "bitcast" "[" <Expression> ":" <TypeSpecifier> "]"
type BitCastExpression struct {
	nodeOperand

	Target Expression
	Type   TypeSpec
}

var _ Expression = BitCastExpression{}

func (BitCastExpression) Kind() exn.Kind {
	return exn.BitCast
}

func (e BitCastExpression) Pin() source.Pos {
	return e.Target.Pin()
}

type UnaryOperator struct {
	Pos  source.Pos
	Kind uop.Kind
}

func UnaryOperatorFromToken(tok token.Token) UnaryOperator {
	return UnaryOperator{
		Pos:  tok.Pos,
		Kind: uop.FromToken(tok.Kind),
	}
}

type BinaryOperator struct {
	Pos  source.Pos
	Kind bop.Kind
}

func BinaryOperatorFromToken(tok token.Token) BinaryOperator {
	return BinaryOperator{
		Pos:  tok.Pos,
		Kind: bop.FromToken(tok.Kind),
	}
}

func (o BinaryOperator) Precedence() int {
	return o.Kind.Precedence()
}

func (o BinaryOperator) Power() int {
	return o.Kind.Power()
}
