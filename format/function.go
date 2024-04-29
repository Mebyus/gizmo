package format

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (g *Builder) FunctionDefinition(node ast.FunctionDefinition) {
	g.FunctionDeclaration(node.Head)
	g.Block(node.Body)
}

func (g *Builder) FunctionDeclaration(decl ast.FunctionDeclaration) {
	g.gen(token.Fn)
	g.ss()
	g.idn(decl.Name)
	g.ss()

	g.gen(token.LeftParentheses)
	g.sep()
	g.FunctionParams(decl.Signature.Params)
	g.sep()
	g.gen(token.RightParentheses)
	g.ss()

	if decl.Signature.Never {
		g.gen(token.RightArrow)
		g.ss()
		g.gen(token.Never)
		g.ss()
		return
	}

	if decl.Signature.Result == nil {
		return
	}

	g.gen(token.RightArrow)
	g.ss()
	g.TypeSpecifier(decl.Signature.Result)
	g.ss()
}

func (g *Builder) FunctionParams(params []ast.FieldDefinition) {
	if len(params) == 0 {
		return
	}

	for i := 0; i < len(params)-1; i += 1 {
		p := params[i]

		g.idn(p.Name)
		g.gen(token.Colon)
		g.ss()
		g.TypeSpecifier(p.Type)
		g.gen(token.Comma)
		g.space()
	}

	last := params[len(params)-1]

	g.idn(last.Name)
	g.gen(token.Colon)
	g.ss()
	g.TypeSpecifier(last.Type)
	g.trailComma()
}
