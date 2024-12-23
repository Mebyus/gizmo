package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/enums/tpk"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/stg/scp"
	"github.com/mebyus/gizmo/token"
)

// Scan constructs expression from a given AST. Uses current scope for symbol
// resolution.
func (s *Scope) Scan(ctx *Context, exp ast.Exp) (Exp, error) {
	if exp == nil {
		return nil, nil
	}

	return s.scan(ctx, exp)
}

// ScanInitExp onstructs expression from a given AST expression node.
// Uses current scope for symbol resolution.
//
// This method is intended only for variable and object literal
// field init expressions.
func (s *Scope) ScanInitExp(ctx *Context, exp ast.Exp) (Exp, error) {
	if exp == nil {
		return nil, nil
	}
	if exp.Kind() == exk.Dirty {
		pos := exp.(ast.Dirty).Pos
		return Dirty{Pos: pos}, nil
	}

	return s.scan(ctx, exp)
}

func (s *Scope) scan(ctx *Context, exp ast.Exp) (Exp, error) {
	switch exp.Kind() {
	case exk.Basic:
		return scanBasicLiteral(exp.(ast.BasicLiteral)), nil
	case exk.Symbol:
		return s.scanSymbolExp(ctx, exp.(ast.SymbolExp))
	case exk.Chain:
		return s.scanChainOperand(ctx, exp.(ast.ChainOperand))
	case exk.Unary:
		return s.scanUnaryExp(ctx, exp.(*ast.UnaryExp))
	case exk.Binary:
		return s.scanBinExp(ctx, exp.(ast.BinExp))
	case exk.Paren:
		return s.scanParenExp(ctx, exp.(ast.ParenExp))
	case exk.Cast:
		return s.scanCastExp(ctx, exp.(ast.CastExp))
	case exk.Tint:
		return s.scanTintExp(ctx, exp.(ast.TintExp))
	case exk.IncompName:
		return s.scanIncompNameExp(ctx, exp.(ast.IncompNameExp))
	case exk.MemCast:
		return s.scanMemCastExp(ctx, exp.(ast.MemCastExpression))
	// case exk.BitCast:
	// 	// g.BitCastExpression(expr.(ast.BitCastExpression))
	// case exk.Object:
	// 	// g.ObjectLiteral(expr.(ast.ObjectLiteral))
	default:
		panic(fmt.Sprintf("not implemented for %s expression", exp.Kind().String()))
	}
}

func (s *Scope) scanMemCastExp(ctx *Context, exp ast.MemCastExpression) (*MemCastExp, error) {
	// TODO: check that types have the same size

	target, err := s.scan(ctx, exp.Target)
	if err != nil {
		return nil, err
	}
	t, err := s.Types.lookup(ctx, exp.Type)
	if err != nil {
		return nil, err
	}

	// TODO: perform types compatibility check

	return &MemCastExp{
		Target:   target,
		DestType: t,
	}, nil
}

func (s *Scope) scanIncompNameExp(ctx *Context, exp ast.IncompNameExp) (*EnumExp, error) {
	name := exp.Identifier.Lit
	pos := exp.Identifier.Pos

	if ctx.enum == nil {
		return nil, fmt.Errorf("%s: incomplete name \".%s\" used outside of enum context", pos, name)
	}

	enum := ctx.enum.Def.(*Type)
	def := enum.Def.(CustomTypeDef).Base.Def.(*EnumTypeDef)
	entry := def.Entry(name)
	if entry == nil {
		return nil, fmt.Errorf("%s: enum \"%s\" does not have entry named \".%s\"",
			pos, ctx.enum.Name, name)
	}

	return &EnumExp{
		Pos:   pos,
		Enum:  enum,
		Entry: entry,
	}, nil
}

