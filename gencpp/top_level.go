package gencpp

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
	"github.com/mebyus/gizmo/ast/tps"
)

func (g *Builder) TopLevel(node ast.TopLevel) {
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
	default:
		g.write(fmt.Sprintf("<%s node not implemented>", node.Kind().String()))
	}
}

func (g *Builder) TopType(top ast.TopType) {
	if top.Spec.Kind() == tps.Struct {
		g.TopStructType(top.Name, top.Spec.(ast.StructType))
		return
	}

	g.write(fmt.Sprintf("<top %s type not implemented>", top.Spec.Kind()))
}

func (g *Builder) TopStructType(name ast.Identifier, spec ast.StructType) {
	g.write("struct ")
	g.Identifier(name)
	g.space()
	g.structFieldsWithDirtyConstructor(spec.Fields, name.Lit)
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
