package ast

import (
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

// <Identifier> = <Word>
type Identifier struct {
	Pos source.Pos
	Lit string
}

func (n Identifier) Token() token.Token {
	return token.Token{
		Kind: token.Identifier,
		Pos:  n.Pos,
		Lit:  n.Lit,
	}
}

func (n Identifier) String() string {
	if len(n.Lit) == 0 {
		return "<nil>"
	}

	return n.Lit
}
