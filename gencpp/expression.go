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
	case exn.Address:
		g.AddressExpression(expr.(ast.AddressExpression))
	case exn.Cast:
		g.CastExpression(expr.(ast.CastExpression))
	case exn.Instance:
		g.InstanceExpression(expr.(ast.InstanceExpression))
	case exn.Index:
		g.IndexExpression(expr.(ast.IndexExpression))
	case exn.Slice:
		g.SliceExpression(expr.(ast.SliceExpression))
	case exn.BitCast:
		g.BitCastExpression(expr.(ast.BitCastExpression))
	case exn.Object:
		g.ObjectLiteral(expr.(ast.ObjectLiteral))
	default:
		g.write(fmt.Sprintf("<%s expr>", expr.Kind().String()))
	}
}

func (g *Builder) ObjectField(field ast.ObjectField) {
	g.write(".")
	g.Identifier(field.Name)
	g.write(" = ")
	g.Expression(field.Value)
}

func (g *Builder) ObjectLiteral(lit ast.ObjectLiteral) {
	if len(lit.Fields) == 0 {
		g.write("{}")
		return
	}

	g.write("{")
	g.nl()

	g.inc()

	g.indent()
	g.ObjectField(lit.Fields[0])
	g.write(",")
	g.nl()
	for _, field := range lit.Fields[1:] {
		g.indent()
		g.ObjectField(field)
		g.write(",")
		g.nl()
	}
	g.dec()

	g.indent()
	g.write("}")
}

func (g *Builder) SliceExpression(expr ast.SliceExpression) {
	g.Expression(expr.Target)
	if expr.Start == nil && expr.End == nil {
		g.write(".slice()")
		return
	}
	if expr.Start != nil && expr.End == nil {
		g.write(".tail(")
		g.Expression(expr.Start)
		g.write(")")
		return
	}
	if expr.Start == nil && expr.End != nil {
		g.write(".head(")
		g.Expression(expr.End)
		g.write(")")
		return
	}

	// start != nil && end != nil
	g.write(".slice(")
	g.Expression(expr.Start)
	g.write(", ")
	g.Expression(expr.End)
	g.write(")")
}

func (g *Builder) IndexExpression(expr ast.IndexExpression) {
	g.Expression(expr.Target)
	g.write("[")
	g.Expression(expr.Index)
	g.write("]")
}

func (g *Builder) InstanceExpression(expr ast.InstanceExpression) {
	g.ScopedIdentifier(expr.Target)
	g.templateArgs(expr.Args)
}

func (g *Builder) templateArgs(args []ast.TypeSpecifier) {
	g.write("<")
	g.TypeSpecifier(args[0])
	for _, arg := range args[1:] {
		g.write(", ")
		g.TypeSpecifier(arg)
	}
	g.write(">")
}

func (g *Builder) CastExpression(expr ast.CastExpression) {
	g.write("(")
	g.TypeSpecifier(expr.Type)
	g.write(")(")
	g.Expression(expr.Target)
	g.write(")")
}

func (g *Builder) BitCastExpression(expr ast.BitCastExpression) {
	g.write("bit_cast(")
	g.TypeSpecifier(expr.Type)
	g.write(", ")
	g.Expression(expr.Target)
	g.write(")")
}

func (g *Builder) AddressExpression(expr ast.AddressExpression) {
	g.write("&")
	g.Expression(expr.Target)
}

const maxInt64 = 0x7fff_ffff_ffff_ffff

func (g *Builder) BasicLiteral(lit ast.BasicLiteral) {
	if lit.Token.Kind == token.Nil {
		g.write("nullptr")
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
		g.write(strconv.FormatUint(lit.Token.Val, 10))
		g.write(")")
		return
	}
	if lit.Token.Kind == token.OctalInteger {
		g.write("0" + strconv.FormatUint(lit.Token.Val, 8))
		return
	}
	if lit.Token.Kind == token.DecimalInteger && lit.Token.Val > maxInt64 {
		g.write(strconv.FormatUint(lit.Token.Val, 10))
		g.write("ul")
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
