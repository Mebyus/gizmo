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
	AddAssign
	Return
	Const
	Var
	If
	Expr
	For
	ForCond
	Match
	Jump
)

var text = [...]string{
	empty: "<nil>",

	Block:     "block",
	Assign:    "assign",
	AddAssign: "add_assign",
	Return:    "return",
	Const:     "const",
	Var:       "var",
	If:        "if",
	Expr:      "expr",
	For:       "for",
	ForCond:   "for_cond",
	Match:     "match",
	Jump:      "jump",
}

func (k Kind) String() string {
	return text[k]
}
