package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/ast/lbl"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) parseStatement() (statement ast.Statement, err error) {
	switch p.tok.Kind {
	case token.LeftCurly:
		return p.block()
	case token.Const:
		return p.constStatement()
	case token.Let:
		return p.letStatement()
	case token.Var:
		return p.varStatement()
	case token.If:
		return p.ifStatement()
	case token.Return:
		return p.returnStatement()
	case token.For:
		return p.forStatement()
	case token.Match:
		return p.matchStatement()
	case token.Jump:
		return p.jumpStatement()
	case token.Defer:
		return p.deferStatement()
	default:
		return p.otherStatement()
	}
}

func (p *Parser) deferStatement() (ast.DeferStatement, error) {
	pos := p.pos()
	p.advance() // skip "defer"

	if p.tok.Kind != token.Identifier && p.tok.Kind != token.Receiver {
		return ast.DeferStatement{}, p.unexpected(p.tok)
	}

	var cop ast.ChainOperand
	var err error
	if p.tok.Kind == token.Identifier {
		identifier := p.idn()
		p.advance() // skip identifier

		cop, err = p.chainOperand(ast.ChainStart{Identifier: identifier})
		if err != nil {
			return ast.DeferStatement{}, err
		}
	} else {
		rvPos := p.pos()
		p.advance() // skip "rv"

		cop, err = p.chainOperand(ast.Receiver{Pos: rvPos})
		if err != nil {
			return ast.DeferStatement{}, err
		}
	}

	if cop.Kind() != exn.Call {
		return ast.DeferStatement{}, fmt.Errorf("%s: only call statements can be deferred", pos.String())
	}

	if p.tok.Kind != token.Semicolon {
		return ast.DeferStatement{}, p.unexpected(p.tok)
	}
	p.advance() // skip ";"

	return ast.DeferStatement{
		Pos:  pos,
		Call: cop.(ast.CallExpression),
	}, nil
}

func (p *Parser) jumpStatement() (ast.JumpStatement, error) {
	p.advance() // slip "jump"

	label, err := p.label()
	if err != nil {
		return ast.JumpStatement{}, err
	}

	err = p.expect(token.Semicolon)
	if err != nil {
		return ast.JumpStatement{}, err
	}
	p.advance() // consume ";"

	return ast.JumpStatement{Label: label}, nil
}

func (p *Parser) label() (ast.Label, error) {
	pos := p.pos()

	switch p.tok.Kind {
	case token.LabelNext:
		p.advance()
		return ast.ReservedLabel{Pos: pos, ResKind: lbl.Next}, nil
	case token.LabelEnd:
		p.advance()
		return ast.ReservedLabel{Pos: pos, ResKind: lbl.End}, nil
	default:
		panic("not implemented for label " + p.tok.Kind.String())
	}
}

func (p *Parser) matchStatement() (ast.MatchStatement, error) {
	pos := p.pos()

	p.advance() // skip "match"

	expr, err := p.expr()
	if err != nil {
		return ast.MatchStatement{}, err
	}

	if p.tok.Kind != token.LeftCurly {
		return ast.MatchStatement{}, p.unexpected(p.tok)
	}
	p.advance() // skip "{"

	var cases []ast.MatchCase
	for {
		switch p.tok.Kind {
		case token.Case:
			c, err := p.matchCase()
			if err != nil {
				return ast.MatchStatement{}, err
			}
			cases = append(cases, c)
		case token.Else:
			block, err := p.matchElse()
			if err != nil {
				return ast.MatchStatement{}, err
			}

			// closing curly of match statement
			if p.tok.Kind != token.RightCurly {
				return ast.MatchStatement{}, p.unexpected(p.tok)
			}
			p.advance() // skip "}"

			return ast.MatchStatement{
				Pos:        pos,
				Expression: expr,
				Cases:      cases,
				Else:       block,
			}, nil
		default:
			return ast.MatchStatement{}, p.unexpected(p.tok)
		}
	}
}

