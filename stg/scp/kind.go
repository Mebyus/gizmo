package scp

// Kind indicates scope kind
type Kind uint8

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	// Scope with global program-wide access. Holds builtin language symbols
	Global

	// Scope with unit-wide access (across all unit atoms)
	Unit

	// Scope created by top-level function or method
	Top

	// Scope created by function literal (with or without closure).
	FnLit

	Block
	Loop
	If
	IfElse
	Else
	Case
)

var text = [...]string{
	empty: "<nil>",

	Global: "global",
	Unit:   "unit",
	Top:    "top",
	FnLit:  "fnlit",
	Block:  "block",
	Loop:   "loop",
	If:     "if",
	IfElse: "ifelse",
	Else:   "else",
	Case:   "case",
}

func (k Kind) String() string {
	return text[k]
}
