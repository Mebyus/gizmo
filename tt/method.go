package tt

type MethodDef struct {
	Body Block

	// Function parameters, equals nil if function has no parameters.
	Params []*Symbol

	// List of top-level unit symbols which are used in method.
	Refs []*Symbol

	// Method return type. Equals nil if function returns nothing or never returns.
	Result *Type

	// Equals true for functions which never return.
	Never bool
}
