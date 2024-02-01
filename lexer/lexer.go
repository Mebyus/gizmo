package lexer

import (
	"io"

	"github.com/mebyus/gizmo/source"
)

type Lexer struct {
	// buffer for next literal
	buf [maxLiteralLength]byte

	// source text which is scanned by Lexer
	src []byte

	// Lexer position inside source text
	pos source.Pos

	// literal length in buffer
	len int

	// scanning index of source text
	i int

	// previous character code
	prev int

	// character code at current Lexer position
	c int

	// next character code
	next int
}

func FromBytes(b []byte) *Lexer {
	lx := &Lexer{src: b}

	// init lexer buffer
	lx.advance()
	lx.advance()

	lx.pos.Reset()
	return lx
}

func FromString(s string) *Lexer {
	return FromBytes([]byte(s))
}

func FromFile(filename string) (*Lexer, error) {
	src, err := source.Load(filename)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return FromSource(src), nil
}

func FromSource(src *source.File) *Lexer {
	lx := FromBytes(src.Bytes)
	lx.pos.File = src
	return lx
}

func FromReader(r io.Reader) (*Lexer, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromBytes(b), nil
}
