package parser

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) method() (ast.Method, error) {
	p.advance() // skip "method"

	if p.tok.Kind != token.LeftParentheses {
		return ast.Method{}, p.unexpected(p.tok)
	}
	p.advance() // skip "("

	receiver, err := p.typeSpecifier()
	if err != nil {
		return ast.Method{}, err
	}

	if p.tok.Kind != token.RightParentheses {
		return ast.Method{}, p.unexpected(p.tok)
	}
	p.advance() // skip ")"

	if p.tok.Kind != token.Identifier {
		return ast.Method{}, p.unexpected(p.tok)
	}
	name := p.idn()
	p.advance() // skip method name

	signature, err := p.functionSignature()
	if err != nil {
		return ast.Method{}, err
	}

	if p.tok.Kind != token.LeftCurly {
		return ast.Method{}, p.unexpected(p.tok)
	}

	body, err := p.block()
	if err != nil {
		return ast.Method{}, err
	}

	return ast.Method{
		Receiver:  receiver,
		Name:      name,
		Signature: signature,
		Body:      body,
	}, nil
}
