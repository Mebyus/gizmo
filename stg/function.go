package stg

type FunDef struct {
	nodeSymDef

	Signature

	Body Block

	// List of top-level unit symbols which are used in function.
	Refs []*Symbol

	// Temporary field for transferring context between inspect
	// and full scan phases.
	ctx *Context
}

// Explicit interface implementation check.
var _ SymDef = &FunDef{}

// Signature provides information about a call layout
// (of a function, method, function pointer, etc.).
type Signature struct {
	// Function parameters, equals nil if function has no parameters.
	Params []*Symbol

	// Function return type. Equals nil if function returns nothing or never returns.
	Result *Type

	// Equals true for functions which never return.
	Never bool
}
