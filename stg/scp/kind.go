package scp

// Kind indicates scope kind
type Kind uint8

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	// Scope with global program-wide access. Holds builtin language symbols.
	Global

	// Scope with unit-wide access (across all unit atoms).
	//
	// Does not include unit tests.
	Unit

	// Scope created by unit level function (inside function body).
	Fun

	// Scope created by method (inside method body).
	Method

	// Scope that holds collection of all tests inside a unit.
	UnitTests

	// Scope created by unit test function.
	Test

	// Scope created by function literal (with or without closure).
	FunLit

	Block
	Loop
	If
	Else
	Case
)

var text = [...]string{
	empty: "<nil>",

	Global: "global",
	Unit:   "unit",
	Fun:    "fun",
	FunLit: "fnlit",
	Block:  "block",
	Loop:   "loop",
	If:     "if",
	Else:   "else",
	Case:   "case",

	UnitTests: "unit.tests",
}

func (k Kind) String() string {
	return text[k]
}
