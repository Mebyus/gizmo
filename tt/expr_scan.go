package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/token"
	"github.com/mebyus/gizmo/tt/scp"
)

// Scan constructs expression from a given AST. Uses current scope for symbol
// resolution.
func (s *Scope) Scan(ctx *Context, expr ast.Expression) (Expression, error) {
	if expr == nil {
		return nil, nil
	}

	switch expr.Kind() {
	// case exn.Start:
	// 	// g.ScopedIdentifier(expr.(ast.ChainStart).Identifier)
	case exn.Symbol:
		return s.scanSymbolExpression(ctx, expr.(ast.SymbolExpression))
	case exn.Basic:
		return scanBasicLiteral(expr.(ast.BasicLiteral)), nil
	// case exn.Indirect:
	// 	// g.IndirectExpression(expr.(ast.IndirectExpression))
	case exn.Unary:
		return s.scanUnaryExpression(ctx, expr.(*ast.UnaryExpression))
	case exn.Binary:
		return s.scanBinaryExpression(ctx, expr.(ast.BinaryExpression))
	// case exn.Call:
	// 	// g.CallExpression(expr.(ast.CallExpression))
	// case exn.Indirx:
	// 	// g.IndirectIndexExpression(expr.(ast.IndirectIndexExpression))
	// case exn.Paren:
	// 	// g.ParenthesizedExpression(expr.(ast.ParenthesizedExpression))
	// case exn.Select:
	// 	// g.SelectorExpression(expr.(ast.SelectorExpression))
	// case exn.Address:
	// 	// g.AddressExpression(expr.(ast.AddressExpression))
	// case exn.Cast:
	// 	// g.CastExpression(expr.(ast.CastExpression))
	// case exn.Instance:
	// 	// g.InstanceExpression(expr.(ast.InstanceExpression))
	// case exn.Index:
	// 	// g.IndexExpression(expr.(ast.IndexExpression))
	// case exn.Slice:
	// 	// g.SliceExpression(expr.(ast.SliceExpression))
	// case exn.BitCast:
	// 	// g.BitCastExpression(expr.(ast.BitCastExpression))
	// case exn.Object:
	// 	// g.ObjectLiteral(expr.(ast.ObjectLiteral))
	default:
		panic(fmt.Sprintf("not implemented for %s expression", expr.Kind().String()))
	}
}

func (s *Scope) scanSymbolExpression(ctx *Context, expr ast.SymbolExpression) (*SymbolExpression, error) {
	name := expr.Identifier.Name.Lit
	pos := expr.Identifier.Name.Pos
	symbol := s.Lookup(name, pos.Num)
	if symbol == nil {
		return nil, fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}
	if symbol.Scope.Kind == scp.Unit {
		ctx.ref.Add(symbol)
	}

	return &SymbolExpression{
		Pos: pos,
		Sym: symbol,
	}, nil
}

func (s *Scope) scanUnaryExpression(ctx *Context, expr *ast.UnaryExpression) (*UnaryExpression, error) {
	return &UnaryExpression{}, nil
}

func (s *Scope) scanBinaryExpression(ctx *Context, expr ast.BinaryExpression) (*BinaryExpression, error) {
	return &BinaryExpression{}, nil
}

func scanBasicLiteral(lit ast.BasicLiteral) Literal {
	pos := lit.Token.Pos

	switch lit.Token.Kind {
	case token.True:
		return True{Pos: pos}
	case token.False:
		return False{Pos: pos}
	case token.BinaryInteger, token.OctalInteger, token.DecimalInteger, token.HexadecimalInteger:
		return Integer{Pos: pos, Val: lit.Token.Val}
	case token.String:
		return String{Pos: pos, Val: lit.Token.Lit}
	default:
		panic(fmt.Sprintf("not implemented for %s literal", lit.Token.Kind.String()))
	}
}
