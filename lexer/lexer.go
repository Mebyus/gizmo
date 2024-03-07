package lexer

import (
	"io"

	"github.com/mebyus/gizmo/source"
)

type Lexer struct {
	// source text which is scanned by lexer
	src []byte

	// Lexer position inside source text. Directly used to
	// create tokens
	pos source.Pos

	// Mark index
	//
	// Mark is used to slice input text for token literals
	mark int

	// next byte read index
	i int

	// scanning index of source text (c = src[s])
	s int

	// Byte that was previously at scan position
	prev byte

	// Byte at current scan position
	//
	// This is a cached value. Look into advance() method
	// for details about how this caching algorithm works
	c byte

	// Next byte that will be placed at scan position
	//
	// This is a cached value from previous read
	next byte

	// True if lexer reached end of input (by scan index)
	eof bool
}

func FromBytes(b []byte) *Lexer {
	lx := &Lexer{src: b}
	lx.init()
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
