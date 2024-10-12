package ast

import (
	"github.com/mebyus/gizmo/ast/bop"
	"github.com/mebyus/gizmo/ast/uop"
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

// <Exp> = <PrimaryOperand> | <BinaryExpression>
type Exp interface {
	Node

	// dummy discriminator method
	Expression()

	Kind() exk.Kind
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

func (SymbolExp) Kind() exk.Kind {
	return exk.Symbol
}

func (e SymbolExp) Pin() source.Pos {
	return e.Identifier.Pos
}

type IncompNameExp struct {
	NodeO

	Identifier Identifier
}

var _ Operand = IncompNameExp{}

func (IncompNameExp) Kind() exk.Kind {
	return exk.IncompName
}

func (e IncompNameExp) Pin() source.Pos {
	return e.Identifier.Pos
}

// <ParenExp> = "(" <Expression> ")"
type ParenExp struct {
	NodeO

	Pos source.Pos

	Inner Exp
}

var _ Operand = ParenExp{}

func (ParenExp) Kind() exk.Kind {
	return exk.Paren
}

func (e ParenExp) Pin() source.Pos {
	return e.Pos
}

// <UnaryExp> = <UnaryOperator> <UnaryOperand>
//
// <UnaryOperand> = <Operand> | <UnaryExp>
type UnaryExp struct {
	NodeE

	Operator UnaryOperator
	Inner    Exp
}

var _ Exp = UnaryExp{}

func (UnaryExp) Kind() exk.Kind {
	return exk.Unary
}

func (e UnaryExp) Pin() source.Pos {
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

func (BinExp) Kind() exk.Kind {
	return exk.Binary
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

func (CastExp) Kind() exk.Kind {
	return exk.Cast
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

func (TintExp) Kind() exk.Kind {
	return exk.Tint
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

func (MemCastExpression) Kind() exk.Kind {
	return exk.MemCast
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
