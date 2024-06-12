package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/token"
	"github.com/mebyus/gizmo/tt/scp"
	"github.com/mebyus/gizmo/tt/sym"
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
	case exn.SymbolCall:
		return s.scanSymbolCallExpression(ctx, expr.(ast.SymbolCallExpression))
	// case exn.Indirx:
	// 	// g.IndirectIndexExpression(expr.(ast.IndirectIndexExpression))
	case exn.Paren:
		return s.scanParenthesizedExpression(ctx, expr.(ast.ParenthesizedExpression))
	// case exn.Select:
	// 	// g.SelectorExpression(expr.(ast.SelectorExpression))
	case exn.SymbolAddress:
		return s.scanSymbolAddressExpression(ctx, expr.(ast.SymbolAddressExpression))
	case exn.MemberCall:
		return s.scanMemberCallExpression(ctx, expr.(ast.MemberCallExpression))
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

func (s *Scope) scanMemberCallExpression(ctx *Context, expr ast.MemberCallExpression) (Expression, error) {
	name := expr.Target.Lit
	pos := expr.Target.Pos
	symbol := s.Lookup(name, pos.Num)
	if symbol == nil {
		return nil, fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}

	if symbol.Scope.Kind == scp.Unit && symbol.Kind != sym.Import {
		ctx.ref.Add(symbol)
	}

	if symbol.Kind == sym.Import {
		// call to imported function
		panic("not implemented")
	}

	if symbol.Kind == sym.Type {
		// call to function namespaced within type
		panic("not implemented")
	}

	if !(symbol.Kind == sym.Let || symbol.Kind == sym.Param || symbol.Kind == sym.Var) {
		return nil, fmt.Errorf("%s: %s symbol \"%s\" is not selectable", pos.String(), symbol.Kind.String(), name)
	}

	switch symbol.Type.Base.Kind {
	case typ.Struct:
	case typ.Pointer:
		refType := symbol.Type.Base.Def.(PtrTypeDef).RefType
		if refType.Base.Kind != typ.Struct {
			return nil, fmt.Errorf("%s: symbol \"%s\" is a pointer to %s type which cannot have members", pos.String(),
				name, refType.Base.Kind)
		}
	default:
		return nil, fmt.Errorf("%s: symbol \"%s\" is of %s type which cannot have members", pos.String(),
			name, symbol.Type.Base.Kind.String())
	}

	return nil, nil
}

func (s *Scope) scanSymbolAddressExpression(ctx *Context, expr ast.SymbolAddressExpression) (*SymbolAddressExpression, error) {
	name := expr.Target.Lit
	pos := expr.Target.Pos
	symbol := s.Lookup(name, pos.Num)
	if symbol == nil {
		return nil, fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}

	if !symbol.Kind.Addressable() {
		return nil, fmt.Errorf("%s: symbol \"%s\" (%s) is not addressable", pos.String(), name, symbol.Kind.String())
	}

	if symbol.Scope.Kind == scp.Unit {
		ctx.ref.Add(symbol)
	}

	return &SymbolAddressExpression{
		Pos:    pos,
		Target: symbol,
	}, nil
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

func (s *Scope) scanSymbolCallExpression(ctx *Context, expr ast.SymbolCallExpression) (*SymbolCallExpression, error) {
	name := expr.Callee.Lit
	pos := expr.Callee.Pos
	symbol := s.Lookup(name, pos.Num)
	if symbol == nil {
		return nil, fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}
	if symbol.Kind != sym.Fn {
		return nil, fmt.Errorf("%s: call to symbol \"%s\", which is not a function", pos.String(), name)
	}

	def := symbol.Def.(*FnDef)
	if len(expr.Arguments) < len(def.Params) {
		return nil, fmt.Errorf("%s: not enough arguments (got %d) to call \"%s\" function (want %d)",
			pos.String(), len(expr.Arguments), name, len(def.Params))
	}
	if len(expr.Arguments) > len(def.Params) {
		return nil, fmt.Errorf("%s: too many arguments (got %d) in function \"%s\" call (want %d)",
			pos.String(), len(expr.Arguments), name, len(def.Params))
	}
	args, err := s.scanCallArgs(ctx, def.Params, expr.Arguments)
	if err != nil {
		return nil, err
	}

	if symbol.Scope.Kind == scp.Unit {
		ctx.ref.Add(symbol)
	}

	return &SymbolCallExpression{
		Pos:       pos,
		Callee:    symbol,
		Arguments: args,
	}, nil
}

func (s *Scope) scanCallArgs(ctx *Context, params []*Symbol, exprs []ast.Expression) ([]Expression, error) {
	if len(exprs) != len(params) {
		panic("unreachable due to previous conditions")
	}
	n := len(params)
	if n == 0 {
		return nil, nil
	}

	args := make([]Expression, 0, n)
	for i := 0; i < n; i += 1 {
		expr := exprs[i]
		param := params[i]

		if expr == nil {
			panic("empty argument expression")
		}

		arg, err := s.scan(ctx, expr)
		if err != nil {
			return nil, err
		}

		err = checkCallArgType(param, arg)
		if err != nil {
			return nil, fmt.Errorf("%s: mismatched types of call argument and parameter", arg.Pin())
		}

		args = append(args, arg)
	}
	return args, nil
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
		Pos:        expr.Pos,
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
