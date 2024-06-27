package vm

import "strconv"

type TokKind uint8

const (
	emptyToken TokKind = iota

	EOF

	Comma
	Colon
	Semicolon
	LeftSquare
	RightSquare

	Register
	Mnemonic

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

type Lexer struct {
	// Source text which is scanned by lexer.
	src []byte

	// Mark index
	//
	// Mark is used to slice input text for token literals
	mark int

	// next byte read index
	i int

	// scanning index of source text (c = src[s])
	s int

	line uint32
	col  uint32

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

func (lx *Lexer) Lex() *Token {
	if lx.eof {
		return lx.create(EOF)
	}

	lx.skipWhitespace()
	if lx.eof {
		return lx.create(EOF)
	}

	if lx.c == '@' && lx.next == '.' {
		return lx.label()
	}

	return lx.other()
}

func (lx *Lexer) label() *Token {

}

func (lx *Lexer) other() *Token {

}

func (lx *Lexer) advance() {

}

func (lx *Lexer) skipWhitespace() {

}

func (lx *Lexer) skipLine() {
	for !lx.eof && lx.c != '\n' {
		lx.advance()
	}
	if lx.c == '\n' {
		lx.advance()
	}
}

func (lx *Lexer) create(kind TokKind) *Token {
	return &Token{
		Line: lx.line,
		Col:  lx.col,
		Kind: kind,
	}
}
