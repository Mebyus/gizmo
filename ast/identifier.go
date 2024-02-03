package ast

import "github.com/mebyus/gizmo/source"

// <Identifier> = word
type Identifier struct {
	Pos source.Pos
	Lit string
}

// <ScopedIdentifier> = <Identifier> { "::" <Identifier> }
type ScopedIdentifier struct {
	// Can be nil if <ScopedIdentifier> represents regular <Identifier>
	Scopes []Identifier

	// Last name in chain
	Name Identifier
}
