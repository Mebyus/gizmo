package parser

import (
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

func ParseBytes(b []byte) (ast.UnitAtom, error) {
	p := FromBytes(b)
	return p.parse()
}

func ParseFile(filename string) (ast.UnitAtom, error) {
	p, err := FromFile(filename)
	if err != nil {
		return ast.UnitAtom{}, err
	}
	return p.parse()
}

func ParseSource(src *source.File) (ast.UnitAtom, error) {
	p, err := FromSource(src)
	if err != nil {
		return ast.UnitAtom{}, err
	}
	return p.parse()
}

func Parse(r io.Reader) (ast.UnitAtom, error) {
	p, err := FromReader(r)
	if err != nil {
		return ast.UnitAtom{}, err
	}
	return p.parse()
}

func (p *Parser) parse() (ast.UnitAtom, error) {
	var unit *ast.UnitBlock
	var err error
	if p.tok.Kind == token.Unit {
		unit, err = p.unitBlock()
		if err != nil {
			return ast.UnitAtom{}, err
		}
	}

	var blocks []ast.NamespaceBlock
	for {
		if p.isEOF() {
			return ast.UnitAtom{
				Unit:   unit,
				Blocks: blocks,
			}, nil
		}

		block, err := p.namespaceBlock()
		if err != nil {
			return ast.UnitAtom{}, err
		}
		blocks = append(blocks, block)
	}
}
