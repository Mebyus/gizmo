package tt

import "github.com/mebyus/gizmo/ast"

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
	return nil
}

// Merge is called after all unit atoms were added to merger to build type tree of the unit.
func (m *Merger) Merge() (*Unit, error) {
	return &m.unit, nil
}
