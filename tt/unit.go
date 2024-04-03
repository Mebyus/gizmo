package tt

type Unit struct {
	Name string

	// Scope that holds all top-level symbols from all unit atoms.
	//
	// This field is always not nil and Scope.Kind is always equal to scp.Unit.
	Scope *Scope
}
