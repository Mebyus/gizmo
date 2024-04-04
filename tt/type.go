package tt

import "github.com/mebyus/gizmo/ast"

type Type struct {
	nodeSymDef
}

func (m *Merger) lookupType(spec ast.TypeSpecifier) *Type {
	if spec == nil {
		return nil
	}
	return &Type{}
}