func (p *Parser) matchCase() (ast.MatchCase, error) {
	p.advance() // skip "case"

	expr, err := p.expr()
	if err != nil {
		return ast.MatchCase{}, err
	}

	if p.tok.Kind != token.RightArrow {
		return ast.MatchCase{}, p.unexpected(p.tok)
	}
	p.advance() // skip "=>"

	block, err := p.block()
	if err != nil {
		return ast.MatchCase{}, err
	}

	return ast.MatchCase{
		Expression: expr,
		Body:       block,
	}, nil
}

func (p *Parser) matchElse() (ast.BlockStatement, error) {
	p.advance() // skip "else"

	if p.tok.Kind != token.RightArrow {
		return ast.BlockStatement{}, p.unexpected(p.tok)
	}
	p.advance() // skip "=>"

	block, err := p.block()
	if err != nil {
		return ast.BlockStatement{}, err
	}

	return block, nil
}

func (p *Parser) forStatement() (ast.Statement, error) {
	if p.next.Kind == token.LeftCurly {
		return p.forSimpleStatement()
	}
	if p.next.Kind == token.Let {
		return p.forEachStatement()
	}
	return p.forWithConditionStatement()
}

func (p *Parser) forSimpleStatement() (ast.ForStatement, error) {
	pos := p.pos()

	p.advance() // skip "for"

	body, err := p.block()
	if err != nil {
		return ast.ForStatement{}, err
	}
	if len(body.Statements) == 0 {
		return ast.ForStatement{}, fmt.Errorf("%s for loop without condition cannot have empty body", pos.Short())
	}

	return ast.ForStatement{
		Pos:  pos,
		Body: body,
	}, nil
}

func (p *Parser) forEachStatement() (ast.ForEachStatement, error) {
	pos := p.pos()

	p.advance() // skip "for"
	p.advance() // skip "let"

	if p.tok.Kind != token.Identifier {
		return ast.ForEachStatement{}, p.unexpected(p.tok)
	}

	name := p.idn()
	p.advance() // skip name identifier

	if p.tok.Kind != token.In {
		return ast.ForEachStatement{}, p.unexpected(p.tok)
	}
	p.advance() // skip "in"

	expr, err := p.expr()
	if err != nil {
		return ast.ForEachStatement{}, err
	}

	if p.tok.Kind != token.LeftCurly {
		return ast.ForEachStatement{}, p.unexpected(p.tok)
	}
	body, err := p.block()
	if err != nil {
		return ast.ForEachStatement{}, err
	}

	return ast.ForEachStatement{
		Pos:      pos,
		Name:     name,
		Iterator: expr,
		Body:     body,
	}, nil
}

