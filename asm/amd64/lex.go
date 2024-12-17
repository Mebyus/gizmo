package amd64

import (
	"fmt"
	"io"
	"strconv"

	"github.com/mebyus/gizmo/char"
)

type TokKind uint8

const (
	TokNil TokKind = iota

	TokEOF

	TokComma
	TokColon
	TokPeriod
	LeftSquare
	RightSquare

	TokFun
	TokDef
	TokLet

	tokNoStaticLiteral

	TokRegister
	TokMnemonic
	TokProperty
	TokFlag

	TokLabel
	TokHexInteger
	TokIdentifier
	TokString
	TokIllegal
)

func (k TokKind) hasStaticLiteral() bool {
	return k < tokNoStaticLiteral
}

var tokKindLiteral = [...]string{
	TokNil: "EMPTY",

	TokEOF: "EOF",

	TokComma:    ",",
	TokColon:    ":",
	TokPeriod:   ".",
	LeftSquare:  "[",
	RightSquare: "]",

	TokFun: "fun",
	TokDef: "def",
	TokLet: "let",

	TokRegister: "REG",
	TokMnemonic: "ME",
	TokProperty: "PROP",
	TokFlag:     "FLAG",

	TokIllegal:    "ILG",
	TokLabel:      "LABEL",
	TokHexInteger: "INT.HEX",
	TokIdentifier: "IDN",
	TokString:     "STR",
}

func (k TokKind) String() string {
	return tokKindLiteral[k]
}

type Token struct {
	Text string
	Val  uint64

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
	case TokRegister:
		return Register(t.Val).String()
	case TokProperty:
		return propText[t.Val]
	case TokMnemonic:
		return meText[t.Val]
	case TokFlag:
		return "?." + flagText[t.Val]
	case TokHexInteger:
		return "0x" + strconv.FormatUint(t.Val, 16)
	case TokIdentifier:
		return "$." + t.Text
	case TokString:
		return "\"" + t.Text + "\""
	case TokLabel:
		return "@." + t.Text
	case TokIllegal:
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
		if tok.Kind == TokEOF {
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
		return lx.create(TokEOF)
	}

	lx.SkipWhitespaceAndComments()
	if lx.EOF {
		return lx.create(TokEOF)
	}

	if lx.C == '"' {
		return lx.stringLiteral()
	}

	if char.IsLetter(lx.C) {
		return lx.word()
	}

	if lx.C == '0' && lx.Next == 'x' {
		return lx.hexNumber()
	}

	if lx.C == '?' && lx.Next == '.' {
		return lx.flag()
	}

	if lx.C == '@' && lx.Next == '.' {
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

func (lx *Lexer) stringLiteral() *Token {
	tok := lx.tok()

	lx.Advance() // skip '"'

	if lx.C == '"' {
		// common case of empty string literal
		lx.Advance()
		tok.Kind = TokString
		return tok
	}

	lx.Start()
	lx.SkipString()
	lit, ok := lx.Take()
	if !ok {
		panic("not implemented")
	}

	// TODO: handle EOF edgecases
	lx.Advance() // skip quote

	tok.Kind = TokString
	tok.Text = lit
	return tok
}

func (lx *Lexer) hexNumber() *Token {
	tok := lx.tok()

	lx.Advance() // skip "0"
	lx.Advance() // skip "x"

	lx.Start()
	lx.SkipHexDigits()

	tok.Kind = TokHexInteger
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
	case "fun":
		tok.Kind = TokFun
	case "def":
		tok.Kind = TokDef
	case "let":
		tok.Kind = TokLet
	case "ptr":
		tok.Kind = TokProperty
		tok.Val = PropPtr
	case "len":
		tok.Kind = TokProperty
		tok.Val = PropLen
	case "rax":
		tok.Kind = TokRegister
		tok.Val = uint64(RAX)
	case "rcx":
		tok.Kind = TokRegister
		tok.Val = uint64(RCX)
	case "rdx":
		tok.Kind = TokRegister
		tok.Val = uint64(RDX)
	case "rbx":
		tok.Kind = TokRegister
		tok.Val = uint64(RBX)
	case "rsp":
		tok.Kind = TokRegister
		tok.Val = uint64(RSP)
	case "rbp":
		tok.Kind = TokRegister
		tok.Val = uint64(RBP)
	case "rsi":
		tok.Kind = TokRegister
		tok.Val = uint64(RSI)
	case "rdi":
		tok.Kind = TokRegister
		tok.Val = uint64(RDI)
	case "r8":
		tok.Kind = TokRegister
		tok.Val = uint64(R8)
	case "r9":
		tok.Kind = TokRegister
		tok.Val = uint64(R9)
	case "r10":
		tok.Kind = TokRegister
		tok.Val = uint64(R10)
	case "r11":
		tok.Kind = TokRegister
		tok.Val = uint64(R11)
	case "r12":
		tok.Kind = TokRegister
		tok.Val = uint64(R12)
	case "r13":
		tok.Kind = TokRegister
		tok.Val = uint64(R13)
	case "r14":
		tok.Kind = TokRegister
		tok.Val = uint64(R14)
	case "r15":
		tok.Kind = TokRegister
		tok.Val = uint64(R15)
	default:
		me, ok := meLookup[lit]
		if ok {
			tok.Kind = TokMnemonic
			tok.Val = uint64(me)
		} else {
			tok.Kind = TokIllegal
			tok.Text = lit
		}
	}

	return tok
}

func (lx *Lexer) flag() *Token {
	tok := lx.tok()

	lx.Advance() // skip "?"
	lx.Advance() // skip "."

	lx.Start()
	lx.SkipWord()

	lit, ok := lx.Take()
	if !ok {
		panic("not implemented")
	}

	switch lit {
	case "z":
		tok.Val = FlagZero
	case "nz":
		tok.Val = FlagNotZero
	default:
		tok.Kind = TokIllegal
		return tok
	}

	tok.Kind = TokFlag
	return tok
}

func (lx *Lexer) label() *Token {
	tok := lx.tok()

	lx.Advance() // skip "@"
	lx.Advance() // skip "."

	lx.Start()
	lx.SkipWord()

	lit, ok := lx.Take()
	if !ok {
		panic("not implemented")
	}

	tok.Kind = TokLabel
	tok.Text = lit
	return tok
}

func (lx *Lexer) other() *Token {
	switch lx.C {
	case '[':
		return lx.oneByteToken(LeftSquare)
	case ']':
		return lx.oneByteToken(RightSquare)
	case '.':
		return lx.oneByteToken(TokPeriod)
	case ',':
		return lx.oneByteToken(TokComma)
	case ':':
		return lx.oneByteToken(TokColon)
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
	tok.Kind = TokIllegal
	lx.Advance()
	return tok
}

func (lx *Lexer) create(kind TokKind) *Token {
	tok := lx.tok()
	tok.Kind = kind
	return tok
}
