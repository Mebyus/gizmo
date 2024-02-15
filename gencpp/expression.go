package gencpp

import (
	"fmt"
	"strconv"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/bop"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/ast/uop"
	"github.com/mebyus/gizmo/token"
)

func (g *Builder) Expression(expr ast.Expression) {
	switch expr.Kind() {
	case exn.Start:
		g.ScopedIdentifier(expr.(ast.ChainStart).Identifier)
	case exn.Subs:
		g.ScopedIdentifier(expr.(ast.SubsExpression).Identifier)
	case exn.Basic:
		g.BasicLiteral(expr.(ast.BasicLiteral))
	case exn.Indirect:
		g.IndirectExpression(expr.(ast.IndirectExpression))
	case exn.Binary:
		g.BinaryExpression(expr.(ast.BinaryExpression))
	case exn.Call:
		g.CallExpression(expr.(ast.CallExpression))
	case exn.Unary:
		g.UnaryExpression(expr.(*ast.UnaryExpression))
	case exn.Indirx:
		g.IndirectIndexExpression(expr.(ast.IndirectIndexExpression))
	case exn.Paren:
		g.ParenthesizedExpression(expr.(ast.ParenthesizedExpression))
	case exn.Select:
		g.SelectorExpression(expr.(ast.SelectorExpression))
	default:
		g.write(fmt.Sprintf("<%s expr>", expr.Kind().String()))
	}
}

func (g *Builder) BasicLiteral(lit ast.BasicLiteral) {
	if lit.Token.Kind == token.Nil {
		g.write("0")
		return
	}
	if lit.Token.Kind == token.String {
		if len(lit.Token.Lit) == 0 {
			g.write("empty_string")
			return
		}

		g.write("make_static_string(")
		g.write(lit.Token.Literal())
		g.write(", ")
		g.write(strconv.FormatInt(int64(len(lit.Token.Lit)), 10))
		g.write(")")
		return
	}
	if lit.Token.Kind == token.OctalInteger {
		g.write("0" + strconv.FormatUint(lit.Token.Val, 8))
		return
	}

	g.write(lit.Token.Literal())
}

func (g *Builder) IndirectExpression(expr ast.IndirectExpression) {
	g.write("*(")
	g.Expression(expr.Target)
	g.write(")")
}

func (g *Builder) BinaryExpression(expr ast.BinaryExpression) {
	g.Expression(expr.Left)
	g.space()
	g.BinaryOperator(expr.Operator)
	g.space()
	g.Expression(expr.Right)
}

func (g *Builder) BinaryOperator(op ast.BinaryOperator) {
	if op.Kind == bop.BitwiseAndNot {
		g.write("&~")
		return
	}
	g.write(op.Kind.String())
}

func (g *Builder) UnaryExpression(expr *ast.UnaryExpression) {
	g.UnaryOperator(expr.Operator)
	g.Expression(expr.Inner)
}

func (g *Builder) UnaryOperator(op ast.UnaryOperator) {
	if op.Kind == uop.BitwiseNot {
		g.write("~")
		return
	}
	g.write(op.Kind.String())
}

func (g *Builder) CallExpression(expr ast.CallExpression) {
	g.Expression(expr.Callee)
	g.callArguments(expr.Arguments)
}

func (g *Builder) callArguments(args []ast.Expression) {
	if len(args) == 0 {
		g.write("()")
		return
	}

	g.wb('(')

	g.Expression(args[0])
	for _, arg := range args[1:] {
		g.write(", ")
		g.Expression(arg)
	}

	g.wb(')')
}

func (g *Builder) IndirectIndexExpression(expr ast.IndirectIndexExpression) {
	g.Expression(expr.Target)
	g.write("[")
	g.Expression(expr.Index)
	g.write("]")
}

func (g *Builder) ParenthesizedExpression(expr ast.ParenthesizedExpression) {
	g.write("(")
	g.Expression(expr.Inner)
	g.write(")")
}

func (g *Builder) SelectorExpression(expr ast.SelectorExpression) {
	g.Expression(expr.Target)
	g.write(".")
	g.Identifier(expr.Selected)
}
