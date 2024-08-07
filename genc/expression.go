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

func (g *Builder) expr(exp stg.Expression) {
	switch exp.Kind() {

	case exn.Integer:
		g.Integer(exp.(stg.Integer))
	case exn.True:
		g.True()
	case exn.False:
		g.False()
	case exn.Symbol:
		g.SymbolExp(exp.(*stg.SymbolExpression))
	case exn.Receiver:
		g.puts("g")
	case exn.Binary:
		g.BinaryExp(exp.(*stg.BinaryExpression))
	case exn.Unary:
		g.UnaryExp(exp.(*stg.UnaryExpression))
	case exn.Paren:
		g.ParenExp(exp.(*stg.ParenthesizedExpression))
	case exn.Cast:
		g.CastExp(exp.(*stg.CastExpression))
	case exn.Tint:
		g.TintExp(exp.(*stg.TintExp))
	case exn.Chain, exn.Member, exn.Address, exn.Indirect, exn.IndirectIndex,
		exn.Call, exn.IndirectMember, exn.ChunkMember, exn.ChunkIndex, exn.ChunkSlice, exn.ArraySlice:

		g.ChainOperand(exp.(stg.ChainOperand))

	default:
		panic(fmt.Sprintf("%s expression not implemented", exp.Kind().String()))
	}
}

func (g *Builder) True() {
	g.puts("true")
}

func (g *Builder) False() {
	g.puts("false")
}

func (g *Builder) TintExp(exp *stg.TintExp) {
	g.puts("(")
	g.TypeSpec(exp.DestType)
	g.puts(")(")
	g.expr(exp.Target)
	g.puts(")")
}

func (g *Builder) CastExp(expr *stg.CastExpression) {
	g.puts("(")
	g.TypeSpec(expr.DestType)
	g.puts(")(")
	g.expr(expr.Target)
	g.puts(")")
}

func (g *Builder) ChainOperand(exp stg.ChainOperand) {
	switch exp.Kind() {
	case exn.Chain:
		g.ChainSymbol(exp.(*stg.ChainSymbol))
	case exn.Member:
		g.MemberExp(exp.(*stg.MemberExpression))
	case exn.Address:
		g.AddressExp(exp.(*stg.AddressExpression))
	case exn.Indirect:
		g.IndirectExp(exp.(*stg.IndirectExpression))
	case exn.IndirectIndex:
		g.IndirectIndexExp(exp.(*stg.IndirectIndexExpression))
	case exn.Call:
		g.CallExp(exp.(*stg.CallExpression))
	case exn.IndirectMember:
		g.IndirectMemberExp(exp.(*stg.IndirectMemberExpression))
	case exn.ChunkMember:
		g.ChunkMemberExp(exp.(*stg.ChunkMemberExpression))
	case exn.ChunkIndex:
		g.ChunkIndexExp(exp.(*stg.ChunkIndexExpression))
	case exn.ArraySlice:
		g.ArraySliceExp(exp.(*stg.ArraySliceExp))
	default:
		panic(fmt.Sprintf("%s operand not implemented", exp.Kind().String()))
	}
}

func (g *Builder) ChainSymbol(exp *stg.ChainSymbol) {
	if exp.Sym == nil {
		g.puts("g")
		return
	}
	g.SymbolName(exp.Sym)
}

func (g *Builder) ArraySliceExp(exp *stg.ArraySliceExp) {
	if exp.Start == nil && exp.End == nil {
		// full array slice
		panic("not implemented")
	}
	if exp.Start == nil && exp.End != nil {
		g.arrayHeadSliceExp(exp)
		return
	}
	if exp.Start != nil && exp.End == nil {
		// tail slice
		panic("not implemented")
	}
	if exp.Start != nil && exp.End != nil {
		// regular slice
		panic("not implemented")
	}

	panic("impossible condition")
}

func (g *Builder) arrayHeadSliceExp(exp *stg.ArraySliceExp) {
	g.ArrayTypeHeadSliceMethodName(exp.Target.Type())
	g.puts("(&")
	g.expr(exp.Target)
	g.puts(", ")
	g.expr(exp.End)
	g.puts(")")
}

func (g *Builder) ChunkIndexExp(exp *stg.ChunkIndexExpression) {
	g.ChunkTypeIndexMethodName(exp.Type())
	g.puts("(")
	g.expr(exp.Target)
	g.puts(", ")
	g.expr(exp.Index)
	g.puts(")")
}

func (g *Builder) ChunkIndirectElemExp(exp *stg.ChunkIndexExpression) {
	g.puts("*(")
	g.ChunkTypeElemMethodName(exp.Type())
	g.puts("(")
	g.expr(exp.Target)
	g.puts(", ")
	g.expr(exp.Index)
	g.puts("))")
}

func (g *Builder) ArrayIndirectElemExp(exp *stg.ArrayIndexExp) {
	g.puts("*(")
	g.ArrayTypeElemMethodName(exp.Target.Type())
	g.puts("(&")
	g.expr(exp.Target)
	g.puts(", ")
	g.expr(exp.Index)
	g.puts("))")
}

func (g *Builder) ChunkMemberExp(exp *stg.ChunkMemberExpression) {
	g.expr(exp.Target)
	g.puts(".")
	g.puts(exp.Name)
}

func (g *Builder) ParenExp(exp *stg.ParenthesizedExpression) {
	g.puts("(")
	g.expr(exp.Inner)
	g.puts(")")
}

func (g *Builder) UnaryExp(exp *stg.UnaryExpression) {
	g.puts(exp.Operator.Kind.String())
	g.expr(exp.Inner)
}

func (g *Builder) IndirectIndexExp(exp *stg.IndirectIndexExpression) {
	g.ChainOperand(exp.Target)
	g.puts("[")
	g.expr(exp.Index)
	g.puts("]")
}

func (g *Builder) IndirectExp(exp *stg.IndirectExpression) {
	g.puts("*(")
	g.ChainOperand(exp.Target)
	g.puts(")")
}

func (g *Builder) AddressExp(exp *stg.AddressExpression) {
	g.puts("&")
	g.ChainOperand(exp.Target)
}

func (g *Builder) MemberExp(exp *stg.MemberExpression) {
	g.ChainOperand(exp.Target)
	g.puts(".")
	g.puts(exp.Member.Name)
}

func (g *Builder) IndirectMemberExp(exp *stg.IndirectMemberExpression) {
	g.ChainOperand(exp.Target)
	g.puts("->")
	g.puts(exp.Member.Name)
}

func (g *Builder) Integer(exp stg.Integer) {
	if exp.Neg {
		g.puts("-")
	}
	g.puts(strconv.FormatUint(exp.Val, 10))
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

func (g *Builder) CallExp(expr *stg.CallExpression) {
	g.ChainOperand(expr.Callee)
	g.CallArgs(expr.Arguments)
}

func (g *Builder) BinaryExp(expr *stg.BinaryExpression) {
	g.expr(expr.Left)
	g.space()
	g.puts(expr.Operator.Kind.String())
	g.space()
	g.expr(expr.Right)
}

func (g *Builder) SymbolExp(expr *stg.SymbolExpression) {
	g.SymbolName(expr.Sym)
}
