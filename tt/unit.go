package tt

type Unit struct {
	// List of all top-level symbols defined in this unit
	Top []*Symbol

	// Symbol map. Maps top-level symbol name to the symbol itself
	sm map[string]*Symbol
}
