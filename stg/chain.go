package stg

import (
	"github.com/mebyus/gizmo/ast/exn"
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
type nodeChain struct{ nodeOperand }

func (nodeChain) ChainOperand() {}

type ChainSymbol struct {
	nodeChain

	Pos source.Pos

	// can be nil for receiver
	Sym *Symbol

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ChainSymbol{}

func (s *ChainSymbol) Kind() exn.Kind {
	return exn.Chain
}

func (s *ChainSymbol) Pin() source.Pos {
	return s.Pos
}

func (s *ChainSymbol) Type() *Type {
	return s.typ
}

type IndirectExpression struct {
	nodeChain

	Pos source.Pos

	Target ChainOperand

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &IndirectExpression{}

func (*IndirectExpression) Kind() exn.Kind {
	return exn.Indirect
}

func (e *IndirectExpression) Pin() source.Pos {
	return e.Pos
}

func (e *IndirectExpression) Type() *Type {
	return e.typ
}

type MemberExpression struct {
	nodeChain

	Pos source.Pos

	Target ChainOperand

	Member *Member
}

// Explicit interface implementation check.
var _ ChainOperand = &MemberExpression{}

func (*MemberExpression) Kind() exn.Kind {
	return exn.Member
}

func (e *MemberExpression) Pin() source.Pos {
	return e.Pos
}

func (e *MemberExpression) Type() *Type {
	return e.Member.Type
}

type IndirectMemberExpression struct {
	nodeChain

	Pos source.Pos

	Target ChainOperand

	Member *Member
}

// Explicit interface implementation check.
var _ ChainOperand = &IndirectMemberExpression{}

func (*IndirectMemberExpression) Kind() exn.Kind {
	return exn.IndirectMember
}

func (e *IndirectMemberExpression) Pin() source.Pos {
	return e.Pos
}

func (e *IndirectMemberExpression) Type() *Type {
	return e.Member.Type
}

type CallExpression struct {
	nodeChain

	Pos source.Pos

	Arguments []Expression

	Callee ChainOperand

	// Return type of the call.
	typ *Type

	// true for calls that never return.
	never bool
}

// Explicit interface implementation check
var _ ChainOperand = &CallExpression{}

func (*CallExpression) Kind() exn.Kind {
	return exn.Call
}

func (e *CallExpression) Pin() source.Pos {
	return e.Pos
}

func (e *CallExpression) Type() *Type {
	return e.typ
}

type AddressExpression struct {
	nodeChain

	Pos source.Pos

	Target ChainOperand

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &AddressExpression{}

func (*AddressExpression) Kind() exn.Kind {
	return exn.Address
}

func (e *AddressExpression) Pin() source.Pos {
	return e.Pos
}

func (e *AddressExpression) Type() *Type {
	return e.typ
}

type IndirectIndexExpression struct {
	nodeChain

	Pos source.Pos

	Target ChainOperand

	Index Expression

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &IndirectIndexExpression{}

func (*IndirectIndexExpression) Kind() exn.Kind {
	return exn.IndirectIndex
}

func (e *IndirectIndexExpression) Pin() source.Pos {
	return e.Pos
}

func (e *IndirectIndexExpression) Type() *Type {
	return e.typ
}

type ChunkIndexExpression struct {
	nodeChain

	Pos source.Pos

	Target ChainOperand

	Index Expression

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ChunkIndexExpression{}

func (*ChunkIndexExpression) Kind() exn.Kind {
	return exn.ChunkIndex
}

func (e *ChunkIndexExpression) Pin() source.Pos {
	return e.Pos
}

func (e *ChunkIndexExpression) Type() *Type {
	return e.typ
}

type ArrayIndexExp struct {
	nodeChain

	Pos source.Pos

	Target ChainOperand

	Index Expression

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ArrayIndexExp{}

func (*ArrayIndexExp) Kind() exn.Kind {
	return exn.ArrayIndex
}

func (e *ArrayIndexExp) Pin() source.Pos {
	return e.Pos
}

func (e *ArrayIndexExp) Type() *Type {
	return e.typ
}

type ArraySliceExp struct {
	nodeChain

	Pos source.Pos

	Target ChainOperand

	// Can be nil if expression is omitted.
	Start Expression

	// Can be nil if expression is omitted.
	End Expression

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ArraySliceExp{}

func (*ArraySliceExp) Kind() exn.Kind {
	return exn.ArraySlice
}

func (e *ArraySliceExp) Pin() source.Pos {
	return e.Pos
}

func (e *ArraySliceExp) Type() *Type {
	return e.typ
}

type ChunkSliceExp struct {
	nodeChain

	Pos source.Pos

	Target ChainOperand

	// Can be nil if expression is omitted.
	Start Expression

	// Can be nil if expression is omitted.
	End Expression

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ChunkSliceExp{}

func (*ChunkSliceExp) Kind() exn.Kind {
	return exn.ChunkSlice
}

func (e *ChunkSliceExp) Pin() source.Pos {
	return e.Pos
}

func (e *ChunkSliceExp) Type() *Type {
	return e.typ
}

type ChunkMemberExpression struct {
	nodeChain

	Pos source.Pos

	Target ChainOperand

	// property name:
	//	- len
	//	- ptr
	Name string

	typ *Type
}

// Explicit interface implementation check.
var _ ChainOperand = &ChunkMemberExpression{}

func (*ChunkMemberExpression) Kind() exn.Kind {
	return exn.ChunkMember
}

func (e *ChunkMemberExpression) Pin() source.Pos {
	return e.Pos
}

func (e *ChunkMemberExpression) Type() *Type {
	return e.typ
}
