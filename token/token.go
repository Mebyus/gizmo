package token

import (
	"fmt"
	"strconv"

	"github.com/mebyus/gizmo/source"
)

type Token struct {
	Pos source.Pos

	// Not empty only for tokens which do not have static literal
	//
	// Examples: identifiers, strings, illegal tokens
	//
	// For tokens obtained from regular string literals (as in "some string")
	// this field contains original byte sequence from source file (without surrounding quotes).
	// That means escape sequences are not decoded and bytes from inside quotes in source file
	// are placed in this field "as is"
	Lit string

	// Meaning of this value is dependant on token Kind
	//
	//	Integer:	parsed integer value (if it fits into 64 bits)
	//	Rune:		integer value of code point
	//	String:		string raw byte size (as represented in memory, not in source code)
	//	EOF:		error code (can be 0)
	//	Illegal:	error code (always not 0)
	Val uint64

	Kind Kind
}

func (t Token) Pin() source.Pos {
	return t.Pos
}

func (t Token) Len() uint32 {
	if t.Kind.IsEOF() {
		return 0
	}
	if t.Kind.hasStaticLiteral() {
		return uint32(len(t.Kind.String()))
	}
	return uint32(len(t.Lit))
}

func (t Token) IsEOF() bool {
	return t.Kind.IsEOF()
}

func (t Token) IsLit() bool {
	return t.Kind.IsLit()
}

func (t Token) IsIdent() bool {
	return t.Kind.IsIdent()
}

func (t Token) IsLeftPar() bool {
	return t.Kind.IsLeftPar()
}

func (t Token) Literal() string {
	switch t.Kind {
	case Identifier:
		return t.Lit
	case Illegal:
		if t.Val == LengthOverflow {
			return "<overflow>"
		}
		return "=[ " + t.Lit + " ]="
	case BinaryInteger:
		return "0b" + strconv.FormatUint(t.Val, 2)
	case OctalInteger:
		return "0o" + strconv.FormatUint(t.Val, 8)
	case DecimalInteger:
		return strconv.FormatUint(t.Val, 10)
	case HexadecimalInteger:
		return "0x" + strconv.FormatUint(t.Val, 16)
	case DecimalFloat:
		return t.Lit
	case Rune:
		if t.Lit != "" {
			return "'" + t.Lit + "'"
		}
		switch t.Val {
		case '\\':
			return `'\\'`
		case '\'':
			return `'\''`
		case '\n':
			return `'\n'`
		case '\t':
			return `'\t'`
		case '\r':
			return `'\r'`
		}
		return "'" + string(rune(t.Val)) + "'"
	case String:
		return "\"" + t.Lit + "\""
	case Nil:
		return "nil"
	case True:
		return "true"
	case False:
		return "false"
	default:
		panic("must not be invoked with static literal tokens: " + t.Kind.String())
	}
}

func (t Token) Short() string {
	if t.Kind.hasStaticLiteral() {
		return fmt.Sprintf("%-6d%-12s%-12s%s", t.Pos.Num+1, t.Pos.Short(), ".", t.Kind.String())
	}

	return fmt.Sprintf("%-6d%-12s%-12s%s", t.Pos.Num+1, t.Pos.Short(), t.Kind.String(), t.Literal())
}

func (t Token) String() string {
	if t.Kind.hasStaticLiteral() {
		return fmt.Sprintf("%s%12s", t.Pos.String(), t.Kind.String())
	}

	return fmt.Sprintf("%s%-12s%s", t.Pos.String(), t.Kind.String(), t.Literal())
}
