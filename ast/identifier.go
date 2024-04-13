package ast

import (
	"github.com/mebyus/gizmo/source"
)

// <Identifier> = word
type Identifier struct {
	Pos source.Pos
	Lit string
}

func (n Identifier) String() string {
	if len(n.Lit) == 0 {
		return "<nil>"
	}

	return n.Lit
}
