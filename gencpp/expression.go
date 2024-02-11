package gencpp

import "github.com/mebyus/gizmo/ast"

func (g *Builder) Expression(expr ast.Expression) {
	g.write("<expr>")
}
