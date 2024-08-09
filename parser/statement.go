package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/aop"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/ast/lbl"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) Statement() (ast.Statement, error) {
	switch p.tok.Kind {
	case token.LeftCurly:
		return p.Block()
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
	case token.Receiver:
		return p.receiverStartStatement()
	case token.Identifier:
		return p.identifierStartStatement()
	default:
		return nil, p.unexpected(p.tok)
	}
}

func (p *Parser) deferStatement() (ast.DeferStatement, error) {
	pos := p.pos()
	p.advance() // skip "defer"

	if p.tok.Kind != token.Identifier && p.tok.Kind != token.Receiver {
		return ast.DeferStatement{}, p.unexpected(p.tok)
	}

	var chain ast.ChainOperand
	if p.tok.Kind == token.Identifier {
		identifier := p.idn()
		p.advance() // skip identifier
		chain = identifier.AsChainOperand()
	} else {
		r := p.receiver()
		chain = r.AsChainOperand()
	}

	err := p.chainOperand(&chain)
	if err != nil {
		return ast.DeferStatement{}, err
	}
	if chain.Last() != exn.Call {
		return ast.DeferStatement{}, fmt.Errorf("%s: only call statements can be deferred", pos.String())
	}

	if p.tok.Kind != token.Semicolon {
		return ast.DeferStatement{}, p.unexpected(p.tok)
	}
	p.advance() // skip ";"

	return ast.DeferStatement{
		Pos:  pos,
		Call: chain,
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

	block, err := p.Block()
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

	block, err := p.Block()
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

	body, err := p.Block()
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
	body, err := p.Block()
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
	body, err := p.Block()
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
		body, err = p.Block()
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
	body, err := p.Block()
	if err != nil {
		return
	}

	return ast.IfClause{
		Pos:       pos,
		Condition: expression,
		Body:      body,
	}, nil
}

func (p *Parser) receiverStartStatement() (ast.Statement, error) {
	r := p.receiver()
	return p.chainStartStatement(r.AsIdentifier())
}

func (p *Parser) identifierStartStatement() (ast.Statement, error) {
	idn := p.idn()
	p.advance() // skip identifier
	return p.chainStartStatement(idn)
}

func (p *Parser) chainStartStatement(identifier ast.Identifier) (ast.Statement, error) {
	chain := identifier.AsChainOperand()
	err := p.chainOperand(&chain)
	if err != nil {
		return nil, err
	}
	op, ok := aop.FromToken(p.tok.Kind)
	if ok {
		p.advance() // skip assign operator token
		return p.assignStatement(op, chain)
	}

	return p.callStatement(chain)
}

func (p *Parser) assignStatement(op aop.Kind, target ast.ChainOperand) (ast.AssignStatement, error) {
	switch target.Last() {
	case exn.Call, exn.Address, exn.Slice, exn.Receiver:
		return ast.AssignStatement{}, fmt.Errorf("%s: cannot assign to %s operand",
			target.Pin().String(), target.Last().String())
	}

	expr, err := p.expr()
	if err != nil {
		return ast.AssignStatement{}, err
	}
	err = p.expect(token.Semicolon)
	if err != nil {
		return ast.AssignStatement{}, err
	}
	p.advance() // skip ";"

	return ast.AssignStatement{
		Operator:   op,
		Target:     target,
		Expression: expr,
	}, nil
}

func (p *Parser) callStatement(chain ast.ChainOperand) (ast.CallStatement, error) {
	if chain.Last() != exn.Call {
		return ast.CallStatement{},
			fmt.Errorf("%s: standalone expression in statement must be call expression",
				chain.Pin().String())
	}

	if p.tok.Kind != token.Semicolon {
		return ast.CallStatement{}, p.unexpected(p.tok)
	}
	p.advance() // consume ";"

	return ast.CallStatement{Call: chain}, nil
}

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
			Var: ast.Var{
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
		Var: ast.Var{
			Pos:  pos,
			Name: name,
			Type: specifier,
			Exp:  expr,
		},
	}, nil
}

func (p *Parser) letWalrusStatement() (statement ast.LetStatement, err error) {
	name := p.idn()
	p.advance() // skip let name identifier
	p.advance() // skip ":="

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
		Let: ast.Let{
			Name: name,
			Exp:  expression,
		},
	}, nil
}

func (p *Parser) letStatement() (statement ast.LetStatement, err error) {
	p.advance() // skip "let"
	err = p.expect(token.Identifier)
	if err != nil {
		return
	}
	if p.next.Kind == token.Walrus {
		return p.letWalrusStatement()
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
		Let: ast.Let{
			Name: name,
			Type: specifier,
			Exp:  expression,
		},
	}, nil
}

func (p *Parser) Block() (block ast.BlockStatement, err error) {
	block.Pos = p.pos()
	p.skip(token.LeftCurly)
	for {
		if p.tok.Kind == token.RightCurly {
			p.advance() // consume "}"
			return
		}

		var statement ast.Statement
		statement, err = p.Statement()
		if err != nil {
			return
		}
		block.Statements = append(block.Statements, statement)
	}
}
