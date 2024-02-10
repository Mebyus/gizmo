package treeview

import "github.com/mebyus/gizmo/ast"

func formatIdentifier(identifier ast.Identifier) string {
	if len(identifier.Lit) == 0 {
		return "<nil>"
	}

	return identifier.Lit
}

func formatScopedIdentifier(identifier ast.ScopedIdentifier) string {
	if len(identifier.Scopes) == 0 && len(identifier.Name.Lit) == 0 {
		return "<nil>"
	}

	var s string
	for _, name := range identifier.Scopes {
		s += name.Lit + "::"
	}
	s += identifier.Name.Lit
	return s
}
