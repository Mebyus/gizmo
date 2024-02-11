package gencpp

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
)

func (g *Builder) TopLevel(node ast.TopLevel) {
	switch node.Kind() {
	case toplvl.Fn:
		g.FunctionDefinition(node.(ast.TopFunctionDefinition).Definition)
	default:
		g.write(fmt.Sprintf("<%s node not implemented>", node.Kind().String()))
	}
}
