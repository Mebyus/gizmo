package ast

import (
	"github.com/mebyus/gizmo/ast/tps"
	"github.com/mebyus/gizmo/source"
)

// <TypeSpec> = <TypeName> | <TypeLiteral>
type TypeSpec interface {
	Node

	// dummy discriminator method
	TypeSpec()

	Kind() tps.Kind
}

// Provides quick, easy to use implementation of discriminator TypeSpec() method
//
// Used for embedding into other (non-dummy) type specifier nodes.
type nodeTypeSpec struct{}

func (nodeTypeSpec) TypeSpec() {}

type TypeName struct {
	nodeTypeSpec

	Name Identifier
}

// Explicit interface implementation check
var _ TypeSpec = TypeName{}

func (TypeName) Kind() tps.Kind {
	return tps.Name
}

func (t TypeName) Pin() source.Pos {
	return t.Name.Pos
}

// <StructTypeLiteral> = "struct" "{" { <FieldDefinition> "," } "}"
type StructType struct {
	nodeTypeSpec

	Pos source.Pos

	Fields []FieldDefinition
}

// Explicit interface implementation check
var _ TypeSpec = StructType{}

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
	nodeTypeSpec

	Pos source.Pos

	RefType TypeSpec
}

// Explicit interface implementation check
var _ TypeSpec = PointerType{}

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
	nodeTypeSpec

	Pos source.Pos

	ElemType TypeSpec
}

// Explicit interface implementation check
var _ TypeSpec = ArrayPointerType{}

func (ArrayPointerType) Kind() tps.Kind {
	return tps.ArrayPointer
}

func (t ArrayPointerType) Pin() source.Pos {
	return t.Pos
}

type ChunkType struct {
	nodeTypeSpec

	Pos source.Pos

	ElemType TypeSpec
}

// Explicit interface implementation check
var _ TypeSpec = ChunkType{}

func (ChunkType) Kind() tps.Kind {
	return tps.Chunk
}

func (t ChunkType) Pin() source.Pos {
	return t.Pos
}

type ArrayType struct {
	nodeTypeSpec

	Pos source.Pos

	ElemType TypeSpec

	Size Expression
}

// Explicit interface implementation check
var _ TypeSpec = ArrayType{}

func (ArrayType) Kind() tps.Kind {
	return tps.Array
}

func (t ArrayType) Pin() source.Pos {
	return t.Pos
}

type EnumType struct {
	nodeTypeSpec

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
var _ TypeSpec = EnumType{}

func (EnumType) Kind() tps.Kind {
	return tps.Enum
}

func (t EnumType) Pin() source.Pos {
	return t.Pos
}

// <FunctionType> = "fn" <FunctionSignature>
type FunctionType struct {
	nodeTypeSpec

	Pos source.Pos

	Signature Signature
}

// Explicit interface implementation check
var _ TypeSpec = FunctionType{}

func (FunctionType) Kind() tps.Kind {
	return tps.Function
}

func (t FunctionType) Pin() source.Pos {
	return t.Pos
}

// <UnionType> = "union" "{" { <FieldDefinition> "," } "}"
type UnionType struct {
	nodeTypeSpec

	Pos source.Pos

	Fields []FieldDefinition
}

// Explicit interface implementation check
var _ TypeSpec = UnionType{}

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
	Result TypeSpec

	Name Identifier

	// Equals true if function never returns
	Never bool
}

// <BagType> = <bag> "{" { <BagMethodSpec> "," } "}"
type BagType struct {
	nodeTypeSpec

	Pos source.Pos

	Methods []BagMethodSpec
}

// Explicit interface implementation check
var _ TypeSpec = UnionType{}

func (BagType) Kind() tps.Kind {
	return tps.Bag
}

func (t BagType) Pin() source.Pos {
	return t.Pos
}

// <TupleType> = "(" { <FieldDefinition> "," } ")"
type TupleType struct {
	nodeTypeSpec

	Pos source.Pos

	Fields []FieldDefinition
}

// Explicit interface implementation check
var _ TypeSpec = UnionType{}

func (TupleType) Kind() tps.Kind {
	return tps.Tuple
}

func (t TupleType) Pin() source.Pos {
	return t.Pos
}
