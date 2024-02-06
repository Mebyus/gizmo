package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) parseStatement() (statement ast.Statement, err error) {
	if p.tok.Kind == token.LeftCurly {
		return p.block()
	}
	if p.tok.Kind == token.Const {
		return p.constStatement()
	}
	// if p.tok.Kind == token.Var {
	// 	statement, err = p.parseVariableStatement()
	// 	return
	// }
	// if p.tok.Kind == token.If {
	// 	statement, err = p.parseIfStatement()
	// 	return
	// }
	if p.tok.Kind == token.Return {
		return p.returnStatement()
	}
	// if p.tok.Kind == token.For {
	// 	statement, err = p.loop()
	// 	return
	// }
	return p.parseExpressionStartStatement()
}

func (p *Parser) returnStatement() (statement ast.ReturnStatement, err error) {
	pos := p.tok.Pos
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

// func (p *Parser) parseIfStatement() (statement ast.IfStatement, err error) {
// 	ifClause, err := p.parseIfClause()
// 	if err != nil {
// 		return
// 	}

// 	var elseClause *ast.ElseClause
// 	if p.tok.Kind == token.Else {
// 		p.advance() // skip "else"
// 		err = p.expect(token.LeftCurly)
// 		if err != nil {
// 			return
// 		}
// 		var body ast.BlockStatement
// 		body, err = p.block()
// 		if err != nil {
// 			return
// 		}
// 		elseClause = &ast.ElseClause{
// 			Body: body,
// 		}
// 	}
// 	statement = ast.IfStatement{
// 		If:   ifClause,
// 		Else: elseClause,
// 	}
// 	return
// }

// func (p *Parser) parseIfClause() (clause ast.IfClause, err error) {
// 	p.advance() // skip "if"
// 	expression, err := p.expr()
// 	if err != nil {
// 		return
// 	}
// 	err = p.expect(token.LeftCurly)
// 	if err != nil {
// 		return
// 	}
// 	body, err := p.block()
// 	if err != nil {
// 		return
// 	}
// 	clause = ast.IfClause{
// 		Condition: expression,
// 		Body:      body,
// 	}
// 	return
// }

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

	if p.tok.IsIdent() && p.next.Kind == token.Assign {
		target := p.idn()
		p.advance() // skip identifier
		p.advance() // skip "="
		var expr ast.Expression
		expr, err = p.expr()
		if err != nil {
			return
		}
		err = p.expect(token.Semicolon)
		if err != nil {
			return
		}
		p.advance() // skip ";"
		statement = ast.AssignStatement{
			Target:     target,
			Expression: expr,
		}
		return
	}

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

func (p *Parser) parseExpressionStartStatement() (statement ast.Statement, err error) {
	statement, err = p.tryParseConstStatement()
	if err != nil {
		return nil, err
	}
	if statement != nil {
		return statement, nil
	}

	// expression, err := p.expr()
	// if err != nil {
	// 	return
	// }
	// if p.tok.Kind == token.Semicolon {
	// 	p.advance() // consume ";"
	// 	statement = ast.ExpressionStatement{
	// 		Expression: expression,
	// 	}
	// 	return
	// }
	// if p.tok.Kind == token.Comma {
	// 	statement, err = p.continueParseMultipleAssignStatement(expression)
	// 	return
	// }

	return nil, fmt.Errorf("other statement types not implemented")
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

// func (p *Parser) parseVariableStatement() (statement ast.VarStatement, err error) {
// 	p.advance() // skip "var"
// 	err = p.expect(token.Identifier)
// 	if err != nil {
// 		return
// 	}
// 	name := p.ident()
// 	p.advance() // skip identifier

// 	if p.tok.Kind == token.ShortAssign {
// 		p.advance() // skip ":="
// 		var expr ast.Expression
// 		expr, err = p.expr()
// 		if err != nil {
// 			return
// 		}
// 		err = p.expect(token.Semicolon)
// 		if err != nil {
// 			return
// 		}
// 		p.advance() // consume ";"
// 		statement = ast.VarShortAssign{
// 			Name:       name,
// 			Expression: expr,
// 		}
// 		return
// 	}

// 	err = p.expect(token.Colon)
// 	if err != nil {
// 		return
// 	}
// 	p.advance() // skip ":"
// 	specifier, err := p.parseTypeSpecifier()
// 	if err != nil {
// 		return
// 	}
// 	if p.tok.Kind == token.Semicolon {
// 		p.advance() // skip ";"
// 		statement = ast.VarDeclaration{
// 			Name: name,
// 			Type: specifier,
// 		}
// 		return
// 	}
// 	err = p.expect(token.Assign)
// 	if err != nil {
// 		return
// 	}
// 	p.advance() // skip "="
// 	if p.tok.Kind == token.Dirty {
// 		p.advance() // skip "dirty"

// 		err = p.expect(token.Semicolon)
// 		if err != nil {
// 			return
// 		}
// 		p.advance() // consume ";"
// 		statement = ast.VarDirty{
// 			Name: name,
// 			Type: specifier,
// 		}
// 		return
// 	}
// 	expression, err := p.expr()
// 	if err != nil {
// 		return
// 	}
// 	err = p.expect(token.Semicolon)
// 	if err != nil {
// 		return
// 	}
// 	p.advance() // consume ";"
// 	statement = ast.VarDefinition{
// 		Declaration: ast.VarDeclaration{
// 			Name: name,
// 			Type: specifier,
// 		},
// 		Expression: expression,
// 	}
// 	return
// }

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
