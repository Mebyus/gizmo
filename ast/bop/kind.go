package bop

import (
	"fmt"

	"github.com/mebyus/gizmo/token"
)

// Kind indicates binary operator kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left inspecified
	empty Kind = iota

	Equal    // ' == '
	NotEqual // ' != '

	Less           // ' < '
	Greater        // ' > '
	LessOrEqual    // ' <= '
	GreaterOrEqual // ' >= '

	And // ' && '
	Or  // ' || '

	Add // ' + '
	Sub // ' - '
	Mul // ' * '
	Div // ' / '
	Mod // ' % '

	Xor           // ' ^ '
	BitwiseAnd    // ' & '
	BitwiseOr     // ' | '
	BitwiseAndNot // ' &^ '
	LeftShift     // ' << '
	RightShift    // ' >> '
)

var text = [...]string{
	empty: "<nil>",

	Equal:    "==",
	NotEqual: "!=",

	Less:           "<",
	Greater:        ">",
	LessOrEqual:    "<=",
	GreaterOrEqual: ">=",

	And: "&&",
	Or:  "||",

	Add: "+",
	Sub: "-",
	Mul: "*",
	Div: "/",
	Mod: "%",

	Xor:           "^",
	BitwiseAnd:    "&",
	BitwiseOr:     "|",
	BitwiseAndNot: "&^",
	LeftShift:     "<<",
	RightShift:    ">>",
}

func (k Kind) String() string {
	return text[k]
}

var precedence = [...]int{
	empty: 0,

	Mul:           1,
	Div:           1,
	Mod:           1,
	LeftShift:     1,
	RightShift:    1,
	BitwiseAnd:    1,
	BitwiseAndNot: 1,

	Add:       2,
	Sub:       2,
	Xor:       2,
	BitwiseOr: 2,

	Equal:          3,
	NotEqual:       3,
	Less:           3,
	Greater:        3,
	LessOrEqual:    3,
	GreaterOrEqual: 3,

	And: 4,

	Or: 5,
}

// Precedence gives binary operator precedence.
// Greater values mean later binding
//
//	1 :  *  /  %  <<  >>  &  &^
//	2 :  +  -  ^  |
//	3 :  ==  !=  <  <=  >  >=
//	4 :  &&
//	5 :  ||
func (k Kind) Precedence() int {
	return precedence[k]
}

func FromToken(k token.Kind) Kind {
	switch k {
	case token.Equal:
		return Equal
	case token.NotEqual:
		return NotEqual
	case token.LeftAngle:
		return Less
	case token.RightAngle:
		return Greater
	case token.LessOrEqual:
		return LessOrEqual
	case token.GreaterOrEqual:
		return GreaterOrEqual
	case token.LogicalAnd:
		return And
	case token.LogicalOr:
		return Or
	case token.Plus:
		return Add
	case token.Minus:
		return Sub
	case token.Asterisk:
		return Mul
	case token.Slash:
		return Div
	case token.Percent:
		return Mod
	case token.Caret:
		return Xor
	case token.Ampersand:
		return BitwiseAnd
	case token.Pipe:
		return BitwiseOr
	case token.BitwiseAndNot:
		return BitwiseAndNot
	case token.LeftShift:
		return LeftShift
	case token.RightShift:
		return RightShift
	default:
		panic(fmt.Sprintf("token ' %s ' cannot be binary operator", k.String()))
	}
}

/*

a=target variable, b=bit number to act upon 0-n

#define BIT_SET(a,b) ((a) |= (1ULL<<(b)))
#define BIT_CLEAR(a,b) ((a) &= ~(1ULL<<(b)))
#define BIT_FLIP(a,b) ((a) ^= (1ULL<<(b)))
#define BIT_CHECK(a,b) (!!((a) & (1ULL<<(b))))        // '!!' to make sure this returns 0 or 1

#define BITMASK_SET(x, mask) ((x) |= (mask))
#define BITMASK_CLEAR(x, mask) ((x) &= (~(mask)))
#define BITMASK_FLIP(x, mask) ((x) ^= (mask))
#define BITMASK_CHECK_ALL(x, mask) (!(~(x) & (mask)))
#define BITMASK_CHECK_ANY(x, mask) ((x) & (mask))

*/
