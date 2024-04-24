package format

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/token"
)

func (g *Builder) Expression(expr ast.Expression) {
	switch expr.Kind() {
	case exn.Symbol:
		g.SymbolExpression(expr.(ast.SymbolExpression))
	case exn.Receiver:
		g.ReceiverExpression(expr.(ast.Receiver))
	case exn.Basic:
		g.BasicLiteral(expr.(ast.BasicLiteral))
	case exn.Index:
		g.IndexExpression(expr.(ast.IndexExpression))
	case exn.Binary:
		g.BinaryExpression(expr.(ast.BinaryExpression))
	case exn.Call:
		g.CallExpression(expr.(ast.CallExpression))
	case exn.Address:
		g.AddressExpression(expr.(ast.AddressExpression))
	case exn.Paren:
		g.ParenthesizedExpression(expr.(ast.ParenthesizedExpression))
	case exn.Indirect:
		g.IndirectExpression(expr.(ast.IndirectExpression))
	case exn.Start:
		g.ChainStart(expr.(ast.ChainStart))
	default:
		panic(fmt.Sprintf("%s expression node not implemented", expr.Kind().String()))
	}
}

func (g *Builder) ChainStart(expr ast.ChainStart) {
	g.idn(expr.Identifier)
}

func (g *Builder) IndirectExpression(expr ast.IndirectExpression) {
	g.Expression(expr.Target)
	g.genpos(token.Indirect, expr.Pos)
}

func (g *Builder) AddressExpression(expr ast.AddressExpression) {
	g.Expression(expr.Target)
	g.gen(token.Address)
}

func (g *Builder) CallExpression(expr ast.CallExpression) {
	g.Expression(expr.Callee)

	args := expr.Arguments
	if len(args) == 0 {
		g.gen(token.LeftParentheses)
		g.gen(token.RightParentheses)
		return
	}

	g.gen(token.LeftParentheses)
	g.sep()

	for i := 0; i < len(args)-1; i += 1 {
		arg := args[i]
		g.Expression(arg)
		g.gen(token.Comma)
		g.space()
	}

	last := args[len(args)-1]
	g.Expression(last)
	g.trailComma()

	g.sep()
	g.gen(token.RightParentheses)
}

func (g *Builder) IndexExpression(expr ast.IndexExpression) {
	g.Expression(expr)
	g.gen(token.LeftSquare)
	g.Expression(expr.Index)
	g.gen(token.RightSquare)
}

func (g *Builder) BinaryExpression(expr ast.BinaryExpression) {
	g.Expression(expr.Left)
	g.space()
	g.bop(expr.Operator)
	g.space()
	g.Expression(expr.Right)
}

func (g *Builder) BasicLiteral(lit ast.BasicLiteral) {
	g.tok(lit.Token)
}

func (g *Builder) ReceiverExpression(expr ast.Receiver) {
	g.genpos(token.Receiver, expr.Pos)
}

func (g *Builder) SymbolExpression(expr ast.SymbolExpression) {
	g.idn(expr.Identifier)
}

func (g *Builder) ParenthesizedExpression(expr ast.ParenthesizedExpression) {
	g.genpos(token.LeftParentheses, expr.Pos)
	g.sep()
	g.Expression(expr.Inner)
	g.sep()
	g.gen(token.RightParentheses)
}
