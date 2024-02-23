package uop

import (
	"fmt"

	"github.com/mebyus/gizmo/token"
)

// Kind indicates unary operator kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left unspecified
	empty Kind = iota

	Not // ' ! '

	Plus  // ' + '
	Minus // ' - '

	BitwiseNot // ' ^ '
)

var text = [...]string{
	empty: "<nil>",

	Not:        "!",
	Plus:       "+",
	Minus:      "-",
	BitwiseNot: "^",
}

func (k Kind) String() string {
	return text[k]
}

func FromToken(k token.Kind) Kind {
	switch k {
	case token.Not:
		return Not
	case token.Plus:
		return Plus
	case token.Minus:
		return Minus
	case token.Caret:
		return BitwiseNot
	default:
		panic(fmt.Sprintf("token ' %s ' cannot be unary operator", k.String()))
	}
}
