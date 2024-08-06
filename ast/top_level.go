package ast

import (
	"github.com/mebyus/gizmo/source"
)

// Top level node.
//
// <Top> = <Fun> | <Method> | <Type> | <Var> | <Con>
type Top interface {
	Node

	// dummy discriminator method
	Top()
}

// Provides quick, easy to use implementation of discriminator Top() method.
// Used for embedding into other (non-dummy) top level nodes.
type nodeTop struct{}

func (nodeTop) Top() {}

// Trait container object for passing around top node
// attributes and properties.
type Traits struct {
	// List of node's properties.
	Props *[]Prop

	// True for public nodes.
	Pub bool
}

// TopDec top level node with function declaration.
//
// <TopDec> = [ "pub" ] "fn" <Name> <Signature>
type TopDec struct {
	nodeTop

	Signature Signature

	Name Identifier

	Traits
}

var _ Top = TopDec{}

func (t TopDec) Pin() source.Pos {
	return t.Name.Pos
}

// TopFun top level node with function definition.
//
// <TopFun> = <FunDec> <Body>
type TopFun struct {
	nodeTop

	Signature Signature

	Name Identifier

	Body BlockStatement

	Traits
}

var _ Top = TopFun{}

func (t TopFun) Pin() source.Pos {
	return t.Name.Pos
}

// <TopLet> = [ "pub" ] <Let>
type TopLet struct {
	nodeTop

	Let

	Traits
}

var _ Top = TopLet{}

func (t TopLet) Pin() source.Pos {
	return t.Name.Pos
}

// <TopType> = [ "pub" ] "type" <Name> <TypeSpecifier>
//
// <Name> = <Identifier>
type TopType struct {
	nodeTop

	Name Identifier

	Spec TypeSpec

	Traits
}

var _ Top = TopType{}

func (t TopType) Pin() source.Pos {
	return t.Name.Pos
}

// <TopVar> = [ "pub" ] <Var>
type TopVar struct {
	nodeTop

	Var

	Traits
}

var _ Top = TopVar{}

func (t TopVar) Pin() source.Pos {
	return t.Pos
}

// <ReceiverTypeSpec> = [ "*" ] <Identifier>
type ReceiverTypeSpec struct {
	Name Identifier
	Ptr  bool
}

// <Method> = "fn" "[" <ReceiverTypeSpec> "]" <Name> <Signature> <Body>
//
// <Receiver> = <TypeSpecifier>
// <Name> = <Identifier>
// <Signature> = <FunctionSignature>
// <Body> = <BlockStatement>
type Method struct {
	nodeTop

	Name Identifier

	Receiver ReceiverTypeSpec

	Signature Signature

	Body BlockStatement

	Traits
}

var _ Top = Method{}

func (m Method) Pin() source.Pos {
	return m.Name.Pos
}
