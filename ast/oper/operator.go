package oper

import (
	"strconv"

	"github.com/mebyus/gizmo/token"
)

type Operator interface {
	Precedence() uint8
}

// Token.Kind is Plus, Minus or Not
type Unary token.Token

func NewUnary(tok token.Token) Unary {
	tok.Val = uint64(getUnaryPrecedence(tok.Kind))
	return Unary(tok)
}

func getUnaryPrecedence(k token.Kind) uint8 {
	switch k {
	case token.Plus, token.Minus:
		return 4
	case token.Not:
		return 12
	default:
		panic("unexpected token kind " + strconv.FormatInt(int64(k), 10))
	}
}

func (u Unary) Precedence() uint8 {
	return uint8(u.Val)
}

// Token.Kind is Less, Greater, Equal, LessEqual, GreaterEqual, NotEqual, And, Or, Plus, Minus, Mult, Div or Mod
type Binary token.Token

func NewBinary(tok token.Token) Binary {
	tok.Val = uint64(getBinaryPrecedence(tok.Kind))
	return Binary(tok)
}

func getBinaryPrecedence(k token.Kind) uint8 {
	switch k {
	case token.Asterisk, token.Slash, token.Percent:
		return 6
	case token.Plus, token.Minus:
		return 7
	case token.Less, token.Greater, token.Equal, token.LessOrEqual, token.GreaterOrEqual, token.NotEqual:
		return 10
	case token.LogicalAnd:
		return 13
	case token.LogicalOr:
		return 14
	default:
		panic("unexpected token kind " + strconv.FormatInt(int64(k), 10))
	}
}

func (b Binary) Precedence() uint8 {
	return uint8(b.Val)
}
