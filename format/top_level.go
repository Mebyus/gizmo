package format

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
)

func (g *Builder) TopLevel(node ast.TopLevel) {
	switch node.Kind() {
	case toplvl.Fn:
		g.TopFn(node.(ast.TopFunctionDefinition))
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
		panic(fmt.Sprintf("top-level %s node not implemented", node.Kind().String()))
	}
}

func (g *Builder) TopFn(top ast.TopFunctionDefinition) {
	g.FunctionDefinition(top.Definition)
}

func (g *Builder) TopConst(top ast.TopConst) {
	// f.ConstInit(top.ConstInit)
}

func (g *Builder) TopDeclare(top ast.TopFunctionDeclaration) {
	// f.FunctionDeclaration(top.Declaration)
}

func (g *Builder) TopVar(top ast.TopVar) {
	// f.VarInit(top.VarInit)
}

func (g *Builder) Method(top ast.Method) {
	g.g.Str("<method>\n")
}

func (g *Builder) TopType(top ast.TopType) {

}
