package format

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/token"
)

func (g *Noder) Expression(expr ast.Exp) {
	switch expr.Kind() {
	case exn.Symbol:
		g.SymbolExpression(expr.(ast.SymbolExp))
	case exn.Basic:
		g.BasicLiteral(expr.(ast.BasicLiteral))
	case exn.Binary:
		g.BinaryExpression(expr.(ast.BinExp))
	case exn.Paren:
		g.ParenthesizedExpression(expr.(ast.ParenthesizedExpression))
	case exn.Chain:
		g.ChainOperand(expr.(ast.ChainOperand))
	default:
		panic(fmt.Sprintf("%s expression node not implemented", expr.Kind().String()))
	}
}

func (g *Noder) ChainOperand(expr ast.ChainOperand) {
	g.idn(expr.Identifier)

	for _, part := range expr.Parts {
		g.chainPart(part)
	}
}

func (g *Noder) chainPart(part ast.ChainPart) {
	switch part.Kind() {
	case exn.Call:
		g.callPart(part.(ast.CallPart))
	case exn.Address:
		g.addressPart(part.(ast.AddressPart))
	case exn.Index:
		g.indexPart(part.(ast.IndexPart))
	case exn.Indirect:
		g.indirectPart(part.(ast.IndirectPart))
	default:
		panic(fmt.Sprintf("%s chain part node not implemented", part.Kind().String()))
	}
}

func (g *Noder) indirectPart(expr ast.IndirectPart) {
	g.genpos(token.Indirect, expr.Pos)
}

func (g *Noder) addressPart(expr ast.AddressPart) {
	g.genpos(token.Address, expr.Pos)
}

func (g *Noder) callPart(expr ast.CallPart) {
	args := expr.Args
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

func (g *Noder) indexPart(expr ast.IndexPart) {
	g.gen(token.LeftSquare)
	g.Expression(expr.Index)
	g.gen(token.RightSquare)
}

func (g *Noder) BinaryExpression(expr ast.BinExp) {
	g.Expression(expr.Left)
	g.space()
	g.bop(expr.Operator)
	g.space()
	g.Expression(expr.Right)
}

func (g *Noder) BasicLiteral(lit ast.BasicLiteral) {
	g.tok(lit.Token)
}

func (g *Noder) SymbolExpression(expr ast.SymbolExp) {
	g.idn(expr.Identifier)
}

func (g *Noder) ParenthesizedExpression(expr ast.ParenthesizedExpression) {
	g.genpos(token.LeftParentheses, expr.Pos)
	g.sep()
	g.Expression(expr.Inner)
	g.sep()
	g.gen(token.RightParentheses)
}
