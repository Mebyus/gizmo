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
type nodeTypeSpecifier struct{ uidHolder }

func (nodeTypeSpecifier) TypeSpecifier() {}

type TypeName struct {
	nodeTypeSpecifier

	Name ScopedIdentifier
}

// interface implementation check
var _ TypeSpecifier = TypeName{}

func (TypeName) Kind() tps.Kind {
	return tps.Name
}

func (t TypeName) Pin() source.Pos {
	return t.Name.Pos()
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

// <TemplateInstanceType> = <TemplateName> "<$" { <TypeSpecifier> "," } ">"
//
// <TemplateName> = <ScopedIdentifier>
type TemplateInstanceType struct {
	nodeTypeSpecifier

	// Cannot be nil
	Params []TypeSpecifier

	Name ScopedIdentifier
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

// <TypeLiteral> = <ArrayTypeLiteral> | <PointerTypeLiteral> | <ChunkTypeLiteral> |
// <ArrayPointerTypeLiteral> | <StructTypeLiteral>
type TypeLiteral any
