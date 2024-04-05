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

	// Scope's nesting level. Starts from 0 for global scope. Language structure
	// implies that first levels are dependant on Kind:
	//
	//	// kind => level
	//	- global => 0
	//	- unit   => 1
	//	- top    => 2
	//
	// Subsequent levels are created inside function and method bodies by means of
	// various language constructs.
	Level uint32

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
		Level:  parent.nextLevel(),

		sm: make(map[string]*Symbol),
	}
}

func NewUnitScope(global *Scope) *Scope {
	return NewScope(scp.Unit, global, 0)
}

func NewTopScope(unit *Scope) *Scope {
	return NewScope(scp.Top, unit, 0)
}

func (s *Scope) nextLevel() uint32 {
	if s == nil {
		return 0
	}
	return s.Level + 1
}

// Lookup a symbol by its name and source code position inside scope
// and (if not found) all of its parents. If symbol is found its usage
// count is increased.
func (s *Scope) Lookup(name string, pos uint32) *Symbol {
	sym := s.lookup(name, pos)
	if sym == nil {
		return nil
	}

	sym.RefNum += 1
	return sym
}

// Same as Lookup, but usage count for found symbol is not increased.
func (s *Scope) lookup(name string, pos uint32) *Symbol {
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

// lookup symbol in local scope without position check.
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
	sym.Scope = s
}
