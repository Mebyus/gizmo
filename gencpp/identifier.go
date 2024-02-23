package gencpp

import (
	"strings"

	"github.com/mebyus/gizmo/ast"
)

func (g *Builder) Identifier(identifier ast.Identifier) {
	g.write(identifier.Lit)
}

func (g *Builder) ScopedIdentifier(identifier ast.ScopedIdentifier) {
	for _, s := range identifier.Scopes {
		g.Identifier(s)
		g.write("::")
	}
	g.Identifier(identifier.Name)
}

func (g *Builder) symName(name ast.Identifier) string {
	return strings.Join(append(g.currentScopes, name.Lit), "::")
}
