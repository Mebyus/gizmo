package ast

import (
	"github.com/mebyus/gizmo/ast/tps"
	"github.com/mebyus/gizmo/source"
)

// <TypeSpecifier> = <TypeName> | <TypeLiteral>
//
// <TypeName> = <ScopedIdentifier>
type TypeSpecifier interface {
	Node

	// dummy discriminator method
	TypeSpecifier()

	Kind() tps.Kind
}

// Dummy provides quick, easy to use implementation of discriminator TypeSpecifier() method
//
// Used for embedding into other (non-dummy) type specifier nodes
type nodeTypeSpecifier struct{}

func (nodeTypeSpecifier) TypeSpecifier() {}

type TypeName struct {
	nodeTypeSpecifier

	Name Identifier
}

// Explicit interface implementation check
var _ TypeSpecifier = TypeName{}

func (TypeName) Kind() tps.Kind {
	return tps.Name
}

func (t TypeName) Pin() source.Pos {
	return t.Name.Pos
}

// <StructTypeLiteral> = "struct" "{" { <FieldDefinition> "," } "}"
type StructType struct {
	nodeTypeSpecifier

	Pos source.Pos

	Fields []FieldDefinition
}

// Explicit interface implementation check
var _ TypeSpecifier = StructType{}

func (StructType) Kind() tps.Kind {
	return tps.Struct
}

func (t StructType) Pin() source.Pos {
	return t.Pos
}

// <PointerType> = "*" <ReferencedType>
//
// <ReferencedType> = <TypeSpecifier>
type PointerType struct {
	nodeTypeSpecifier

	Pos source.Pos

	RefType TypeSpecifier
}

// Explicit interface implementation check
var _ TypeSpecifier = PointerType{}

func (PointerType) Kind() tps.Kind {
	return tps.Pointer
}

func (t PointerType) Pin() source.Pos {
	return t.Pos
}

// <ArrayPointerType> = "[*]" <ElemType>
//
// <ElemType> = <TypeSpecifier>
type ArrayPointerType struct {
	nodeTypeSpecifier

	Pos source.Pos

	ElemType TypeSpecifier
}

// Explicit interface implementation check
var _ TypeSpecifier = ArrayPointerType{}

func (ArrayPointerType) Kind() tps.Kind {
	return tps.ArrayPointer
}

func (t ArrayPointerType) Pin() source.Pos {
	return t.Pos
}

type ChunkType struct {
	nodeTypeSpecifier

	Pos source.Pos

	ElemType TypeSpecifier
}

// Explicit interface implementation check
var _ TypeSpecifier = ChunkType{}

func (ChunkType) Kind() tps.Kind {
	return tps.Chunk
}

func (t ChunkType) Pin() source.Pos {
	return t.Pos
}

type ArrayType struct {
	nodeTypeSpecifier

	Pos source.Pos

	ElemType TypeSpecifier

	Size Expression
}

// Explicit interface implementation check
var _ TypeSpecifier = ArrayType{}

func (ArrayType) Kind() tps.Kind {
	return tps.Array
}

func (t ArrayType) Pin() source.Pos {
	return t.Pos
}

type EnumType struct {
	nodeTypeSpecifier

	Pos source.Pos

	// Base integer type
	Base TypeName

	Entries []EnumEntry
}

type EnumEntry struct {
	Name Identifier

	// Can be nil if entry does not have explicit assigned value
	Expression Expression
}

// Explicit interface implementation check
var _ TypeSpecifier = EnumType{}

func (EnumType) Kind() tps.Kind {
	return tps.Enum
}

func (t EnumType) Pin() source.Pos {
	return t.Pos
}

// <FunctionType> = "fn" <FunctionSignature>
type FunctionType struct {
	nodeTypeSpecifier

	Pos source.Pos

	Signature FunctionSignature
}

// Explicit interface implementation check
var _ TypeSpecifier = FunctionType{}

func (FunctionType) Kind() tps.Kind {
	return tps.Function
}

func (t FunctionType) Pin() source.Pos {
	return t.Pos
}

// <UnionType> = "union" "{" { <FieldDefinition> "," } "}"
type UnionType struct {
	nodeTypeSpecifier

	Pos source.Pos

	Fields []FieldDefinition
}

// Explicit interface implementation check
var _ TypeSpecifier = UnionType{}

func (UnionType) Kind() tps.Kind {
	return tps.Union
}

func (t UnionType) Pin() source.Pos {
	return t.Pos
}

type TypeParam struct {
	Name Identifier

	// Equals nil if there is no constraint on this param
	Constraint any
}

// <BagMethodSpec> = <Name> <Parameters> [ "=>" ( <Result> | "never" ) ]
//
// <Name> = <Identifier>
//
// <Parameters> = "(" { <FieldDefinition> "," } ")"
type BagMethodSpec struct {
	// Equals nil if there are no parameters in signature
	Params []FieldDefinition

	// Equals nil if function returns nothing or never returns
	Result TypeSpecifier

	Name Identifier

	// Equals true if function never returns
	Never bool
}

// <BagType> = <bag> "{" { <BagMethodSpec> "," } "}"
type BagType struct {
	nodeTypeSpecifier

	Pos source.Pos

	Methods []BagMethodSpec
}

// Explicit interface implementation check
var _ TypeSpecifier = UnionType{}

func (BagType) Kind() tps.Kind {
	return tps.Bag
}

func (t BagType) Pin() source.Pos {
	return t.Pos
}

// <TupleType> = "(" { <FieldDefinition> "," } ")"
type TupleType struct {
	nodeTypeSpecifier

	Pos source.Pos

	Fields []FieldDefinition
}

// Explicit interface implementation check
var _ TypeSpecifier = UnionType{}

func (TupleType) Kind() tps.Kind {
	return tps.Tuple
}

func (t TupleType) Pin() source.Pos {
	return t.Pos
}
