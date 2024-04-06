package tt

// Context carries information and objects which are needed to perform
// statements and expressions transformations (from AST to TT), type checking,
// symbol usage counting, etc.
type Context struct {
	m *Merger

	// gathers top-level symbols used in function, method, unit constant or variable.
	ref SymSet

	// return type of a function or method, equals nil if it returns nothing or never.
	ret *Type

	// true for function or method that does not return.
	never bool
}

func (m *Merger) newFnCtx(def *FnDef) *Context {
	return &Context{
		m:     m,
		ref:   NewSymSet(),
		ret:   def.Result,
		never: def.Never,
	}
}
