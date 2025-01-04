package parser

import (
	"fmt"
	"io"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/token"
)

// Parser keeps track of atom parsing state. In most cases parser usage
// should be logically close to following:
//
//	src, err := source.Load(filename)
//	if err != nil {
//		// handle error
//	}
//	p := FromSource(src)
//	header, err := p.Header()
//	// analyze atom header
//	atom, err := p.Parse()
//	// process atom tree
type Parser struct {
	// stored output to transfer already parsed data
	// between external methods calls
	atom ast.Atom

	// saved properties, will be attached to next object/symbol/field
	props []ast.Prop

	lx lexer.Stream

	// token at current parser position
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

func FromSource(src *source.File) *Parser {
	lx := lexer.FromSource(src)
	p := New(lx)
	return p
}

func ParseBytes(b []byte) (*ast.Atom, error) {
	p := FromBytes(b)
	return p.FullParse()
}

func ParseFile(filename string) (*ast.Atom, error) {
	p, err := FromFile(filename)
	if err != nil {
		return nil, err
	}
	return p.FullParse()
}

func ParseSource(src *source.File) (*ast.Atom, error) {
	p := FromSource(src)
	return p.FullParse()
}

// Header returns list of all import paths specified in parsed atom.
// Should be called before calling Parse method. Under the hood this method
// parses first blocks in atom until it encounteres non-import top-level block.
//
// Resulting slice may be empty or contain non-unique entries.
func (p *Parser) Header() (*ast.AtomHeader, error) {
	err := p.header()
	if err != nil {
		return nil, err
	}
	return &p.atom.Header, nil
}

// Parse continues atom parsing from the state where Header method left.
// Resulting atom tree is complete and also contains info gathered by
// Header method.
func (p *Parser) Parse() (*ast.Atom, error) {
	err := p.parse()
	if err != nil {
		return nil, err
	}
	return &p.atom, nil
}

// FullParse combines Header and Parse into one method.
// This method is mostly a convenience for testing and prototyping.
func (p *Parser) FullParse() (*ast.Atom, error) {
	_, err := p.Header()
	if err != nil {
		return nil, err
	}
	return p.Parse()
}

func (p *Parser) header() error {
	var unit *ast.UnitClause
	var err error
	if p.tok.Kind == token.Unit {
		unit, err = p.unit()
		if err != nil {
			return err
		}
	}
	p.atom.Header.Unit = unit

	blocks, err := p.imports()
	if err != nil {
		return err
	}
	p.atom.Header.Imports.Blocks = blocks

	var paths []origin.Path
	for _, block := range blocks {
		for _, spec := range block.Specs {
			paths = append(paths, origin.Path{
				Origin: block.Origin,
				ImpStr: spec.String.Lit,
			})
		}
	}
	p.atom.Header.Imports.Paths = paths
	return nil
}

func (p *Parser) parse() error {
	for {
		if p.isEOF() {
			if p.tok.Lit != "" {
				return fmt.Errorf(p.tok.Lit)
			}
			return nil
		}

		err := p.top()
		if err != nil {
			return err
		}
	}
}
