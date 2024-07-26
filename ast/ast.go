package ast

import "github.com/mebyus/gizmo/source/origin"

// <UnitClause> = "unit" <UnitName>
type UnitClause struct {
	Props []Prop

	Name Identifier
}

// Atom smallest piece of processed source code inside a unit. In most
// cases this represents a file with source code. Exceptions may include
// source code generated at compile time.
//
// <Atom> = [ <UnitClause> ] { <TopNode> }
//
// All top nodes inside atom are listed in order they appear in source code.
type Atom struct {
	Header AtomHeader

	// List of top custom type definition nodes.
	Types []TopType

	// List of top constant definition nodes.
	Lets []TopLet

	// List of top variable definition nodes.
	Vars []TopVar

	// List of top function definition nodes.
	Funs []TopFun

	// List of top function declaration nodes.
	Decs []TopDec

	// List of method nodes.
	Meds []Method
}

// AtomHeader stores info about atom that affects build-time decisions.
// Includes:
//
//   - unit clause (explicit unit name + build props)
//   - units imported in this atom
type AtomHeader struct {
	Imports AtomImports

	// Can be nil in case unit clause is not present.
	Unit *UnitClause
}

type AtomImports struct {
	Blocks []ImportBlock

	// Preprocessed data from all import blocks.
	// Builder uses this list to access all imports at once.
	Paths []origin.Path
}
