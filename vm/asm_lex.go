package vm

import (
	"io"
	"strconv"
	"time"

	"github.com/mebyus/gizmo/char"
)

type TokKind uint8

const (
	emptyToken TokKind = iota

	EOF

	Comma
	Colon
	Semicolon
	LeftSquare
	RightSquare

	Function

	Register
	Mnemonic

	Illegal
	Label
	HexInteger
	Identifier
)

// Mnemonic values
const (
	MeNop = iota
	MeHalt
	MeSysCall
	MeLoad
	MeStore
	MeAdd
)

var meText = [...]string{
	MeNop:     "nop",
	MeHalt:    "halt",
	MeSysCall: "syscall",
	MeLoad:    "load",
	MeStore:   "store",
	MeAdd:     "add",
}

var meWord = map[string]int{
	"nop":     MeNop,
	"halt":    MeHalt,
	"syscall": MeSysCall,
	"load":    MeLoad,
	"store":   MeStore,
	"add":     MeAdd,
}

type Token struct {
	Lit string
	Val uint64

	Line uint32
	Col  uint32

	Kind TokKind
}

func (t *Token) Pos() string {
	return strconv.FormatUint(uint64(t.Line+1), 10) + ":" + strconv.FormatUint(uint64(t.Col+1), 10)
}

func ListTokens(w io.Writer, lx *Lexer) error {
	for {
		tok := lx.Lex()
		err := RenderToken(w, tok)
		if err != nil {
			return err
		}
		err = putNewline(w)
		if err != nil {
			return err
		}
		if tok.Kind == EOF {
			return nil
		}
		// TODO: remove this debug sleep
		time.Sleep(100 * time.Millisecond)
	}
}

func RenderToken(w io.Writer, tok *Token) error {
	_, err := io.WriteString(w, tok.Pos())
	return err
}

func putNewline(w io.Writer) error {
	_, err := io.WriteString(w, "\n")
	return err
}

type Lexer struct {
	char.Chopper
}

func NewLexer(text []byte) *Lexer {
	lx := Lexer{}
	lx.Init(text)
	return &lx
}

func (lx *Lexer) Lex() *Token {
	if lx.EOF {
		return lx.create(EOF)
	}

	lx.SkipWhitespace()
	if lx.EOF {
		return lx.create(EOF)
	}

	if char.IsLetter(lx.C) {
		return lx.word()
	}

	if lx.C == '.' {
		return lx.label()
	}

	return lx.other()
}

// start new token at current position
func (lx *Lexer) tok() *Token {
	return &Token{
		Line: lx.Line,
		Col:  lx.Col,
	}
}

func (lx *Lexer) word() *Token {
	tok := lx.tok()

	lx.Start()
	lx.SkipWord()

	lit, ok := lx.Take()
	if !ok {
		panic("not implemented")
	}
	me, ok := meWord[lit]
	if ok {
		tok.Kind = Mnemonic
		tok.Val = uint64(me)
		return tok
	}

	tok.Kind = Identifier
	tok.Lit = lit
	return tok
}

func (lx *Lexer) label() *Token {
	return &Token{}
}

func (lx *Lexer) other() *Token {
	return &Token{}
}

func (lx *Lexer) create(kind TokKind) *Token {
	tok := lx.tok()
	tok.Kind = kind
	return tok
}
