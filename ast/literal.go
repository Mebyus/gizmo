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
	nodeExpression

	Token token.Token
}

var _ Expression = BasicLiteral{}

func (BasicLiteral) Kind() exn.Kind {
	return exn.Basic
}

func (l BasicLiteral) Pin() source.Pos {
	return l.Token.Pos
}

// <ListLiteral> = ".[" { <Expression> "," } "]"
type ListLiteral struct {
	nodeExpression

	Pos source.Pos

	Elems []Expression
}

func (ListLiteral) Kind() exn.Kind {
	return exn.List
}

func (l ListLiteral) Pin() source.Pos {
	return l.Pos
}
