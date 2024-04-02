package sym

// Kind indicates symbol kind.
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	// Function (produced from declaration or definition)
	Fn

	// Method (produced from declaration or definition)
	Method

	// Type definition
	Type

	// Buildtime constant definition (name + type + value)
	Const

	// Runtime constant definition (name + type + value)
	Let

	// Variable definition (name + type + initial value)
	Var

	// Blueprint, aka "function template"
	Blue

	// Prototype, aka "type template"
	Proto

	// Prototype method bluepint, aka "method template"
	Pmb

	// Symbol created by importing other unit
	Import
)

var text = [...]string{
	empty: "<nil>",

	Fn:     "fn",
	Method: "method",
	Type:   "type",
	Const:  "const",
	Let:    "let",
	Var:    "var",
	Blue:   "blue",
	Proto:  "proto",
	Pmb:    "pmb",
	Import: "import",
}

func (k Kind) String() string {
	return text[k]
}
