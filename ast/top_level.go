package ast

import (
	"github.com/mebyus/gizmo/ast/toplvl"
	"github.com/mebyus/gizmo/source"
)

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
type nodeTopLevel struct{}

func (nodeTopLevel) TopLevel() {}

// <TopFunctionDeclaration> = [ "pub" ] "declare" <FunctionDeclaration>
type TopFunctionDeclaration struct {
	nodeTopLevel

	Declaration FunctionDeclaration

	Props []Prop

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

	Props []Prop

	Public bool
}

var _ TopLevel = TopFunctionDefinition{}

func (TopFunctionDefinition) Kind() toplvl.Kind {
	return toplvl.Fn
}

func (t TopFunctionDefinition) Pin() source.Pos {
	return t.Definition.Head.Name.Pos
}

// <TopConst> = [ "pub" ] <ConstInit>
type TopConst struct {
	nodeTopLevel

	ConstInit

	Public bool
}

var _ TopLevel = TopConst{}

func (TopConst) Kind() toplvl.Kind {
	return toplvl.Const
}

func (t TopConst) Pin() source.Pos {
	return t.Pos
}

// <TopType> = [ "pub" ] "type" <Name> <TypeSpecifier>
//
// <Name> = <Identifier>
type TopType struct {
	nodeTopLevel

	Name Identifier

	Spec TypeSpecifier

	Public bool
}

var _ TopLevel = TopType{}

func (TopType) Kind() toplvl.Kind {
	return toplvl.Type
}

func (t TopType) Pin() source.Pos {
	return t.Name.Pos
}

// <TopVar> = [ "pub" ] <VarInit>
type TopVar struct {
	nodeTopLevel

	VarInit

	Public bool
}

var _ TopLevel = TopVar{}

func (TopVar) Kind() toplvl.Kind {
	return toplvl.Var
}

func (t TopVar) Pin() source.Pos {
	return t.Pos
}

// <TopBlueprint> = [ "pub" ] "fn" <Name> <TypeParams> <FunctionSignature> <Body>
type TopBlueprint struct {
	nodeTopLevel

	Name Identifier

	// Contains at least one element
	Params []TypeParam

	Signature FunctionSignature

	Body BlockStatement

	Public bool
}

var _ TopLevel = TopBlueprint{}

func (TopBlueprint) Kind() toplvl.Kind {
	return toplvl.Blue
}

func (t TopBlueprint) Pin() source.Pos {
	return t.Name.Pos
}

// <TopPrototype> = [ "pub" ] "type" <Name> <TypeParams> <TypeSpecifier>
type TopPrototype struct {
	nodeTopLevel

	Name Identifier

	// Contains at least one element
	Params []TypeParam

	// Must be struct type
	Spec TypeSpecifier

	Public bool
}

var _ TopLevel = TopPrototype{}

func (TopPrototype) Kind() toplvl.Kind {
	return toplvl.Proto
}

func (t TopPrototype) Pin() source.Pos {
	return t.Name.Pos
}
