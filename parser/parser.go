package parser

import (
	"fmt"
	"io"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

type Parser struct {
	lx lexer.Stream

	// previous token
	prev token.Token

	// token at current Parser position
	tok token.Token

	// next token
	next token.Token
}

func New(lx lexer.Stream) *Parser {
	p := &Parser{lx: lx}

	// init parser buffer
	p.advance()
	p.advance()
	return p
}

func FromReader(r io.Reader) (*Parser, error) {
	lx, err := lexer.FromReader(r)
	if err != nil {
		return nil, err
	}
	return New(lx), nil
}

func FromBytes(b []byte) *Parser {
	return New(lexer.FromBytes(b))
}

func FromFile(filename string) (*Parser, error) {
	lx, err := lexer.FromFile(filename)
	if err != nil {
		return nil, err
	}
	return New(lx), nil
}

func FromSource(src *source.File) (p *Parser, err error) {
	lx := lexer.FromSource(src)
	p = New(lx)
	return p, nil
}

func ParseBytes(b []byte) (ast.UnitBlock, error) {
	p := FromBytes(b)
	return p.parse()
}

func ParseFile(filename string) (ast.UnitBlock, error) {
	p, err := FromFile(filename)
	if err != nil {
		return ast.UnitBlock{}, err
	}
	return p.parse()
}

func ParseSource(src *source.File) (ast.UnitBlock, error) {
	p, err := FromSource(src)
	if err != nil {
		return ast.UnitBlock{}, err
	}
	return p.parse()
}

func Parse(r io.Reader) (ast.UnitBlock, error) {
	p, err := FromReader(r)
	if err != nil {
		return ast.UnitBlock{}, err
	}
	return p.parse()
}

func (p *Parser) parse() (ast.UnitBlock, error) {
	unit, err := p.parseUnitBlock()
	if err != nil {
		return ast.UnitBlock{}, err
	}
	if unit.Name.Kind.IsEmpty() {
		return ast.UnitBlock{}, fmt.Errorf("expected unit block")
	}
	return unit, nil
}
