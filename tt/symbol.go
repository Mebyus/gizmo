package tt

import (
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/tt/sym"
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
	Def SymDef

	// Always an alphanumerical word. Always not empty.
	Name string

	// Always not nil in a completed tree. Can be nil during tree construction.
	Type *Type

	// Scope where this symbol was declared.
	Scope *Scope

	Kind sym.Kind

	// Number of times symbol is referenced (used in source code outside its origin)
	RefNum uint32

	// Used only for top-level symbols in unit. If true than other units, which import unit where this
	// symbol is declared, can use this symbol in their code.
	Public bool
}

type SymDef interface {
	// dummy discriminator method
	SymDef()
}

// Dummy provides quick, easy to use implementation of discriminator SymDef() method
//
// Used for embedding into other (non-dummy) symbol definition nodes
type nodeSymDef struct{}

func (nodeSymDef) SymDef() {}
