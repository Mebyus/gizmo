package toplvl

// Kind indicates top level construct kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left unspecified
	empty Kind = iota

	// Function definition
	Fn

	// Function declaration
	Declare

	// Method definition
	Method

	// Type definition
	Type

	// Type definition of evaluated type
	TypeEval

	// Constant definition (name + type + value)
	Const

	// Variable definition (name + type + initial value)
	Var

	// Blueprint, aka "function template"
	Blue

	// Prototype, aka "type template"
	Proto

	// Prototype method bluepint, aka "method template"
	Pmb
)

var text = [...]string{
	empty: "<nil>",

	Fn:       "fn",
	Declare:  "declare",
	Method:   "method",
	Type:     "type",
	TypeEval: "type_eval",
	Const:    "const",
	Var:      "var",
	Blue:     "blue",
	Proto:    "proto",
	Pmb:      "pmb",
}

func (k Kind) String() string {
	return text[k]
}
