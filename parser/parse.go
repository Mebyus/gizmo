package parser

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

// create identifier from current token
func (p *Parser) word() ast.Identifier {
	return ast.Identifier{
		Pos: p.tok.Pos,
		Lit: p.tok.Lit,
	}
}

func (p *Parser) identifier() (ast.Identifier, error) {
	if p.tok.Kind != token.Identifier {
		return ast.Identifier{}, p.unexpected(p.tok)
	}
	identifier := p.word()
	p.advance()
	return identifier, nil
}

// gives a copy of current parser position
func (p *Parser) pos() source.Pos {
	return p.tok.Pos
}

func (p *Parser) basic() ast.BasicLiteral {
	return ast.BasicLiteral{Token: p.tok}
}

func (p *Parser) expect(k token.Kind) error {
	if p.tok.Kind == k {
		return nil
	}
	return p.unexpected(p.tok)
}

func (p *Parser) unit() (*ast.UnitClause, error) {
	p.advance() // skip "unit"
	if p.tok.Kind != token.Identifier {
		return nil, p.unexpected(p.tok)
	}
	name := p.word()
	p.advance() // skip unit name identifier

	return &ast.UnitClause{Name: name}, nil
}

func (p *Parser) top() error {
	err := p.gatherProps()
	if err != nil {
		return err
	}
	traits := ast.Traits{Props: p.takeProps()}

	switch p.tok.Kind {
	case token.Type:
		return p.topType(traits)
	case token.Var:
		return p.topVar(traits)
	case token.Fun:
		return p.topFun(traits)
	case token.Let:
		return p.topConstant(traits)
	case token.Test:
		return p.topTest(traits)
	case token.Pub:
		traits.Pub = true
		return p.topPub(traits)
	default:
		return p.unexpected(p.tok)
	}
}

func (p *Parser) topVar(traits ast.Traits) error {
	s, err := p.varStatement()
	if err != nil {
		return err
	}

	v := ast.TopVar{
		Var:    s.Var,
		Traits: traits,
	}
	p.atom.Vars = append(p.atom.Vars, v)
	return nil
}

func (p *Parser) topConstant(traits ast.Traits) error {
	s, err := p.letStatement()
	if err != nil {
		return err
	}

	c := ast.TopLet{
		Let:    s.Let,
		Traits: traits,
	}
	p.atom.Constants = append(p.atom.Constants, c)
	return nil
}

func (p *Parser) topPub(traits ast.Traits) error {
	p.advance() // skip "pub"

	switch p.tok.Kind {
	case token.Type:
		return p.topType(traits)
	case token.Fun:
		return p.topFun(traits)
	case token.Let:
		return p.topConstant(traits)
	default:
		return p.unexpected(p.tok)
	}
}
