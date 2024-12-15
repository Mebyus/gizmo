package ast

import (
	"github.com/mebyus/gizmo/source"
)

// <Identifier> = <Word>
type Identifier struct {
	Pos source.Pos
	Lit string
}

func (n Identifier) AsChainOperand() ChainOperand {
	return ChainOperand{Start: n}
}

func (n Identifier) String() string {
	if len(n.Lit) == 0 {
		return "<nil>"
	}

	return n.Lit
}
