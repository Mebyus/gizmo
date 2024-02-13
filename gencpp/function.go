package gencpp

import "github.com/mebyus/gizmo/ast"

func (g *Builder) FunctionDefinition(definition ast.FunctionDefinition) {
	g.FunctionDeclaration(definition.Head)
	g.space()
	g.Block(definition.Body)
	g.nl()
}

func (g *Builder) FunctionDeclaration(declaration ast.FunctionDeclaration) {
	g.functionReturnType(declaration.Signature.Result)
	g.space()

	g.Identifier(declaration.Name)
	g.functionParams(declaration.Signature.Params)
	g.write(" noexcept")
}

func (g *Builder) functionParams(params []ast.FieldDefinition) {
	if len(params) == 0 {
		g.write("()")
		return
	}

	g.wb('(')

	g.functionParam(params[0])
	for _, param := range params[1:] {
		g.write(", ")
		g.functionParam(param)
	}

	g.wb(')')
}

func (g *Builder) functionParam(param ast.FieldDefinition) {
	g.TypeSpecifier(param.Type)
	g.space()
	g.Identifier(param.Name)
}

func (g *Builder) functionReturnType(spec ast.TypeSpecifier) {
	if spec == nil {
		g.write("void")
		return
	}

	g.TypeSpecifier(spec)
}
