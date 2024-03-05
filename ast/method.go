package ast

import (
	"github.com/mebyus/gizmo/ast/toplvl"
	"github.com/mebyus/gizmo/source"
)

// <Method> = "method" "(" <Receiver> ")" <Name> <Signature> <Body>
//
// <Receiver> = <Identifier>
// <Name> = <Identifier>
// <Signature> = <FunctionSignature>
// <Body> = <BlockStatement>
type Method struct {
	nodeTopLevel

	Receiver Identifier

	Name Identifier

	Signature FunctionSignature

	Body BlockStatement
}

var _ TopLevel = Method{}

func (Method) Kind() toplvl.Kind {
	return toplvl.Method
}

func (m Method) Pin() source.Pos {
	return m.Name.Pos
}

// <Method> = "method" "(" <Receiver> <TypeParams> ")" <Name> <Signature> <Body>
//
// <Receiver> = <Identifier>
// <Name> = <Identifier>
// <Signature> = <FunctionSignature>
// <Body> = <BlockStatement>
type MethodTemplate struct {
	nodeTopLevel

	Receiver Identifier

	// Contains at least one element
	TypeParams []Identifier

	Name Identifier

	Signature FunctionSignature

	Body BlockStatement
}

var _ TopLevel = Method{}

func (MethodTemplate) Kind() toplvl.Kind {
	return toplvl.MethodTemplate
}

func (m MethodTemplate) Pin() source.Pos {
	return m.Name.Pos
}
