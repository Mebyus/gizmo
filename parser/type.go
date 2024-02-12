package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) typeSpecifier() (ast.TypeSpecifier, error) {
	if p.tok.Kind == token.Identifier {
		name, err := p.scopedIdentifier()
		if err != nil {
			return nil, err
		}
		return ast.TypeName{Name: name}, nil
	}
	if p.tok.Kind == token.Asterisk {
		return p.pointerType()
	}
	return nil, fmt.Errorf("other type specifiers not implemented %s", p.tok.Short())
}

func (p *Parser) pointerType() (ast.PointerType, error) {
	pos := p.pos()

	p.advance() // skip "*"

	ref, err := p.typeSpecifier()
	if err != nil {
		return ast.PointerType{}, err
	}

	return ast.PointerType{
		Pos:     pos,
		RefType: ref,
	}, nil
}
