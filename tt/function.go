package tt

type FnDef struct {
	nodeSymDef

	// Function parameters, equals nil if function has no parameters.
	Params []*Symbol

	// List of top-level package symbols which are used in function.
	Refs []*Symbol

	// Body Block

	// Function return type. Equals nil if function returns nothing or never returns.
	Result *Type

	// Equals true for functions which never return.
	Never bool
}

// Explicit interface implementation check.
var _ SymDef = &FnDef{}
