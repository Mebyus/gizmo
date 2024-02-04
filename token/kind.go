package token

type Kind uint32

const (
	// Zero value of Kind. Do not use explicitly
	empty Kind = iota

	// Special tokens
	EOF

	// Operators and/or punctuators

	Underscore // _

	Address  // .&
	Indirect // .@

	Plus      // +
	Minus     // -
	Asterisk  // *
	Slash     // /
	Percent   // %
	Ampersand // &
	Quest     // ?

	Pipe          // |
	Caret         // ^
	LeftShift     // <<
	RightShift    // >>
	BitwiseAndNot // &^

	Semicolon   // ;
	Period      // .
	Colon       // :
	Comma       // ,
	DoubleColon // ::

	Equal          // ==
	NotEqual       // !=
	LessOrEqual    // <=
	GreaterOrEqual // >=
	Less           // <
	Greater        // >
	Not            // !

	Assign          // =
	ShortAssign     // :=
	AddAssign       // +=
	SubtractAssign  // -=
	MultiplyAssign  // *=
	QuotientAssign  // /=
	RemainderAssign // %=

	LogicalAnd // &&
	LogicalOr  // ||
	LeftArrow  // <-
	RightArrow // =>

	// Brackets

	LeftCurly        // {
	RightCurly       // }
	LeftSquare       // [
	RightSquare      // ]
	LeftParentheses  // (
	RightParentheses // )

	Compound // .{
	List     // .[
	Insist   // .!
	Chain    // .?
	Slice    // []
	AutoLen  // [_]

	Nillable      // ?|
	NillableChunk // [?]

	// Keywords

	Break
	Case
	Continue

	If
	Else
	In
	For

	Defer
	Fn
	Import
	Declare

	Bag
	Bind

	Unit
	Return

	Enum
	Struct
	Union

	Switch
	Type
	Var

	Pub
	Atr

	// Special literals

	Never
	Dirty
	Nil
	True
	False

	noStaticLiteral

	Illegal // any byte sequence unknown to lexer

	// Identifiers and basic type literals
	Identifier         // myvar, main, Line, print
	BinaryInteger      // 0b1101100001
	OctalInteger       // 0o43671
	DecimalInteger     // 5367, 43432, 1000097
	HexadecimalInteger // 0x43da1
	DecimalFloat       // 123.45
	Character          // 'a', '\t', 'p'
	String             // "abc", "", "\t\n  42Hello\n"
)

func (k Kind) String() string {
	return Literal[k]
}

func (k Kind) IsEmpty() bool {
	return k == empty
}

func (k Kind) IsEOF() bool {
	return k == EOF
}

func (k Kind) hasStaticLiteral() bool {
	return k < noStaticLiteral
}

func (k Kind) IsIdent() bool {
	return k == Identifier
}

func (k Kind) IsLeftPar() bool {
	return k == LeftParentheses
}

func (k Kind) IsLit() bool {
	switch k {
	case
		String, Character, DecimalInteger, DecimalFloat, BinaryInteger,
		OctalInteger, HexadecimalInteger, True, False, Nil:

		return true
	default:
		return false
	}
}

func (k Kind) IsUnaryOperator() bool {
	switch k {
	case Plus, Minus, Not:
		return true
	default:
		return false
	}
}

func (k Kind) IsBinaryOperator() bool {
	switch k {
	case
		Less, Greater, Equal, LessOrEqual, GreaterOrEqual, NotEqual, LogicalAnd,
		LogicalOr, Plus, Minus, Asterisk, Slash, Percent:

		return true
	default:
		return false
	}
}
