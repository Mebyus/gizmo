package tt

type ConstDef struct {
	nodeSymDef

	Expr Expression

	// List of top-level unit symbols which are used in constant definition.
	// Only applicable for top-level unit constants. For constants inside
	// functions and methods this field is always nil.
	Refs []*Symbol

	// Declared constant type. Equals nil if constant is untyped.
	Type *Type
}
