package ast

import "github.com/mebyus/gizmo/source/origin"

// <UnitBlock> = "unit" <UnitName> "{" { <Statement> } "}"
type UnitBlock struct {
	Props []Prop

	Name Identifier

	Block BlockStatement
}

// UnitAtom smallest piece of processed source code inside a unit. In most
// cases this represents a file with source code. Exceptions may include
// source code generated at compile time
//
// <UnitAtom> = [ <UnitBlock> ] { <Namespace> }
type UnitAtom struct {
	Header AtomHeader

	// Saved in order they appear in source code
	Blocks []NamespaceBlock
}

// AtomHeader stores info about atom that affects build-time decisions.
// Includes:
//
//   - unit clause (explicit unit name + build props)
//   - units imported in this atom
type AtomHeader struct {
	Imports AtomImports

	// Can be nil in case UnitClause is not present
	Unit *UnitBlock
}

type AtomImports struct {
	ImportBlocks []ImportBlock

	// preprocessed data from blocks
	ImportPaths []origin.Path
}
