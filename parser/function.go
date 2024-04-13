package parser

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) topLevelFn() (ast.TopLevel, error) {
	if p.next.Kind == token.LeftSquare {
		return p.topLevelMethod()
	}

	p.advance() // consume "fn"
	if p.tok.Kind != token.Identifier {
		return nil, p.unexpected(p.tok)
	}

	declaration := ast.FunctionDeclaration{
		Name: p.idn(),
	}
	p.advance() // consume function name identifier

	signature, err := p.functionSignature()
	if err != nil {
		return nil, err
	}
	declaration.Signature = signature

	if p.tok.Kind != token.LeftCurly {
		return ast.TopFunctionDeclaration{
			Declaration: declaration,
			Props:       p.takeProps(),
		}, nil
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	definition := ast.FunctionDefinition{
		Head: declaration,
		Body: body,
	}
	return ast.TopFunctionDefinition{
		Definition: definition,
		Props:      p.takeProps(),
	}, nil
}

func (p *Parser) functionSignature() (ast.FunctionSignature, error) {
	params, err := p.functionParams()
	if err != nil {
		return ast.FunctionSignature{}, err
	}

	if p.tok.Kind != token.RightArrow {
		return ast.FunctionSignature{Params: params}, nil
	}

	p.advance() // skip "=>"

	if p.tok.Kind == token.Never {
		p.advance() // skip "never"
		return ast.FunctionSignature{
			Params: params,
			Never:  true,
		}, nil
	}

	result, err := p.typeSpecifier()
	if err != nil {
		return ast.FunctionSignature{}, err
	}

	return ast.FunctionSignature{
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
			return nil, p.unexpected(p.tok)
		}
	}
}

func (p *Parser) field() (field ast.FieldDefinition, err error) {
	err = p.expect(token.Identifier)
	if err != nil {
		return
	}
	field.Name = p.idn()
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
