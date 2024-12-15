package ast

import (
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/source"
)

// <ChainOperand> = <Identifier> | <SelectExp> | <IndirectExp> | <IndirectIndexExp> | <IndirectIndexExp>
type ChainOperand struct {
	NodeO

	// Start of chain operand.
	Start Identifier

	// Parts are arranged from left to right as in source text.
	Parts []ChainPart
}

// Explicit interface implementation check.
var _ Operand = ChainOperand{}

func (o ChainOperand) Pin() source.Pos {
	return o.Start.Pos
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

func (o ChainOperand) LastPartPin() source.Pos {
	if len(o.Parts) == 0 {
		return o.Start.Pos
	}
	return o.Parts[len(o.Parts)-1].Pin()
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

// <CallExp> = <ChainOperand> "(" { <Exp> "," } ")"
type CallExp struct {
	NodeO

	Args []Exp

	Callee ChainOperand
}

// Explicit interface implementation check.
var _ Operand = CallExp{}

func (CallExp) Kind() exk.Kind {
	return exk.Call
}

func (p CallExp) Pin() source.Pos {
	return p.Callee.LastPartPin()
}

// <IndexPart> = "[" <Exp> "]"
type IndexPart struct {
	NodeP

	Index Exp
}

// Explicit interface implementation check.
var _ ChainPart = IndexPart{}

func (IndexPart) Kind() exk.Kind {
	return exk.Index
}

func (p IndexPart) Pin() source.Pos {
	return p.Index.Pin()
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

// <AddressExp> = <ChainOperand> ".&"
type AddressExp struct {
	NodeO

	Chain ChainOperand
}

// Explicit interface implementation check.
var _ Operand = AddressExp{}

func (AddressExp) Kind() exk.Kind {
	return exk.Address
}

func (p AddressExp) Pin() source.Pos {
	return p.Chain.LastPartPin()
}

// <IndirectIndexPart> => ".[" <Index> "]"
//
// <Index> => <Exp>
type IndirectIndexPart struct {
	NodeP

	Index Exp
}

// Explicit interface implementation check.
var _ ChainPart = IndirectIndexPart{}

func (IndirectIndexPart) Kind() exk.Kind {
	return exk.IndirectIndex
}

func (p IndirectIndexPart) Pin() source.Pos {
	return p.Index.Pin()
}

// <SliceExp> = "[" [ <Start> ] ":" [ <End> ] "]"
type SliceExp struct {
	NodeO

	Chain ChainOperand

	// Part before ":". Can be nil if expression is omitted.
	Start Exp

	// Part after ":". Can be nil if expression is omitted.
	End Exp
}

// Explicit interface implementation check.
var _ Operand = SliceExp{}

func (SliceExp) Kind() exk.Kind {
	return exk.Slice
}

func (p SliceExp) Pin() source.Pos {
	return p.Chain.LastPartPin()
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
