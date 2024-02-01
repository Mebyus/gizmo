package lexer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/mebyus/gizmo/token"
)

type Stream interface {
	Lex() token.Token
}

// Parrot implements Stream by yielding tokens from supplied list
type Parrot struct {
	toks []token.Token
	i    int
}

func FromTokens(toks []token.Token) *Parrot {
	return &Parrot{
		toks: toks,
	}
}

func (p *Parrot) Lex() token.Token {
	if p.i >= len(p.toks) {
		tok := token.Token{Kind: token.EOF}
		if len(p.toks) == 0 {
			return tok
		}
		pos := p.toks[len(p.toks)-1].Pos
		pos.Line++
		pos.Col = 0
		tok.Pos = pos
		return tok
	}
	tok := p.toks[p.i]
	p.i++
	return tok
}

// Eraser wraps a Stream removing position information from tokens
//
// Commonly used for testing
type Eraser struct {
	s Stream
}

func NoPos(s Stream) Eraser {
	return Eraser{s: s}
}

func (e Eraser) Lex() token.Token {
	tok := e.s.Lex()
	tok.Pos.Clear()
	return tok
}

func List(w io.Writer, s Stream) error {
	for {
		tok := s.Lex()
		err := putToken(w, tok)
		if err != nil {
			return err
		}
		err = putNewline(w)
		if err != nil {
			return err
		}
		if tok.IsEOF() {
			return nil
		}
	}
}

func putToken(w io.Writer, tok token.Token) error {
	_, err := w.Write([]byte(tok.Short()))
	return err
}

func putNewline(w io.Writer) error {
	_, err := w.Write([]byte("\n"))
	return err
}

func ParseList(r io.Reader) ([]token.Token, error) {
	sc := bufio.NewScanner(r)
	var toks []token.Token
	i := 0
	for sc.Scan() {
		i++
		line := sc.Text()
		if line == "" {
			continue
		}
		tok, err := token.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("line %d, parse token: %w", i, err)
		}
		toks = append(toks, tok)
	}
	err := sc.Err()
	if err != nil {
		return nil, err
	}
	return toks, nil
}
