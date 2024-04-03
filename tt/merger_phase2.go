package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
)

// Phase 2 merger indexing coagulates info from AST nodes that were skipped during
// phase 1. It also performs additional scans on the tree to determine:
//
//   - top-level symbol resolution
//   - symbol hoisting information inside the unit
//   - used and defined types indexing
//   - type checking of various operations
func (m *Merger) runPhaseTwo() error {
	err := m.bindNodes()
	if err != nil {
		return err
	}
	return nil
}

func (m *Merger) bindNodes() error {
	for _, top := range m.symBindNodes {
		var err error

		switch top.Kind() {
		case toplvl.Method:
			err = m.bindMethod(top.(ast.Method))
		case toplvl.Pmb:
		default:
			panic(fmt.Sprintf("unexpected top-level %s node", top.Kind().String()))
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Merger) bindMethod(method ast.Method) error {
	rvName := method.Receiver.Lit
	rvSym := m.unit.Scope.sym(rvName)
	if rvSym == nil {
		return fmt.Errorf("%s: unresolved symbol \"%s\" in method receiver", method.Receiver.Pos.String(), rvName)
	}

	// TODO: add method to list of receiver's type methods
	return nil
}
