package treeview

import "github.com/mebyus/gizmo/ast"

func formatScopedIdentifier(identifier ast.ScopedIdentifier) string {
	var s string
	for _, name := range identifier.Scopes {
		s += name.Lit + "::"
	}
	s += identifier.Name.Lit
	return s
}
