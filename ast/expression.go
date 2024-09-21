package ast

import (
	"github.com/mebyus/gizmo/ast/bop"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/ast/uop"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

// <Exp> = <PrimaryOperand> | <BinaryExpression>
type Exp interface {
	Node

	// dummy discriminator method
	Expression()

	Kind() exn.Kind
}

type NodeE struct{}

func (NodeE) Expression() {}

// <PrimaryOperand> = <Operand> | <UnaryExpression>
type PrimaryOperand any

// <Operand> = <Literal> | <SubsExpression> | <ParenthesizedExpression> | <SelectorExpression> |
// <IndexExpression> | <CallExpression> | <AddressExpression> | <CastExpression>
type Operand interface {
	Exp

	// dummy discriminator method
	Operand()
}

// Dummy operand node provides quick, easy to use implementation of discriminator Operand() method
//
// Used for embedding into other (non-dummy) operand nodes
type NodeO struct{ NodeE }

func (NodeO) Operand() {}

// <SymbolExp> = <Identifier>
type SymbolExp struct {
	NodeO

	Identifier Identifier
}

var _ Operand = SymbolExp{}

func (SymbolExp) Kind() exn.Kind {
	return exn.Symbol
}

func (e SymbolExp) Pin() source.Pos {
	return e.Identifier.Pos
}

type IncompNameExp struct {
	NodeO

	Identifier Identifier
}

var _ Operand = IncompNameExp{}

func (IncompNameExp) Kind() exn.Kind {
	return exn.IncompName
}

func (e IncompNameExp) Pin() source.Pos {
	return e.Identifier.Pos
}

// <ParenthesizedExpression> = "(" <Expression> ")"
type ParenthesizedExpression struct {
	NodeO

	Pos source.Pos

	Inner Exp
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
	NodeE

	Operator UnaryOperator
	Inner    Exp
}

var _ Exp = UnaryExpression{}

func (UnaryExpression) Kind() exn.Kind {
	return exn.Unary
}

func (e UnaryExpression) Pin() source.Pos {
	return e.Operator.Pos
}

// <BinExp> = <Expression> <BinaryOperator> <Expression>
type BinExp struct {
	NodeE

	Operator BinaryOperator
	Left     Exp
	Right    Exp
}

var _ Exp = BinExp{}

func (BinExp) Kind() exn.Kind {
	return exn.Binary
}

func (e BinExp) Pin() source.Pos {
	return e.Left.Pin()
}

// <CastExp> = "cast" "(" <TypeSpec> "," <Exp> ")"
type CastExp struct {
	NodeO

	Pos    source.Pos
	Target Exp
	Type   TypeSpec
}

// Explicit interface implementation check.
var _ Exp = CastExp{}

func (CastExp) Kind() exn.Kind {
	return exn.Cast
}

func (e CastExp) Pin() source.Pos {
	return e.Pos
}

// <TintExp> = "tint" "(" <TypeSpec> "," <Exp> ")"
type TintExp struct {
	NodeO

	Pos    source.Pos
	Target Exp
	Type   TypeSpec
}

var _ Exp = TintExp{}

func (TintExp) Kind() exn.Kind {
	return exn.Tint
}

func (e TintExp) Pin() source.Pos {
	return e.Pos
}

// <MemCastExpression> = "mcast" "(" <TypeSpecifier> "," <Expression> ")"
type MemCastExpression struct {
	NodeO

	Target Exp
	Type   TypeSpec
}

var _ Exp = MemCastExpression{}

func (MemCastExpression) Kind() exn.Kind {
	return exn.MemCast
}

func (e MemCastExpression) Pin() source.Pos {
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
