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
	// Examples: identifiers, illegal tokens
	Lit string

	// Meaning of this value is dependant on token Kind
	//
	//	Integer:	parsed integer value
	//	Character:	integer value of code point
	//	EOF:		error code
	//	Illegal:	error code
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
	case BinaryInteger:
		return "0b" + strconv.FormatUint(t.Val, 2)
	case OctalInteger:
		return "0o" + strconv.FormatUint(t.Val, 8)
	case DecimalInteger:
		return strconv.FormatUint(t.Val, 10)
	case HexadecimalInteger:
		return "0x" + strconv.FormatUint(t.Val, 16)
	case Character:
		return fmt.Sprintf("'%c'", t.Val)
	case String:
		return "\"" + t.Lit + "\""
	case Nil:
		return "nil"
	default:
		panic("must not be invoked with static literal tokens: " + t.Kind.String())
	}
}

func (t Token) Short() string {
	if t.Kind.hasStaticLiteral() {
		return fmt.Sprintf("%-12s%s", t.Pos.Short(), t.Kind.String())
	}

	switch t.Kind {
	case BinaryInteger:
		return fmt.Sprintf("%-12s%-12s0b%b", t.Pos.Short(), t.Kind.String(), t.Val)
	case OctalInteger:
		return fmt.Sprintf("%-12s%-12s0o%o", t.Pos.Short(), t.Kind.String(), t.Val)
	case DecimalInteger:
		return fmt.Sprintf("%-12s%-12s%d", t.Pos.Short(), t.Kind.String(), t.Val)
	case HexadecimalInteger:
		return fmt.Sprintf("%-12s%-12s0x%X", t.Pos.Short(), t.Kind.String(), t.Val)
	case Character:
		return fmt.Sprintf("%-12s%-12s'%c'", t.Pos.Short(), t.Kind.String(), t.Val)
	case String:
		return fmt.Sprintf("%-12s%-12s\"%s\"", t.Pos.Short(), t.Kind.String(), t.Lit)
	}

	return fmt.Sprintf("%-12s%-12s%s", t.Pos.Short(), t.Kind.String(), t.Lit)
}

func (t Token) String() string {
	if t.Kind.hasStaticLiteral() {
		return fmt.Sprintf("%s%12s", t.Pos.String(), t.Kind.String())
	}

	switch t.Kind {
	case BinaryInteger:
		return fmt.Sprintf("%s%-12s0b%b", t.Pos.String(), t.Kind.String(), t.Val)
	case OctalInteger:
		return fmt.Sprintf("%s%-12s0o%o", t.Pos.String(), t.Kind.String(), t.Val)
	case DecimalInteger:
		return fmt.Sprintf("%s%-12s%d", t.Pos.String(), t.Kind.String(), t.Val)
	case HexadecimalInteger:
		return fmt.Sprintf("%s%-12s0x%X", t.Pos.String(), t.Kind.String(), t.Val)
	case Character:
		return fmt.Sprintf("%s%-12s'%c'", t.Pos.String(), t.Kind.String(), t.Val)
	case String:
		return fmt.Sprintf("%s%-12s\"%s\"", t.Pos.String(), t.Kind.String(), t.Lit)
	}

	return fmt.Sprintf("%s%-12s%s", t.Pos.String(), t.Kind.String(), t.Lit)
}
