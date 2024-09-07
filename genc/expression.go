package genc

import (
	"fmt"
	"strconv"

	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/stg"
)

func (g *Builder) Exp(expr stg.Expression) {
	if expr == nil {
		panic("nil expression")
	}

	g.exp(expr)
}

func (g *Builder) exp(exp stg.Expression) {
	switch exp.Kind() {

	case exn.Integer:
		g.Integer(exp.(stg.Integer))
	case exn.String:
		g.String(exp.(stg.String))
	case exn.True:
		g.True()
	case exn.False:
		g.False()
	case exn.Symbol:
		g.SymbolExp(exp.(*stg.SymbolExpression))
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
	case exn.Enum:
		g.EnumExp(exp.(*stg.EnumExp))

	case exn.Chain, exn.Member, exn.Address, exn.Indirect, exn.IndirectIndex,
		exn.Call, exn.IndirectMember, exn.ChunkMember, exn.ChunkIndex, exn.ChunkSlice, exn.ArraySlice:

		g.ChainOperand(exp.(stg.ChainOperand))

	default:
		panic(fmt.Sprintf("%s expression not implemented", exp.Kind().String()))
	}
}

func (g *Builder) String(exp stg.String) {
	g.puts("ku_static_string((const u8*)(u8\"")
	g.puts(exp.Val)
	g.puts("\"), ")
	g.putn(exp.Size)
	g.puts(")")
}

func (g *Builder) True() {
	g.puts("true")
}

func (g *Builder) False() {
	g.puts("false")
}

func (g *Builder) EnumExp(exp *stg.EnumExp) {
	symbol := exp.Enum.Def.(stg.CustomTypeDef).Symbol
	name := getEnumEntryName(symbol.Name, exp.Entry.Name)
	g.puts(name)
}

func (g *Builder) TintExp(exp *stg.TintExp) {
	g.puts("(")
	g.TypeSpec(exp.DestType)
	g.puts(")(")
	g.exp(exp.Target)
	g.puts(")")
}

func (g *Builder) CastExp(expr *stg.CastExpression) {
	g.puts("(")
	g.TypeSpec(expr.DestType)
	g.puts(")(")
	g.exp(expr.Target)
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
	g.exp(exp.Target)
	g.puts(", ")
	g.exp(exp.End)
	g.puts(")")
}

func (g *Builder) ChunkIndexExp(exp *stg.ChunkIndexExpression) {
	g.ChunkTypeIndexMethodName(exp.Type())
	g.puts("(")
	g.exp(exp.Target)
	g.puts(", ")
	g.exp(exp.Index)
	g.puts(")")
}

func (g *Builder) ChunkIndirectElemExp(exp *stg.ChunkIndexExpression) {
	g.puts("*(")
	g.ChunkTypeElemMethodName(exp.Type())
	g.puts("(")
	g.exp(exp.Target)
	g.puts(", ")
	g.exp(exp.Index)
	g.puts("))")
}

func (g *Builder) ArrayIndirectElemExp(exp *stg.ArrayIndexExp) {
	g.puts("*(")
	g.ArrayTypeElemMethodName(exp.Target.Type())
	g.puts("(&")
	g.exp(exp.Target)
	g.puts(", ")
	g.exp(exp.Index)
	g.puts("))")
}

func (g *Builder) ChunkMemberExp(exp *stg.ChunkMemberExpression) {
	g.exp(exp.Target)
	g.puts(".")
	g.puts(exp.Name)
}

func (g *Builder) ParenExp(exp *stg.ParenthesizedExpression) {
	g.puts("(")
	g.exp(exp.Inner)
	g.puts(")")
}

func (g *Builder) UnaryExp(exp *stg.UnaryExpression) {
	g.puts(exp.Operator.Kind.String())
	g.exp(exp.Inner)
}

func (g *Builder) IndirectIndexExp(exp *stg.IndirectIndexExpression) {
	g.ChainOperand(exp.Target)
	g.puts("[")
	g.exp(exp.Index)
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
	g.exp(args[0])
	for _, arg := range args[1:] {
		g.puts(", ")
		g.exp(arg)
	}
	g.puts(")")
}

func (g *Builder) CallExp(expr *stg.CallExpression) {
	g.ChainOperand(expr.Callee)
	g.CallArgs(expr.Arguments)
}

func (g *Builder) BinaryExp(expr *stg.BinaryExpression) {
	g.exp(expr.Left)
	g.space()
	g.puts(expr.Operator.Kind.String())
	g.space()
	g.exp(expr.Right)
}

func (g *Builder) SymbolExp(expr *stg.SymbolExpression) {
	g.SymbolName(expr.Sym)
}
