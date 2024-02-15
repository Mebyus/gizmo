package ast

// <UnitBlock> = "unit" <UnitName> "{" { <Statement> } "}"
type UnitBlock struct {
	Name Identifier

	Block BlockStatement
}

// UnitAtom smallest piece of processed source code inside a unit. In most
// cases this represents a file with source code. Exceptions may include
// source code generated at compile time
//
// <UnitAtom> = [ <UnitBlock> ] { <Namespace> }
type UnitAtom struct {
	// Can be nil in case UnitBlock is not present
	Unit *UnitBlock

	// Saved in order they appear in source code
	Blocks []NamespaceBlock
}
