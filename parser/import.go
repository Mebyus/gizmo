package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) imports() ([]ast.ImportBlock, error) {
	var blocks []ast.ImportBlock
	for {
		if p.tok.Kind == token.Import {
			block, err := p.topLevelImport()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, block)
		} else if p.tok.Kind == token.Pub && p.next.Kind == token.Import {
			p.advance() // skip "pub"
			block, err := p.topLevelImport()
			if err != nil {
				return nil, err
			}
			block.Pub = true
			blocks = append(blocks, block)
		} else {
			return blocks, nil
		}
	}
}

func (p *Parser) topLevelImport() (ast.ImportBlock, error) {
	pos := p.pos()
	p.advance() // skip "import"

	origin, err := p.importOrigin()
	if err != nil {
		return ast.ImportBlock{}, err
	}

	err = p.expect(token.LeftCurly)
	if err != nil {
		return ast.ImportBlock{}, err
	}
	p.advance() // skip "{"

	var specs []ast.ImportSpec
	for {
		if p.tok.Kind == token.RightCurly {
			p.advance() // skip "}"
			return ast.ImportBlock{
				Pos:    pos,
				Specs:  specs,
				Origin: origin,
			}, nil
		}

		spec, err := p.importSpec()
		if err != nil {
			return ast.ImportBlock{}, err
		}
		specs = append(specs, spec)
	}
}

func (p *Parser) importOrigin() (origin.Origin, error) {
	if p.tok.Kind != token.Identifier {
		return origin.Loc, nil
	}

	name := p.idn()
	p.advance() //  skip import origin name

	origin, ok := origin.Parse(name.Lit)
	if !ok {
		return 0, fmt.Errorf("unknown import origin \"%s\" at %s", name.Lit, name.Pos.String())
	}
	return origin, nil
}

func (p *Parser) importSpec() (ast.ImportSpec, error) {
	err := p.expect(token.Identifier)
	if err != nil {
		return ast.ImportSpec{}, err
	}
	name := p.idn()
	p.advance() // skip import name identifier

	err = p.expect(token.RightArrow)
	if err != nil {
		return ast.ImportSpec{}, err
	}
	p.advance() // skip "=>"

	err = p.expect(token.String)
	if err != nil {
		return ast.ImportSpec{}, err
	}
	if p.tok.Lit == "" {
		return ast.ImportSpec{}, fmt.Errorf("empty import string at %s", p.tok.Pos.String())
	}
	str := ast.ImportString{
		Pos: p.tok.Pos,
		Lit: p.tok.Lit,
	}
	p.advance() // skip import string

	return ast.ImportSpec{
		Name:   name,
		String: str,
	}, nil
}
