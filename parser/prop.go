package parser

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) gatherProps() error {
	for {
		if p.tok.Kind != token.PropStart {
			return nil
		}

		prop, err := p.prop()
		if err != nil {
			return err
		}
		p.props = append(p.props, prop)
	}
}

func (p *Parser) prop() (ast.Prop, error) {
	pos := p.pos()
	p.advance() // skip "#["

	key, err := p.propKey()
	if err != nil {
		return ast.Prop{}, err
	}

	if p.tok.Kind != token.Assign {
		return ast.Prop{}, p.unexpected(p.tok)
	}
	p.advance() // skip "="

	value, err := p.propValue()
	if err != nil {
		return ast.Prop{}, err
	}

	if p.tok.Kind != token.RightSquare {
		return ast.Prop{}, p.unexpected(p.tok)
	}
	p.advance() // skip "]"

	return ast.Prop{
		Pos:   pos,
		Key:   key,
		Value: value,
	}, nil
}

func (p *Parser) propKey() (string, error) {
	if p.tok.Kind != token.Identifier {
		return "", p.unexpected(p.tok)
	}
	identifier := p.idn()
	p.advance() // skip identifier

	s := identifier.Lit

	for {
		if p.tok.Kind != token.Period {
			return s, nil
		}
		p.advance() // skip "."

		if p.tok.Kind != token.Identifier {
			return "", p.unexpected(p.tok)
		}
		identifier := p.idn()
		p.advance() // skip identifier

		s += "." + identifier.Lit
	}
}

func (p *Parser) propValue() (ast.PropValue, error) {
	tok := p.tok
	p.advance()

	switch tok.Kind {
	case token.DecimalInteger:
		return ast.PropValueInteger{Pos: tok.Pos, Val: tok.Val}, nil
	case token.String:
		return ast.PropValueString{Pos: tok.Pos, Val: tok.Lit}, nil
	case token.True:
		return ast.PropValueBool{Pos: tok.Pos, Val: true}, nil
	case token.False:
		return ast.PropValueBool{Pos: tok.Pos, Val: false}, nil
	default:
		return nil, p.unexpected(tok)
	}
}

func (p *Parser) takeProps() []ast.Prop {
	props := p.props
	p.props = nil
	return props
}
