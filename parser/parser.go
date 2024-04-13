package parser

import (
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

	// tokens saved by advance backup
	buf CycleTokenBuffer

	// saved properties, will be attached to next object/symbol/field
	props []ast.Prop

	lx lexer.Stream

	// token at current parser position
	tok token.Token

	// next token
	next token.Token

	// how many tokens must be taken from backtrack buffer
	back uint
}

type CycleTokenBuffer struct {
	buf [16]token.Token

	// push index
	tip uint

	// number of stored elements
	len uint
}

func (b *CycleTokenBuffer) Push(tok token.Token) {
	if b.len >= 16 {
		panic("full")
	}
	b.buf[b.tip] = tok
	b.tip = (b.tip + 1) & 0xF // (pos + 1) mod 16
	b.len += 1
}

func (b *CycleTokenBuffer) Pop() token.Token {
	if b.len == 0 {
		panic("empty")
	}
	i := (b.tip - b.len) & 0xF // (pos - len) mod 16
	b.len -= 1
	return b.buf[i]
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

func ParseBytes(b []byte) (ast.Atom, error) {
	p := FromBytes(b)
	return p.parse()
}

func ParseFile(filename string) (ast.Atom, error) {
	p, err := FromFile(filename)
	if err != nil {
		return ast.Atom{}, err
	}
	return p.parse()
}

func ParseSource(src *source.File) (ast.Atom, error) {
	p := FromSource(src)
	return p.parse()
}

func Parse(r io.Reader) (ast.Atom, error) {
	p, err := FromReader(r)
	if err != nil {
		return ast.Atom{}, err
	}
	return p.parse()
}

// Header returns list of all import paths specified in parsed atom.
// Should be called before calling Parse method. Under the hood this method
// parses first blocks in atom until it encounteres non-import top-level block
//
// Resulting slice may be empty or contain non-unique entries
func (p *Parser) Header() (ast.AtomHeader, error) {
	var unit *ast.UnitBlock
	var err error
	if p.tok.Kind == token.Unit {
		unit, err = p.unitBlock()
		if err != nil {
			return ast.AtomHeader{}, err
		}
	}
	p.atom.Header.Unit = unit

	blocks, err := p.imports()
	if err != nil {
		return ast.AtomHeader{}, err
	}
	p.atom.Header.Imports.ImportBlocks = blocks

	var paths []origin.Path
	for _, block := range blocks {
		for _, spec := range block.Specs {
			paths = append(paths, origin.Path{
				Origin: block.Origin,
				ImpStr: spec.String.Lit,
			})
		}
	}
	p.atom.Header.Imports.ImportPaths = paths

	return p.atom.Header, nil
}

// Parse continues atom parsing from the state where Header method left.
// Resulting atom tree is complete and also contains info gathered by
// Header method
func (p *Parser) Parse() (ast.Atom, error) {
	return p.parse()
}

func (p *Parser) parse() (ast.Atom, error) {
	var nodes []ast.TopLevel
	for {
		if p.isEOF() {
			p.atom.Nodes = nodes
			return p.atom, nil
		}

		top, err := p.topLevel()
		if err != nil {
			return ast.Atom{}, err
		}
		nodes = append(nodes, top)
	}
}
