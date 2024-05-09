package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/token"
	"github.com/mebyus/gizmo/tt/scp"
	"github.com/mebyus/gizmo/tt/typ"
)

// Scan constructs expression from a given AST. Uses current scope for symbol
// resolution.
func (s *Scope) Scan(ctx *Context, expr ast.Expression) (Expression, error) {
	if expr == nil {
		return nil, nil
	}
	if expr.Kind() == exn.Start {
		panic("chain start cannot appear as standalone expression")
	}

	return s.scan(ctx, expr)
}

func (s *Scope) scan(ctx *Context, expr ast.Expression) (Expression, error) {
	switch expr.Kind() {
	case exn.Basic:
		return scanBasicLiteral(expr.(ast.BasicLiteral)), nil
	case exn.Symbol:
		return s.scanSymbolExpression(ctx, expr.(ast.SymbolExpression))
	case exn.Start:
		return s.scanChainStart(ctx, expr.(ast.ChainStart))
	case exn.Indirect:
		return s.scanIndirectExpression(ctx, expr.(ast.IndirectExpression))
	case exn.Unary:
		return s.scanUnaryExpression(ctx, expr.(*ast.UnaryExpression))
	case exn.Binary:
		return s.scanBinaryExpression(ctx, expr.(ast.BinaryExpression))
	case exn.Call:
		return s.scanCallExpression(ctx, expr.(ast.CallExpression))
	// case exn.Indirx:
	// 	// g.IndirectIndexExpression(expr.(ast.IndirectIndexExpression))
	case exn.Paren:
		return s.scanParenthesizedExpression(ctx, expr.(ast.ParenthesizedExpression))
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

func (s *Scope) scanParenthesizedExpression(ctx *Context, expr ast.ParenthesizedExpression) (*ParenthesizedExpression, error) {
	pos := expr.Pos
	inner, err := s.scan(ctx, expr.Inner)
	if err != nil {
		return nil, err
	}

	return &ParenthesizedExpression{
		Pos:   pos,
		Inner: inner,
	}, nil
}

func (s *Scope) scanSymbolExpression(ctx *Context, expr ast.SymbolExpression) (*SymbolExpression, error) {
	name := expr.Identifier.Lit
	pos := expr.Identifier.Pos
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
	inner, err := s.scan(ctx, expr.Inner)
	if err != nil {
		return nil, err
	}

	return &UnaryExpression{
		Operator: UnaryOperator(expr.Operator),
		Inner:    inner,
	}, nil
}

func (s *Scope) scanBinaryExpression(ctx *Context, expr ast.BinaryExpression) (*BinaryExpression, error) {
	left, err := s.scan(ctx, expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := s.scan(ctx, expr.Right)
	if err != nil {
		return nil, err
	}

	return &BinaryExpression{
		Operator: BinaryOperator(expr.Operator),
		Left:     left,
		Right:    right,
	}, nil
}

func (s *Scope) scanChainStart(ctx *Context, start ast.ChainStart) (*ChainStart, error) {
	name := start.Identifier.Lit
	pos := start.Identifier.Pos
	symbol := s.Lookup(name, pos.Num)
	if symbol == nil {
		return nil, fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}
	if symbol.Scope.Kind == scp.Unit {
		ctx.ref.Add(symbol)
	}

	return &ChainStart{
		Pos: pos,
		Sym: symbol,
	}, nil
}

func (s *Scope) scanIndirectExpression(ctx *Context, expr ast.IndirectExpression) (*IndirectExpression, error) {
	tg, err := s.scan(ctx, expr.Target)
	if err != nil {
		return nil, err
	}
	target := tg.(ChainOperand)
	targetType := target.Type()
	if targetType.Kind != typ.Pointer {
		return nil, fmt.Errorf("%s: invalid operation (indirect on non-pointer type)", expr.Pos)
	}
	return &IndirectExpression{
		Pos:        expr.Pos,
		Target:     target,
		ChainDepth: expr.ChainDepth,

		typ: targetType.Def.(PtrTypeDef).RefType,
	}, nil
}

func (s *Scope) scanCallExpression(ctx *Context, expr ast.CallExpression) (*CallExpression, error) {
	return &CallExpression{
		ChainDepth: expr.ChainDepth,
	}, nil
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
