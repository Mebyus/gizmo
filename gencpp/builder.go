package gencpp

import (
	"io"
)

// Builder is used to generate formatted output (C++ code)
//
// Output is stored in memory
type Builder struct {
	// Stored output
	out []byte

	// Identation buffer
	//
	// Stores sequence of bytes which is used for indenting current line
	// in output. When a new line starts this buffer is used to add identation
	ibuf []byte

	// Output reader position
	pos int

	// list of scopes in current namespace, equals nil inside default namespace
	currentScopes []string

	cfg *Config
}

type Config struct {
	DefaultNamespace string

	GlobalNamespacePrefix string

	// Initial preallocated memory capacity for storing output
	Size int

	SourceLocationComments bool
}

// Explicit interface implementation check
var _ io.Reader = &Builder{}

// NewBuilder create a ready to use builder with specified initial
// preallocated memory capacity
func NewBuilder(cfg *Config) *Builder {
	return &Builder{
		out: make([]byte, 0, cfg.Size),
		cfg: cfg,
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

// Reset builder to a clean state, but keep underlying memory
// for further usage
func (g *Builder) Reset() {
	g.out = g.out[:0]
	g.ibuf = g.ibuf[:0]
	g.pos = 0
}

// write a string directly to output without formatting
func (g *Builder) write(s string) {
	g.out = append(g.out, []byte(s)...)
}

// writes a single byte to output
func (g *Builder) wb(b byte) {
	g.out = append(g.out, b)
}

// increment indentation by one level
func (g *Builder) inc() {
	g.ibuf = append(g.ibuf, '\t')
}

// decrement indentation by one level
func (g *Builder) dec() {
	g.ibuf = g.ibuf[:len(g.ibuf)-1]
}

func (g *Builder) comment(s string) {
	g.write("// ")
	g.write(s)
	g.nl()
}

func (g *Builder) space() {
	g.wb(' ')
}

func (g *Builder) semi() {
	g.wb(';')
}

// start new line
func (g *Builder) nl() {
	g.wb('\n')
}

// add indentation to current line
func (g *Builder) indent() {
	g.out = append(g.out, g.ibuf...)
}
