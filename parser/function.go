package parser

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) topLevelFunction() (ast.TopFunctionDefinition, error) {
	p.advance() // consume "fn"
	if p.tok.Kind != token.Identifier {
		return ast.TopFunctionDefinition{}, p.unexpected(p.tok)
	}

	declaration := ast.FunctionDeclaration{
		Name: p.idn(),
	}
	p.advance() // consume function name identifier

	signature, err := p.functionSignature()
	if err != nil {
		return ast.TopFunctionDefinition{}, err
	}
	declaration.Signature = signature

	if p.tok.Kind != token.LeftCurly {
		return ast.TopFunctionDefinition{}, p.unexpected(p.tok)
	}
	body, err := p.block()
	if err != nil {
		return ast.TopFunctionDefinition{}, err
	}
	definition := ast.FunctionDefinition{
		Head: declaration,
		Body: body,
	}
	return ast.TopFunctionDefinition{
		Definition: definition,
	}, nil
}

func (p *Parser) functionSignature() (ast.FunctionSignature, error) {
	params, err := p.functionParams()
	if err != nil {
		return ast.FunctionSignature{}, err
	}

	var result ast.TypeSpecifier
	if p.tok.Kind == token.RightArrow {
		p.advance() // skip "=>"

		if p.tok.Kind == token.Never {
			p.advance() // skip "never"
			return ast.FunctionSignature{
				Params: params,
				Result: result,
				Never:  true,
			}, nil
		}

		result, err = p.typeSpecifier()
		if err != nil {
			return ast.FunctionSignature{}, err
		}
	}

	return ast.FunctionSignature{
		Params: params,
		Result: result,
	}, nil
}

func (p *Parser) functionParams() (params []ast.FieldDefinition, err error) {
	err = p.expect(token.LeftParentheses)
	if err != nil {
		return
	}
	p.advance() // skip "("

	first := true
	comma := false
	for {
		if p.tok.Kind == token.RightParentheses {
			p.advance() // skip ")"
			return params, nil
		}

		if first {
			first = false
		} else if comma {
			comma = false
		} else {
			return nil, p.unexpected(p.tok)
		}

		var field ast.FieldDefinition
		field, err := p.field()
		if err != nil {
			return nil, err
		}
		params = append(params, field)
		if p.tok.Kind == token.Comma {
			comma = true
			p.advance() // skip ","
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

func (p *Parser) topLevelDeclare() (ast.TopFunctionDeclaration, error) {
	p.advance() // consume "declare"
	if p.tok.Kind != token.Fn {
		return ast.TopFunctionDeclaration{}, p.unexpected(p.tok)
	}

	p.advance() // consume "fn"
	if p.tok.Kind != token.Identifier {
		return ast.TopFunctionDeclaration{}, p.unexpected(p.tok)
	}

	declaration := ast.FunctionDeclaration{
		Name: p.idn(),
	}
	p.advance() // consume function name identifier

	signature, err := p.functionSignature()
	if err != nil {
		return ast.TopFunctionDeclaration{}, err
	}
	declaration.Signature = signature

	return ast.TopFunctionDeclaration{
		Declaration: declaration,
	}, nil
}
