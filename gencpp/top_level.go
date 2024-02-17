package gencpp

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
	"github.com/mebyus/gizmo/ast/tps"
)

func (g *Builder) TopLevel(node ast.TopLevel) {
	g.comment("gizmo.source = " + node.Pin().String())

	switch node.Kind() {
	case toplvl.Fn:
		g.FunctionDefinition(node.(ast.TopFunctionDefinition).Definition)
	case toplvl.Type:
		g.TopType(node.(ast.TopType))
	case toplvl.Const:
		g.TopConst(node.(ast.TopConst))
	case toplvl.Declare:
		g.TopDeclare(node.(ast.TopFunctionDeclaration))
	case toplvl.Var:
		g.TopVar(node.(ast.TopVar))
	case toplvl.Method:
		g.Method(node.(ast.Method))
	default:
		g.write(fmt.Sprintf("<%s node not implemented>", node.Kind().String()))
	}
}

func (g *Builder) TopType(top ast.TopType) {
	switch top.Spec.Kind() {
	case tps.Struct:
		g.TopStructType(top.Name, top.Spec.(ast.StructType))
	case tps.Enum:
		g.TopEnumType(top.Name, top.Spec.(ast.EnumType))
	default:
		g.write(fmt.Sprintf("<top %s type not implemented>", top.Spec.Kind()))
	}
}

func (g *Builder) TopStructType(name ast.Identifier, spec ast.StructType) {
	g.write("struct ")
	g.Identifier(name)
	g.space()
	g.structFields(spec.Fields)
	g.semi()
	g.nl()
}

func (g *Builder) TopConst(top ast.TopConst) {
	g.ConstInit(top.ConstInit)
	g.semi()
	g.nl()
}

func (g *Builder) TopDeclare(top ast.TopFunctionDeclaration) {
	g.FunctionDeclaration(top.Declaration)
	g.semi()
	g.nl()
}

func (g *Builder) TopVar(top ast.TopVar) {
	g.VarInit(top.VarInit)
	g.semi()
	g.nl()
}

func (g *Builder) Method(top ast.Method) {
	g.TypeSpecifier(top.Signature.Result)
	g.space()
	g.Identifier(top.Receiver)
	g.write("::")
	g.Identifier(top.Name)
	g.functionParams(top.Signature.Params)
	g.write(" noexcept ")
	g.Block(top.Body)
	g.nl()
}

func (g *Builder) TopEnumType(name ast.Identifier, spec ast.EnumType) {
	g.write("enum struct ")
	g.Identifier(name)
	g.space()

	if len(spec.Entries) == 0 {
		g.write("{};")
		g.nl()
		return
	}

	g.write("{")
	g.nl()
	
	g.inc()
	for _, entry := range spec.Entries {
		g.indent()

		g.Identifier(entry.Name)
		
		if entry.Expression != nil {
			g.write(" = ")
			g.Expression(entry.Expression)
		}

		g.write(",")
		g.nl()
	}
	g.dec()
	
	g.write("};")
	g.nl()
}
