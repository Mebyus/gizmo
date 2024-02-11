package gencpp

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/tps"
)

func (g *Builder) TypeSpecifier(spec ast.TypeSpecifier) {
	switch spec.Kind() {
	case tps.Name:
		g.ScopedIdentifier(spec.(ast.TypeName).Name)
	default:
		g.write(fmt.Sprintf("<%s type specifier not implemented>", spec.Kind().String()))
	}
}
