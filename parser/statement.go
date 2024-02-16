package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) parseStatement() (statement ast.Statement, err error) {
	switch p.tok.Kind {
	case token.LeftCurly:
		return p.block()
	case token.Const:
		return p.constStatement()
	case token.Var:
		return p.varStatement()
	case token.If:
		return p.ifStatement()
	case token.Return:
		return p.returnStatement()
	case token.For:
		return p.forStatement()
	default:
		return p.identifierStartStatement()
	}
}

func (p *Parser) forStatement() (ast.Statement, error) {
	if p.next.Kind == token.LeftCurly {
		return p.forSimpleStatement()
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

func (p *Parser) tryParseConstStatement() (statement ast.Statement, err error) {
	// if p.tok.IsIdent() && p.next.Kind == token.ShortAssign {
	// 	target := p.ident()
	// 	p.advance() // skip identifier
	// 	p.advance() // skip ":="
	// 	var expr ast.Expression
	// 	expr, err = p.expr()
	// 	if err != nil {
	// 		return
	// 	}
	// 	err = p.expect(token.Semicolon)
	// 	if err != nil {
	// 		return
	// 	}
	// 	p.advance() // skip ";"
	// 	statement = ast.ShortAssign{
	// 		Name:       target,
	// 		Expression: expr,
	// 	}
	// 	return
	// }

	// if p.tok.IsIdent() && p.next.Kind == token.Assign {
	// 	target := p.idn()
	// 	p.advance() // skip identifier
	// 	p.advance() // skip "="
	// 	var expr ast.Expression
	// 	expr, err = p.expr()
	// 	if err != nil {
	// 		return
	// 	}
	// 	err = p.expect(token.Semicolon)
	// 	if err != nil {
	// 		return
	// 	}
	// 	p.advance() // skip ";"
	// 	statement = ast.AssignStatement{
	// 		Target:     target,
	// 		Expression: expr,
	// 	}
	// 	return
	// }

	// if p.tok.IsIdent() && p.next.Kind == token.AddAssign {
	// 	target := p.ident()
	// 	p.advance() // skip identifier
	// 	p.advance() // skip "+="
	// 	var expr ast.Expression
	// 	expr, err = p.expr()
	// 	if err != nil {
	// 		return
	// 	}
	// 	err = p.expect(token.Semicolon)
	// 	if err != nil {
	// 		return
	// 	}
	// 	p.advance() // skip ";"
	// 	statement = ast.AddAssign{
	// 		Target:     target,
	// 		Expression: expr,
	// 	}
	// 	return
	// }

	// if p.tok.IsIdent() && p.next.Kind == token.Colon {
	// 	name := p.ident()
	// 	p.advance() // skip identifier
	// 	p.advance() // skip ":"
	// 	var spec ast.TypeSpecifier
	// 	spec, err = p.parseTypeSpecifier()
	// 	if err != nil {
	// 		return
	// 	}
	// 	err = p.expect(token.Assign)
	// 	if err != nil {
	// 		return
	// 	}
	// 	p.advance() // skip "="
	// 	var expr ast.Expression
	// 	expr, err = p.expr()
	// 	if err != nil {
	// 		return
	// 	}
	// 	err = p.expect(token.Semicolon)
	// 	if err != nil {
	// 		return
	// 	}
	// 	p.advance() // skip ";"
	// 	statement = ast.TypedAssign{
	// 		Name:       name,
	// 		Type:       spec,
	// 		Expression: expr,
	// 	}
	// 	return
	// }

	return
}

func (p *Parser) tryAssignStatement() (ast.Statement, error) {
	if p.tok.Kind != token.Identifier {
		return nil, nil
	}

	identifier, err := p.scopedIdentifier()
	if err != nil {
		return nil, err
	}

	target, err := p.chainOperand(ast.ChainStart{Identifier: identifier})
	if err != nil {
		return nil, err
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

func (p *Parser) continueExpressionStatement(operand ast.Operand) (ast.ExpressionStatement, error) {
	expr, err := p.continueExpressionFromOperand(operand)
	if err != nil {
		return ast.ExpressionStatement{}, err
	}

	kind := expr.Kind()
	if kind != exn.Call {
		return ast.ExpressionStatement{},
			fmt.Errorf("standalone expression (%s) in statement must be call expression: %s", expr.Pin(), kind.String())
	}

	if p.tok.Kind != token.Semicolon {
		return ast.ExpressionStatement{}, p.unexpected(p.tok)
	}

	p.advance() // consume ";"
	return ast.ExpressionStatement{Expression: expr}, nil
}

func (p *Parser) identifierStartStatement() (ast.Statement, error) {
	return p.tryAssignStatement()
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
	pos := p.tok.Pos

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
