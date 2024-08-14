package token

type Kind uint32

const (
	// Zero value of Kind. Do not use explicitly
	empty Kind = iota

	// Special tokens
	EOF

	// Operators and/or punctuators

	Underscore // _

	Address       // .&
	Indirect      // .@
	IndirectIndex // .[

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
	LeftAngle      // <
	RightAngle     // >
	Not            // !

	Assign    // =
	Walrus    // :=
	AddAssign // +=
	SubAssign // -=
	MulAssign // *=
	DivAssign // /=
	RemAssign // %=

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

	PropStart // #[

	Compound     // .{
	Insist       // .!
	Chain        // .?
	Chunk        // []
	AutoLen      // [_]
	ArrayPointer // [*]
	CapBuffer    // [^]

	Nillable      // ?|
	NillableChunk // [?]

	// Keywords

	Case
	Jump

	If
	Else
	In
	For

	Defer
	Fun
	Import

	Bag

	Unit
	Return

	Enum
	Struct
	Union

	Match
	Type
	Let
	Var

	Pub

	// Special literals

	Never
	Dirty
	Nil
	True
	False

	Cast
	Tint    // truncate (cast with storage size change) integer
	BitCast // TODO: rename this to memcast

	Receiver // method receiver, consists of one letter "g"

	LabelNext // @.next
	LabelEnd  // @.end

	noStaticLiteral

	Illegal // any byte sequence unknown to lexer

	// Identifiers and basic type literals
	Identifier         // myvar, main, Line, print
	BinaryInteger      // 0b1101100001
	OctalInteger       // 0o43671
	DecimalInteger     // 5367, 43432, 1000097
	HexadecimalInteger // 0x43da1
	DecimalFloat       // 123.45
	Rune               // 'a', '\t', 'p'
	String             // "abc", "", "\t\n  42Hello\n"

	// Comments
	LineComment  // Line comment starts with //
	BlockComment // Comment inside /* comment */ block
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

func (k Kind) IsComment() bool {
	return k == LineComment || k == BlockComment
}

func (k Kind) IsLeftPar() bool {
	return k == LeftParentheses
}

func (k Kind) IsLit() bool {
	switch k {
	case
		String, Rune, DecimalInteger, DecimalFloat, BinaryInteger,
		OctalInteger, HexadecimalInteger, True, False, Nil:

		return true
	default:
		return false
	}
}

func (k Kind) IsUnaryOperator() bool {
	switch k {
	case Plus, Minus, Not, Caret:
		return true
	default:
		return false
	}
}

func (k Kind) IsBinaryOperator() bool {
	switch k {
	case
		LeftAngle, RightAngle, Equal, LessOrEqual, GreaterOrEqual, NotEqual, LogicalAnd,
		LogicalOr, Plus, Minus, Asterisk, Slash, Percent, LeftShift, RightShift,
		Pipe, Ampersand, BitwiseAndNot, Caret:

		return true
	default:
		return false
	}
}
