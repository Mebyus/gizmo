package format

import "github.com/mebyus/gizmo/ast"

func (g *Builder) FunctionDefinition(node ast.FunctionDefinition) {
	g.Block(node.Body)
}
