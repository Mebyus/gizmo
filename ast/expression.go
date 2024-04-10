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

// <SymbolExpression> = <ScopedIdentifier>
type SymbolExpression struct {
	nodeOperand

	Identifier ScopedIdentifier
}

var _ Operand = SymbolExpression{}

func (SymbolExpression) Kind() exn.Kind {
	return exn.Symbol
}

func (e SymbolExpression) Pin() source.Pos {
	return e.Identifier.Pin()
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

// <ChainOperand> = <Receiver> | <ChainStart> | <CallExpression> | <SelectorExpression> | <IndexExpression> | <IndirectExpression>
type ChainOperand interface {
	Operand

	ChainOperand()

	// Depth of chain operand. Starts from zero for the first operand in chain.
	Depth() uint32
}

type nodeChainOperand struct{ nodeOperand }

func (nodeChainOperand) ChainOperand() {}

// <Receiver> = "rv"
type Receiver struct {
	nodeChainOperand

	Pos source.Pos
}

var _ ChainOperand = Receiver{}

func (Receiver) Kind() exn.Kind {
	return exn.Receiver
}

func (r Receiver) Pin() source.Pos {
	return r.Pos
}

func (r Receiver) Depth() uint32 {
	return 0
}

// <ChainStart> = <ScopedIdentifier>
type ChainStart struct {
	nodeChainOperand

	Identifier ScopedIdentifier
}

var _ ChainOperand = ChainStart{}

func (ChainStart) Kind() exn.Kind {
	return exn.Start
}

func (s ChainStart) Pin() source.Pos {
	return s.Identifier.Pin()
}

func (s ChainStart) Depth() uint32 {
	return 0
}

// <CallExpression> = <CallableExpression> "(" { <Expression> "," } ")"
type CallExpression struct {
	nodeChainOperand

	Callee    ChainOperand
	Arguments []Expression

	ChainDepth uint32
}

var _ ChainOperand = CallExpression{}

func (CallExpression) Kind() exn.Kind {
	return exn.Call
}

func (e CallExpression) Pin() source.Pos {
	return e.Callee.Pin()
}

func (e CallExpression) Depth() uint32 {
	return e.ChainDepth
}

// <SelectorExpression> = <SelectableExpression> "." <Identifier>
type SelectorExpression struct {
	nodeChainOperand

	Target   ChainOperand
	Selected Identifier

	ChainDepth uint32
}

var _ ChainOperand = SelectorExpression{}

func (SelectorExpression) Kind() exn.Kind {
	return exn.Select
}

func (e SelectorExpression) Pin() source.Pos {
	return e.Target.Pin()
}

func (e SelectorExpression) Depth() uint32 {
	return e.ChainDepth
}

// <IndexExpression> = <IndexableExpression> "[" <Expression> "]"
type IndexExpression struct {
	nodeChainOperand

	Target ChainOperand
	Index  Expression

	ChainDepth uint32
}

var _ ChainOperand = IndexExpression{}

func (IndexExpression) Kind() exn.Kind {
	return exn.Index
}

func (e IndexExpression) Pin() source.Pos {
	return e.Target.Pin()
}

func (e IndexExpression) Depth() uint32 {
	return e.ChainDepth
}

// <IndirectExpression> = <ChainOperand> ".@"
type IndirectExpression struct {
	nodeChainOperand

	Target ChainOperand

	ChainDepth uint32
}

var _ ChainOperand = IndirectExpression{}

func (IndirectExpression) Kind() exn.Kind {
	return exn.Indirect
}

func (e IndirectExpression) Pin() source.Pos {
	return e.Target.Pin()
}

func (e IndirectExpression) Depth() uint32 {
	return e.ChainDepth
}

// <AddressExpression> = <ChainOperand> ".&"
type AddressExpression struct {
	nodeChainOperand

	Target ChainOperand

	ChainDepth uint32
}

var _ ChainOperand = AddressExpression{}

func (AddressExpression) Kind() exn.Kind {
	return exn.Address
}

func (e AddressExpression) Pin() source.Pos {
	return e.Target.Pin()
}

func (e AddressExpression) Depth() uint32 {
	return e.ChainDepth
}

// <IndirectIndexExpression> = <Target> ".[" <Index> "]"
//
// <Target> = <ChainOperand>
//
// <Index> = <Expression>
type IndirectIndexExpression struct {
	nodeChainOperand

	Target ChainOperand
	Index  Expression

	ChainDepth uint32
}

var _ ChainOperand = IndirectIndexExpression{}

func (IndirectIndexExpression) Kind() exn.Kind {
	return exn.Indirx
}

func (e IndirectIndexExpression) Pin() source.Pos {
	return e.Target.Pin()
}

func (e IndirectIndexExpression) Depth() uint32 {
	return e.ChainDepth
}

// <SliceExpression> = <Target> "[" [ <Start> ] ":" [ <End> ] "]"
type SliceExpression struct {
	nodeChainOperand

	Target ChainOperand

	// Part before ":". Can be nil if expression is omitted
	Start Expression

	// Part after ":". Can be nil if expression is omitted
	End Expression

	ChainDepth uint32
}

var _ ChainOperand = SliceExpression{}

func (SliceExpression) Kind() exn.Kind {
	return exn.Slice
}

func (e SliceExpression) Pin() source.Pos {
	return e.Target.Pin()
}

func (e SliceExpression) Depth() uint32 {
	return e.ChainDepth
}

// <CastExpression> = "cast" "[" <Expression> ":" <TypeSpecifier> "]"
type CastExpression struct {
	nodeOperand

	Target Expression
	Type   TypeSpecifier
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
	Type   TypeSpecifier
}

var _ Expression = BitCastExpression{}

func (BitCastExpression) Kind() exn.Kind {
	return exn.BitCast
}

func (e BitCastExpression) Pin() source.Pos {
	return e.Target.Pin()
}

// <InstanceExpression> = <ScopedIdentifier> "[[" <Args> "]]"
type InstanceExpression struct {
	nodeChainOperand

	Target ScopedIdentifier

	// Always has at least one element
	Args []TypeSpecifier
}

var _ Expression = InstanceExpression{}

func (InstanceExpression) Kind() exn.Kind {
	return exn.Instance
}

func (e InstanceExpression) Pin() source.Pos {
	return e.Target.Name.Pos
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
