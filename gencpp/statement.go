package gencpp

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/stm"
)

// Same as BlockStatement method, but indentation formatting is different
// to start block on the same line
func (g *Builder) Block(block ast.BlockStatement) {
	g.wb('{')
	g.nl()
	g.inc()

	for _, statement := range block.Statements {
		g.Statement(statement)
	}

	g.dec()
	g.indent()
	g.wb('}')
}

func (g *Builder) BlockStatement(block ast.BlockStatement) {
	g.indent()
	g.Block(block)
	g.nl()
}

func (g *Builder) Statement(statement ast.Statement) {
	switch statement.Kind() {
	case stm.Block:
		g.BlockStatement(statement.(ast.BlockStatement))
	case stm.Return:
		g.ReturnStatement(statement.(ast.ReturnStatement))
	case stm.Const:
		g.ConstStatement(statement.(ast.ConstStatement))
	case stm.Var:
		g.VarStatement(statement.(ast.VarStatement))
	case stm.If:
		g.IfStatement(statement.(ast.IfStatement))
	case stm.Expr:
		g.ExpressionStatement(statement.(ast.ExpressionStatement))
	case stm.Assign:
		g.AssignStatement(statement.(ast.AssignStatement))
	default:
		g.indent()
		g.write(fmt.Sprintf("<%s statement not implemented>", statement.Kind().String()))
		g.nl()
	}
}

func (g *Builder) ExpressionStatement(statement ast.ExpressionStatement) {
	g.indent()

	g.Expression(statement.Expression)

	g.semi()
	g.nl()
}

func (g *Builder) AssignStatement(statement ast.AssignStatement) {
	g.indent()

	g.Expression(statement.Target)
	g.write(" = ")
	g.Expression(statement.Expression)

	g.semi()
	g.nl()
}

func (g *Builder) ReturnStatement(statement ast.ReturnStatement) {
	g.indent()
	g.write("return")

	if statement.Expression != nil {
		g.space()
		g.Expression(statement.Expression)
	}

	g.semi()
	g.nl()
}

func (g *Builder) ConstStatement(statement ast.ConstStatement) {
	g.indent()
	g.write("const")

	g.space()
	g.TypeSpecifier(statement.Type)

	g.space()
	g.Identifier(statement.Name)

	g.write(" = ")
	g.Expression(statement.Expression)

	g.semi()
	g.nl()
}

func (g *Builder) VarStatement(statement ast.VarStatement) {
	g.indent()

	g.TypeSpecifier(statement.Type)

	g.space()
	g.Identifier(statement.Name)

	if statement.Expression != nil {
		g.write(" = ")
		g.Expression(statement.Expression)
	}

	g.semi()
	g.nl()
}

func (g *Builder) IfStatement(statement ast.IfStatement) {
	g.indent()
	g.write("if (")
	g.Expression(statement.If.Condition)
	g.write(") ")
	g.Block(statement.If.Body)

	for _, clause := range statement.ElseIf {
		g.write(" else if (")
		g.Expression(clause.Condition)
		g.write(") ")
		g.Block(clause.Body)
	}

	if statement.Else != nil {
		g.write(" else ")
		g.Block(statement.Else.Body)
	}

	g.nl()
}
