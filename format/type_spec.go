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
	case tps.Name:
		g.TypeName(spec.(ast.TypeName))
	case tps.Pointer:
		g.PointerType(spec.(ast.PointerType))
	default:
		panic(fmt.Sprintf("%s type node not implemented", spec.Kind().String()))
	}
}

func (g *Builder) PointerType(spec ast.PointerType) {
	g.genpos(token.Asterisk, spec.Pos)
	g.TypeSpecifier(spec.RefType)
}

func (g *Builder) TypeName(spec ast.TypeName) {
	g.idn(spec.Name)
}

func (g *Builder) StructType(spec ast.StructType) {
	g.genpos(token.Struct, spec.Pos)
	g.ss()

	g.gen(token.LeftCurly)
	g.inc()

	for i := 0; i < len(spec.Fields); i += 1 {
		f := spec.Fields[i]

		g.nli()
		g.idn(f.Name)
		g.gen(token.Colon)
		g.ss()
		g.TypeSpecifier(f.Type)
		g.gen(token.Comma)
	}

	g.dec()
	g.nl()
	g.gen(token.RightCurly)
}

func (g *Builder) EnumType(spec ast.EnumType) {
	g.genpos(token.Enum, spec.Pos)
	g.ss()
	g.idn(spec.Base.Name)
	g.ss()

	// TODO: enum members
}
