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
	case tps.ArrayPointer:
		g.ArrayPointerType(spec.(ast.ArrayPointerType))
	default:
		g.write(fmt.Sprintf("<%s type specifier not implemented>", spec.Kind().String()))
	}
}

func (g *Builder) PointerType(spec ast.PointerType) {
	g.TypeSpecifier(spec.RefType)
	g.wb('*')
}

func (g *Builder) ArrayPointerType(spec ast.ArrayPointerType) {
	g.TypeSpecifier(spec.ElemType)
	g.wb('*')
}

func (g *Builder) structFields(fields []ast.FieldDefinition) {
	if len(fields) == 0 {
		g.write("<empty struct not implemented>")
		return
	}

	g.write("{")
	g.nl()
	g.inc()

	for _, field := range fields {
		g.indent()
		g.structField(field)
		g.semi()
		g.nl()
	}

	g.dec()
	g.write("};")
}

func (g *Builder) structField(field ast.FieldDefinition) {
	g.TypeSpecifier(field.Type)
	g.space()
	g.Identifier(field.Name)
}
