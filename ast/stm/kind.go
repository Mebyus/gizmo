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
	Let
	Var
	If
	Expr
	For
	ForCond
	ForEach
	Match
	Jump

	// Defer function or method call
	Defer

	// Defer block statement
	DeferBlock

	// Example:
	//
	//	x = 10 + a;
	SymbolAssign

	// Example:
	//
	//	x.@ = 10 + a;
	IndirectAssign
)

var text = [...]string{
	empty: "<nil>",

	Block:     "block",
	Assign:    "assign",
	AddAssign: "add_assign",
	Return:    "return",
	Const:     "const",
	Let:       "let",
	Var:       "var",
	If:        "if",
	Expr:      "expr",
	For:       "for",
	ForCond:   "for_cond",
	ForEach:   "for_each",
	Match:     "match",
	Jump:      "jump",
	Defer:     "defer",

	DeferBlock:     "defer_block",
	
	SymbolAssign:   "symbol_assign",
	IndirectAssign: "indirect_assign",
}

func (k Kind) String() string {
	return text[k]
}
