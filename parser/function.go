package parser

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) topFun(traits ast.Traits) error {
	if p.next.Kind == token.LeftSquare {
		return p.method(traits)
	}

	p.advance() // consume "fun"

	if p.tok.Kind != token.Identifier {
		return p.unexpected()
	}
	name := p.word()
	p.advance() // consume function name identifier

	signature, err := p.functionSignature()
	if err != nil {
		return err
	}

	if p.tok.Kind != token.LeftCurly {
		d := ast.TopDec{
			Signature: signature,
			Name:      name,
			Traits:    traits,
		}
		p.atom.Decs = append(p.atom.Decs, d)
		return nil
	}

	body, err := p.Block()
	if err != nil {
		return err
	}
	f := ast.TopFun{
		Signature: signature,
		Name:      name,
		Body:      body,
		Traits:    traits,
	}
	p.atom.Funs = append(p.atom.Funs, f)
	return nil
}

func (p *Parser) topTest(traits ast.Traits) error {
	p.advance() // consume "test"

	if p.tok.Kind != token.Identifier {
		return p.unexpected()
	}
	name := p.word()
	p.advance() // consume test name identifier

	signature, err := p.functionSignature()
	if err != nil {
		return err
	}

	if p.tok.Kind != token.LeftCurly {
		return p.unexpected()
	}

	body, err := p.Block()
	if err != nil {
		return err
	}
	f := ast.TopFun{
		Signature: signature,
		Name:      name,
		Body:      body,
		Traits:    traits,
	}
	p.atom.Tests = append(p.atom.Tests, f)
	return nil
}

func (p *Parser) functionSignature() (ast.Signature, error) {
	params, err := p.functionParams()
	if err != nil {
		return ast.Signature{}, err
	}

	if p.tok.Kind != token.RightArrow {
		return ast.Signature{Params: params}, nil
	}

	p.advance() // skip "=>"

	if p.tok.Kind == token.Never {
		p.advance() // skip "never"
		return ast.Signature{
			Params: params,
			Never:  true,
		}, nil
	}

	result, err := p.typeSpecifier()
	if err != nil {
		return ast.Signature{}, err
	}

	return ast.Signature{
		Params: params,
		Result: result,
	}, nil
}

func (p *Parser) functionParams() ([]ast.FieldDefinition, error) {
	err := p.expect(token.LeftParentheses)
	if err != nil {
		return nil, err
	}
	p.advance() // skip "("

	var params []ast.FieldDefinition
	for {
		if p.tok.Kind == token.RightParentheses {
			p.advance() // skip ")"
			return params, nil
		}

		param, err := p.field()
		if err != nil {
			return nil, err
		}
		params = append(params, param)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.RightParentheses {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected()
		}
	}
}

func (p *Parser) field() (field ast.FieldDefinition, err error) {
	err = p.expect(token.Identifier)
	if err != nil {
		return
	}
	field.Name = p.word()
	p.advance() // skip identifier

	err = p.expect(token.Colon)
	if err != nil {
		return
	}
	p.advance() // consume ":"
	spec, err := p.typeSpecifier()
	if err != nil {
		return
	}
	field.Type = spec
	return
}
