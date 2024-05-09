package format

import (
	"fmt"
	"io"

	"github.com/mebyus/gizmo/token"
)

type Printer struct {
	buf Buffer

	tokens []token.Token

	nodes []Node
}

func NewPrinter(tokens []token.Token, nodes []Node) *Printer {
	return &Printer{
		tokens: tokens,
		nodes:  nodes,
	}
}

func Print(tokens []token.Token, nodes []Node) []byte {
	p := NewPrinter(tokens, nodes)
	return p.print()
}

func (p *Printer) print() []byte {
	for i := 0; i < len(p.nodes); i += 1 {
		node := p.nodes[i]

		switch node.Kind {
		case GenNode:
			p.gen(token.Kind(node.Val))
		case TokNode:
			p.tok(node.Val)
		case StrictSpaceNode:
			p.ss()
		case StartNode:
			p.start()
		case IncNode:
			p.inc()
		case DecNode:
			p.dec()
		case StartBlockNode:
			p.sb()
		case EndBlockNode:
			p.eb()
		case SepNode:
			// TODO: logic for non default SepNode behaviour
		case TrailCommaNode:
			// TODO: logic for non default TrailCommaNode behaviour
		case SpaceNode:
			p.space()
		case BlankNode:
			p.blank()
		case NewlineNode:
			p.nl()
		case NewlineIndentNode:
			p.nli()
		default:
			panic(fmt.Sprintf("unknown node %s", node.Kind.String()))
		}
	}

	return p.buf.Bytes()
}

func (p *Printer) gen(kind token.Kind) {
	p.buf.write(kind.String())
}

func (p *Printer) tok(num uint32) {
	tok := p.tokens[num]

	if tok.Kind == token.LineComment {
		p.lcom(tok.Lit)
		return
	}

	p.buf.write(tok.Literal())
}

func (p *Printer) inc() {
	p.buf.inc()
}

func (p *Printer) dec() {
	p.buf.dec()
}

func (p *Printer) space() {
	// TODO: logic for non default SpaceNode behaviour
	p.ss()
}

func (p *Printer) lcom(s string) {
	p.buf.write("// ")
	p.buf.write(s)
	p.buf.nl()
}

func (p *Printer) blank() {
	p.buf.nl()
	p.buf.nl()
}

func (p *Printer) nl() {
	p.buf.nl()
}

func (p *Printer) nli() {
	p.buf.nl()
	p.buf.indent()
}

func (p *Printer) ss() {
	p.buf.ss()
}

func (p *Printer) start() {
	p.buf.nl()
	p.buf.indent()
}

func (p *Printer) sb() {
	p.buf.put('{')
	p.buf.inc()
}

func (p *Printer) eb() {
	p.buf.dec()
	p.buf.nl()
	p.buf.indent()
	p.buf.put('}')
}

// Buffer is used to accumulate formatted source code output.
//
// Output is stored in memory.
type Buffer struct {
	// Stored output.
	out []byte

	// Indentation buffer.
	//
	// Stores sequence of bytes which is used for indenting current line
	// in output. When a new line starts this buffer is used to add identation.
	ib []byte

	// Output reader position.
	pos int
}

// Explicit interface implementation check.
var _ io.Reader = &Buffer{}
var _ io.WriterTo = &Buffer{}

// Bytes returns stored output in form of raw bytes.
func (g *Buffer) Bytes() []byte {
	return g.out
}

// Read implements io.Reader from stored output.
func (g *Buffer) Read(p []byte) (int, error) {
	n := copy(p, g.out[g.pos:])

	if n == 0 {
		return 0, io.EOF
	}

	g.pos += n

	return n, nil
}

// WriteTo implements io.WriterTo.
func (g *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(g.Bytes())
	return int64(n), err
}

// Reset buffer to a clean state, but keep underlying memory.
// for further usage.
func (g *Buffer) Reset() {
	g.out = g.out[:0]
	g.ib = g.ib[:0]
	g.pos = 0
}

// write a string directly to output without additional formatting.
func (g *Buffer) write(s string) {
	g.out = append(g.out, []byte(s)...)
}

// write a single byte to output.
func (g *Buffer) put(b byte) {
	g.out = append(g.out, b)
}

// increment indentation by one level.
func (g *Buffer) inc() {
	g.ib = append(g.ib, '\t')
}

// decrement indentation by one level
func (g *Buffer) dec() {
	g.ib = g.ib[:len(g.ib)-1]
}

// write space character to output.
func (g *Buffer) ss() {
	g.put(' ')
}

// start new line in output.
func (g *Buffer) nl() {
	g.put('\n')
}

// add indentation to current line.
func (g *Buffer) indent() {
	g.out = append(g.out, g.ib...)
}
