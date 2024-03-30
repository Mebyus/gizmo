package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
)

// Merger is a high-level algorithm driver that gathers multiple ASTs of unit's atoms
// to produce that unit's type tree.
type Merger struct {
	ctx Context

	// Unit that is currently built by merger
	unit Unit
}

func New(ctx Context) *Merger {
	return &Merger{
		ctx: ctx,
		unit: Unit{
			sm: make(map[string]*Symbol),
		},
	}
}

// Context is a reference data structure that contains type and symbol information
// about imported units. It also holds build conditions under which unit compilation
// is performed.
type Context struct {
}

func (m *Merger) Add(atom ast.UnitAtom) error {
	for _, block := range atom.Blocks {
		// TODO: remove namespace blocks support
		if !block.Default {
			continue
		}

		for _, top := range block.Nodes {
			switch top.Kind() {
			case toplvl.Fn:
				fmt.Println("fn", top.(ast.TopFunctionDefinition).Definition.Head.Name.Lit)
			case toplvl.Type:
				fmt.Println("type", top.(ast.TopType).Name.Lit)
			}
		}
	}
	return nil
}

// Merge is called after all unit atoms were added to merger to build type tree of the unit.
func (m *Merger) Merge() (*Unit, error) {
	return &m.unit, nil
}
