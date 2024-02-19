package parser

import (
	"fmt"
	"sort"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/token"
)

// ParseExpression is a pure function for usage in unit tests
func ParseExpression(str string) (ast.Expression, error) {
	p := New(lexer.NoPos(lexer.FromString(str)))
	return p.expr()
}

func (p *Parser) expr() (ast.Expression, error) {
	operand, err := p.primary()
	if err != nil {
		return nil, err
	}
	return p.continueExpressionFromOperand(operand)
}

func (p *Parser) continueExpressionFromOperand(operand ast.Expression) (ast.Expression, error) {
	if !p.tok.Kind.IsBinaryOperator() {
		return operand, nil
	}

	operands := []ast.Expression{operand}
	var operators []ast.BinaryOperator
	for p.tok.Kind.IsBinaryOperator() {
		opr := ast.BinaryOperatorFromToken(p.tok)
		p.advance() // skip binary operator

		expr, err := p.primary()
		if err != nil {
			return nil, err
		}

		operands = append(operands, expr)
		operators = append(operators, opr)
	}

	// handle common cases with 1, 2 or 3 operators by hand
	switch len(operators) {
	case 0:
		panic("unreachable: slice must contain at least one operator")
	case 1:
		return composeBinaryExpressionWithOneOperator(operators, operands), nil
	case 2:
		return composeBinaryExpressionWithTwoOperators(operators, operands), nil
	case 3:
		return composeBinaryExpressionWithThreeOperators(operators, operands), nil
	default:
		return composeBinaryExpression(operators, operands), nil
	}
}

type positionedOperator struct {
	op  ast.BinaryOperator
	pos int
}

type positionedOperators []positionedOperator

func (p positionedOperators) Len() int {
	return len(p)
}

func (p positionedOperators) Less(i, j int) bool {
	return (p[i].op.Precedence() < p[j].op.Precedence()) ||
		(p[i].op.Precedence() == p[j].op.Precedence() && p[i].pos < p[j].pos)
}

func (p positionedOperators) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// sortOperators returns slice of operators sorted in order of execution
// along with original position information
func sortOperators(ops []ast.BinaryOperator) positionedOperators {
	p := make(positionedOperators, 0, len(ops))
	for i, o := range ops {
		p = append(p, positionedOperator{op: o, pos: i})
	}
	sort.Sort(p)
	return p
}

// composeBinaryExpression текущий алгоритм сборки дерева работает за O(n^2) по времени
//
// Входные слайсы из операторов и операндов должны удовлетворять условию len(ops) + 1 = len(operands)
//
// Слайс операндов будет использован для in-place мутаций и получения итогового выражения
func composeBinaryExpression(ops []ast.BinaryOperator, nds []ast.Expression) ast.BinaryExpression {
	p := sortOperators(ops)

	// iterate over each operator in order of execution
	for i := 0; i < len(p); i++ {
		o := p[i]

		// collapse adjacent operands
		lopr := nds[o.pos]
		ropr := nds[o.pos+1]
		b := ast.BinaryExpression{
			Operator: o.op,
			Left:     lopr,
			Right:    ropr,
		}
		nds[o.pos] = b

		// необходимо сдвинуть все операнды начиная с o.pos+2 влево на 1
		l := len(nds) - i // текущее количество операндов в срезе, с учетом предыдущих сверток
		for j := o.pos + 1; j < l-1; j++ {
			nds[j] = nds[j+1]
		}

		// теперь нужно сделать поправку позиций операторов, которые находились правее места свертки
		for j := i + 1; j < len(p); j++ {
			if p[j].pos > o.pos {
				p[j].pos--
			}
		}
	}

	return nds[0].(ast.BinaryExpression)
}

func bex(o ast.BinaryOperator, l, r ast.Expression) ast.BinaryExpression {
	return ast.BinaryExpression{
		Operator: o,
		Left:     l,
		Right:    r,
	}
}

func composeBinaryExpressionWithOneOperator(o []ast.BinaryOperator, n []ast.Expression) ast.BinaryExpression {
	return bex(o[0], n[0], n[1])
}

func composeBinaryExpressionWithTwoOperators(o []ast.BinaryOperator, n []ast.Expression) ast.BinaryExpression {
	if o[0].Precedence() <= o[1].Precedence() {
		// a + b + c = ((a + b) + c)
		return bex(o[1], bex(o[0], n[0], n[1]), n[2])
	}

	// a + b * c = (a + (b * c))
	return bex(o[0], n[0], bex(o[1], n[1], n[2]))
}

