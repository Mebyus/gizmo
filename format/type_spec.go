package format

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/tps"
	"github.com/mebyus/gizmo/token"
)

func (g *Builder) TypeSpecifier(spec ast.TypeSpecifier) {
	switch spec.Kind() {
	case tps.Struct:
		g.StructType(spec.(ast.StructType))
	case tps.Enum:
		g.EnumType(spec.(ast.EnumType))
	default:
		panic(fmt.Sprintf("%s type node not implemented", spec.Kind().String()))
	}
}

func (g *Builder) StructType(spec ast.StructType) {
	g.genpos(token.Struct, spec.Pos)
	g.ss()

	// TODO: struct fields
}

func (g *Builder) EnumType(spec ast.EnumType) {
	g.genpos(token.Enum, spec.Pos)
	g.ss()
	g.idn(spec.Base.Name)
	g.ss()

	// TODO: enum members
}
