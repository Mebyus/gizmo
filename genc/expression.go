package genc

import (
	"fmt"
	"strconv"

	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/stg"
)

func (g *Builder) Exp(expr stg.Exp) {
	if expr == nil {
		panic("nil expression")
	}

	g.exp(expr)
}

func (g *Builder) exp(exp stg.Exp) {
	switch exp.Kind() {

	case exk.Integer:
		g.Integer(exp.(stg.Integer))
	case exk.String:
		g.String(exp.(stg.String))
	case exk.True:
		g.True()
	case exk.False:
		g.False()
	case exk.Symbol:
		g.SymbolExp(exp.(*stg.SymbolExp))
	case exk.Binary:
		g.BinaryExp(exp.(*stg.BinExp))
	case exk.Unary:
		g.UnaryExp(exp.(*stg.UnaryExp))
	case exk.Paren:
		g.ParenExp(exp.(*stg.ParenExp))
	case exk.Cast:
		g.CastExp(exp.(*stg.CastExp))
	case exk.MemCast:
		g.MemCastExp(exp.(*stg.MemCastExp))
	case exk.Tint:
		g.TintExp(exp.(*stg.TintExp))
	case exk.Enum:
		g.EnumExp(exp.(*stg.EnumExp))

	case exk.Chain, exk.Member, exk.Address, exk.Indirect, exk.IndirectIndex,
		exk.Call, exk.IndirectMember, exk.ChunkMember, exk.ChunkIndex, exk.ChunkSlice, exk.ArraySlice:

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

func (g *Builder) CastExp(exp *stg.CastExp) {
	if exp.DestType == stg.StrType {
		g.CastStringFromChunk(exp)
		return
	}

	g.puts("(")
	g.TypeSpec(exp.DestType)
	g.puts(")(")
	g.exp(exp.Target)
	g.puts(")")
}

func (g *Builder) CastStringFromChunk(exp *stg.CastExp) {
	g.puts("ku_str_from_chunk(")
	g.exp(exp.Target)
	g.puts(")")
}

func (g *Builder) MemCastExp(exp *stg.MemCastExp) {
	g.puts("__builtin_bit_cast(") // TODO: this does not work in C
	g.TypeSpec(exp.DestType)
	g.puts(", ")
	g.exp(exp.Target)
	g.puts(")")
}

func (g *Builder) ChainOperand(exp stg.ChainOperand) {
	switch exp.Kind() {
	case exk.Chain:
		g.ChainSymbol(exp.(*stg.ChainSymbol))
	case exk.Member:
		g.MemberExp(exp.(*stg.MemberExp))
	case exk.Address:
		g.AddressExp(exp.(*stg.AddressExp))
	case exk.Indirect:
		g.IndirectExp(exp.(*stg.IndirectExp))
	case exk.IndirectIndex:
		g.IndirectIndexExp(exp.(*stg.IndirectIndexExp))
	case exk.Call:
		g.CallExp(exp.(*stg.CallExp))
	case exk.IndirectMember:
		g.IndirectMemberExp(exp.(*stg.IndirectMemberExp))
	case exk.ChunkMember:
		g.ChunkMemberExp(exp.(*stg.ChunkMemberExp))
	case exk.ChunkIndex:
		g.ChunkIndexExp(exp.(*stg.ChunkIndexExp))
	case exk.ArraySlice:
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
		g.arrayFullSliceExp(exp)
		return
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

func (g *Builder) arrayFullSliceExp(exp *stg.ArraySliceExp) {
	g.ArrayTypeFullSliceMethodName(exp.Target.Type())
	g.puts("(&")
	g.exp(exp.Target)
	g.puts(")")
}

func (g *Builder) ChunkIndexExp(exp *stg.ChunkIndexExp) {
	g.ChunkTypeIndexMethodName(exp.Type())
	g.puts("(")
	g.exp(exp.Target)
	g.puts(", ")
	g.exp(exp.Index)
	g.puts(")")
}

func (g *Builder) ChunkIndirectElemExp(exp *stg.ChunkIndexExp) {
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

func (g *Builder) ChunkMemberExp(exp *stg.ChunkMemberExp) {
	g.exp(exp.Target)
	g.puts(".")
	g.puts(exp.Name)
}

func (g *Builder) ParenExp(exp *stg.ParenExp) {
	g.puts("(")
	g.exp(exp.Inner)
	g.puts(")")
}

func (g *Builder) UnaryExp(exp *stg.UnaryExp) {
	g.puts(exp.Operator.Kind.String())
	g.exp(exp.Inner)
}

func (g *Builder) IndirectIndexExp(exp *stg.IndirectIndexExp) {
	g.ChainOperand(exp.Target)
	g.puts("[")
	g.exp(exp.Index)
	g.puts("]")
}

func (g *Builder) IndirectExp(exp *stg.IndirectExp) {
	g.puts("*(")
	g.ChainOperand(exp.Target)
	g.puts(")")
}

func (g *Builder) AddressExp(exp *stg.AddressExp) {
	g.puts("&")
	g.ChainOperand(exp.Target)
}

func (g *Builder) MemberExp(exp *stg.MemberExp) {
	g.ChainOperand(exp.Target)
	g.puts(".")
	g.puts(exp.Member.Name)
}

func (g *Builder) IndirectMemberExp(exp *stg.IndirectMemberExp) {
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

func (g *Builder) CallArgs(args []stg.Exp) {
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

func (g *Builder) CallExp(expr *stg.CallExp) {
	g.ChainOperand(expr.Callee)
	g.CallArgs(expr.Arguments)
}

func (g *Builder) BinaryExp(expr *stg.BinExp) {
	g.exp(expr.Left)
	g.space()
	g.puts(expr.Operator.Kind.String())
	g.space()
	g.exp(expr.Right)
}

func (g *Builder) SymbolExp(expr *stg.SymbolExp) {
	g.SymbolName(expr.Sym)
}
