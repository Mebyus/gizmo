package source

import (
	"io"
)

// Builder is used to generate formatted source code output.
//
// Output is stored in memory.
type Builder struct {
	// Stored output.
	out []byte

	// Identation buffer
	//
	// Stores sequence of bytes which is used for indenting current line
	// in output. When a new line starts this buffer is used to add identation.
	ibuf []byte

	// Output reader position.
	pos int
}

// Explicit interface implementation check
var _ io.Reader = &Builder{}
var _ io.WriterTo = &Builder{}

// NewBuilder create a ready to use builder with specified initial
// preallocated memory capacity
func NewBuilder(size int) *Builder {
	return &Builder{
		out: make([]byte, 0, size),
	}
}

// Bytes returns stored output in form of raw bytes
func (g *Builder) Bytes() []byte {
	return g.out
}

// Read implements io.Reader from stored output
func (g *Builder) Read(p []byte) (int, error) {
	n := copy(p, g.out[g.pos:])

	if n == 0 {
		return 0, io.EOF
	}

	g.pos += n

	return n, nil
}

func (g *Builder) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(g.Bytes())
	return int64(n), err
}

// Reset builder to a clean state, but keep underlying memory.
// for further usage
func (g *Builder) Reset() {
	g.out = g.out[:0]
	g.ibuf = g.ibuf[:0]
	g.pos = 0
}

// Str a string directly to output without formatting.
func (g *Builder) Str(s string) {
	g.out = append(g.out, []byte(s)...)
}

// Byte write a single byte to output.
func (g *Builder) Byte(b byte) {
	g.out = append(g.out, b)
}

// Inc increment indentation by one level.
func (g *Builder) Inc() {
	g.ibuf = append(g.ibuf, '\t')
}

// Dec decrement indentation by one level
func (g *Builder) Dec() {
	g.ibuf = g.ibuf[:len(g.ibuf)-1]
}

// LineComment add line comment and start new line.
func (g *Builder) LineComment(s string) {
	g.Str("// ")
	g.Str(s)
	g.Line()
}

// Space write single space character to output.
func (g *Builder) Space() {
	g.Byte(' ')
}

// Semi write single semicolon character to output.
func (g *Builder) Semi() {
	g.Byte(';')
}

// Line start new line in output.
func (g *Builder) Line() {
	g.Byte('\n')
}

// Indent add indentation to current line.
func (g *Builder) Indent() {
	g.out = append(g.out, g.ibuf...)
}
