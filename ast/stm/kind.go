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
	ForCond
	ForEach
	Match
	MatchBool
	Jump
	Never

	// Defer function or method call
	Defer

	// Defer block statement
	DeferBlock

	// If without else and else-if clauses.
	SimpleIf
)

var text = [...]string{
	empty: "<nil>",

	Block:   "block",
	Assign:  "assign",
	Return:  "return",
	Const:   "const",
	Let:     "let",
	Var:     "var",
	If:      "if",
	Call:    "call",
	For:     "for",
	ForCond: "for_cond",
	ForEach: "for_each",
	Match:   "match",
	Defer:   "defer",
	Jump:    "jump",
	Never:   "never",

	MatchBool: "match.bool",

	DeferBlock: "defer.block",

	SimpleIf: "if.simple",
}

func (k Kind) String() string {
	return text[k]
}
