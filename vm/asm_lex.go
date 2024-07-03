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
	Period
	LeftSquare
	RightSquare

	Fn
	Def
	Let

	noStaticLiteral

	SysReg
	Register
	Mnemonic
	Property
	Flag

	Illegal
	Label
	HexInteger
	Identifier
	String
)

const (
	SysRegIP = 0
	SysRegSP = 1
	SysRegFP = 2
	SysRegSC = 3
	SysRegCF = 4
)

var sysRegText = [...]string{
	SysRegIP: "ip",
	SysRegSP: "sp",
	SysRegFP: "fp",
	SysRegSC: "sc",
	SysRegCF: "cf",
}

const (
	PropPtr = 0
	PropLen = 1
)

var propText = [...]string{
	PropPtr: "ptr",
	PropLen: "len",
}

const (
	FlagZero    = 0
	FlagNotZero = 1
)

var flagText = [...]string{
	FlagZero:    "z",
	FlagNotZero: "nz",
}

func (k TokKind) hasStaticLiteral() bool {
	return k < noStaticLiteral
}

var tokKindLiteral = [...]string{
	emptyToken: "EMPTY",

	EOF: "EOF",

	Comma:       ",",
	Colon:       ":",
	Period:      ".",
	LeftSquare:  "[",
	RightSquare: "]",

	Fn:  "fn",
	Def: "def",
	Let: "let",

	SysReg:   "REG.SYS",
	Register: "REG",
	Mnemonic: "ME",
	Property: "PROP",
	Flag:     "FLAG",

	Illegal:    "ILG",
	Label:      "LABEL",
	HexInteger: "INT.HEX",
	Identifier: "IDN",
	String:     "STR",
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
	MeJump
	MeTest
	MeInc
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
	MeJump:    "jump",
	MeTest:    "test",
	MeInc:     "inc",
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
	"jump":    MeJump,
	"test":    MeTest,
	"inc":     MeInc,
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
	case SysReg:
		return sysRegText[t.Val]
	case Register:
		return "r" + strconv.FormatUint(t.Val, 10)
	case Property:
		return propText[t.Val]
	case Mnemonic:
		return meText[t.Val]
	case Flag:
		return "?." + flagText[t.Val]
	case HexInteger:
		return "0x" + strconv.FormatUint(t.Val, 16)
	case Identifier:
		return t.Lit
	case String:
		return "\"" + t.Lit + "\""
	case Label:
		return "@." + t.Lit
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

	lx.SkipWhitespaceAndComments()
	if lx.EOF {
		return lx.create(EOF)
	}

	if lx.C == '"' {
		return lx.stringLiteral()
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
		tok.Kind = String
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

	tok.Kind = String
	tok.Lit = lit
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
	case "let":
		tok.Kind = Let
		return tok
	case "ptr":
		tok.Kind = Property
		tok.Val = PropPtr
		return tok
	case "len":
		tok.Kind = Property
		tok.Val = PropLen
		return tok
	case "ip":
		tok.Kind = SysReg
		tok.Val = SysRegIP
		return tok
	case "sp":
		tok.Kind = SysReg
		tok.Val = SysRegSP
		return tok
	case "fp":
		tok.Kind = SysReg
		tok.Val = SysRegFP
		return tok
	case "sc":
		tok.Kind = SysReg
		tok.Val = SysRegSC
		return tok
	case "cf":
		tok.Kind = SysReg
		tok.Val = SysRegCF
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
		tok.Kind = Illegal
		return tok
	}

	tok.Kind = Flag
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
	case '.':
		return lx.oneByteToken(Period)
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
