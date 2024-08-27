package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/enums/tpk"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/stg/scp"
)

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

	// Types which are visible inside this scope.
	Types *TypeIndex

	// Symbol map. Maps name to its local symbol.
	sm map[string]*Symbol

	// Local position of scope inclusion point in parent scope.
	// Determined by Pos of scope block start. Irrelevant for
	// global, unit and top scopes.
	Pos *source.Pos

	// Unit where this scope is located.
	// Always nil for global scope.
	Unit *Unit

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

	// How many loops deep this scope resides (within a function or method).
	// Starts from 0 for top scope (where there are no surrounding loops).
	LoopLevel uint32

	// How many logic branches deep (excluding loops) this scope resides (within a function or method).
	// Starts from 0 for top scope (where there are no surrounding branches).
	BranchLevel uint32

	Kind scp.Kind
}

func NewScope(kind scp.Kind, parent *Scope, pos *source.Pos) *Scope {
	if kind == scp.Global && parent != nil {
		panic("global scope cannot have a parent")
	}
	if kind == scp.Unit && parent.Kind != scp.Global {
		panic("parent of unit scope must always be global scope")
	}
	if kind == scp.Top && parent.Kind != scp.Unit {
		panic("parent of top-level scope must always be unit scope")
	}
	if kind >= scp.Top && pos == nil {
		panic("scope parent inclusion position is not specified")
	}
	s := &Scope{
		Kind:   kind,
		Pos:    pos,
		Parent: parent,
		Level:  parent.nextLevel(),

		sm: make(map[string]*Symbol),
	}
	s.Types = &TypeIndex{
		scope: s,
		tm:    inheritTypeMap(parent),
	}
	if parent != nil {
		s.Unit = parent.Unit
		s.LoopLevel = parent.LoopLevel
		s.BranchLevel = parent.BranchLevel

		switch kind {
		case scp.Loop:
			s.LoopLevel += 1
		case scp.If, scp.Case, scp.Else:
			s.BranchLevel += 1
		}
	}
	return s
}

func inheritTypeMap(parent *Scope) map[Stable]*Type {
	if parent == nil {
		return make(map[Stable]*Type)
	}
	return parent.Types.tm
}

func NewUnitScope(unit *Unit, global *Scope) *Scope {
	s := NewScope(scp.Unit, global, nil)
	s.Unit = unit
	return s
}

func NewTopScope(parent *Scope, pos *source.Pos) *Scope {
	return NewScope(scp.Top, parent, pos)
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
//
// Note that for global and unit scopes position is irrelevant, since
// symbols created there are subject to hoisting.
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
func (s *Scope) Bind(symbol *Symbol) {
	s.bind(symbol)
}

func (s *Scope) bind(symbol *Symbol) {
	s.Symbols = append(s.Symbols, symbol)
	s.sm[symbol.Name] = symbol
	symbol.Scope = s
}

func (s *Scope) BindTypeSymbol(symbol *Symbol) {
	if symbol.Kind != smk.Type {
		panic("method must be called only with symbols representing a type")
	}

	t := symbol.Def.(*Type)

	if t.Kind == tpk.Custom {
		panic(fmt.Sprintf("unexpected type kind: %s", t.Kind.String()))
	}

	s.bind(symbol)
}

// CheckUsage scans symbols usage count in this scope. Returns error if there are
// declared, but not used symbols.
func (s *Scope) CheckUsage(ctx *Context) error {
	var list []*Symbol
	for _, symbol := range s.Symbols {
		if symbol.RefNum == 0 {
			if symbol.Kind == smk.Param {
				ctx.m.warn(symbol.Pos, fmt.Sprintf("unused function parameter \"%s\"", symbol.Name))
			} else {
				list = append(list, symbol)
			}
		}
	}
	if len(list) == 0 {
		return nil
	}
	sym := list[0]
	return fmt.Errorf("%s: symbol \"%s\" declared and not used", sym.Pos.String(), sym.Name)
}

// WarnUnused acts much like CheckUsage, but does not produce an error, only warnings
// about unused symbols.
func (s *Scope) WarnUnused(ctx *Context) {
	for _, symbol := range s.Symbols {
		if !symbol.Pub && symbol.RefNum == 0 {
			ctx.m.warn(symbol.Pos, fmt.Sprintf("symbol \"%s\" has no references inside unit", symbol.Name))
		}
	}
}
