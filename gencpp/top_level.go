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
	g.structFields(spec.Fields)
	g.nl()
}
