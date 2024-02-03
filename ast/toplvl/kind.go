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
	Decl

	// Method definition
	Method

	// Type definition
	Type

	// Constant definition (name + type + value)
	Const

	// Variable definition (name + type + initial value)
	Var

	Template
)

var text = [...]string{
	empty: "<nil>",

	Fn:       "fn",
	Decl:     "decl",
	Method:   "method",
	Type:     "type",
	Const:    "const",
	Var:      "var",
	Template: "template",
}

func (k Kind) String() string {
	return text[k]
}
