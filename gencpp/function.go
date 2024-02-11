package gencpp

import "github.com/mebyus/gizmo/ast"

func (g *Builder) FunctionDefinition(definition ast.FunctionDefinition) {
	g.FunctionDeclaration(definition.Head)
	g.Block(definition.Body)
	g.nl()
}

func (g *Builder) FunctionDeclaration(declaration ast.FunctionDeclaration) {
	g.functionReturnType(declaration.Signature.Result)
	g.space()

	g.Identifier(declaration.Name)
	g.wb('(')

	for _, param := range declaration.Signature.Params {
		g.TypeSpecifier(param.Type)
		g.space()
		g.Identifier(param.Name)
		g.write(", ")
	}

	g.wb(')')
	g.write(" noexcept ")
}

func (g *Builder) functionReturnType(spec ast.TypeSpecifier) {
	if spec == nil {
		g.write("void")
		return
	}

	g.TypeSpecifier(spec)
}
