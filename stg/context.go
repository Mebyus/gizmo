package stg

// Context carries information and objects which are needed to perform
// statements and expressions transformations (from AST to STG), type checking,
// symbol usage counting, etc.
type Context struct {
	m *Merger

	// gathers top-level symbols used in function, method, unit constant or variable.
	ref SymSet

	// return type of a function or method, equals nil if it returns nothing or never.
	ret *Type

	// receiver type of a method, equals nil for regular functions
	rv *Type

	// Context for resolving incomplete enum names.
	enum *Symbol

	// true for function or method that does not return.
	never bool

	// TODO: add fields for analyzing defers usage
	// most likely it will be
	//	- returns counter
	//	- slice of defer info (conditional/unconditional, arg expressions, etc.)
}

func (c *Context) pushEnum(s *Symbol) *Symbol {
	old := c.enum
	c.enum = s
	return old
}

func (m *Merger) newFunCtx() *Context {
	return &Context{
		m:   m,
		ref: NewSymSet(),
	}
}

func (m *Merger) newMethodCtx() *Context {
	return &Context{
		m:   m,
		ref: NewSymSet(),
	}
}

func (m *Merger) newTypeCtx() *Context {
	return &Context{
		m:   m,
		ref: NewSymSet(),
	}
}

func (m *Merger) newConstCtx() *Context {
	return &Context{
		m:   m,
		ref: NewSymSet(),
	}
}
