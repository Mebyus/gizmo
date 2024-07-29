package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/enums/tpk"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/stg/scp"
	"github.com/mebyus/gizmo/token"
)

// Scan constructs expression from a given AST. Uses current scope for symbol
// resolution.
func (s *Scope) Scan(ctx *Context, expr ast.Expression) (Expression, error) {
	if expr == nil {
		return nil, nil
	}

	return s.scan(ctx, expr)
}

func (s *Scope) scan(ctx *Context, expr ast.Expression) (Expression, error) {
	switch expr.Kind() {
	case exn.Basic:
		return scanBasicLiteral(expr.(ast.BasicLiteral)), nil
	case exn.Symbol:
		return s.scanSymbolExpression(ctx, expr.(ast.SymbolExpression))
	case exn.Chain:
		return s.scanChainOperand(ctx, expr.(ast.ChainOperand))
	case exn.Unary:
		return s.scanUnaryExpression(ctx, expr.(*ast.UnaryExpression))
	case exn.Binary:
		return s.scanBinaryExpression(ctx, expr.(ast.BinaryExpression))
	case exn.Paren:
		return s.scanParenthesizedExpression(ctx, expr.(ast.ParenthesizedExpression))
	case exn.Cast:
		return s.scanCastExpression(ctx, expr.(ast.CastExpression))
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

func (s *Scope) scanChainOperand(ctx *Context, expr ast.ChainOperand) (ChainOperand, error) {
	name := expr.Identifier.Lit
	pos := expr.Identifier.Pos

	var symbol *Symbol
	var t *Type
	if name != "" {
		symbol = s.Lookup(name, pos.Num)
		if symbol == nil {
			return nil, fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
		}
		// TODO: probably check symbol kind here
		if symbol.Scope.Kind == scp.Unit {
			ctx.ref.Add(symbol)
		}

		t = symbol.Type
	} else {
		if ctx.rv == nil {
			return nil, fmt.Errorf("%s: receiver usage in regular function", pos.String())
		}

		// chain starts from receiver
		t = ctx.rv
	}

	chain := &ChainSymbol{
		Pos: pos,
		Sym: symbol,
		typ: t,
	}
	var tip ChainOperand = chain

	var err error
	for _, part := range expr.Parts {
		tip, err = s.scanChainPart(ctx, tip, part)
		if err != nil {
			return nil, err
		}
	}

	return tip, nil
}

func (s *Scope) scanChainPart(ctx *Context, tip ChainOperand, part ast.ChainPart) (ChainOperand, error) {
	switch part.Kind() {
	case exn.Address:
		return s.scanAddressPart(ctx, tip, part.(ast.AddressPart))
	case exn.Indirect:
		return s.scanIndirectPart(ctx, tip, part.(ast.IndirectPart))
	case exn.Member:
		return s.scanMemberPart(ctx, tip, part.(ast.MemberPart))
	case exn.Call:
		return s.scanCallPart(ctx, tip, part.(ast.CallPart))
	case exn.IndirectIndex:
		return s.scanIndirectIndexPart(ctx, tip, part.(ast.IndirectIndexPart))
	case exn.Index:
		return s.scanIndexPart(ctx, tip, part.(ast.IndexPart))
	default:
		panic(fmt.Sprintf("not implemented for %s expression", part.Kind().String()))
	}
}

func (s *Scope) scanIndexPart(ctx *Context, tip ChainOperand, part ast.IndexPart) (ChainOperand, error) {
	t := tip.Type()
	pos := part.Pos
	switch t.Kind {
	case tpk.Custom:
		panic("not implemented")
	case tpk.Chunk:
		index, err := s.scan(ctx, part.Index)
		if err != nil {
			return nil, err
		}
		if !index.Type().IsIntegerType() {
			return nil, fmt.Errorf("%s: type %s cannot be used as index", pos.String(), index.Type().Kind.String())
		}
		return &ChunkIndexExpression{
			Pos:    pos,
			Target: tip,
			Index:  index,
			typ:    t.Def.(ChunkTypeDef).ElemType,
		}, nil
	default:
		panic(fmt.Sprintf("not implemented for %s types", t.Kind))
	}
}

func (s *Scope) scanAddressPart(ctx *Context, tip ChainOperand, part ast.AddressPart) (ChainOperand, error) {
	return &AddressExpression{
		Pos:    part.Pos,
		Target: tip,
		typ:    s.Types.storePointer(tip.Type()),
	}, nil
}

func (s *Scope) scanIndirectIndexPart(ctx *Context, tip ChainOperand, part ast.IndirectIndexPart) (ChainOperand, error) {
	t := tip.Type()
	// TODO: implement for custom type
	if t.Kind != tpk.ArrayPointer {
		return nil, fmt.Errorf("%s: cannot indirect index %s operand of %s type",
			part.Pos.String(), tip.Kind().String(), t.Kind.String())
	}
	index, err := s.scan(ctx, part.Index)
	if err != nil {
		return nil, err
	}

	// TODO: check that index is of integer type

	return &IndirectIndexExpression{
		Pos:    part.Pos,
		Target: tip,
		Index:  index,
		typ:    t.Def.(ArrayPointerTypeDef).RefType,
	}, nil
}

func (s *Scope) scanIndirectPart(ctx *Context, tip ChainOperand, part ast.IndirectPart) (ChainOperand, error) {
	t := tip.Type()
	// TODO: implement for custom type
	if t.Kind != tpk.Pointer {
		return nil, fmt.Errorf("%s: cannot indirect %s operand of %s type",
			part.Pos.String(), tip.Kind().String(), t.Kind.String())
	}
	return &IndirectExpression{
		Pos:    part.Pos,
		Target: tip,
		typ:    t.Def.(PointerTypeDef).RefType,
	}, nil
}

func (s *Scope) scanMemberPart(ctx *Context, tip ChainOperand, part ast.MemberPart) (ChainOperand, error) {
	tt := tip.Type()

	// TODO: think up a better way to lookup members on types,
	// perhaps we should add a dedicated Type method for this

	var t *Type
	if tt.Kind == tpk.Custom {
		// TODO: first search here for possible methods
		t = tt.Def.(CustomTypeDef).Base
	} else {
		t = tt
	}

	pos := part.Member.Pos
	name := part.Member.Lit

	switch t.Kind {
	case tpk.Struct:
		def := t.Def.(*StructTypeDef)
		m := def.Members.Find(name)
		if m == nil {
			return nil, fmt.Errorf("%s: type %s no member \"%s\"",
				pos.String(), t.Kind.String(), name)
		}
		return &MemberExpression{
			Pos:    pos,
			Target: tip,
			Member: m,
		}, nil
	case tpk.Pointer:
		ref := t.Def.(PointerTypeDef).RefType
		if ref.Kind != tpk.Custom && ref.Def.(CustomTypeDef).Base.Kind != tpk.Struct {
			return nil, fmt.Errorf("%s: cannot select a member from %s type",
				pos.String(), ref.Def.(CustomTypeDef).Base.Kind)
		}
		def := ref.Def.(CustomTypeDef).Base.Def.(*StructTypeDef)
		m := def.Members.Find(name)
		if m == nil {
			return nil, fmt.Errorf("%s: type %s no member \"%s\"",
				pos.String(), t.Kind.String(), name)
		}
		return &IndirectMemberExpression{
			Pos:    pos,
			Target: tip,
			Member: m,
		}, nil
	case tpk.Chunk:
		name := part.Member.Lit
		switch name {
		case "len":
			return &ChunkMemberExpression{
				Pos:    pos,
				Target: tip,
				Name:   "len",
				typ:    UintType,
			}, nil
		case "ptr":
			return &ChunkMemberExpression{
				Pos:    pos,
				Target: tip,
				Name:   "ptr",

				// TODO: we probably need to construct this type differntly
				// based when it is custom chunk type
				typ: s.Types.storeArrayPointer(t.Def.(*ChunkTypeDef).ElemType),
			}, nil
		default:
			return nil, fmt.Errorf("%s: chunks do not have \"%s\" member", pos.String(), name)
		}
	case tpk.Signed, tpk.Unsigned, tpk.Boolean, tpk.Float:
		return nil, fmt.Errorf("%s: cannot select a member from %s type",
			pos.String(), t.Kind.String())
	default:
		panic(fmt.Sprintf("%s types not implemented", t.Kind.String()))
	}
}

func getSignatureByChainSymbol(c *ChainSymbol) (*Signature, error) {
	s := c.Sym
	if s == nil {
		panic("receiver not implemented")
	}

	switch s.Kind {
	case smk.Fun:
		return &s.Def.(*FunDef).Signature, nil
	case smk.Type:
		return nil, fmt.Errorf("%s: cannot call type \"%s\"", c.Pos.String(), s.Name)
	case smk.Import:
		return nil, fmt.Errorf("%s: cannot call imported unit \"%s\"", c.Pos.String(), s.Name)
	case smk.Method:
		panic(fmt.Sprintf("bad symbol name \"%s\" resolution", s.Name))
	default:
		panic(fmt.Sprintf("%s symbols not implemented", s.Kind.String()))
	}
}

func getSignature(o ChainOperand) (*Signature, error) {
	switch o.Kind() {
	case exn.Chain:
		return getSignatureByChainSymbol(o.(*ChainSymbol))
	default:
		panic(fmt.Sprintf("%s: %s operands not implemented", o.Pin().String(), o.Kind().String()))
	}
}

func (s *Scope) scanCallPart(ctx *Context, tip ChainOperand, part ast.CallPart) (ChainOperand, error) {
	sig, err := getSignature(tip)
	if err != nil {
		return nil, err
	}

	pos := part.Pos
	params := sig.Params
	args := part.Args

	aa, err := s.scanCallArgs(ctx, pos, params, args)
	if err != nil {
		return nil, err
	}
	return &CallExpression{
		Pos:       pos,
		Arguments: aa,
		Callee:    tip,

		typ:   sig.Result,
		never: sig.Never,
	}, nil
}

func (s *Scope) scanCastExpression(ctx *Context, expr ast.CastExpression) (*CastExpression, error) {
	target, err := s.scan(ctx, expr.Target)
	if err != nil {
		return nil, err
	}
	t, err := s.Types.lookup(expr.Type)
	if err != nil {
		return nil, err
	}

	return &CastExpression{
		// Pos: expr.,
		Target:          target,
		DestinationType: t,
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

func (s *Scope) scanCallArgs(ctx *Context, pos source.Pos, params []*Symbol, args []ast.Expression) ([]Expression, error) {
	if len(args) < len(params) {
		return nil, fmt.Errorf("%s: not enough arguments (got %d) to call \"%s\" function (want %d)",
			pos.String(), len(args), "???", len(params)) // TODO: provide call context to supply a name here
	}
	if len(args) > len(params) {
		return nil, fmt.Errorf("%s: too many arguments (got %d) in function \"%s\" call (want %d)",
			pos.String(), len(args), "???", len(params))
	}

	if len(args) != len(params) {
		panic("unreachable due to previous checks")
	}
	n := len(params)
	if n == 0 {
		return nil, nil
	}

	aa := make([]Expression, 0, n)
	for i := 0; i < n; i += 1 {
		arg := args[i]
		param := params[i]

		if arg == nil {
			panic("empty argument expression")
		}

		a, err := s.scan(ctx, arg)
		if err != nil {
			return nil, err
		}

		err = checkCallArgType(param, a)
		if err != nil {
			return nil, err
		}

		aa = append(aa, a)
	}
	return aa, nil
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
