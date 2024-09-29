package ast

import (
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

// <Dirty> => "dirty"
//
// Dirty literal can only appear as a complete expression when initializing
// a variable or a field in object literal. It cannot be a part of a bigger
// expression. For example construct
//
//	(dirty)
//
// is illegal.
type Dirty struct {
	NodeE

	Pos source.Pos
}

var _ Exp = Dirty{}

func (Dirty) Kind() exk.Kind {
	return exk.Dirty
}

func (d Dirty) Pin() source.Pos {
	return d.Pos
}

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

func (BasicLiteral) Kind() exk.Kind {
	return exk.Basic
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

func (ListLiteral) Kind() exk.Kind {
	return exk.List
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

func (ObjectLiteral) Kind() exk.Kind {
	return exk.Object
}

func (l ObjectLiteral) Pin() source.Pos {
	return l.Pos
}
