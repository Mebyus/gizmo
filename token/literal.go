package token

// Literal maps Kind to token static literal
var Literal = [...]string{
	empty: "<nil>",

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
	Walrus:           ":=",
	AddAssign:        "+=",
	SubAssign:        "-=",
	MulAssign:        "*=",
	DivAssign:        "/=",
	RemAssign:        "%=",
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
	LeftAngle:        "<",
	RightAngle:       ">",
	LeftCurly:        "{",
	RightCurly:       "}",
	LeftSquare:       "[",
	RightSquare:      "]",
	LeftParentheses:  "(",
	RightParentheses: ")",
	PropStart:        "#[",
	Compound:         ".{",
	IndirectIndex:    ".[",
	Insist:           ".!",
	Chain:            ".?",
	Chunk:            "[]",
	AutoLen:          "[_]",
	ArrayPointer:     "[*]",
	CapBuffer:        "[^]",
	Nillable:         "?|",
	NillableChunk:    "[?]",

	// Keywords

	Import: "import",
	Fun:    "fun",
	Jump:   "jump",
	Return: "return",
	For:    "for",
	Else:   "else",
	If:     "if",
	Defer:  "defer",
	Bag:    "bag",
	In:     "in",
	Var:    "var",
	Type:   "type",
	Test:   "test",
	Enum:   "enum",
	Struct: "struct",
	Union:  "union",
	Pub:    "pub",
	Unit:   "unit",
	Let:    "let",

	// Special literals

	Never: "never",
	Dirty: "dirty",
	Nil:   "nil",
	True:  "true",
	False: "false",

	Cast:    "cast",
	Tint:    "tint",
	MemCast: "mcast",
	MemSize: "msize",

	Any: "any",

	LabelNext: "@.next",
	LabelOut:  "@.out",

	DirIf:    "#if",
	DirBuild: "#build",

	// Non static literals

	Illegal:    "ILG",
	Identifier: "IDN",
	String:     "STR",
	RawString:  "STR.RAW",
	FillString: "STR.FILL",
	Rune:       "RUNE",
	BinInteger: "INT.BIN",
	OctInteger: "INT.OCT",
	DecInteger: "INT.DEC",
	HexInteger: "INT.HEX",
	DecFloat:   "FLT.DEC",
	Macro:      "MACRO",

	// Comments

	LineComment:  "COM.LINE",
	BlockComment: "COM.BLOCK",
}

// ScanStringByteSize determines how many bytes are needed to represent a given
// string literal (as written in source code) in memory. Handles utf-8 encoding
// and escape sequences. It also returns ok flag, if ok == false then string contains
// bad escape sequence
func ScanStringByteSize(s string) (uint64, bool) {
	var size uint64
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			i += 1
			switch s[i] {
			case '\\', 'n', 't', 'r', '"':
			default:
				return 0, false
			}
		}
		size += 1
	}
	return size, true
}
