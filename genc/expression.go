package genc

import (
	"fmt"
	"strconv"

	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/stg"
)

func (g *Builder) Expression(expr stg.Expression) {
	if expr == nil {
		panic("nil expression")
	}

	g.expr(expr)
}

func (g *Builder) expr(expr stg.Expression) {
	switch expr.Kind() {

	case exn.Integer:
		g.Integer(expr.(stg.Integer))
	case exn.True:
		g.True()
	case exn.False:
		g.False()
	case exn.Symbol:
		g.SymbolExpression(expr.(*stg.SymbolExpression))
	case exn.Binary:
		g.BinaryExpression(expr.(*stg.BinaryExpression))
	case exn.Unary:
		g.UnaryExpression(expr.(*stg.UnaryExpression))
	case exn.Paren:
		g.ParenthesizedExpression(expr.(*stg.ParenthesizedExpression))
	case exn.Cast:
		g.CastExpression(expr.(*stg.CastExpression))
	case exn.Chain, exn.Member, exn.Address, exn.Indirect, exn.IndirectIndex,
		exn.Call, exn.IndirectMember, exn.ChunkMember, exn.ChunkIndex:

		g.ChainOperand(expr.(stg.ChainOperand))

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

func (g *Builder) CastExpression(expr *stg.CastExpression) {
	g.puts("(")
	g.TypeSpec(expr.DestinationType)
	g.puts(")(")
	g.expr(expr.Target)
	g.puts(")")
}

func (g *Builder) ChainOperand(expr stg.ChainOperand) {
	switch expr.Kind() {
	case exn.Chain:
		g.ChainSymbol(expr.(*stg.ChainSymbol))
	case exn.Member:
		g.MemberExpression(expr.(*stg.MemberExpression))
	case exn.Address:
		g.AddressExpression(expr.(*stg.AddressExpression))
	case exn.Indirect:
		g.IndirectExpression(expr.(*stg.IndirectExpression))
	case exn.IndirectIndex:
		g.IndirectIndexExpression(expr.(*stg.IndirectIndexExpression))
	case exn.Call:
		g.CallExpression(expr.(*stg.CallExpression))
	case exn.IndirectMember:
		g.IndirectMemberExpression(expr.(*stg.IndirectMemberExpression))
	case exn.ChunkMember:
		g.ChunkMemberExpression(expr.(*stg.ChunkMemberExpression))
	case exn.ChunkIndex:
		g.ChunkIndexExpression(expr.(*stg.ChunkIndexExpression))
	default:
		panic(fmt.Sprintf("%s operand not implemented", expr.Kind().String()))
	}
}

func (g *Builder) ChainSymbol(expr *stg.ChainSymbol) {
	if expr.Sym == nil {
		g.puts("g")
		return
	}
	g.SymbolName(expr.Sym)
}

func (g *Builder) ChunkIndexExpression(expr *stg.ChunkIndexExpression) {
	g.ChunkTypeIndexMethodName(expr.Type())
	g.puts("(")
	g.expr(expr.Target)
	g.puts(", ")
	g.expr(expr.Index)
	g.puts(")")
}

func (g *Builder) ChunkMemberExpression(expr *stg.ChunkMemberExpression) {
	g.expr(expr.Target)
	g.puts(".")
	g.puts(expr.Name)
}

func (g *Builder) ParenthesizedExpression(expr *stg.ParenthesizedExpression) {
	g.puts("(")
	g.expr(expr.Inner)
	g.puts(")")
}

func (g *Builder) UnaryExpression(expr *stg.UnaryExpression) {
	g.puts(expr.Operator.Kind.String())
	g.expr(expr.Inner)
}

func (g *Builder) IndirectIndexExpression(expr *stg.IndirectIndexExpression) {
	g.ChainOperand(expr.Target)
	g.puts("[")
	g.expr(expr.Index)
	g.puts("]")
}

func (g *Builder) IndirectExpression(expr *stg.IndirectExpression) {
	g.puts("*(")
	g.ChainOperand(expr.Target)
	g.puts(")")
}

func (g *Builder) AddressExpression(expr *stg.AddressExpression) {
	g.puts("&")
	g.ChainOperand(expr.Target)
}

func (g *Builder) MemberExpression(expr *stg.MemberExpression) {
	g.ChainOperand(expr.Target)
	g.puts(".")
	g.puts(expr.Member.Name)
}

func (g *Builder) IndirectMemberExpression(expr *stg.IndirectMemberExpression) {
	g.ChainOperand(expr.Target)
	g.puts("->")
	g.puts(expr.Member.Name)
}

func (g *Builder) Integer(expr stg.Integer) {
	g.puts(strconv.FormatUint(expr.Val, 10))
}

func (g *Builder) CallArgs(args []stg.Expression) {
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

func (g *Builder) CallExpression(expr *stg.CallExpression) {
	g.ChainOperand(expr.Callee)
	g.CallArgs(expr.Arguments)
}

func (g *Builder) BinaryExpression(expr *stg.BinaryExpression) {
	g.expr(expr.Left)
	g.space()
	g.puts(expr.Operator.Kind.String())
	g.space()
	g.expr(expr.Right)
}

func (g *Builder) SymbolExpression(expr *stg.SymbolExpression) {
	g.SymbolName(expr.Sym)
}
