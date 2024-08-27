package genc

import (
	"strconv"

	"github.com/mebyus/gizmo/stg"
)

type Builder struct {
	// Mangled unit name. Index corresponds to Unit.Index.
	unames []string

	buf []byte

	// Indentation buffer.
	//
	// Stores sequence of bytes which is used for indenting current line
	// in output. When a new line starts this buffer is used to add identation.
	ib []byte

	// global name prefix for generated symbols.
	prefix string

	// cached type specs.
	specs map[*stg.Type]string

	// from type index.
	arrays map[*stg.Type][]*stg.Type

	// from type index.
	chunks map[*stg.Type]*stg.Type
}

func (g *Builder) Bytes() []byte {
	return g.buf
}

func (g *Builder) LineComment(s string) {
	g.puts("// ")
	g.puts(s)
	g.nl()
}

func (g *Builder) BlockTitle(unit string, s string) {
	g.puts("/* ===== ")
	g.puts(unit)
	g.puts(": ")
	g.puts(s)
	g.puts(" ===== */")
	g.nl()
}

// put decimal formatted integer into output buffer
func (g *Builder) putn(n uint64) {
	g.puts(strconv.FormatUint(n, 10))
}

// put string into output buffer
func (g *Builder) puts(s string) {
	g.buf = append(g.buf, s...)
}

// put single byte into output buffer
func (g *Builder) putb(b byte) {
	g.buf = append(g.buf, b)
}

func (g *Builder) put(b []byte) {
	g.buf = append(g.buf, b...)
}

func (g *Builder) nl() {
	g.putb('\n')
}

func (g *Builder) space() {
	g.putb(' ')
}

func (g *Builder) semi() {
	g.putb(';')
}

// increment indentation by one level.
func (g *Builder) inc() {
	g.ib = append(g.ib, '\t')
}

// decrement indentation by one level.
func (g *Builder) dec() {
	g.ib = g.ib[:len(g.ib)-1]
}

// add indentation to current line.
func (g *Builder) indent() {
	g.put(g.ib)
}
