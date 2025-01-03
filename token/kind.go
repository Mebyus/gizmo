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

	Compound       // .{
	Insist         // .!
	Chain          // .?
	Chunk          // []
	AutoLen        // [_]
	ArrayPointer   // [*]
	CapBuffer      // [^]
	IndirectSelect // .@.

	Nillable      // ?|
	NillableChunk // [?]

	// Keywords

	Jump

	If
	Else
	In
	For

	Defer
	Fun
	Import
	Test

	Bag

	Unit
	Return

	Enum
	Struct
	Union

	Type
	Let
	Var

	Pub

	// Special literals

	Never
	Stub
	Dirty
	Nil
	True
	False

	Cast
	Tint // truncate (cast with storage size change) integer
	MemCast
	MemSize

	Any // designator to use as *any (void* analog)

	LabelNext // @.next
	LabelOut  // @.out

	DirName    // #name
	DirInclude // #include
	DirLink    // #link
	DirIf      // #if

	noStaticLiteral

	Illegal // any byte sequence unknown to lexer

	// Identifiers and basic type literals
	Identifier // myvar, main, Line, print
	BinInteger // 0b1101100001
	OctInteger // 0o43671
	DecInteger // 5367, 43432, 1000097
	HexInteger // 0x43da1
	DecFloat   // 123.45
	Rune       // 'a', '\t', 'p'
	String     // "abc", "", "\t\n  42Hello\n"
	RawString  // #"raw string literal"
	FillString // "string with ${10 + 1} interpolated ${a - b} expressions"
	Macro      // #:MACRO_NAME

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
		String, Rune, DecInteger, DecFloat, BinInteger,
		OctInteger, HexInteger, True, False, Nil:

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
