package ast

import (
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/source"
)

// <ChainOperand> = <Identifier> | <CallExpression> | <SelectorExpression> | <IndexExpression> | <IndirectExpression>
type ChainOperand struct {
	NodeO

	// Start of chain operand. Field Identifier.Lit can be empty in this
	// context, if this is the case, then it is receiver, not identifier.
	Identifier Identifier

	// Parts are arranged from left to right as in source text.
	Parts []ChainPart
}

// Explicit interface implementation check.
var _ Operand = ChainOperand{}

func (o ChainOperand) Pin() source.Pos {
	return o.Identifier.Pos
}

func (o ChainOperand) Kind() exk.Kind {
	return exk.Chain
}

func (o ChainOperand) Last() exk.Kind {
	if len(o.Parts) == 0 {
		return exk.Symbol
	}
	return o.Parts[len(o.Parts)-1].Kind()
}

type ChainPart interface {
	Node

	ChainPart()

	Kind() exk.Kind
}

// Dummy chain part node provides quick, easy to use implementation
// of discriminator ChainPart() method.
//
// Used for embedding into other (non-dummy) chain part nodes.
type NodeP struct{}

func (NodeP) ChainPart() {}

// <SelectPart> = "." <Name>
//
// <Name> = <Identifier>
type SelectPart struct {
	NodeP

	Name Identifier
}

// Explicit interface implementation check.
var _ ChainPart = SelectPart{}

func (SelectPart) Kind() exk.Kind {
	return exk.Select
}

func (p SelectPart) Pin() source.Pos {
	return p.Name.Pos
}

// <CallPart> = "(" { <Expression> "," } ")"
type CallPart struct {
	NodeP

	Pos source.Pos

	Args []Exp
}

// Explicit interface implementation check.
var _ ChainPart = CallPart{}

func (CallPart) Kind() exk.Kind {
	return exk.Call
}

func (p CallPart) Pin() source.Pos {
	return p.Pos
}

// <IndexPart> = "[" <Expression> "]"
type IndexPart struct {
	NodeP

	Pos   source.Pos
	Index Exp
}

// Explicit interface implementation check.
var _ ChainPart = IndexPart{}

func (IndexPart) Kind() exk.Kind {
	return exk.Index
}

func (p IndexPart) Pin() source.Pos {
	return p.Pos
}

// <IndirectPart> = ".@"
type IndirectPart struct {
	NodeP

	Pos source.Pos
}

// Explicit interface implementation check.
var _ ChainPart = IndirectPart{}

func (IndirectPart) Kind() exk.Kind {
	return exk.Indirect
}

func (p IndirectPart) Pin() source.Pos {
	return p.Pos
}

// <AddressPart> = ".&"
type AddressPart struct {
	NodeP

	Pos source.Pos
}

// Explicit interface implementation check.
var _ ChainPart = AddressPart{}

func (AddressPart) Kind() exk.Kind {
	return exk.Address
}

func (p AddressPart) Pin() source.Pos {
	return p.Pos
}

// <IndirectIndexPart> = ".[" <Index> "]"
//
// <Index> = <Expression>
type IndirectIndexPart struct {
	NodeP

	Pos source.Pos

	Index Exp
}

// Explicit interface implementation check.
var _ ChainPart = IndirectIndexPart{}

func (IndirectIndexPart) Kind() exk.Kind {
	return exk.IndirectIndex
}

func (p IndirectIndexPart) Pin() source.Pos {
	return p.Pos
}

// <SlicePart> = "[" [ <Start> ] ":" [ <End> ] "]"
type SlicePart struct {
	NodeP

	Pos source.Pos

	// Part before ":". Can be nil if expression is omitted.
	Start Exp

	// Part after ":". Can be nil if expression is omitted.
	End Exp
}

// Explicit interface implementation check.
var _ ChainPart = SlicePart{}

func (SlicePart) Kind() exk.Kind {
	return exk.Slice
}

func (p SlicePart) Pin() source.Pos {
	return p.Pos
}

// <TestPart> => "." "test" "." <Name>
// <Name> => <Word>
type TestPart struct {
	NodeP

	Name Identifier
}

// Explicit interface implementation check.
var _ ChainPart = TestPart{}

func (TestPart) Kind() exk.Kind {
	return exk.Test
}

func (p TestPart) Pin() source.Pos {
	return p.Name.Pos
}
