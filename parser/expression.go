package parser

import (
	"sort"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/oper"
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
	if !p.tok.Kind.IsBinaryOperator() {
		return operand, nil
	}

	operands := []ast.PrimaryOperand{operand}
	var operators []oper.Binary
	for p.tok.Kind.IsBinaryOperator() {
		opr := oper.NewBinary(p.tok)
		p.advance()

		operand, err = p.primary()
		if err != nil {
			return nil, err
		}

		operands = append(operands, operand)
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
	op  oper.Binary
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
func sortOperators(ops []oper.Binary) positionedOperators {
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
func composeBinaryExpression(ops []oper.Binary, nds []ast.PrimaryOperand) ast.BinaryExpression {
	p := sortOperators(ops)

	// iterate over each operator in order of execution
	for i := 0; i < len(p); i++ {
		o := p[i]

		// collapse adjacent operands
		lopr := nds[o.pos]
		ropr := nds[o.pos+1]
		b := ast.BinaryExpression{
			Operator:  o.op,
			LeftSide:  lopr,
			RightSide: ropr,
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

func bexpr(op oper.Binary, left, right ast.Expression) ast.BinaryExpression {
	return ast.BinaryExpression{
		Operator:  op,
		LeftSide:  left,
		RightSide: right,
	}
}

func composeBinaryExpressionWithOneOperator(ops []oper.Binary, nds []ast.PrimaryOperand) ast.BinaryExpression {
	return bexpr(ops[0], nds[0], nds[1])
}

func composeBinaryExpressionWithTwoOperators(ops []oper.Binary, nds []ast.PrimaryOperand) ast.BinaryExpression {
	if ops[0].Precedence() <= ops[1].Precedence() {
		// a + b + c = ((a + b) + c)
		return bexpr(ops[1], bexpr(ops[0], nds[0], nds[1]), nds[2])
	}

	// a + b * c = (a + (b * c))
	return bexpr(ops[0], nds[0], bexpr(ops[1], nds[1], nds[2]))
}

func composeBinaryExpressionWithThreeOperators(ops []oper.Binary, operands []ast.PrimaryOperand) ast.BinaryExpression {
	switch {

	case ops[0].Precedence() <= ops[1].Precedence() && ops[1].Precedence() <= ops[2].Precedence():
		// a + b + c + d = (((a + b) + c) + d)
		return ast.BinaryExpression{
			Operator: ops[2],
			LeftSide: ast.BinaryExpression{
				Operator: ops[1],
				LeftSide: ast.BinaryExpression{
					Operator:  ops[0],
					LeftSide:  operands[0],
					RightSide: operands[1],
				},
				RightSide: operands[2],
			},
			RightSide: operands[3],
		}

	case ops[0].Precedence() <= ops[2].Precedence() && ops[2].Precedence() < ops[1].Precedence():
		// a * b + c * d = ((a * b) + (c * d))
		return ast.BinaryExpression{
			Operator: ops[1],
			LeftSide: ast.BinaryExpression{
				Operator:  ops[0],
				LeftSide:  operands[0],
				RightSide: operands[1],
			},
			RightSide: ast.BinaryExpression{
				Operator:  ops[2],
				LeftSide:  operands[2],
				RightSide: operands[3],
			},
		}

	case ops[1].Precedence() < ops[0].Precedence() && ops[0].Precedence() <= ops[2].Precedence():
		// a + b * c + d =  ((a + (b * c)) + d)
		return ast.BinaryExpression{
			Operator: ops[2],
			LeftSide: ast.BinaryExpression{
				Operator: ops[0],
				LeftSide: operands[0],
				RightSide: ast.BinaryExpression{
					Operator:  ops[1],
					LeftSide:  operands[1],
					RightSide: operands[2],
				},
			},
			RightSide: operands[3],
		}

	case ops[1].Precedence() <= ops[2].Precedence() && ops[2].Precedence() < ops[0].Precedence():
		// a + b * c * d = (a + ((b * c) * d))
		return ast.BinaryExpression{
			Operator: ops[1],
			LeftSide: operands[0],
			RightSide: ast.BinaryExpression{
				Operator: ops[2],
				LeftSide: ast.BinaryExpression{
					Operator:  ops[1],
					LeftSide:  operands[1],
					RightSide: operands[2],
				},
				RightSide: operands[3],
			},
		}

	case ops[2].Precedence() < ops[0].Precedence() && ops[0].Precedence() <= ops[1].Precedence():
		// a + b + c * d = ((a + b) + (c * d))
		return ast.BinaryExpression{
			Operator: ops[1],
			LeftSide: ast.BinaryExpression{
				Operator:  ops[0],
				LeftSide:  operands[0],
				RightSide: operands[1],
			},
			RightSide: ast.BinaryExpression{
				Operator:  ops[2],
				LeftSide:  operands[2],
				RightSide: operands[3],
			},
		}

	case ops[2].Precedence() < ops[1].Precedence() && ops[1].Precedence() < ops[0].Precedence():
		// a < b + c * d = (a < (b + (c * d)))
		return ast.BinaryExpression{
			Operator: ops[0],
			LeftSide: operands[0],
			RightSide: ast.BinaryExpression{
				Operator: ops[1],
				LeftSide: operands[1],
				RightSide: ast.BinaryExpression{
					Operator:  ops[2],
					LeftSide:  operands[2],
					RightSide: operands[3],
				},
			},
		}

	default:
		panic("must cover all cases")
	}
}

func (p *Parser) primary() (ast.PrimaryOperand, error) {
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
		Operator: oper.NewUnary(p.tok),
	}
	p.advance()

	tipExp := topExp
	for p.tok.Kind.IsUnaryOperator() {
		nextExp := &ast.UnaryExpression{
			Operator: oper.NewUnary(p.tok),
		}
		p.advance()
		tipExp.UnaryOperand = nextExp
		tipExp = nextExp
	}
	operand, err := p.operand()
	if err != nil {
		return nil, err
	}
	tipExp.UnaryOperand = operand
	return topExp, nil
}

func (p *Parser) operand() (ast.Operand, error) {
	operand, err := p.tryOperand()
	if err != nil {
		return nil, err
	}
	if operand == nil {
		return nil, p.unexpected(p.tok)
	}
	return operand, nil
}

func (p *Parser) tryOperand() (ast.Operand, error) {
	if p.tok.IsLit() {
		lit := p.basic()
		p.advance()
		return lit, nil
	}

	if p.tok.IsIdent() {
		return p.identStartOperand()
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

	if p.tok.Kind == token.List {
		p.advance() // skip ".["
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

// SelectorExpression, IndexExpression or CallExpression
func (p *Parser) identStartOperand() (ast.Operand, error) {
	var tip ast.Operand
	tip = p.ident()
	p.advance() // skip identifier

	for {
		switch p.tok.Kind {
		// case token.LeftParentheses:
		// 	tuple, err := p.tupleLiteral()
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	tip = ast.CallExpression{
		// 		Callee:    tip,
		// 		Arguments: tuple,
		// 	}
		case token.Period:
			p.advance() // skip "."
			err := p.expect(token.Identifier)
			if err != nil {
				return nil, err
			}
			selected := p.ident()
			p.advance() // skip identifier
			tip = ast.SelectorExpression{
				Target:   tip,
				Selected: selected,
			}
		case token.LeftSquare:
			p.advance() // skip "["
			if p.tok.Kind == token.Colon {
				p.advance() // skip ":"
				if p.tok.Kind == token.RightSquare {
					p.advance() // skip "]"
					tip = ast.SliceExpression{
						Target: tip,
					}
				} else {
					expr, err := p.expr()
					if err != nil {
						return nil, err
					}
					err = p.expect(token.RightSquare)
					if err != nil {
						return nil, err
					}
					p.advance() // skip "]"
					tip = ast.SliceExpression{
						Target: tip,
						End:    expr,
					}
				}
			} else {
				expr, err := p.expr()
				if err != nil {
					return nil, err
				}
				if p.tok.Kind == token.Colon {
					p.advance() // skip ":"
					err = p.expect(token.RightSquare)
					if err != nil {
						return nil, err
					}
					p.advance() // skip "]"
					tip = ast.SliceExpression{
						Target: tip,
						Start:  expr,
					}
				} else {
					err = p.expect(token.RightSquare)
					if err != nil {
						return nil, err
					}
					p.advance() // skip "]"
					tip = ast.IndexExpression{
						Target: tip,
						Index:  expr,
					}
				}
			}
		default:
			return tip, nil
		}
	}
}
