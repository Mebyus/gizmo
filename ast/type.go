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

// <StructTypeLiteral> = "struct" "{" { <FieldDefinition> "," } "}"
type StructType struct {
	nodeTypeSpecifier

	Pos source.Pos

	Fields []FieldDefinition
}

// interface implementation check
var _ TypeSpecifier = StructType{}

func (t StructType) Kind() tps.Kind {
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

// <TypeLiteral> = <ArrayTypeLiteral> | <PointerTypeLiteral> | <ChunkTypeLiteral> |
// <ArrayPointerTypeLiteral> | <StructTypeLiteral>
type TypeLiteral any