func (s *Scope) scanChainOperand(ctx *Context, exp ast.ChainOperand) (ChainOperand, error) {
	name := exp.Start.Lit
	pos := exp.Start.Pos

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
		if symbol.Kind == smk.Import {
			return s.scanImportExp(ctx, pos, symbol, exp.Parts)
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
	return s.scanChainParts(ctx, chain, exp.Parts)
}

func (s *Scope) scanChainParts(ctx *Context, tip ChainOperand, parts []ast.ChainPart) (ChainOperand, error) {
	var err error
	for _, part := range parts {
		tip, err = s.scanChainPart(ctx, tip, part)
		if err != nil {
			return nil, err
		}
	}
	return tip, nil
}

func (s *Scope) scanImportExp(ctx *Context, pos source.Pos, imp *Symbol, parts []ast.ChainPart) (ChainOperand, error) {
	if len(parts) == 0 {
		return nil, fmt.Errorf("%s: import symbol \"%s\" cannot be used as standalone expression",
			pos, imp.Name)
	}

	var chain *ChainSymbol

	part := parts[0]
	switch part.Kind() {
	case exk.Select:
		unit := imp.Def.(ImportSymDef).Unit
		member := part.(ast.SelectPart).Name
		name := member.Lit
		pos := member.Pos
		symbol := unit.Scope.sym(name)
		if symbol == nil {
			return nil, fmt.Errorf("%s: unit %s has no \"%s\" symbol", pos, unit.Name, name)
		}
		if !symbol.Pub {
			return nil, fmt.Errorf("%s: imported symbol \"%s.%s\" is not public", pos, unit.Name, name)
		}
		chain = &ChainSymbol{
			Pos: pos,
			Sym: symbol,
			typ: symbol.Type,
		}
	default:
		return nil, fmt.Errorf("%s: %s chain cannot be used on \"%s\" import symbol",
			part.Pin(), part.Kind(), imp.Name)
	}

	return s.scanChainParts(ctx, chain, parts[1:])
}

func (s *Scope) scanChainPart(ctx *Context, tip ChainOperand, part ast.ChainPart) (ChainOperand, error) {
	switch part.Kind() {
	case exk.Address:
		return nil, nil
		// return s.scanAddressPart(ctx, tip, part.(ast.AddressPart))
	case exk.Indirect:
		return s.scanIndirectPart(ctx, tip, part.(ast.IndirectPart))
	case exk.Select:
		return s.scanSelectPart(ctx, tip, part.(ast.SelectPart))
	case exk.Call:
		return nil, nil
		// return s.scanCallPart(ctx, tip, part.(ast.CallPart))
	case exk.IndirectIndex:
		return s.scanIndirectIndexPart(ctx, tip, part.(ast.IndirectIndexPart))
	case exk.Index:
		return s.scanIndexPart(ctx, tip, part.(ast.IndexPart))
	case exk.Slice:
		return nil, nil
		// return s.scanSlicePart(ctx, tip, part.(ast.SlicePart))
	default:
		panic(fmt.Sprintf("not implemented for %s expression", part.Kind().String()))
	}
}

func (s *Scope) scanSlicePart(ctx *Context, tip ChainOperand, part ast.SliceExp) (ChainOperand, error) {
	start, err := s.Scan(ctx, part.Start)
	if err != nil {
		return nil, err
	}
	if start != nil && !start.Type().IsIntegerType() {
		return nil, fmt.Errorf("%s: type %s cannot be used as index", start.Pin(), start.Type().Kind)
	}

	end, err := s.Scan(ctx, part.End)
	if err != nil {
		return nil, err
	}
	if end != nil && !end.Type().IsIntegerType() {
		return nil, fmt.Errorf("%s: type %s cannot be used as index", end.Pin(), end.Type().Kind)
	}

	pos := part.Start.Pin()
	t := tip.Type()

	switch t.Kind {
	case tpk.Custom:
		panic("not implemented")
	case tpk.Chunk:
		return &ChunkSliceExp{
			Pos:    pos,
			Target: tip,
			Start:  start,
			End:    end,
			typ:    s.Types.storeChunk(t.Def.(ChunkTypeDef).ElemType),
		}, nil
	case tpk.Array:
		return &ArraySliceExp{
			Pos:    pos,
			Target: tip,
			Start:  start,
			End:    end,
			typ:    s.Types.storeChunk(t.Def.(ArrayTypeDef).ElemType),
		}, nil
	default:
		panic(fmt.Sprintf("not implemented for %s types", t.Kind))
	}
}

func (s *Scope) scanIndexPart(ctx *Context, tip ChainOperand, part ast.IndexPart) (ChainOperand, error) {
	index, err := s.scan(ctx, part.Index)
	if err != nil {
		return nil, err
	}
	if !index.Type().IsIntegerType() {
		return nil, fmt.Errorf("%s: type %s cannot be used as index", index.Pin(), index.Type().Kind)
	}

	pos := part.Index.Pin()
	t := tip.Type()
	switch t.Kind {
	case tpk.Custom:
		panic("not implemented")
	case tpk.Chunk:
		return &ChunkIndexExp{
			Pos:    pos,
			Target: tip,
			Index:  index,
			typ:    t.Def.(ChunkTypeDef).ElemType,
		}, nil
	case tpk.Array:
		return &ArrayIndexExp{
			Pos:    pos,
			Target: tip,
			Index:  index,
			typ:    t.Def.(ArrayTypeDef).ElemType,
		}, nil
	default:
		panic(fmt.Sprintf("not implemented for %s types", t.Kind))
	}
}

func (s *Scope) scanAddressPart(ctx *Context, tip ChainOperand, part ast.AddressExp) (ChainOperand, error) {
	return &AddressExp{
		Pos:    part.Pin(),
		Target: tip,
		typ:    s.Types.storePointer(tip.Type()),
	}, nil
}

func (s *Scope) scanIndirectIndexPart(ctx *Context, tip ChainOperand, part ast.IndirectIndexPart) (ChainOperand, error) {
	t := tip.Type()
	// TODO: implement for custom type
	if t.Kind != tpk.ArrayPointer {
		return nil, fmt.Errorf("%s: cannot indirect index %s operand of %s type",
			part.Pin(), tip.Kind().String(), t.Kind.String())
	}
	index, err := s.scan(ctx, part.Index)
	if err != nil {
		return nil, err
	}

	// TODO: check that index is of integer type

	return &IndirectIndexExp{
		Pos:    part.Pin(),
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
	return &IndirectExp{
		Pos:    part.Pos,
		Target: tip,
		typ:    t.Def.(PointerTypeDef).RefType,
	}, nil
}

func (s *Scope) scanSelectPart(ctx *Context, tip ChainOperand, part ast.SelectPart) (ChainOperand, error) {
	// TODO: add type and symdef for import symbols
	tt := tip.Type()

	// TODO: think up a better way to lookup members on types,
	// perhaps we should add a dedicated Type method for this

	name := part.Name

	switch tt.Kind {
	case tpk.Custom:
		member, err := tt.Member(name)
		if err != nil {
			return nil, err
		}

		switch m := member.(type) {
		case Field:
			return &FieldExp{
				Pos:    name.Pos,
				Target: tip,
				Field:  &m,
			}, nil
		case Method:
			return &BoundMethodExp{
				Pos:      name.Pos,
				Receiver: tip,
				Symbol:   m.Symbol,
				Pointer:  m.Symbol.Def.(*MethodDef).Receiver.Kind == tpk.Pointer,
			}, nil
		default:
			panic(fmt.Sprintf("%T unexpected member type", m))
		}
	case tpk.Pointer:
		rt := tt.Def.(PointerTypeDef).RefType
		member, err := rt.Member(name)
		if err != nil {
			return nil, err
		}

		switch m := member.(type) {
		case Field:
			return &IndirectFieldExp{
				Pos:    name.Pos,
				Target: tip,
				Field:  &m,
			}, nil
		case Method:
			return &BoundMethodExp{
				Pos:      name.Pos,
				Receiver: tip,
				Symbol:   m.Symbol,
			}, nil
		default:
			panic(fmt.Sprintf("%T unexpected member type", m))
		}
	case tpk.Chunk:
		switch name.Lit {
		case "len":
			return &ChunkMemberExp{
				Pos:    name.Pos,
				Target: tip,
				Name:   "len",
				typ:    UintType,
			}, nil
		case "ptr":
			return &ChunkMemberExp{
				Pos:    name.Pos,
				Target: tip,
				Name:   "ptr",

				// TODO: we probably need to construct this type differently
				// when it is custom chunk type
				typ: s.Types.storeArrayPointer(tt.Def.(ChunkTypeDef).ElemType),
			}, nil
		default:
			return nil, fmt.Errorf("%s: chunks do not have \"%s\" member", name.Pos, name.Lit)
		}
	case tpk.Integer, tpk.Boolean, tpk.Float:
		return nil, fmt.Errorf("%s: cannot select a member from %s type",
			name.Pos, tt.Kind)
	default:
		panic(fmt.Sprintf("%s types not implemented", tt.Kind))
	}
}

type CallDetails struct {
	Signature

	Args []Exp

	// Not nil only for method calls.
	Receiver Exp

	// What symbol is being called. One of:
	//
	//	- function
	//	- method
	//	- variable of function type
	Symbol *Symbol
}

func getCallDetailsByChainSymbol(c *ChainSymbol) (*CallDetails, error) {
	s := c.Sym
	if s == nil {
		panic("receiver not implemented")
	}

	switch s.Kind {
	case smk.Fun:
		return &CallDetails{
			Signature: s.Def.(*FunDef).Signature,
			Symbol:    s,
		}, nil
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

func (s *Scope) getCallDetailsByBoundMethod(m *BoundMethodExp) (*CallDetails, error) {
	var r Exp
	if m.Pointer {
		r = &AddressExp{
			Pos:    m.Receiver.Pin(),
			Target: m.Receiver,
			typ:    s.Types.storePointer(m.Receiver.Type()),
		}
	} else {
		r = m.Receiver
	}
	return &CallDetails{
		Signature: Signature{},
		Symbol:    m.Symbol,
		Receiver:  r,
	}, nil
}

// TODO: refactor this to accept CallExpression.
func (s *Scope) getCallDetails(o ChainOperand) (*CallDetails, error) {
	switch o.Kind() {
	case exk.Chain:
		return getCallDetailsByChainSymbol(o.(*ChainSymbol))
	case exk.BoundMethod:
		return s.getCallDetailsByBoundMethod(o.(*BoundMethodExp))
	default:
		panic(fmt.Sprintf("%s: %s operands not implemented", o.Pin().String(), o.Kind().String()))
	}
}

func (s *Scope) scanCallPart(ctx *Context, tip ChainOperand, part ast.CallExp) (ChainOperand, error) {
	details, err := s.getCallDetails(tip)
	if err != nil {
		return nil, err
	}

	pos := part.Pin()
	params := details.Params
	args := part.Args

	aa, err := s.scanCallArgs(ctx, pos, params, args)
	if err != nil {
		return nil, err
	}
	return &CallExp{
		Pos:       pos,
		Arguments: aa,
		Callee:    tip,

		typ:   details.Result,
		never: details.Never,
	}, nil
}

func (s *Scope) scanTintExp(ctx *Context, exp ast.TintExp) (*TintExp, error) {
	target, err := s.scan(ctx, exp.Target)
	if err != nil {
		return nil, err
	}
	t, err := s.Types.lookup(ctx, exp.Type)
	if err != nil {
		return nil, err
	}

	if !t.IsIntegerType() {
		return nil, fmt.Errorf("%s: destination must be integer type", exp.Type.Pin())
	}
	if !target.Type().IsIntegerType() {
		return nil, fmt.Errorf("%s: target is not an integer", target.Pin())
	}

	pos := exp.Pos
	return &TintExp{
		Pos:    pos,
		Target: target,

		DestType: t,
	}, nil
}

func (s *Scope) scanCastExp(ctx *Context, expr ast.CastExp) (*CastExp, error) {
	target, err := s.scan(ctx, expr.Target)
	if err != nil {
		return nil, err
	}
	t, err := s.Types.lookup(ctx, expr.Type)
	if err != nil {
		return nil, err
	}

	// TODO: perform types compatibility check

	return &CastExp{
		// Pos: expr.,
		Target:   target,
		DestType: t,
	}, nil
}

func (s *Scope) scanParenExp(ctx *Context, expr ast.ParenExp) (*ParenExp, error) {
	pos := expr.Pos
	inner, err := s.scan(ctx, expr.Inner)
	if err != nil {
		return nil, err
	}

	return &ParenExp{
		Pos:   pos,
		Inner: inner,
	}, nil
}

func (s *Scope) scanCallArgs(ctx *Context, pos source.Pos, params []*Symbol, args []ast.Exp) ([]Exp, error) {
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

	aa := make([]Exp, 0, n)
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

		err = typeCheckExp(param.Type, a)
		if err != nil {
			return nil, err
		}

		aa = append(aa, a)
	}
	return aa, nil
}

func (s *Scope) scanSymbolExp(ctx *Context, expr ast.SymbolExp) (*SymbolExp, error) {
	name := expr.Identifier.Lit
	pos := expr.Identifier.Pos
	symbol := s.Lookup(name, pos.Num)
	if symbol == nil {
		return nil, fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}
	if symbol.Scope.Kind == scp.Unit {
		ctx.ref.Add(symbol)
	}

	return &SymbolExp{
		Pos: pos,
		Sym: symbol,
	}, nil
}

func (s *Scope) scanUnaryExp(ctx *Context, expr *ast.UnaryExp) (*UnaryExp, error) {
	inner, err := s.scan(ctx, expr.Inner)
	if err != nil {
		return nil, err
	}

	return &UnaryExp{
		Operator: UnaryOperator(expr.Operator),
		Inner:    inner,
	}, nil
}

func (s *Scope) scanBinExp(ctx *Context, exp ast.BinExp) (Exp, error) {
	left, err := s.scan(ctx, exp.Left)
	if err != nil {
		return nil, err
	}
	right, err := s.scan(ctx, exp.Right)
	if err != nil {
		return nil, err
	}
	bin := &BinExp{
		Operator: BinaryOperator(exp.Operator),
		Left:     left,
		Right:    right,
	}

	e, err := typeCheckDesugarBinExp(bin)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func scanBasicLiteral(lit ast.BasicLiteral) Literal {
	pos := lit.Token.Pos

	switch lit.Token.Kind {
	case token.True:
		return True{Pos: pos}
	case token.False:
		return False{Pos: pos}
	case token.BinInteger, token.OctInteger, token.DecInteger, token.HexInteger:
		return Integer{Pos: pos, Val: lit.Token.Val, typ: StaticIntegerType}
	case token.String:
		return String{Pos: pos, Val: lit.Token.Lit, Size: lit.Token.Val}
	case token.Rune:
		// TODO: separate basic literals into classes in AST level,
		// parser should also verify and transform literals for this nodes
		if lit.Token.Lit != "" {
			panic(fmt.Sprintf("complex runes (%s) not implemented", lit.Token.Lit))
		}
		return Integer{Pos: pos, Val: lit.Token.Val, typ: StaticIntegerType}
	default:
		panic(fmt.Sprintf("not implemented for %s literal", lit.Token.Kind.String()))
	}
}