func (p *Parser) forWithConditionStatement() (ast.ForConditionStatement, error) {
	pos := p.pos()

	p.advance() // skip "for"

	condition, err := p.expr()
	if err != nil {
		return ast.ForConditionStatement{}, err
	}
	if p.tok.Kind != token.LeftCurly {
		return ast.ForConditionStatement{}, p.unexpected(p.tok)
	}
	body, err := p.block()
	if err != nil {
		return ast.ForConditionStatement{}, err
	}

	return ast.ForConditionStatement{
		Pos:       pos,
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) returnStatement() (statement ast.ReturnStatement, err error) {
	pos := p.pos()

	p.advance() // consume "return"

	if p.tok.Kind == token.Semicolon {
		p.advance() // consume ";"
		return ast.ReturnStatement{Pos: pos}, nil
	}

	expression, err := p.expr()
	if err != nil {
		return
	}
	err = p.expect(token.Semicolon)
	if err != nil {
		return
	}
	p.advance() // consume ";"
	statement = ast.ReturnStatement{
		Pos:        pos,
		Expression: expression,
	}
	return
}

// func (p *Parser) loop() (ast.Statement, error) {
// 	if p.next.Kind == token.LeftCurly {
// 		return p.forEver()
// 	}

// 	p.advance() // skip "for"

// 	condition, err := p.expr()
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = p.expect(token.LeftCurly)
// 	if err != nil {
// 		return nil, err
// 	}
// 	body, err := p.block()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ast.While{
// 		Body:      body,
// 		Condition: condition,
// 	}, nil
// }

// func (p *Parser) forEver() (statement ast.ForEver, err error) {
// 	p.advance() // consume "for"
// 	body, err := p.block()
// 	if err != nil {
// 		return
// 	}
// 	statement = ast.ForEver{Body: body}
// 	return
// }

func (p *Parser) ifStatement() (statement ast.IfStatement, err error) {
	ifClause, err := p.ifClause()
	if err != nil {
		return
	}

	var elseIf []ast.ElseIfClause
	for {
		if p.tok.Kind == token.Else && p.next.Kind == token.If {
			p.advance() // skip "else"
		} else {
			break
		}

		clause, err := p.ifClause()
		if err != nil {
			return ast.IfStatement{}, err
		}
		elseIf = append(elseIf, ast.ElseIfClause(clause))
	}

	var elseClause *ast.ElseClause
	if p.tok.Kind == token.Else {
		p.advance() // skip "else"
		err = p.expect(token.LeftCurly)
		if err != nil {
			return
		}
		var body ast.BlockStatement
		body, err = p.block()
		if err != nil {
			return
		}
		elseClause = &ast.ElseClause{
			Body: body,
		}
	}

	return ast.IfStatement{
		If:     ifClause,
		ElseIf: elseIf,
		Else:   elseClause,
	}, nil
}

func (p *Parser) ifClause() (clause ast.IfClause, err error) {
	pos := p.tok.Pos

	p.advance() // skip "if"
	expression, err := p.expr()
	if err != nil {
		return
	}
	err = p.expect(token.LeftCurly)
	if err != nil {
		return
	}
	body, err := p.block()
	if err != nil {
		return
	}

	return ast.IfClause{
		Pos:       pos,
		Condition: expression,
		Body:      body,
	}, nil
}

// ExpressionStatement or AssignStatement
func (p *Parser) otherStatement() (ast.Statement, error) {
	if p.tok.Kind == token.Identifier && p.next.Kind == token.LeftParentheses {
		name := p.idn()
		p.advance() // skip identifier

		args, err := p.callArguments()
		if err != nil {
			return nil, err
		}

		if p.tok.Kind != token.Semicolon {
			panic("chain after call not implemented")
		}
		p.advance() // skip ";"

		return ast.SymbolCallStatement{
			Callee:    name,
			Arguments: args,
		}, nil
	}

	if p.tok.Kind != token.Identifier && p.tok.Kind != token.Receiver {
		return nil, p.unexpected(p.tok)
	}

	if p.tok.Kind == token.Identifier && p.next.Kind == token.Assign {
		return p.symbolAssignStatement()
	}

	var err error
	var target ast.ChainOperand

	if p.tok.Kind == token.Identifier && p.next.Kind == token.Indirect {
		idn := p.idn()
		p.advance() // skip identifier
		p.advance() // skip ".@"

		if p.tok.Kind == token.Assign {
			return p.indirectAssignStatement(idn)
		}
		target, err = p.chainOperand(ast.IndirectExpression{
			Target:     ast.ChainStart{Identifier: idn},
			ChainDepth: 1,
		})
		if err != nil {
			return nil, err
		}
	} else if p.tok.Kind == token.Identifier {
		identifier := p.idn()
		p.advance() // skip identifier

		target, err = p.chainOperand(ast.ChainStart{Identifier: identifier})
		if err != nil {
			return nil, err
		}
	} else if p.tok.Kind == token.Receiver {
		pos := p.pos()
		p.advance() // skip "g"

		if p.tok.Kind == token.Assign {
			return nil, fmt.Errorf("%s: receiver cannot be assigned to", p.tok.Kind.String())
		}

		target, err = p.chainOperand(ast.Receiver{Pos: pos})
		if err != nil {
			return nil, err
		}
	} else {
		panic(fmt.Sprintf("unexpected token %s at %s", p.tok.Kind.String(), p.tok.Pos.String()))
	}

	switch p.tok.Kind {
	case token.Assign, token.AddAssign:
		// continue parse assign statement
	default:
		return p.continueExpressionStatement(target)
	}

	kind := target.Kind()
	if kind == exn.Call {
		return nil, fmt.Errorf("cannot assign to call (%s) result value", target.Pin().String())
	}

	assignTokenKind := p.tok.Kind
	p.advance() // skip assign token

	expr, err := p.expr()
	if err != nil {
		return nil, err
	}

	if p.tok.Kind != token.Semicolon {
		return nil, p.unexpected(p.tok)
	}
	p.advance() // skip ";"

	switch assignTokenKind {
	case token.Assign:
		return ast.AssignStatement{
			Target:     target,
			Expression: expr,
		}, nil
	case token.AddAssign:
		return ast.AddAssignStatement{
			Target:     target,
			Expression: expr,
		}, nil
	default:
		panic("unexpected assign token: " + assignTokenKind.String())
	}
}

func (p *Parser) symbolAssignStatement() (ast.SymbolAssignStatement, error) {
	target := p.idn()
	p.advance() // skip target identifier

	p.advance() // skip "="

	expr, err := p.expr()
	if err != nil {
		return ast.SymbolAssignStatement{}, err
	}

	if p.tok.Kind != token.Semicolon {
		return ast.SymbolAssignStatement{}, p.unexpected(p.tok)
	}
	p.advance() // skip ";"

	return ast.SymbolAssignStatement{
		Target:     target,
		Expression: expr,
	}, nil
}

func (p *Parser) indirectAssignStatement(target ast.Identifier) (ast.IndirectAssignStatement, error) {
	p.advance() // skip "="

	expr, err := p.expr()
	if err != nil {
		return ast.IndirectAssignStatement{}, err
	}

	if p.tok.Kind != token.Semicolon {
		return ast.IndirectAssignStatement{}, p.unexpected(p.tok)
	}
	p.advance() // skip ";"

	return ast.IndirectAssignStatement{
		Target:     target,
		Expression: expr,
	}, nil
}

func (p *Parser) continueExpressionStatement(operand ast.Operand) (ast.ExpressionStatement, error) {
	expr, err := p.continueExpressionFromOperand(operand)
	if err != nil {
		return ast.ExpressionStatement{}, err
	}

	kind := expr.Kind()
	if kind != exn.Call {
		return ast.ExpressionStatement{},
			fmt.Errorf("%s: standalone expression in statement must be call expression: %s",
				expr.Pin().String(), kind.String())
	}

	if p.tok.Kind != token.Semicolon {
		return ast.ExpressionStatement{}, p.unexpected(p.tok)
	}

	p.advance() // consume ";"
	return ast.ExpressionStatement{Expression: expr}, nil
}

// func (p *Parser) continueParseMultipleAssignStatement(
// 	first ast.AssignableExpression) (statement ast.MultipleAssignStatement, err error) {

// 	targets := []ast.AssignableExpression{first}
// 	for {
// 		if p.tok.Kind == token.ShortAssign {
// 			break
// 		}
// 		if p.tok.Kind == token.Comma {
// 			p.advance() // consume ","
// 		} else {
// 			err = fmt.Errorf("missing \",\" inside multiple assign statement, got {%v}", p.tok)
// 			return
// 		}
// 		var target ast.AssignableExpression
// 		target, err = p.expr()
// 		if err != nil {
// 			return
// 		}
// 		targets = append(targets, target)
// 	}

// 	operator := p.aop()
// 	p.advance() // consume assign operator

// 	expressions := []ast.Expression{}

// 	for {
// 		var expression ast.Expression
// 		expression, err = p.expr()
// 		if err != nil {
// 			return
// 		}
// 		expressions = append(expressions, expression)
// 		if p.tok.Kind == token.Comma {
// 			p.advance() // consume ","
// 		} else if p.tok.Kind == token.Semicolon {
// 			break
// 		} else {
// 			err = fmt.Errorf("unexpected token at the end of multiple assignment {%v}", p.tok)
// 			return
// 		}
// 	}
// 	p.advance() // consume terminator

// 	if len(targets) != len(expressions) {
// 		err = fmt.Errorf("left and right sides of assignment have different number of expressions: %d and %d at %v",
// 			len(targets), len(expressions), p.tok.Pos)
// 		return
// 	}
// 	statement = ast.MultipleAssignStatement{
// 		Targets:     targets,
// 		Operator:    operator,
// 		Expressions: expressions,
// 	}
// 	return
// }

func (p *Parser) varStatement() (statement ast.VarStatement, err error) {
	pos := p.pos()

	p.advance() // skip "var"
	err = p.expect(token.Identifier)
	if err != nil {
		return
	}
	name := p.idn()
	p.advance() // skip identifier

	err = p.expect(token.Colon)
	if err != nil {
		return
	}
	p.advance() // skip ":"
	specifier, err := p.typeSpecifier()
	if err != nil {
		return
	}
	err = p.expect(token.Assign)
	if err != nil {
		return
	}
	p.advance() // skip "="
	if p.tok.Kind == token.Dirty {
		p.advance() // skip "dirty"

		err = p.expect(token.Semicolon)
		if err != nil {
			return
		}
		p.advance() // consume ";"

		return ast.VarStatement{
			VarInit: ast.VarInit{
				Pos:  pos,
				Name: name,
				Type: specifier,
			},
		}, nil
	}
	expr, err := p.expr()
	if err != nil {
		return
	}
	err = p.expect(token.Semicolon)
	if err != nil {
		return
	}
	p.advance() // consume ";"

	return ast.VarStatement{
		VarInit: ast.VarInit{
			Pos:        pos,
			Name:       name,
			Type:       specifier,
			Expression: expr,
		},
	}, nil
}

func (p *Parser) constStatement() (statement ast.ConstStatement, err error) {
	pos := p.tok.Pos

	p.advance() // skip "const"
	err = p.expect(token.Identifier)
	if err != nil {
		return
	}
	name := p.idn()
	p.advance() // skip const name identifier

	err = p.expect(token.Colon)
	if err != nil {
		return
	}
	p.advance() // skip ":"
	specifier, err := p.typeSpecifier()
	if err != nil {
		return
	}
	err = p.expect(token.Assign)
	if err != nil {
		return
	}
	p.advance() // skip "="
	expression, err := p.expr()
	if err != nil {
		return
	}
	err = p.expect(token.Semicolon)
	if err != nil {
		return
	}
	p.advance() // consume ";"

	return ast.ConstStatement{
		ConstInit: ast.ConstInit{
			Pos:        pos,
			Name:       name,
			Type:       specifier,
			Expression: expression,
		},
	}, nil
}

func (p *Parser) letStatement() (statement ast.LetStatement, err error) {
	pos := p.tok.Pos

	p.advance() // skip "let"
	err = p.expect(token.Identifier)
	if err != nil {
		return
	}
	name := p.idn()
	p.advance() // skip let name identifier

	err = p.expect(token.Colon)
	if err != nil {
		return
	}
	p.advance() // skip ":"
	specifier, err := p.typeSpecifier()
	if err != nil {
		return
	}
	err = p.expect(token.Assign)
	if err != nil {
		return
	}
	p.advance() // skip "="
	expression, err := p.expr()
	if err != nil {
		return
	}
	err = p.expect(token.Semicolon)
	if err != nil {
		return
	}
	p.advance() // consume ";"

	return ast.LetStatement{
		Pos:        pos,
		Name:       name,
		Type:       specifier,
		Expression: expression,
	}, nil
}

func (p *Parser) block() (block ast.BlockStatement, err error) {
	block.Pos = p.tok.Pos
	p.advance() // consume "{"
	for {
		if p.tok.Kind == token.RightCurly {
			p.advance() // consume "}"
			return
		}

		var statement ast.Statement
		statement, err = p.parseStatement()
		if err != nil {
			return
		}
		block.Statements = append(block.Statements, statement)
	}
}
