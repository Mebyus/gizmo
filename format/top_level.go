package format

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
	"github.com/mebyus/gizmo/token"
)

func (g *Noder) TopLevel(node ast.TopLevel) {
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

func (g *Noder) TopFn(top ast.TopFunctionDefinition) {
	if top.Public {
		g.pub()
	}

	g.FunctionDefinition(top.Definition)
}

func (g *Noder) TopConst(top ast.TopConst) {
	// f.ConstInit(top.ConstInit)
}

func (g *Noder) TopDeclare(top ast.TopFunctionDeclaration) {
	// f.FunctionDeclaration(top.Declaration)
}

func (g *Noder) TopVar(top ast.TopVar) {
	// f.VarInit(top.VarInit)
}

func (g *Noder) Method(top ast.Method) {

}

func (g *Noder) TopType(top ast.TopType) {
	if top.Public {
		g.pub()
	}

	g.gen(token.Type)
	g.ss()
	g.idn(top.Name)
	g.ss()
	g.TypeSpecifier(top.Spec)
}
