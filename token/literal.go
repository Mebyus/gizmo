package token

// Literal maps Kind to token static literal
var Literal = [...]string{
	empty: "EMPTY",

	EOF: "EOF",

	// Operators/punctuators

	Underscore:       "_",
	Address:          ".&",
	Indirect:         ".@",
	Plus:             "+",
	Minus:            "-",
	LogicalAnd:       "&&",
	LogicalOr:        "||",
	Equal:            "==",
	NotEqual:         "!=",
	LessOrEqual:      "<=",
	GreaterOrEqual:   ">=",
	LeftArrow:        "<-",
	RightArrow:       "=>",
	ShortAssign:      ":=",
	AddAssign:        "+=",
	SubtractAssign:   "-=",
	MultiplyAssign:   "*=",
	QuotientAssign:   "/=",
	RemainderAssign:  "%=",
	Pipe:             "|",
	Caret:            "^",
	LeftShift:        "<<",
	RightShift:       ">>",
	BitwiseAndNot:    "&^",
	Assign:           "=",
	Colon:            ":",
	DoubleColon:      "::",
	Semicolon:        ";",
	Asterisk:         "*",
	Quest:            "?",
	Ampersand:        "&",
	Not:              "!",
	Slash:            "/",
	Percent:          "%",
	Period:           ".",
	Comma:            ",",
	Less:             "<",
	Greater:          ">",
	LeftCurly:        "{",
	RightCurly:       "}",
	LeftSquare:       "[",
	RightSquare:      "]",
	LeftParentheses:  "(",
	RightParentheses: ")",
	Compound:         ".{",
	List:             ".[",
	Insist:           ".!",
	Chain:            ".?",
	Slice:            "[]",
	AutoLen:          "[_]",
	Nillable:         "?|",
	NillableChunk:    "[?]",

	// Keywords

	Import:   "import",
	Fn:       "fn",
	Continue: "continue",
	Return:   "return",
	Break:    "break",
	Case:     "case",
	For:      "for",
	Else:     "else",
	If:       "if",
	Defer:    "defer",
	Bag:      "bag",
	Bind:     "bind",
	In:       "in",
	Var:      "var",
	Type:     "type",
	Switch:   "switch",
	Enum:     "enum",
	Struct:   "struct",
	Union:    "union",
	Pub:      "pub",
	Unit:     "unit",
	Atr:      "atr",
	Declare:  "declare",

	// Special literals

	Never: "never",
	Dirty: "dirty",
	Nil:   "nil",
	True:  "true",
	False: "false",

	// Non static literals

	Illegal:            "ILLEGAL",
	Identifier:         "IDENT",
	String:             "STR",
	Character:          "CHAR",
	BinaryInteger:      "INT.BIN",
	OctalInteger:       "INT.OCT",
	DecimalInteger:     "INT.DEC",
	HexadecimalInteger: "INT.HEX",
	DecimalFloat:       "FLT.DEC",
}
