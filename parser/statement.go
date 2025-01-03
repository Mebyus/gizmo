package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/aop"
	"github.com/mebyus/gizmo/ast/lbl"
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/source"
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
		return p.forStm()
	case token.Jump:
		return p.jumpStatement()
	case token.Never:
		return p.neverStatement()
	case token.Stub:
		return p.stubStatement()
	case token.Defer:
		return p.deferStatement()
	case token.Identifier:
		return p.identifierStartStatement()
	default:
		return nil, p.unexpected()
	}
}

func (p *Parser) deferStatement() (ast.DeferStatement, error) {
	pos := p.pos()
	p.advance() // skip "defer"

	if p.tok.Kind != token.Identifier {
		return ast.DeferStatement{}, p.unexpected()
	}

	start := p.word()
	p.advance() // skip identifier

	exp, err := p.chainExp(start)
	if err != nil {
		return ast.DeferStatement{}, err
	}
	if exp.Kind() != exk.Call {
		return ast.DeferStatement{}, fmt.Errorf("%s: only call statements can be deferred", pos.String())
	}

	if p.tok.Kind != token.Semicolon {
		return ast.DeferStatement{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.DeferStatement{
		Pos:  pos,
		Call: exp.(ast.CallExp),
	}, nil
}

func (p *Parser) stubStatement() (ast.StubStatement, error) {
	pos := p.pos()
	p.advance() // skip "stub"

	if p.tok.Kind != token.Semicolon {
		return ast.StubStatement{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.StubStatement{Pos: pos}, nil
}

func (p *Parser) neverStatement() (ast.NeverStatement, error) {
	pos := p.pos()
	p.advance() // skip "never"

	if p.tok.Kind != token.Semicolon {
		return ast.NeverStatement{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.NeverStatement{Pos: pos}, nil
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
	case token.LabelOut:
		p.advance()
		return ast.ReservedLabel{Pos: pos, ResKind: lbl.Out}, nil
	default:
		panic("not implemented for label " + p.tok.Kind.String())
	}
}

func (p *Parser) matchBool() (ast.MatchBoolStatement, error) {
	pos := p.pos()
	p.advance() // skip "if"

	var node ast.MatchBoolStatement
	err := p.fillMatchBoolCases(&node)
	if err != nil {
		return ast.MatchBoolStatement{}, err
	}

	node.Pos = pos
	return node, nil
}

func (p *Parser) fillMatchBoolCases(node *ast.MatchBoolStatement) error {
	for {
		switch p.tok.Kind {
		case token.RightArrow:
			c, err := p.matchBoolCase()
			if err != nil {
				return err
			}
			node.Cases = append(node.Cases, c)
		case token.Else:
			block, err := p.matchElse()
			if err != nil {
				return err
			}
			node.Else = &block
			// else case is always last in match statement
			return nil
		default:
			return nil
		}
	}
}

func (p *Parser) fillMatchCases(node *ast.MatchStatement) error {
	for {
		switch p.tok.Kind {
		case token.RightArrow:
			c, err := p.matchCase()
			if err != nil {
				return err
			}
			node.Cases = append(node.Cases, c)
		case token.Else:
			block, err := p.matchElse()
			if err != nil {
				return err
			}
			node.Else = &block
			// else case is always last in match statement
			return nil
		default:
			return nil
		}
	}
}

// parse list of expressions in case
//
//	case <Exp>, <Exp>, <Exp>, ... {}
func (p *Parser) caseExpList() ([]ast.Exp, error) {
	p.advance() // skip "=>"

	var list []ast.Exp
	for {
		if p.tok.Kind == token.LeftCurly {
			if len(list) == 0 {
				return nil, fmt.Errorf("%s: case with no expressions", p.tok.Pos)
			}
			return list, nil
		}

		expr, err := p.exp()
		if err != nil {
			return nil, err
		}
		list = append(list, expr)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.LeftCurly {
			// will cause return at next iteration
		} else {
			return nil, p.unexpected()
		}
	}
}

func (p *Parser) matchCase() (ast.MatchCase, error) {
	pos := p.pos()

	list, err := p.caseExpList()
	if err != nil {
		return ast.MatchCase{}, err
	}

	block, err := p.Block()
	if err != nil {
		return ast.MatchCase{}, err
	}

	return ast.MatchCase{
		Pos:     pos,
		ExpList: list,
		Body:    block,
	}, nil
}

func (p *Parser) matchStatement(pos source.Pos, exp ast.Exp) (ast.MatchStatement, error) {
	var node ast.MatchStatement
	err := p.fillMatchCases(&node)
	if err != nil {
		return ast.MatchStatement{}, err
	}

	node.Pos = pos
	node.Exp = exp
	return node, nil
}

func (p *Parser) matchBoolCase() (ast.MatchBoolCase, error) {
	pos := p.pos()
	p.advance() // skip "case"

	exp, err := p.exp()
	if err != nil {
		return ast.MatchBoolCase{}, err
	}

	block, err := p.Block()
	if err != nil {
		return ast.MatchBoolCase{}, err
	}

	return ast.MatchBoolCase{
		Pos:  pos,
		Exp:  exp,
		Body: block,
	}, nil
}

func (p *Parser) matchElse() (ast.Block, error) {
	p.advance() // skip "else"
	return p.Block()
}

func (p *Parser) forStm() (ast.Statement, error) {
	if p.next.Kind == token.LeftCurly {
		return p.forSimple()
	}

	p.advance() // skip "for"
	if p.tok.Kind == token.Identifier && (p.next.Kind == token.In || p.next.Kind == token.Colon) {
		return p.forIn()
	}
	return p.forIf()
}

func (p *Parser) forSimple() (ast.For, error) {
	pos := p.pos()

	p.advance() // skip "for"

	body, err := p.Block()
	if err != nil {
		return ast.For{}, err
	}
	if len(body.Statements) == 0 {
		return ast.For{}, fmt.Errorf("%s for loop without condition cannot have empty body", pos.Short())
	}

	return ast.For{
		Pos:  pos,
		Body: body,
	}, nil
}

func (p *Parser) forIn() (ast.Statement, error) {
	if p.tok.Kind != token.Identifier {
		return nil, p.unexpected()
	}

	name := p.word()
	p.advance() // skip name identifier

	if p.tok.Kind == token.Colon {
		panic("for in loop with custom types not implemented")
	}

	if p.tok.Kind != token.In {
		return nil, p.unexpected()
	}
	p.advance() // skip "in"

	if p.tok.Kind == token.Identifier && p.tok.Lit == "range" {
		return p.forRange(name)
	}

	exp, err := p.exp()
	if err != nil {
		return nil, err
	}

	if p.tok.Kind != token.LeftCurly {
		return nil, p.unexpected()
	}
	body, err := p.Block()
	if err != nil {
		return nil, err
	}

	return ast.ForIn{
		Name:     name,
		Iterator: exp,
		Body:     body,
	}, nil
}

func (p *Parser) forRange(name ast.Identifier) (ast.ForRange, error) {
	p.advance() // skip "range"

	if p.tok.Kind != token.LeftParentheses {
		return ast.ForRange{}, p.unexpected()
	}
	p.advance() // skip "("

	exp, err := p.exp()
	if err != nil {
		return ast.ForRange{}, nil
	}

	if p.tok.Kind != token.RightParentheses {
		return ast.ForRange{}, p.unexpected()
	}
	p.advance() // skip ")"

	body, err := p.Block()
	if err != nil {
		return ast.ForRange{}, nil
	}

	return ast.ForRange{
		Name:  name,
		Range: exp,
		Body:  body,
	}, nil
}

func (p *Parser) forIf() (ast.ForIf, error) {
	condition, err := p.exp()
	if err != nil {
		return ast.ForIf{}, err
	}
	if p.tok.Kind != token.LeftCurly {
		return ast.ForIf{}, p.unexpected()
	}
	body, err := p.Block()
	if err != nil {
		return ast.ForIf{}, err
	}

	return ast.ForIf{
		If:   condition,
		Body: body,
	}, nil
}

func (p *Parser) returnStatement() (statement ast.ReturnStatement, err error) {
	pos := p.pos()

	p.advance() // consume "return"

	if p.tok.Kind == token.Semicolon {
		p.advance() // consume ";"
		return ast.ReturnStatement{Pos: pos}, nil
	}

	expression, err := p.exp()
	if err != nil {
		return
	}
	err = p.expect(token.Semicolon)
	if err != nil {
		return
	}
	p.advance() // consume ";"
	statement = ast.ReturnStatement{
		Pos: pos,
		Exp: expression,
	}
	return
}

func (p *Parser) ifStatement() (statement ast.Statement, err error) {
	switch p.next.Kind {
	case token.RightArrow, token.Else:
		return p.matchBool()
	}

	pos := p.tok.Pos
	p.advance() // skip "if"

	exp, err := p.exp()
	if err != nil {
		return
	}

	switch p.tok.Kind {
	case token.RightArrow, token.Else:
		return p.matchStatement(pos, exp)
	case token.LeftCurly:
		// continue regular if statement
	default:
		return nil, p.unexpected()
	}

	body, err := p.Block()
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
		var body ast.Block
		body, err = p.Block()
		if err != nil {
			return
		}
		elseClause = &ast.ElseClause{
			Body: body,
		}
	}

	return ast.IfStatement{
		If: ast.IfClause{
			Condition: exp,
			Pos:       pos,
			Body:      body,
		},
		ElseIf: elseIf,
		Else:   elseClause,
	}, nil
}

func (p *Parser) ifClause() (clause ast.IfClause, err error) {
	pos := p.tok.Pos

	p.advance() // skip "if"
	expression, err := p.exp()
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

func (p *Parser) identifierStartStatement() (ast.Statement, error) {
	if p.next.Kind == token.Walrus {
		return p.shortInitStatement()
	}

	word := p.word()
	p.advance() // skip identifier

	return p.chainStartStatement(word)
}

func (p *Parser) shortInitStatement() (ast.ShortInitStatement, error) {
	name := p.word()
	p.advance() // skip name identifier

	p.advance() // skip ":="

	exp, err := p.exp()
	if err != nil {
		return ast.ShortInitStatement{}, err
	}

	if p.tok.Kind != token.Semicolon {
		return ast.ShortInitStatement{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.ShortInitStatement{
		Name: name,
		Init: exp,
	}, nil
}

func (p *Parser) chainStartStatement(word ast.Identifier) (ast.Statement, error) {
	exp, err := p.chainExp(word)
	if err != nil {
		return nil, err
	}
	op, ok := aop.FromToken(p.tok.Kind)
	if ok {
		if exp.Kind() != exk.Chain {
			return nil, fmt.Errorf("%s: cannot assign to %s expression", exp.Pin(), exp.Kind())
		}

		p.advance() // skip assign operator token
		return p.assignStatement(op, exp.(ast.ChainOperand))
	}

	if exp.Kind() != exk.Call {
		return nil, fmt.Errorf("%s: standalone expression in statement must be call expression", exp.Pin())
	}

	return p.callStatement(exp.(ast.CallExp))
}

func (p *Parser) assignStatement(op aop.Kind, chain ast.ChainOperand) (ast.AssignStatement, error) {
	exp, err := p.exp()
	if err != nil {
		return ast.AssignStatement{}, err
	}
	err = p.expect(token.Semicolon)
	if err != nil {
		return ast.AssignStatement{}, err
	}
	p.advance() // skip ";"

	return ast.AssignStatement{
		Operator: op,
		Chain:    chain,
		Exp:      exp,
	}, nil
}

func (p *Parser) callStatement(exp ast.CallExp) (ast.CallStatement, error) {
	if p.tok.Kind != token.Semicolon {
		return ast.CallStatement{}, p.unexpected()
	}
	p.advance() // consume ";"

	return ast.CallStatement{Call: exp}, nil
}

func (p *Parser) varStatement() (statement ast.VarStatement, err error) {
	pos := p.pos()

	p.advance() // skip "var"
	err = p.expect(token.Identifier)
	if err != nil {
		return
	}
	name := p.word()
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

	if p.tok.Kind == token.Semicolon {
		// empty init expression

		p.advance() // skip ";"
		return ast.VarStatement{
			Var: ast.Var{
				Pos:  pos,
				Name: name,
				Type: specifier,
			},
		}, nil
	}

	err = p.expect(token.Assign)
	if err != nil {
		return
	}
	p.advance() // skip "="

	exp, err := p.InitExp()
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
			Exp:  exp,
		},
	}, nil
}

func (p *Parser) letWalrusStatement() (statement ast.LetStatement, err error) {
	name := p.word()
	p.advance() // skip let name identifier
	p.advance() // skip ":="

	expression, err := p.exp()
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
	name := p.word()
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
	expression, err := p.exp()
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

func (p *Parser) Block() (block ast.Block, err error) {
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
