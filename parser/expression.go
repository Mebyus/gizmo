package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/token"
)

// ParseExpression is a pure function for usage in unit tests
func ParseExpression(str string) (ast.Exp, error) {
	p := New(lexer.NoPos(lexer.FromString(str)))
	return p.exp()
}

// InitExp parse variable or object literal field init expression.
func (p *Parser) InitExp() (ast.Exp, error) {
	if p.tok.Kind == token.Dirty {
		pos := p.pos()
		p.advance()
		return ast.Dirty{Pos: pos}, nil
	}

	return p.exp()
}

// Parse arbitrary expression (no expression will result in error).
//
// Parsing is done via Pratt's recursive descent algorithm variant.
func (p *Parser) exp() (ast.Exp, error) {
	return p.pratt(0)
}

func (p *Parser) pratt(power int) (ast.Exp, error) {
	left, err := p.primary()
	if err != nil {
		return nil, err
	}

	for p.tok.Kind.IsBinaryOperator() {
		o := ast.BinaryOperatorFromToken(p.tok)
		if o.Power() <= power {
			break
		}
		p.advance() // skip binary operator

		right, err := p.pratt(o.Power())
		if err != nil {
			return nil, err
		}

		left = bex(o, left, right)
	}

	return left, nil
}

func bex(op ast.BinaryOperator, left, right ast.Exp) ast.BinExp {
	return ast.BinExp{
		Operator: op,
		Left:     left,
		Right:    right,
	}
}

func (p *Parser) primary() (ast.Exp, error) {
	if p.tok.Kind.IsUnaryOperator() {
		unary, err := p.unary()
		if err != nil {
			return nil, err
		}
		return unary, nil
	}
	return p.operand()
}

func (p *Parser) unary() (*ast.UnaryExpression, error) {
	topExp := &ast.UnaryExpression{
		Operator: ast.UnaryOperatorFromToken(p.tok),
	}
	p.advance()

	tipExp := topExp
	for p.tok.Kind.IsUnaryOperator() {
		nextExp := &ast.UnaryExpression{
			Operator: ast.UnaryOperatorFromToken(p.tok),
		}
		p.advance()
		tipExp.Inner = nextExp
		tipExp = nextExp
	}
	operand, err := p.operand()
	if err != nil {
		return nil, err
	}
	tipExp.Inner = operand
	return topExp, nil
}

func (p *Parser) tint() (exp ast.TintExp, err error) {
	pos := p.pos()
	p.advance() // skip "tint"

	err = p.expect(token.LeftParentheses)
	if err != nil {
		return
	}
	p.advance() // skip "("

	spec, err := p.typeSpecifier()
	if err != nil {
		return
	}

	err = p.expect(token.Comma)
	if err != nil {
		return
	}
	p.advance() // skip ","

	target, err := p.exp()
	if err != nil {
		return
	}

	err = p.expect(token.RightParentheses)
	if err != nil {
		return
	}
	p.advance() // skip ")"

	return ast.TintExp{
		Pos:    pos,
		Target: target,
		Type:   spec,
	}, nil
}

func (p *Parser) cast() (exp ast.CastExp, err error) {
	pos := p.pos()
	p.advance() // skip "cast"

	err = p.expect(token.LeftParentheses)
	if err != nil {
		return
	}
	p.advance() // skip "("

	spec, err := p.typeSpecifier()
	if err != nil {
		return
	}

	err = p.expect(token.Comma)
	if err != nil {
		return
	}
	p.advance() // skip ","

	target, err := p.exp()
	if err != nil {
		return
	}

	err = p.expect(token.RightParentheses)
	if err != nil {
		return
	}
	p.advance() // skip ")"

	return ast.CastExp{
		Pos:    pos,
		Target: target,
		Type:   spec,
	}, nil
}

func (p *Parser) memcast() (exp ast.MemCastExpression, err error) {
	p.advance() // skip "mcast"

	err = p.expect(token.LeftParentheses)
	if err != nil {
		return
	}
	p.advance() // skip "("

	spec, err := p.typeSpecifier()
	if err != nil {
		return
	}

	err = p.expect(token.Comma)
	if err != nil {
		return
	}
	p.advance() // skip ","

	target, err := p.exp()
	if err != nil {
		return
	}

	err = p.expect(token.RightParentheses)
	if err != nil {
		return
	}
	p.advance() // skip ")"

	return ast.MemCastExpression{
		Target: target,
		Type:   spec,
	}, nil
}

