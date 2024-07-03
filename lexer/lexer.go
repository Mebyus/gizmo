package lexer

import (
	"io"

	"github.com/mebyus/gizmo/char"
	"github.com/mebyus/gizmo/source"
)

type Lexer struct {
	char.Chopper

	file *source.File

	// current token number
	num uint32
}

func FromBytes(b []byte) *Lexer {
	lx := Lexer{}
	lx.Init(b)
	return &lx
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
	lx.file = src
	return lx
}

func FromReader(r io.Reader) (*Lexer, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromBytes(b), nil
}

func (lx *Lexer) Stats() *char.Stats {
	return &char.Stats{
		Lines:     lx.Line + 1, // counter if zero-based
		HardLines: lx.HardLines,
		Tokens:    lx.num - 1, // exclude EOF token
		Size:      uint32(lx.Pos),
	}
}
