package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/source"
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

	if p.tok.Kind == token.RightSquare {
		return ast.Prop{}, fmt.Errorf("empty prop at %s", pos.String())
	}

	if p.tok.Kind == token.Identifier && (p.next.Kind == token.RightSquare || p.next.Kind == token.Comma) {
		return p.tagsProp(pos)
	}

	key, err := p.propKey()
	if err != nil {
		return ast.Prop{}, err
	}

	if p.tok.Kind != token.Assign {
		return ast.Prop{}, p.unexpected()
	}
	p.advance() // skip "="

	value, err := p.propValue()
	if err != nil {
		return ast.Prop{}, err
	}

	if p.tok.Kind != token.RightSquare {
		return ast.Prop{}, p.unexpected()
	}
	p.advance() // skip "]"

	return ast.Prop{
		Pos:   pos,
		Key:   key,
		Value: value,
	}, nil
}

func (p *Parser) tagsProp(pos source.Pos) (ast.Prop, error) {
	list, err := p.tagsList()
	if err != nil {
		return ast.Prop{}, err
	}
	p.advance() // skip "]"

	return ast.Prop{
		Pos:   pos,
		Value: ast.PropValueTags{Tags: list},
	}, nil
}

func (p *Parser) tagsList() ([]ast.Identifier, error) {
	var tags []ast.Identifier
	for {
		if p.tok.Kind == token.RightSquare {
			if len(tags) == 0 {
				panic("no tags in list")
			}
			return tags, nil
		}

		tag, err := p.identifier()
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.RightSquare {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected()
		}
	}
}

func (p *Parser) propKey() (string, error) {
	if p.tok.Kind != token.Identifier {
		return "", p.unexpected()
	}
	identifier := p.word()
	p.advance() // skip identifier

	s := identifier.Lit

	for {
		if p.tok.Kind != token.Period {
			return s, nil
		}
		p.advance() // skip "."

		if p.tok.Kind != token.Identifier {
			return "", p.unexpected()
		}
		identifier := p.word()
		p.advance() // skip identifier

		s += "." + identifier.Lit
	}
}

func (p *Parser) propValue() (ast.PropValue, error) {
	tok := p.tok
	p.advance()

	switch tok.Kind {
	case token.DecInteger:
		return ast.PropValueInteger{Pos: tok.Pos, Val: tok.Val}, nil
	case token.String:
		return ast.PropValueString{Pos: tok.Pos, Val: tok.Lit}, nil
	case token.True:
		return ast.PropValueBool{Pos: tok.Pos, Val: true}, nil
	case token.False:
		return ast.PropValueBool{Pos: tok.Pos, Val: false}, nil
	default:
		return nil, p.unexpected()
	}
}

func (p *Parser) takeProps() *[]ast.Prop {
	if len(p.props) == 0 {
		return nil
	}
	props := p.props
	p.props = nil
	return &props
}