func (p *Parser) objectField() (ast.ObjectField, error) {
	if p.tok.Kind != token.Identifier {
		return ast.ObjectField{}, p.unexpected()
	}
	name := p.word()
	p.advance() // skip field name

	if p.tok.Kind != token.Colon {
		return ast.ObjectField{}, p.unexpected()
	}
	p.advance() // skip ":"

	value, err := p.exp()
	if err != nil {
		return ast.ObjectField{}, err
	}

	return ast.ObjectField{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) objectLiteral() (ast.ObjectLiteral, error) {
	pos := p.pos()
	p.advance() // skip "{"

	var fields []ast.ObjectField
	for {
		if p.tok.Kind == token.RightCurly {
			p.advance() // skip "}"
			return ast.ObjectLiteral{
				Pos:    pos,
				Fields: fields,
			}, nil
		}

		field, err := p.objectField()
		if err != nil {
			return ast.ObjectLiteral{}, err
		}
		fields = append(fields, field)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.RightCurly {
			// will be skipped at next iteration
		} else {
			return ast.ObjectLiteral{}, p.unexpected()
		}
	}
}

func (p *Parser) paren() (ast.ParenthesizedExpression, error) {
	pos := p.pos()

	p.advance() // skip "("
	expr, err := p.exp()
	if err != nil {
		return ast.ParenthesizedExpression{}, err
	}
	err = p.expect(token.RightParentheses)
	if err != nil {
		return ast.ParenthesizedExpression{}, err
	}
	p.advance() // skip ")"
	return ast.ParenthesizedExpression{
		Pos:   pos,
		Inner: expr,
	}, nil
}

func (p *Parser) list() (ast.ListLiteral, error) {
	pos := p.pos()
	p.advance() // skip "["

	var list ast.ListLiteral
	for {
		if p.tok.Kind == token.RightSquare {
			p.advance() // skip "]"
			list.Pos = pos
			return list, nil
		}

		expr, err := p.exp()
		if err != nil {
			return ast.ListLiteral{}, err
		}
		list.Elems = append(list.Elems, expr)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.RightSquare {
			// will be skipped at next iteration
		} else {
			return ast.ListLiteral{}, p.unexpected()
		}
	}
}

func (p *Parser) chunkStartOperand() (ast.Operand, error) {
	pos := p.pos()
	p.advance() // skip "[]"

	if p.tok.Kind == token.Chunk || p.tok.Kind == token.Identifier {
		panic("type specifier operands not implemented")
	}
	return ast.ListLiteral{Pos: pos}, nil
}

func (p *Parser) operand() (ast.Operand, error) {
	if p.tok.Kind.IsLit() {
		lit := p.basic()
		p.advance()
		return lit, nil
	}

	switch p.tok.Kind {
	case token.Cast:
		return p.cast()
	case token.Tint:
		return p.tint()
	case token.MemCast:
		return p.memcast()
	case token.LeftCurly:
		return p.objectLiteral()
	case token.Identifier:
		return p.identifierStartOperand()
	case token.LeftParentheses:
		return p.paren()
	case token.LeftSquare:
		return p.list()
	case token.Chunk:
		return p.chunkStartOperand()
	case token.Period:
		return p.incompNameOperand()
	default:
		return nil, p.unexpected()
	}
}

func isChainOperandToken(kind token.Kind) bool {
	switch kind {
	case token.LeftParentheses, token.Period, token.LeftSquare,
		token.Indirect, token.Address, token.IndirectIndex:

		return true
	default:
		return false
	}
}

func (p *Parser) incompNameOperand() (ast.IncompNameExp, error) {
	p.advance() // skip "."

	if p.tok.Kind != token.Identifier {
		return ast.IncompNameExp{}, p.unexpected()
	}

	name := p.word()
	p.advance() // skip name identifier

	return ast.IncompNameExp{Identifier: name}, nil
}

// SymbolExpression, SelectorExpression, IndexExpression, CallExpression or InstanceExpression
func (p *Parser) identifierStartOperand() (ast.Operand, error) {
	idn := p.word()
	p.advance() // skip identifier

	if !isChainOperandToken(p.tok.Kind) {
		return ast.SymbolExp{Identifier: idn}, nil
	}

	chain := ast.ChainOperand{Identifier: idn}
	err := p.chainOperand(&chain)
	if err != nil {
		return nil, err
	}
	return chain, nil
}

func (p *Parser) callPart() (ast.CallPart, error) {
	pos := p.pos()

	args, err := p.callArguments()
	if err != nil {
		return ast.CallPart{}, err
	}

	return ast.CallPart{
		Pos:  pos,
		Args: args,
	}, nil
}

func (p *Parser) memberPart() (ast.MemberPart, error) {
	p.advance() // skip "."
	err := p.expect(token.Identifier)
	if err != nil {
		return ast.MemberPart{}, err
	}
	member := p.word()
	p.advance() // skip identifier

	return ast.MemberPart{Member: member}, nil
}

func (p *Parser) indirectPart() ast.IndirectPart {
	pos := p.pos()
	p.advance() // skip ".@"

	return ast.IndirectPart{Pos: pos}
}

func (p *Parser) chainOperand(chain *ast.ChainOperand) error {
	var parts []ast.ChainPart
	var prev exk.Kind
	for {
		var err error
		var part ast.ChainPart

		switch p.tok.Kind {
		case token.LeftParentheses:
			part, err = p.callPart()
		case token.Period:
			if p.next.Kind == token.Test {
				part, err = p.testPart()
			} else {
				part, err = p.memberPart()
			}
		case token.Indirect:
			part = p.indirectPart()
		case token.Address:
			part, err = p.addressPart(prev)
		case token.IndirectIndex:
			part, err = p.indirectIndexPart()
		case token.LeftSquare:
			part, err = p.leftSquarePart()
		default:
			chain.Parts = parts
			return nil
		}
		if err != nil {
			return err
		}
		prev = part.Kind()
		parts = append(parts, part)
	}
}

func (p *Parser) testPart() (ast.TestPart, error) {
	p.advance() // skip "."
	p.advance() // skip "test"

	if p.tok.Kind != token.Period {
		return ast.TestPart{}, p.unexpected()
	}
	p.advance() // skip "."

	if p.tok.Kind != token.Identifier {
		return ast.TestPart{}, p.unexpected()
	}

	name := p.word()
	p.advance() // skip test name

	return ast.TestPart{Name: name}, nil
}

func (p *Parser) addressPart(prev exk.Kind) (ast.AddressPart, error) {
	pos := p.pos()
	p.advance() // skip ".&"

	switch prev {
	case exk.Call, exk.Slice, exk.Address, exk.Indirect:
		return ast.AddressPart{}, fmt.Errorf("%s: cannot take address of %s expression",
			pos.String(), prev.String())
	}

	return ast.AddressPart{Pos: pos}, nil
}

func (p *Parser) indirectIndexPart() (ast.IndirectIndexPart, error) {
	pos := p.pos()

	p.advance() // skip ".["
	index, err := p.exp()
	if err != nil {
		return ast.IndirectIndexPart{}, err
	}
	if p.tok.Kind != token.RightSquare {
		return ast.IndirectIndexPart{}, p.unexpected()
	}
	p.advance() // skip "]"
	return ast.IndirectIndexPart{
		Pos:   pos,
		Index: index,
	}, nil
}

func (p *Parser) leftSquarePart() (ast.ChainPart, error) {
	pos := p.pos()
	p.advance() // skip "["

	if p.tok.Kind == token.Colon {
		p.advance() // skip ":"
		if p.tok.Kind == token.RightSquare {
			p.advance() // skip "]"
			return ast.SlicePart{Pos: pos}, nil
		}

		expr, err := p.exp()
		if err != nil {
			return nil, err
		}
		err = p.expect(token.RightSquare)
		if err != nil {
			return nil, err
		}
		p.advance() // skip "]"
		return ast.SlicePart{
			Pos: pos,
			End: expr,
		}, nil
	}

	expr, err := p.exp()
	if err != nil {
		return nil, err
	}
	if p.tok.Kind == token.Colon {
		p.advance() // skip ":"
		if p.tok.Kind == token.RightSquare {
			p.advance() // skip "]"
			return ast.SlicePart{
				Pos:   pos,
				Start: expr,
			}, nil
		}
		end, err := p.exp()
		if err != nil {
			return nil, err
		}

		err = p.expect(token.RightSquare)
		if err != nil {
			return nil, err
		}
		p.advance() // skip "]"
		return ast.SlicePart{
			Pos:   pos,
			Start: expr,
			End:   end,
		}, nil
	}

	err = p.expect(token.RightSquare)
	if err != nil {
		return nil, err
	}
	p.advance() // skip "]"
	return ast.IndexPart{
		Pos:   pos,
		Index: expr,
	}, nil
}

func (p *Parser) callArguments() ([]ast.Exp, error) {
	p.advance() // skip "("

	var args []ast.Exp
	for {
		if p.tok.Kind == token.RightParentheses {
			p.advance() // skip ")"
			return args, nil
		}

		expr, err := p.exp()
		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.RightParentheses {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected()
		}
	}
}
