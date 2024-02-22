package toplvl

// Kind indicates top level construct kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left inspecified
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

	// Function template
	FnTemplate
)

var text = [...]string{
	empty: "<nil>",

	Fn:         "fn",
	Declare:    "declare",
	Method:     "method",
	Type:       "type",
	TypeEval:   "type_eval",
	Const:      "const",
	Var:        "var",
	FnTemplate: "fn_template",
}

func (k Kind) String() string {
	return text[k]
}