func composeBinaryExpressionWithThreeOperators(o []ast.BinaryOperator, n []ast.Expression) ast.BinaryExpression {
	switch {

	case o[0].Precedence() <= o[1].Precedence() && o[1].Precedence() <= o[2].Precedence():
		// a + b + c + d = (((a + b) + c) + d)
		return bex(o[2], bex(o[1], bex(o[0], n[0], n[1]), n[2]), n[3])

	case o[0].Precedence() <= o[2].Precedence() && o[2].Precedence() < o[1].Precedence():
		// a * b + c * d = ((a * b) + (c * d))
		return bex(o[1], bex(o[0], n[0], n[1]), bex(o[2], n[2], n[3]))

	case o[1].Precedence() < o[0].Precedence() && o[0].Precedence() <= o[2].Precedence():
		// a + b * c + d = ((a + (b * c)) + d)
		return bex(o[2], bex(o[0], n[0], bex(o[1], n[1], n[2])), n[3])

	case o[1].Precedence() <= o[2].Precedence() && o[2].Precedence() < o[0].Precedence():
		// a + b * c * d = (a + ((b * c) * d))
		return bex(o[0], n[0], bex(o[2], bex(o[1], n[1], n[2]), n[3]))

	case o[2].Precedence() < o[0].Precedence() && o[0].Precedence() <= o[1].Precedence():
		// a + b + c * d = ((a + b) + (c * d))
		return bex(o[1], bex(o[0], n[0], n[1]), bex(o[2], n[2], n[3]))

	case o[2].Precedence() < o[1].Precedence() && o[1].Precedence() < o[0].Precedence():
		// a < b + c * d = (a < (b + (c * d)))
		return bex(o[0], n[0], bex(o[1], n[1], bex(o[2], n[2], n[3])))

	default:
		panic("unreachable: switch must cover all cases")
	}
}

