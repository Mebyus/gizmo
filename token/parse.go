package token

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/mebyus/gizmo/source"
)

func ParseKind(code string) (Kind, bool) {
	k, ok := reverse[code]
	return k, ok
}

var (
	ErrBadPosFormat   = errors.New("bad position spec format")
	ErrBadLineFormat  = errors.New("bad line format")
	ErrBadColFormat   = errors.New("bad column format")
	ErrZeroLineNumber = errors.New("zero line number")
	ErrZeroColNumber  = errors.New("zero column number")
)

func ParsePos(spec string) (source.Pos, error) {
	if len(spec) < 3 {
		return source.Pos{}, ErrBadPosFormat
	}
	j := -1
	for i := 0; i < len(spec); i++ {
		if spec[i] == ':' {
			j = i
			break
		}
	}
	if j == -1 || j == 0 || j == len(spec)-1 {
		return source.Pos{}, ErrBadPosFormat
	}

	linestr := spec[:j]
	colstr := spec[j+1:]

	line, err := strconv.ParseUint(linestr, 10, 32)
	if err != nil {
		return source.Pos{}, ErrBadLineFormat
	}
	col, err := strconv.ParseUint(colstr, 10, 32)
	if err != nil {
		return source.Pos{}, ErrBadColFormat
	}
	if line == 0 {
		return source.Pos{}, ErrZeroLineNumber
	}
	if col == 0 {
		return source.Pos{}, ErrZeroColNumber
	}

	return pos(uint32(line), uint32(col)), nil
}

func pos(line, col uint32) source.Pos {
	return source.Pos{
		Line: line - 1,
		Col:  col - 1,
	}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t'
}

// trims left and right space characters (spaces and tabs)
// and splits string into at most three non-empty strings
// by space characters
//
// returns number of found parts as second value
func splitTokenLine(s string) ([3]string, int) {
	var split [3]string

	i := 0

	// skip leading spaces
	for ; i < len(s) && isSpace(s[i]); i++ {
	}
	if i == len(s) {
		return split, 0
	}

	// first word start
	j := i
	for ; i < len(s) && !isSpace(s[i]); i++ {
	}
	split[0] = s[j:i]
	if i == len(s) {
		return split, 1
	}

	// skip spaces after first word
	for ; i < len(s) && isSpace(s[i]); i++ {
	}
	if i == len(s) {
		return split, 1
	}

	// second word start
	j = i
	for ; i < len(s) && !isSpace(s[i]); i++ {
	}
	split[1] = s[j:i]
	if i == len(s) {
		return split, 2
	}

	// skip spaces after second word
	for ; i < len(s) && isSpace(s[i]); i++ {
	}
	if i == len(s) {
		return split, 2
	}

	// last part start
	j = i

	// search last part end (by skipping spaces from string end)
	for i = len(s) - 1; isSpace(s[i]); i-- {
	}
	split[2] = s[j : i+1]

	return split, 3
}

var (
	ErrBadTokenFormat  = errors.New("bad token format")
	ErrStaticLiteral   = errors.New("token has static literal")
	ErrLiteralNotFound = errors.New("literal not found")
	ErrBadStringFormat = errors.New("bad string format")
)

func Parse(line string) (Token, error) {
	fields, parts := splitTokenLine(line)
	if parts != 2 && parts != 3 {
		return Token{}, ErrBadTokenFormat
	}
	spec := fields[0]
	code := fields[1]
	k, ok := ParseKind(code)
	if !ok {
		return Token{}, fmt.Errorf("unknown kind: %s", code)
	}

	var lit string
	var val uint64
	if k.hasStaticLiteral() {
		if parts != 2 {
			return Token{}, ErrStaticLiteral
		}
	} else {
		if parts != 3 {
			return Token{}, ErrLiteralNotFound
		}
		lit = fields[2]
		var err error
		switch k {
		case DecimalInteger:
			val, err = strconv.ParseUint(lit, 10, 64)
			if err != nil {
				return Token{}, err
			}
			lit = ""
		case String:
			if len(lit) < 2 {
				return Token{}, ErrBadStringFormat
			}
			if lit[0] != '"' || lit[len(lit)-1] != '"' {
				return Token{}, ErrBadStringFormat
			}
			lit = lit[1 : len(lit)-1]
		default:
		}
	}

	pos, err := ParsePos(spec)
	if err != nil {
		return Token{}, err
	}

	return Token{
		Kind: k,
		Lit:  lit,
		Pos:  pos,
		Val:  val,
	}, nil
}

// maps token literal or code literal to its Kind
var reverse = make(map[string]Kind, len(Literal))

func init() {
	for k, lit := range Literal {
		if lit == "" {
			if k < int(noStaticLiteral) {
				panic(fmt.Sprintf("kind = %d should have non-empty static literal", k))
			}
			continue
		}
		kind := Kind(k)
		reverse[lit] = kind
	}
}
