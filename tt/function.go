package tt

type FunDef struct {
	nodeSymDef

	Body Block

	// Function parameters, equals nil if function has no parameters.
	Params []*Symbol

	// List of top-level unit symbols which are used in function.
	Refs []*Symbol

	// Function return type. Equals nil if function returns nothing or never returns.
	Result *Type

	// Equals true for functions which never return.
	Never bool
}

// Explicit interface implementation check.
var _ SymDef = &FunDef{}
