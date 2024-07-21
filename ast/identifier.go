package ast

import (
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

// <Identifier> = <Word>
type Identifier struct {
	Pos source.Pos

	// TODO: make empty string here to be a signal for "g" receiver usage
	// and remove separate Receiver expression after that. Maybe we can
	// integrate this magic directly into lexer.
	Lit string
}

func (n Identifier) AsChainOperand() ChainOperand {
	return ChainOperand{Identifier: n}
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
