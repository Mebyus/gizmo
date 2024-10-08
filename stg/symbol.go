package stg

import (
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/source"
)

// Symbol. Symbols represent objects in source code which we can reference and
// use by using their names (identifiers). Most symbols are created with custom names
// in user's code by using language constructs:
//
//   - var - variable
//   - let - runtime constant, i.e. variable that is read-only, value is assigned to it exactly once
//   - const - build-time constant, i.e. its value must be known or computable at build-time
//   - type - defines a new type or prototype
//   - fn - defines a new function, method, blueprint or prototype blueprint
//   - method - defines a method, can only be used at unit level
//   - import - binds other unit to a local name in current unit
//   - function or prototype parameters
//
// Builtin symbols are created by the compiler and include basic types and primitive functions.
//
// Each symbol has a name, type and parent scope. The latter is determined by the place of origin (declaration)
// of the symbol.
//
// For each unique symbol in a program this struct must be instanced exactly once and then
// passed around in a pointer without creating a dereferenced copy.
type Symbol struct {
	// Source position of symbol origin (where this symbol was declared).
	Pos source.Pos

	// Symbol definition. Actual value stored in this field depends on Kind field.
	// See implementations of SymDef interface for more information.
	//
	// This field can be nil for dried symbol. Typical case for this is separate units compilation.
	//
	// During indexing and type checking this field may contain
	// temporary intermidiate values. For example for most symbols
	// during atoms gathering it is set to symbol index inside the box
	// of corresponding kind.
	Def SymDef

	// Always not empty.
	// Always an alphanumerical word for all symbol types except methods.
	//
	// For methods this field has special format:
	//	"receiver.name"
	//
	// Since other symbol types cannot have period in their names and
	// each custom type method names must be unique this
	// naming scheme cannot have accidental collisions.
	Name string

	// Always not nil in a completed tree. Can be nil during tree construction.
	Type *Type

	// Scope where this symbol was declared.
	Scope *Scope

	Kind smk.Kind

	// Number of times symbol is referenced (used in source code outside its origin).
	RefNum uint32

	// Used only for unit level symbols in unit. If true than other units,
	// which import unit where this symbol is declared, can use this symbol in their code.
	Pub bool
}

func (s *Symbol) Index() astIndexSymDef {
	return s.Def.(astIndexSymDef)
}

// Set of symbols
type SymSet map[*Symbol]struct{}

func NewSymSet() SymSet {
	return make(SymSet)
}

func (s SymSet) Add(sym *Symbol) {
	s[sym] = struct{}{}
}

func (s SymSet) Has(sym *Symbol) bool {
	_, ok := s[sym]
	return ok
}

func (s SymSet) Elems() []*Symbol {
	if len(s) == 0 {
		return nil
	}

	elems := make([]*Symbol, 0, len(s))
	for sym := range s {
		elems = append(elems, sym)
	}
	return elems
}

type SymDef interface {
	// dummy discriminator method
	SymDef()
}

// This is dummy implementation of SymDef interface.
//
// Used for embedding into other (non-dummy) symbol definition nodes.
type nodeSymDef struct{}

func (nodeSymDef) SymDef() {}

// Contains index inside NodesBox's slice of nodes of one
// of the following specific types:
//   - function
//   - type
//   - variable
//   - constant
//   - method
type astIndexSymDef int

func (astIndexSymDef) SymDef() {}

type ImportSymDef struct {
	nodeSymDef

	// Unit imported by symbol.
	Unit *Unit
}