func (p *Parser) primary() (ast.Expression, error) {
	if p.tok.Kind.IsUnaryOperator() {
		unary, err := p.unary()
		if err != nil {
			return nil, err
		}
		return unary, nil
	}
	operand, err := p.tryOperand()
	if err != nil {
		return nil, err
	}
	if operand != nil {
		return operand, nil
	}
	return nil, p.unexpected(p.tok)
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

func (p *Parser) operand() (ast.Expression, error) {
	operand, err := p.tryOperand()
	if err != nil {
		return nil, err
	}
	if operand == nil {
		return nil, p.unexpected(p.tok)
	}
	return operand, nil
}

func (p *Parser) castExpression() (ast.CastExpression, error) {
	p.advance() // skip "cast"

	if p.tok.Kind != token.LeftSquare {
		return ast.CastExpression{}, p.unexpected(p.tok)
	}
	p.advance() // skip "["

	target, err := p.expr()
	if err != nil {
		return ast.CastExpression{}, err
	}

	if p.tok.Kind != token.Colon {
		return ast.CastExpression{}, p.unexpected(p.tok)
	}
	p.advance() // skip ":"

	spec, err := p.typeSpecifier()
	if err != nil {
		return ast.CastExpression{}, p.unexpected(p.tok)
	}

	if p.tok.Kind != token.RightSquare {
		return ast.CastExpression{}, p.unexpected(p.tok)
	}
	p.advance() // skip "]"

	return ast.CastExpression{
		Target: target,
		Type:   spec,
	}, nil
}

func (p *Parser) tryOperand() (ast.Operand, error) {
	if p.tok.Kind == token.Cast {
		return p.castExpression()
	}

	if p.tok.IsLit() {
		lit := p.basic()
		p.advance()
		return lit, nil
	}

	if p.tok.IsIdent() {
		return p.identifierStartOperand()
	}

	if p.tok.IsLeftPar() {
		p.advance() // skip "("
		expr, err := p.expr()
		if err != nil {
			return nil, err
		}
		err = p.expect(token.RightParentheses)
		if err != nil {
			return nil, err
		}
		p.advance() // skip ")"
		return ast.ParenthesizedExpression{Inner: expr}, nil
	}

	if p.tok.Kind == token.LeftSquare {
		p.advance() // skip "["
		var list ast.ListLiteral

		for {
			if p.tok.Kind == token.RightSquare {
				p.advance() // skip "]"
				return list, nil
			}

			expr, err := p.expr()
			if err != nil {
				return nil, err
			}
			list.Elems = append(list.Elems, expr)

			if p.tok.Kind == token.Comma {
				p.advance() // skip ","
			} else if p.tok.Kind == token.RightSquare {
				// will be skipped at next iteration
			} else {
				return nil, p.unexpected(p.tok)
			}
		}
	}

	return nil, nil
}

// SubsExpression, SelectorExpression, IndexExpression or CallExpression
func (p *Parser) identifierStartOperand() (ast.Operand, error) {
	scoped, err := p.scopedIdentifier()
	if err != nil {
		return nil, err
	}

	switch p.tok.Kind {
	case token.Period, token.LeftParentheses, token.LeftSquare, token.Indirect, token.Address:
		return p.chainOperand(ast.ChainStart{Identifier: scoped})
	default:
		return ast.SubsExpression{Identifier: scoped}, nil
	}
}

func (p *Parser) chainOperand(start ast.ChainStart) (ast.ChainOperand, error) {
	// note conversion to interface type
	var tip ast.ChainOperand = start

	for {
		switch p.tok.Kind {
		case token.LeftParentheses:
			args, err := p.callArguments()
			if err != nil {
				return nil, err
			}
			tip = ast.CallExpression{
				Callee:    tip,
				Arguments: args,
			}
		case token.Period:
			p.advance() // skip "."
			err := p.expect(token.Identifier)
			if err != nil {
				return nil, err
			}
			selected := p.idn()
			p.advance() // skip identifier
			tip = ast.SelectorExpression{
				Target:   tip,
				Selected: selected,
			}
		case token.Indirect:
			p.advance() // skip ".@"
			tip = ast.IndirectExpression{Target: tip}
		case token.Address:
			p.advance() // skip ".&"
			if tip.Kind() == exn.Call {
				return nil, fmt.Errorf("cannot take address of a call result %s", tip.Pin().String())
			}
			tip = ast.AddressExpression{Target: tip}
		case token.IndirectIndex:
			p.advance() // skip ".["
			index, err := p.expr()
			if err != nil {
				return nil, err
			}
			if p.tok.Kind != token.RightSquare {
				return nil, p.unexpected(p.tok)
			}
			p.advance() // skip "]"
			tip = ast.IndirectIndexExpression{
				Target: tip,
				Index:  index,
			}
		case token.LeftSquare:
			p.advance() // skip "["
			// if p.tok.Kind == token.Colon {
			// 	p.advance() // skip ":"
			// 	if p.tok.Kind == token.RightSquare {
			// 		p.advance() // skip "]"
			// 		tip = ast.SliceExpression{
			// 			Target: tip,
			// 		}
			// 	} else {
			// 		expr, err := p.expr()
			// 		if err != nil {
			// 			return nil, err
			// 		}
			// 		err = p.expect(token.RightSquare)
			// 		if err != nil {
			// 			return nil, err
			// 		}
			// 		p.advance() // skip "]"
			// 		tip = ast.SliceExpression{
			// 			Target: tip,
			// 			End:    expr,
			// 		}
			// 	}
			// } else {
			// 	expr, err := p.expr()
			// 	if err != nil {
			// 		return nil, err
			// 	}
			// 	if p.tok.Kind == token.Colon {
			// 		p.advance() // skip ":"
			// 		err = p.expect(token.RightSquare)
			// 		if err != nil {
			// 			return nil, err
			// 		}
			// 		p.advance() // skip "]"
			// 		tip = ast.SliceExpression{
			// 			Target: tip,
			// 			Start:  expr,
			// 		}
			// 	} else {
			// 		err = p.expect(token.RightSquare)
			// 		if err != nil {
			// 			return nil, err
			// 		}
			// 		p.advance() // skip "]"
			// 		tip = ast.IndexExpression{
			// 			Target: tip,
			// 			Index:  expr,
			// 		}
			// 	}
			// }
		default:
			return tip, nil
		}
	}
}

func (p *Parser) callArguments() ([]ast.Expression, error) {
	p.advance() // skip "("

	var args []ast.Expression
	for {
		if p.tok.Kind == token.RightParentheses {
			p.advance() // skip ")"
			return args, nil
		}

		expr, err := p.expr()
		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.RightParentheses {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected(p.tok)
		}
	}
}
