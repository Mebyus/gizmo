package ast

import (
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

// Expression which consists of usage of a single basic literal
//
// <BasicLiteral> = <IntegerLiteral> | <FloatLiteral> | <CharacterLiteral> | <StringLiteral> | <NilLiteral> | <TrueLiteral> | <FalseLiteral>
//
// token.Kind is DecimalInteger, DecimalFloat, Character, String, Nil, True, False
type BasicLiteral struct {
	NodeO

	Token token.Token
}

var _ Operand = BasicLiteral{}

func (BasicLiteral) Kind() exn.Kind {
	return exn.Basic
}

func (l BasicLiteral) Pin() source.Pos {
	return l.Token.Pos
}

// <ListLiteral> = "[" [ <Expression> { "," <Expression> } [ "," ] ] "]"
type ListLiteral struct {
	NodeO

	Pos source.Pos

	Elems []Exp
}

func (ListLiteral) Kind() exn.Kind {
	return exn.List
}

func (l ListLiteral) Pin() source.Pos {
	return l.Pos
}

// <ObjectField> = <Name> ":" <Expression>
//
// <Name> = <Identifier>
type ObjectField struct {
	Name  Identifier
	Value Exp
}

// <ObjectLiteral> = "{" [ <ObjectField> { "," <ObjectField> } [ "," ] ] "}"
type ObjectLiteral struct {
	NodeO

	Pos source.Pos

	Fields []ObjectField
}

func (ObjectLiteral) Kind() exn.Kind {
	return exn.Object
}

func (l ObjectLiteral) Pin() source.Pos {
	return l.Pos
}
