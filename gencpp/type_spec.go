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
	case tps.Pointer:
		g.PointerType(spec.(ast.PointerType))
	default:
		g.write(fmt.Sprintf("<%s type specifier not implemented>", spec.Kind().String()))
	}
}

func (g *Builder) PointerType(spec ast.PointerType) {
	g.TypeSpecifier(spec.RefType)
	g.wb('*')
}
