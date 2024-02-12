package ast

import (
	"github.com/mebyus/gizmo/ast/toplvl"
	"github.com/mebyus/gizmo/source"
)

// <NamespaceBlock> = "namespace" "{" { <TopLevel> } "}"
type NamespaceBlock struct {
	Name ScopedIdentifier

	// Saved in order they appear inside namespace block
	Nodes []TopLevel
}

// <TopLevel> = <Function> | <Method> | <Type> | <Var> | <Const> | <Template>
type TopLevel interface {
	Node

	// dummy discriminator method
	TopLevel()

	Kind() toplvl.Kind
}

// Dummy provides quick, easy to use implementation of discriminator TopLevel() method
//
// Used for embedding into other (non-dummy) type specifier nodes
type nodeTopLevel struct{ uidHolder }

func (nodeTopLevel) TopLevel() {}

// <TopFunctionDeclaration> = [ "pub" ] <FunctionDeclaration>
type TopFunctionDeclaration struct {
	nodeTopLevel

	Declaration FunctionDeclaration

	Public bool
}

var _ TopLevel = TopFunctionDeclaration{}

func (TopFunctionDeclaration) Kind() toplvl.Kind {
	return toplvl.Declare
}

func (t TopFunctionDeclaration) Pin() source.Pos {
	return t.Declaration.Name.Pos
}

// <TopFunctionDefinition> = [ "pub" ] <FunctionDefinition>
type TopFunctionDefinition struct {
	nodeTopLevel

	Definition FunctionDefinition

	Public bool
}

var _ TopLevel = TopFunctionDefinition{}

func (TopFunctionDefinition) Kind() toplvl.Kind {
	return toplvl.Fn
}

func (t TopFunctionDefinition) Pin() source.Pos {
	return t.Definition.Head.Name.Pos
}

// <TopConstInit> = [ "pub" ] <ConstInit>
type TopConstInit struct {
	nodeTopLevel

	ConstInit

	Public bool
}

var _ TopLevel = TopConstInit{}

func (TopConstInit) Kind() toplvl.Kind {
	return toplvl.Const
}

func (t TopConstInit) Pin() source.Pos {
	return t.Pos
}

// <TopType> = [ "pub" ] "type" <Name> <TypeSpecifier>
//
// <Name> = <Identifier>
type TopType struct {
	nodeTopLevel

	Name Identifier

	Spec TypeSpecifier
}

var _ TopLevel = TopType{}

func (TopType) Kind() toplvl.Kind {
	return toplvl.Type
}

func (t TopType) Pin() source.Pos {
	return t.Name.Pos
}
