package tt

import "github.com/mebyus/gizmo/tt/scp"

type Scope struct {
	// List of all symbols defined inside this scope. Symbols are
	// listed in order they appear in source code (except for global and unit scopes).
	Symbols []*Symbol

	// First levels of scope hierarchy are fixed:
	//
	//	- global
	//	- unit
	//	- top
	//
	// Next levels may vary based on source code that defines the scope.
	Parent *Scope

	// Symbol map. Maps name to its local symbol.
	sm map[string]*Symbol

	// Types which are visible inside this Scope
	// Types *Tindex

	// Local position of scope inclusion point in parent scope.
	// Determined by Pos.Num of scope block start. Irrelevant for
	// global, unit and top scopes.
	Pos uint32

	Kind scp.Kind
}

func NewScope(kind scp.Kind, parent *Scope, pos uint32) *Scope {
	if kind == scp.Global && parent != nil {
		panic("global scope cannot have a parent")
	}
	if kind == scp.Unit && parent.Kind != scp.Global {
		panic("parent of unit scope must always be global scope")
	}
	if kind == scp.Top && parent.Kind != scp.Unit {
		panic("parent of top-level scope must always be unit scope")
	}
	if kind > scp.Top && pos == 0 {
		panic("scope parent inclusion position is not specified")
	}
	return &Scope{
		Kind:   kind,
		Pos:    pos,
		Parent: parent,
		sm:     make(map[string]*Symbol),
	}
}

func NewGlobalScope() *Scope {
	return NewScope(scp.Global, nil, 0)
}

func NewUnitScope(global *Scope) *Scope {
	return NewScope(scp.Unit, global, 0)
}

func NewTopScope(unit *Scope) *Scope {
	return NewScope(scp.Top, unit, 0)
}

// Lookup a symbol by its name and source code position inside scope
// and (if not found) all of its parents.
func (s *Scope) Lookup(name string, pos uint32) *Symbol {
	// for global and unit scopes lookup position is irrelevant
	if s.Kind == scp.Global {
		return s.sym(name)
	}
	if s.Kind == scp.Unit {
		sym := s.sym(name)
		if sym != nil {
			return sym
		}

		// parent of any unit scope is always the global scope
		return s.Parent.sym(name)
	}

	sym := s.Sym(name, pos)
	if sym != nil {
		return sym
	}
	return s.Parent.Lookup(name, pos)
}

// Sym tries to lookup a symbol by its name inside local scope
// (without going for parent lookup). Returns nil if symbol is not defined
// inside this scope.
func (s *Scope) Sym(name string, pos uint32) *Symbol {
	sym := s.sym(name)
	if sym == nil {
		return nil
	}
	if pos > sym.Pos.Num {
		// if lookup position is after symbol definition position,
		// then lookup is successful
		return sym
	}
	return nil
}

// lookup symbol in local scope without position check
func (s *Scope) sym(name string) *Symbol {
	if len(s.Symbols) == 0 {
		return nil
	}

	// TODO: tweak this constant for optimized lookup perfomance
	if len(s.Symbols) > 16 {
		return s.sm[name]
	}

	for _, sym := range s.Symbols {
		if sym.Name == name {
			return sym
		}
	}
	return nil
}

// Bind adds given symbol to this scope.
func (s *Scope) Bind(sym *Symbol) {
	s.Symbols = append(s.Symbols, sym)
	s.sm[sym.Name] = sym
}
