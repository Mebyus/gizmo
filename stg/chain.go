package stg

import (
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/source"
)

type ChainOperand interface {
	Operand

	ChainOperand()
}

// Dummy chain operand node provides quick, easy to use implementation
// of discriminator ChainOperand() method.
//
// Used for embedding into other (non-dummy) chain part nodes.
type NodeC struct{ NodeO }

func (NodeC) ChainOperand() {}

type ChainSymbol struct {
	NodeC

	Pos source.Pos

	// can be nil for receiver
	Sym *Symbol

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ChainSymbol{}

func (s *ChainSymbol) Kind() exk.Kind {
	return exk.Chain
}

func (s *ChainSymbol) Pin() source.Pos {
	return s.Pos
}

func (s *ChainSymbol) Type() *Type {
	return s.typ
}

type IndirectExp struct {
	NodeC

	Pos source.Pos

	Target ChainOperand

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &IndirectExp{}

func (*IndirectExp) Kind() exk.Kind {
	return exk.Indirect
}

func (e *IndirectExp) Pin() source.Pos {
	return e.Pos
}

func (e *IndirectExp) Type() *Type {
	return e.typ
}

type MemberExp struct {
	NodeC

	Pos source.Pos

	Target ChainOperand

	Member *Member
}

// Explicit interface implementation check.
var _ ChainOperand = &MemberExp{}

func (*MemberExp) Kind() exk.Kind {
	return exk.Member
}

func (e *MemberExp) Pin() source.Pos {
	return e.Pos
}

func (e *MemberExp) Type() *Type {
	return e.Member.Type
}

type IndirectMemberExp struct {
	NodeC

	Pos source.Pos

	Target ChainOperand

	Member *Member
}

// Explicit interface implementation check.
var _ ChainOperand = &IndirectMemberExp{}

func (*IndirectMemberExp) Kind() exk.Kind {
	return exk.IndirectMember
}

func (e *IndirectMemberExp) Pin() source.Pos {
	return e.Pos
}

func (e *IndirectMemberExp) Type() *Type {
	return e.Member.Type
}

type CallExp struct {
	NodeC

	Pos source.Pos

	Arguments []Exp

	Callee ChainOperand

	// Return type of the call.
	typ *Type

	// true for calls that never return.
	never bool
}

// Explicit interface implementation check
var _ ChainOperand = &CallExp{}

func (*CallExp) Kind() exk.Kind {
	return exk.Call
}

func (e *CallExp) Pin() source.Pos {
	return e.Pos
}

func (e *CallExp) Type() *Type {
	return e.typ
}

type AddressExp struct {
	NodeC

	Pos source.Pos

	Target ChainOperand

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &AddressExp{}

func (*AddressExp) Kind() exk.Kind {
	return exk.Address
}

func (e *AddressExp) Pin() source.Pos {
	return e.Pos
}

func (e *AddressExp) Type() *Type {
	return e.typ
}

type IndirectIndexExp struct {
	NodeC

	Pos source.Pos

	Target ChainOperand

	Index Exp

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &IndirectIndexExp{}

func (*IndirectIndexExp) Kind() exk.Kind {
	return exk.IndirectIndex
}

func (e *IndirectIndexExp) Pin() source.Pos {
	return e.Pos
}

func (e *IndirectIndexExp) Type() *Type {
	return e.typ
}

type ChunkIndexExp struct {
	NodeC

	Pos source.Pos

	Target ChainOperand

	Index Exp

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ChunkIndexExp{}

func (*ChunkIndexExp) Kind() exk.Kind {
	return exk.ChunkIndex
}

func (e *ChunkIndexExp) Pin() source.Pos {
	return e.Pos
}

func (e *ChunkIndexExp) Type() *Type {
	return e.typ
}

type ArrayIndexExp struct {
	NodeC

	Pos source.Pos

	Target ChainOperand

	Index Exp

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ArrayIndexExp{}

func (*ArrayIndexExp) Kind() exk.Kind {
	return exk.ArrayIndex
}

func (e *ArrayIndexExp) Pin() source.Pos {
	return e.Pos
}

func (e *ArrayIndexExp) Type() *Type {
	return e.typ
}

type ArraySliceExp struct {
	NodeC

	Pos source.Pos

	Target ChainOperand

	// Can be nil if expression is omitted.
	Start Exp

	// Can be nil if expression is omitted.
	End Exp

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ArraySliceExp{}

func (*ArraySliceExp) Kind() exk.Kind {
	return exk.ArraySlice
}

func (e *ArraySliceExp) Pin() source.Pos {
	return e.Pos
}

func (e *ArraySliceExp) Type() *Type {
	return e.typ
}

type ChunkSliceExp struct {
	NodeC

	Pos source.Pos

	Target ChainOperand

	// Can be nil if expression is omitted.
	Start Exp

	// Can be nil if expression is omitted.
	End Exp

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ChunkSliceExp{}

func (*ChunkSliceExp) Kind() exk.Kind {
	return exk.ChunkSlice
}

func (e *ChunkSliceExp) Pin() source.Pos {
	return e.Pos
}

func (e *ChunkSliceExp) Type() *Type {
	return e.typ
}

type ChunkMemberExp struct {
	NodeC

	Pos source.Pos

	Target ChainOperand

	// property name:
	//	- len
	//	- ptr
	Name string

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ChunkMemberExp{}

func (*ChunkMemberExp) Kind() exk.Kind {
	return exk.ChunkMember
}

func (e *ChunkMemberExp) Pin() source.Pos {
	return e.Pos
}

func (e *ChunkMemberExp) Type() *Type {
	return e.typ
}
