package format

import (
	"fmt"
	"io"

	"github.com/mebyus/gizmo/token"
)

type Stapler struct {
	buf Buffer

	prev token.Token
	next token.Token

	tokens []token.Token

	nodes []Node
}

func NewStapler(tokens []token.Token, nodes []Node) *Stapler {
	return &Stapler{
		tokens: tokens,
		nodes:  nodes,
	}
}

func (s *Stapler) add(node Node) {
	s.nodes = append(s.nodes, node)
}

func (s *Stapler) staple() []byte {
	for i := 0; i < len(s.nodes); i += 1 {
		node := s.nodes[i]

		switch node.Kind {
		case GenNode:
			s.gen(token.Kind(node.Val))
		case TokNode:
			s.tok(node.Val)
		case StrictSpaceNode:
			s.ss()
		case StartNode:
			s.start()
		case IncNode:
			s.inc()
		case DecNode:
			s.dec()
		case StartBlockNode:
			s.sb()
		case EndBlockNode:
			s.eb()
		case SepNode:
			// TODO: logic for non default SepNode behaviour
		case TrailCommaNode:
			// TODO: logic for non default TrailCommaNode behaviour
		case SpaceNode:
			s.space()
		case BlankNode:
			s.blank()
		default:
			panic(fmt.Sprintf("unknown node %s", node.Kind.String()))
		}
	}

	return s.buf.Bytes()
}

func (s *Stapler) gen(kind token.Kind) {
	s.buf.write(kind.String())
}

func (s *Stapler) tok(num uint32) {
	tok := s.tokens[num]
	s.buf.write(tok.Literal())
}

func (s *Stapler) inc() {
	s.buf.inc()
}

func (s *Stapler) dec() {
	s.buf.dec()
}

func (s *Stapler) space() {
	// TODO: logic for non default SpaceNode behaviour
	s.ss()
}

func (s *Stapler) blank() {
	s.buf.nl()
	s.buf.nl()
}

func (s *Stapler) ss() {
	s.buf.ss()
}

func (s *Stapler) start() {
	s.buf.nl()
	s.buf.indent()
}

func (s *Stapler) sb() {
	s.buf.put('{')
	s.buf.inc()
}

func (s *Stapler) eb() {
	s.buf.dec()
	s.buf.nl()
	s.buf.indent()
	s.buf.put('}')
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
