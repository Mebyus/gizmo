package vm

import (
	"fmt"
	"io"
	"strconv"

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

	Fn
	Def

	noStaticLiteral

	Register
	Mnemonic

	Illegal
	Label
	HexInteger
	Identifier
)

func (k TokKind) hasStaticLiteral() bool {
	return k < noStaticLiteral
}

var tokKindLiteral = [...]string{
	emptyToken: "EMPTY",

	EOF: "EOF",

	Comma:       ",",
	Colon:       ":",
	Semicolon:   ";",
	LeftSquare:  "[",
	RightSquare: "]",

	Fn:  "fn",
	Def: "def",

	Register: "REG",
	Mnemonic: "ME",

	Illegal:    "ILG",
	Label:      "LABEL",
	HexInteger: "INT.HEX",
	Identifier: "IDN",
}

func (k TokKind) String() string {
	return tokKindLiteral[k]
}

// Mnemonic values
const (
	MeNop = iota
	MeHalt
	MeSysCall
	MeTrap
	MeClear
	MeLoad
	MeStore
	MeAdd
)

var meText = [...]string{
	MeNop:     "nop",
	MeHalt:    "halt",
	MeSysCall: "syscall",
	MeTrap:    "trap",
	MeClear:   "clear",
	MeLoad:    "load",
	MeStore:   "store",
	MeAdd:     "add",
}

var meWord = map[string]int{
	"nop":     MeNop,
	"halt":    MeHalt,
	"syscall": MeSysCall,
	"trap":    MeTrap,
	"clear":   MeClear,
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

func (t *Token) Short() string {
	if t.Kind.hasStaticLiteral() {
		return fmt.Sprintf("%-12s%-12s%s", t.Pos(), ".", t.Kind.String())
	}

	return fmt.Sprintf("%-12s%-12s%s", t.Pos(), t.Kind.String(), t.Literal())
}

func (t Token) Literal() string {
	if t.Kind.hasStaticLiteral() {
		return t.Kind.String()
	}

	switch t.Kind {
	case Register:
		return "r" + strconv.FormatUint(t.Val, 10)
	case Mnemonic:
		return meText[t.Val]
	case HexInteger:
		return "0x" + strconv.FormatUint(t.Val, 16)
	case Identifier:
		return t.Lit
	case Label:
		return "." + t.Lit
	case Illegal:
		return "."
	default:
		panic("unreachable")
	}
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
	}
}

func RenderToken(w io.Writer, tok *Token) error {
	_, err := io.WriteString(w, tok.Short())
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

	if lx.C == ';' {
		return lx.semiSkipLine()
	}

	if lx.C == 'r' && char.IsDecDigit(lx.Next) {
		return lx.reg()
	}

	if char.IsLetter(lx.C) {
		return lx.word()
	}

	if lx.C == '0' && lx.Next == 'x' {
		return lx.hexNumber()
	}

	if lx.C == '.' && char.IsLetter(lx.Next) {
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

func (lx *Lexer) semiSkipLine() *Token {
	tok := lx.tok()

	lx.Advance() // skip ";"
	lx.SkipLine()

	tok.Kind = Semicolon
	return tok
}

func (lx *Lexer) reg() *Token {
	tok := lx.tok()

	lx.Advance() // skip "r"

	lx.Start()
	lx.SkipDecDigits()

	tok.Kind = Register
	tok.Val = char.ParseDecDigits(lx.View())
	return tok
}

func (lx *Lexer) hexNumber() *Token {
	tok := lx.tok()

	lx.Advance() // skip "0"
	lx.Advance() // skip "x"

	lx.Start()
	lx.SkipHexDigits()

	tok.Kind = HexInteger
	tok.Val = char.ParseHexDigits(lx.View())
	return tok
}

func (lx *Lexer) word() *Token {
	tok := lx.tok()

	lx.Start()
	lx.SkipWord()

	lit, ok := lx.Take()
	if !ok {
		panic("not implemented")
	}

	switch lit {
	case "fn":
		tok.Kind = Fn
		return tok
	case "def":
		tok.Kind = Def
		return tok
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
	tok := lx.tok()

	lx.Advance() // skip "."

	lx.Start()
	lx.SkipWord()

	lit, ok := lx.Take()
	if !ok {
		panic("not implemented")
	}

	tok.Kind = Label
	tok.Lit = lit
	return tok
}

func (lx *Lexer) other() *Token {
	switch lx.C {
	case '[':
		return lx.oneByteToken(LeftSquare)
	case ']':
		return lx.oneByteToken(RightSquare)
	case ',':
		return lx.oneByteToken(Comma)
	case ':':
		return lx.oneByteToken(Colon)
	default:
		return lx.illegalByteToken()
	}
}

func (lx *Lexer) oneByteToken(k TokKind) *Token {
	tok := lx.create(k)
	lx.Advance()
	return tok
}

func (lx *Lexer) illegalByteToken() *Token {
	tok := lx.tok()
	tok.Kind = Illegal
	lx.Advance()
	return tok
}

func (lx *Lexer) create(kind TokKind) *Token {
	tok := lx.tok()
	tok.Kind = kind
	return tok
}
