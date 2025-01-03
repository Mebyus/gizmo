package stm

// Kind indicates statement kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left unspecified
	empty Kind = iota

	Block
	Assign
	Return
	Const
	Let
	Var
	If
	Call
	For
	While
	ForIn
	ForRange
	Match
	MatchBool
	Jump
	Never
	Stub

	ShortInit

	// Defer function or method call
	Defer

	// Defer block statement
	DeferBlock

	// If without else and else-if clauses.
	SimpleIf
)

var text = [...]string{
	empty: "<nil>",

	Block:    "block",
	Assign:   "assign",
	Return:   "return",
	Const:    "const",
	Let:      "let",
	Var:      "var",
	If:       "if",
	Call:     "call",
	For:      "for",
	While:    "while",
	ForIn:    "for.in",
	ForRange: "for.range",
	Match:    "match",
	Defer:    "defer",
	Jump:     "jump",
	Never:    "never",
	Stub:     "stub",

	ShortInit: "var.init.short",

	MatchBool: "match.bool",

	DeferBlock: "defer.block",

	SimpleIf: "if.simple",
}

func (k Kind) String() string {
	return text[k]
}
