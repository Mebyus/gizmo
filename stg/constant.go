package stg

type ConstDef struct {
	nodeSymDef

	// Always not nil.
	Exp Exp

	// List of top-level unit symbols which are used in constant definition.
	// Only applicable for top-level unit constants. For constants inside
	// functions and methods this field is always nil.
	Refs []*Symbol

	// Constant type. If definition does not contain type specifier
	// then type will be inferred from expession type.
	//
	// Always not nil.
	Type *Type
}
