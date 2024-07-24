package parser

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

// returns receiver name and its type params, if method is not a template
// second return value will be nil slice
func (p *Parser) methodReceiver() (ast.ReceiverTypeSpec, error) {
	if p.tok.Kind != token.LeftSquare {
		return ast.ReceiverTypeSpec{}, p.unexpected(p.tok)
	}
	p.advance() // skip "["

	var ptr bool

	if p.tok.Kind == token.Asterisk {
		p.advance() // skip "*"
		ptr = true
	}

	if p.tok.Kind != token.Identifier {
		return ast.ReceiverTypeSpec{}, p.unexpected(p.tok)
	}
	name := p.idn()
	p.advance() // skip receiver type name

	if p.tok.Kind != token.RightSquare {
		return ast.ReceiverTypeSpec{}, p.unexpected(p.tok)
	}
	p.advance() // skip "]"

	return ast.ReceiverTypeSpec{
		Name: name,
		Ptr:  ptr,
	}, nil
}

func (p *Parser) method(traits ast.Traits) error {
	p.advance() // skip "fn"

	receiver, err := p.methodReceiver()
	if err != nil {
		return err
	}

	if p.tok.Kind != token.Identifier {
		return p.unexpected(p.tok)
	}
	name := p.idn()
	p.advance() // skip method name

	signature, err := p.functionSignature()
	if err != nil {
		return err
	}

	if p.tok.Kind != token.LeftCurly {
		return p.unexpected(p.tok)
	}

	body, err := p.Block()
	if err != nil {
		return err
	}

	m := ast.Method{
		Receiver:  receiver,
		Name:      name,
		Signature: signature,
		Body:      body,
		Traits:    traits,
	}
	p.atom.Meds = append(p.atom.Meds, m)
	return nil
}
