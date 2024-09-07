package ast

import (
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/source"
)

// <ChainOperand> = <Receiver> | <Identifier> | <CallExpression> | <SelectorExpression> | <IndexExpression> | <IndirectExpression>
type ChainOperand struct {
	nodeOperand

	// Start of chain operand. Field Identifier.Lit can be empty in this
	// context, if this is the case, then it is receiver, not identifier.
	Identifier Identifier

	// Parts are arranged from left to right as in source text.
	Parts []ChainPart
}

var _ Operand = ChainOperand{}

func (o ChainOperand) Pin() source.Pos {
	return o.Identifier.Pos
}

func (o ChainOperand) Kind() exn.Kind {
	return exn.Chain
}

func (o ChainOperand) Last() exn.Kind {
	if len(o.Parts) == 0 {
		return exn.Symbol
	}
	return o.Parts[len(o.Parts)-1].Kind()
}

type ChainPart interface {
	Node

	ChainPart()

	Kind() exn.Kind
}

// Dummy chain part node provides quick, easy to use implementation
// of discriminator ChainPart() method.
//
// Used for embedding into other (non-dummy) chain part nodes.
type nodeChainPart struct{}

func (nodeChainPart) ChainPart() {}

// <MemberPart> = "." <Member>
//
// <Member> = <Identifier>
type MemberPart struct {
	nodeChainPart

	Member Identifier
}

var _ ChainPart = MemberPart{}

func (MemberPart) Kind() exn.Kind {
	return exn.Member
}

func (p MemberPart) Pin() source.Pos {
	return p.Member.Pos
}

// <CallPart> = "(" { <Expression> "," } ")"
type CallPart struct {
	nodeChainPart

	Pos source.Pos

	Args []Expression
}

var _ ChainPart = CallPart{}

func (CallPart) Kind() exn.Kind {
	return exn.Call
}

func (p CallPart) Pin() source.Pos {
	return p.Pos
}

// <IndexPart> = "[" <Expression> "]"
type IndexPart struct {
	nodeChainPart

	Pos   source.Pos
	Index Expression
}

var _ ChainPart = IndexPart{}

func (IndexPart) Kind() exn.Kind {
	return exn.Index
}

func (p IndexPart) Pin() source.Pos {
	return p.Pos
}

// <IndirectPart> = ".@"
type IndirectPart struct {
	nodeChainPart

	Pos source.Pos
}

var _ ChainPart = IndirectPart{}

func (IndirectPart) Kind() exn.Kind {
	return exn.Indirect
}

func (p IndirectPart) Pin() source.Pos {
	return p.Pos
}

// <AddressPart> = ".&"
type AddressPart struct {
	nodeChainPart

	Pos source.Pos
}

var _ ChainPart = AddressPart{}

func (AddressPart) Kind() exn.Kind {
	return exn.Address
}

func (p AddressPart) Pin() source.Pos {
	return p.Pos
}

// <IndirectIndexPart> = ".[" <Index> "]"
//
// <Index> = <Expression>
type IndirectIndexPart struct {
	nodeChainPart

	Pos source.Pos

	Index Expression
}

var _ ChainPart = IndirectIndexPart{}

func (IndirectIndexPart) Kind() exn.Kind {
	return exn.IndirectIndex
}

func (p IndirectIndexPart) Pin() source.Pos {
	return p.Pos
}

// <SlicePart> = "[" [ <Start> ] ":" [ <End> ] "]"
type SlicePart struct {
	nodeChainPart

	Pos source.Pos

	// Part before ":". Can be nil if expression is omitted
	Start Expression

	// Part after ":". Can be nil if expression is omitted
	End Expression
}

var _ ChainPart = SlicePart{}

func (SlicePart) Kind() exn.Kind {
	return exn.Slice
}

func (p SlicePart) Pin() source.Pos {
	return p.Pos
}
