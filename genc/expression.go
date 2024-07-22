package genc

import (
	"fmt"
	"strconv"

	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/tt"
)

func (g *Builder) Expression(expr tt.Expression) {
	if expr == nil {
		panic("nil expression")
	}

	g.expr(expr)
}

func (g *Builder) expr(expr tt.Expression) {
	switch expr.Kind() {

	case exn.Integer:
		g.Integer(expr.(tt.Integer))
	case exn.True:
		g.True()
	case exn.False:
		g.False()
	case exn.Symbol:
		g.SymbolExpression(expr.(*tt.SymbolExpression))
	case exn.Binary:
		g.BinaryExpression(expr.(*tt.BinaryExpression))
	case exn.Unary:
		g.UnaryExpression(expr.(*tt.UnaryExpression))
	case exn.Paren:
		g.ParenthesizedExpression(expr.(*tt.ParenthesizedExpression))
	case exn.Cast:
		g.CastExpression(expr.(*tt.CastExpression))
	case exn.Chain, exn.Member, exn.Address, exn.Indirect, exn.IndirectIndex:
		g.ChainOperand(expr.(tt.ChainOperand))

	default:
		panic(fmt.Sprintf("%s expression not implemented", expr.Kind().String()))
	}
}

func (g *Builder) True() {
	g.puts("true")
}

func (g *Builder) False() {
	g.puts("false")
}

func (g *Builder) CastExpression(expr *tt.CastExpression) {
	g.puts("(")
	g.TypeSpec(expr.DestinationType)
	g.puts(")(")
	g.expr(expr.Target)
	g.puts(")")
}

func (g *Builder) ChainOperand(expr tt.ChainOperand) {
	switch expr.Kind() {
	case exn.Chain:
		g.ChainSymbol(expr.(*tt.ChainSymbol))
	case exn.Member:
		g.MemberExpression(expr.(*tt.MemberExpression))
	case exn.Address:
		g.AddressExpression(expr.(*tt.AddressExpression))
	case exn.Indirect:
		g.IndirectExpression(expr.(*tt.IndirectExpression))
	case exn.IndirectIndex:
		g.IndirectIndexExpression(expr.(*tt.IndirectIndexExpression))
	default:
		panic(fmt.Sprintf("%s operand not implemented", expr.Kind().String()))
	}
}

func (g *Builder) ChainSymbol(expr *tt.ChainSymbol) {
	if expr.Sym == nil {
		g.puts("g")
		return
	}
	g.SymbolName(expr.Sym)
}

func (g *Builder) ParenthesizedExpression(expr *tt.ParenthesizedExpression) {
	g.puts("(")
	g.expr(expr.Inner)
	g.puts(")")
}

func (g *Builder) UnaryExpression(expr *tt.UnaryExpression) {
	g.puts(expr.Operator.Kind.String())
	g.expr(expr.Inner)
}

func (g *Builder) IndirectIndexExpression(expr *tt.IndirectIndexExpression) {
	g.ChainOperand(expr.Target)
	g.puts("[")
	g.expr(expr.Index)
	g.puts("]")
}

func (g *Builder) IndirectExpression(expr *tt.IndirectExpression) {
	g.puts("*(")
	g.ChainOperand(expr.Target)
	g.puts(")")
}

func (g *Builder) AddressExpression(expr *tt.AddressExpression) {
	g.puts("&")
	g.ChainOperand(expr.Target)
}

func (g *Builder) MemberExpression(expr *tt.MemberExpression) {
	g.ChainOperand(expr.Target)
	g.puts(".")
	g.puts(expr.Member.Name)
}

func (g *Builder) Integer(expr tt.Integer) {
	g.puts(strconv.FormatUint(expr.Val, 10))
}

func (g *Builder) CallArgs(args []tt.Expression) {
	if len(args) == 0 {
		g.puts("()")
		return
	}

	g.puts("(")
	g.expr(args[0])
	for _, arg := range args[1:] {
		g.puts(", ")
		g.expr(arg)
	}
	g.puts(")")
}

func (g *Builder) CallExpression(expr *tt.CallExpression) {
	g.ChainOperand(expr.Callee)
	g.CallArgs(expr.Arguments)
}

func (g *Builder) BinaryExpression(expr *tt.BinaryExpression) {
	g.expr(expr.Left)
	g.space()
	g.puts(expr.Operator.Kind.String())
	g.space()
	g.expr(expr.Right)
}

func (g *Builder) SymbolExpression(expr *tt.SymbolExpression) {
	g.SymbolName(expr.Sym)
}
