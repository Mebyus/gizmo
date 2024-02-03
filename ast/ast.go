package ast

import (
	"github.com/mebyus/gizmo/source"
)

type Node interface {
	source.Pin
}

// <UnitBlock> = "unit" <UnitName> "{" { <Statement> } "}"
type UnitBlock struct {
	Name Identifier

	Block BlockStatement
}

// Smallest piece of processed source code inside a unit. In most
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

// <ReturnStatement> = "return" <Expression> ";"
// type ReturnStatement struct {
// 	Expression Expression
// }

// <IfStatement> = <IfClause> [ <ElseClause> ]
type IfStatement struct {
	If   IfClause
	Else *ElseClause
}

// <IfClause> = "if" <Expression> <BlockStatement>
type IfClause struct {
	Condition Expression
	Body      BlockStatement
}

// <ElseClause> = "else" <BlockStatement>
type ElseClause struct {
	Body BlockStatement
}
