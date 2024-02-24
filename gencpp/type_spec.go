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
	case tps.Chunk:
		g.ChunkType(spec.(ast.ChunkType))
	case tps.Array:
		g.ArrayType(spec.(ast.ArrayType))
	case tps.Instance:
		g.TemplateInstanceType(spec.(ast.TemplateInstanceType))
	default:
		g.write(fmt.Sprintf("<%s type specifier not implemented>", spec.Kind().String()))
	}
}

func (g *Builder) TemplateInstanceType(spec ast.TemplateInstanceType) {
	g.ScopedIdentifier(spec.Name)
	g.write("<")

	g.TypeSpecifier(spec.Params[0])
	for _, param := range spec.Params[1:] {
		g.write(", ")
		g.TypeSpecifier(param)
	}

	g.write(">")
}

func (g *Builder) PointerType(spec ast.PointerType) {
	g.TypeSpecifier(spec.RefType)
	g.wb('*')
}

func (g *Builder) ArrayPointerType(spec ast.ArrayPointerType) {
	g.TypeSpecifier(spec.ElemType)
	g.wb('*')
}

func (g *Builder) structFields(fields []ast.FieldDefinition, methods []ast.FunctionDeclaration) {
	g.structFieldsWithDirtyConstructor(fields, "", methods)
}

func (g *Builder) structFieldsWithDirtyConstructor(fields []ast.FieldDefinition, name string,
	methods []ast.FunctionDeclaration) {

	if len(fields) == 0 {
		g.write("{}")
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

	if name != "" {
		// adds nop constructor to struct definition
		g.nl()
		g.indent()
		g.write(name)
		g.write("() {}")
		g.nl()
	}

	if len(methods) != 0 {
		g.nl()
		for _, method := range methods {
			g.indent()
			g.methodDeclaration(method)
			g.semi()
			g.nl()
		}
	}

	g.dec()
	g.write("}")
}

func (g *Builder) methodDeclaration(method ast.FunctionDeclaration) {
	g.functionReturnType(method.Signature.Result)
	g.space()

	g.Identifier(method.Name)
	g.functionParams(method.Signature.Params)
	g.write(" noexcept")
}

func (g *Builder) structField(field ast.FieldDefinition) {
	g.TypeSpecifier(field.Type)
	g.space()
	g.Identifier(field.Name)
}

func isByteType(spec ast.TypeSpecifier) bool {
	return spec.Kind() == tps.Name && spec.(ast.TypeName).Name.Name.Lit == "u8"
}

func (g *Builder) ChunkType(spec ast.ChunkType) {
	if isByteType(spec.ElemType) {
		g.write("mc")
		return
	}

	g.write("chunk<")
	g.TypeSpecifier(spec.ElemType)
	g.write(">")
}

func (g *Builder) ArrayType(spec ast.ArrayType) {
	g.write("array<")
	g.TypeSpecifier(spec.ElemType)
	g.write(", ")
	g.Expression(spec.Size)
	g.write(">")
}
