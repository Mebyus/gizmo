package format

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (g *Noder) TopFn(top ast.TopFun) {
	if top.Pub {
		g.pub()
	}

	// g.FunctionDefinition(top.Definition)
}

func (g *Noder) TopConst(top ast.TopLet) {
	// f.ConstInit(top.ConstInit)
}

func (g *Noder) TopDeclare(top ast.TopFun) {
	// f.FunctionDeclaration(top.Declaration)
}

func (g *Noder) TopVar(top ast.TopVar) {
	// f.VarInit(top.VarInit)
}

func (g *Noder) Method(top ast.Method) {

}

func (g *Noder) TopType(top ast.TopType) {
	if top.Pub {
		g.pub()
	}

	g.gen(token.Type)
	g.ss()
	g.idn(top.Name)
	g.ss()
	g.TypeSpecifier(top.Spec)
}
