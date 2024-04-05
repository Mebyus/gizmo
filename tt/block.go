package tt

import (
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
)

type Block struct {
	nodeStatement

	// Position where block starts.
	Pos source.Pos

	Nodes []Statement

	Scope *Scope
}

// Explicit interface implementation check
var _ Statement = &Block{}

func (b *Block) Pin() source.Pos {
	return b.Pos
}

func (b *Block) Kind() stm.Kind {
	return stm.Block
}
