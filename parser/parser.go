package parser

import (
	"io"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

type Parser struct {
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
	def := ast.NamespaceBlock{Default: true}
	for {
		if p.isEOF() {
			if len(def.Nodes) != 0 {
				blocks = append(blocks, def)
			}

			return ast.UnitAtom{
				Unit:   unit,
				Blocks: blocks,
			}, nil
		}

		if p.tok.Kind == token.Namespace {
			block, err := p.namespaceBlock()
			if err != nil {
				return ast.UnitAtom{}, err
			}
			if len(def.Nodes) != 0 {
				blocks = append(blocks, def)
				def = ast.NamespaceBlock{Default: true}
			}
			blocks = append(blocks, block)
		} else {
			top, err := p.topLevel()
			if err != nil {
				return ast.UnitAtom{}, err
			}
			def.Nodes = append(def.Nodes, top)
		}
	}
}
