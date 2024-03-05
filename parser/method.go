package parser

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

// returns receiver name and its type params, if method is not a template
// second return value will be nil slice
func (p *Parser) methodReceiver() (ast.Identifier, []ast.Identifier, error) {
	if p.tok.Kind != token.LeftParentheses {
		return ast.Identifier{}, nil, p.unexpected(p.tok)
	}
	p.advance() // skip "("

	if p.tok.Kind != token.Identifier {
		return ast.Identifier{}, nil, p.unexpected(p.tok)
	}
	receiver := p.idn()
	p.advance() // skip receiver name

	if p.tok.Kind == token.RightParentheses {
		// no type params
		p.advance() // skip ")"
		return receiver, nil, nil
	}

	if p.tok.Kind != token.LeftDoubleSquare {
		return ast.Identifier{}, nil, p.unexpected(p.tok)
	}

	params, err := p.templateParams()
	if err != nil {
		return ast.Identifier{}, nil, err
	}

	if p.tok.Kind != token.RightParentheses {
		return ast.Identifier{}, nil, p.unexpected(p.tok)
	}
	p.advance() // skip ")"

	return receiver, params, nil
}

func (p *Parser) method() (ast.TopLevel, error) {
	p.advance() // skip "method"

	receiver, params, err := p.methodReceiver()
	if err != nil {
		return nil, err
	}

	if p.tok.Kind != token.Identifier {
		return nil, p.unexpected(p.tok)
	}
	name := p.idn()
	p.advance() // skip method name

	signature, err := p.functionSignature()
	if err != nil {
		return nil, err
	}

	if p.tok.Kind != token.LeftCurly {
		return nil, p.unexpected(p.tok)
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	if len(params) == 0 {
		return ast.Method{
			Receiver:  receiver,
			Name:      name,
			Signature: signature,
			Body:      body,
		}, nil
	}
	return ast.MethodTemplate{
		Receiver:   receiver,
		TypeParams: params,
		Name:       name,
		Signature:  signature,
		Body:       body,
	}, nil
}
