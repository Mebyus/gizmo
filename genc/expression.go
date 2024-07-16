package genc

import (
	"fmt"
	"strconv"

	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/tt"
)

func (g *Builder) Expression(expr tt.Expression) {
	if expr == nil {
		panic("nit expression")
	}

	g.expr(expr)
}

func (g *Builder) expr(expr tt.Expression) {
	switch expr.Kind() {
	case exn.Integer:
		g.Integer(expr.(tt.Integer))
	case exn.Symbol:
		g.SymbolExpression(expr.(*tt.SymbolExpression))
	case exn.Binary:
		g.BinaryExpression(expr.(*tt.BinaryExpression))
	case exn.SymbolCall:
		g.SymbolCallExpression(expr.(*tt.SymbolCallExpression))
	case exn.Member:
		g.MemberExpression(expr.(*tt.MemberExpression))
	default:
		panic(fmt.Sprintf("%s expression not implemented", expr.Kind().String()))
	}
}

func (g *Builder) MemberExpression(expr *tt.MemberExpression) {
	g.SymbolName(expr.Target)
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

func (g *Builder) SymbolCallExpression(expr *tt.SymbolCallExpression) {
	g.SymbolName(expr.Callee)
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
