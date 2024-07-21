package aop

import (
	"github.com/mebyus/gizmo/token"
)

// Kind indicates assignment operation kind.
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left unspecified
	empty Kind = iota

	Plain // "="

	Add // "+="
	Sub // "-="
	Mul // "*="
	Div // "/="
	Rem // "%="
)

var text = [...]string{
	empty: "<nil>",

	Plain: "=",

	Add: "+=",
	Sub: "-=",
	Mul: "*=",
	Div: "/=",
	Rem: "%=",
}

func (k Kind) String() string {
	return text[k]
}

func FromToken(kind token.Kind) (k Kind, ok bool) {
	switch kind {

	case token.Assign:
		k = Plain
	case token.AddAssign:
		k = Add
	case token.SubAssign:
		k = Sub
	case token.MulAssign:
		k = Mul
	case token.DivAssign:
		k = Div
	case token.RemAssign:
		k = Rem

	default:
		return empty, false
	}

	return k, true
}
