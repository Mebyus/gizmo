package ast

import (
	"strings"

	"github.com/mebyus/gizmo/source"
)

// <Identifier> = word
type Identifier struct {
	Pos source.Pos
	Lit string
}

func (n Identifier) String() string {
	if len(n.Lit) == 0 {
		return "<nil>"
	}

	return n.Lit
}

// <ScopedIdentifier> = <Identifier> { "::" <Identifier> }
type ScopedIdentifier struct {
	// Can be nil if <ScopedIdentifier> represents regular <Identifier>
	Scopes []Identifier

	// Last name in chain
	Name Identifier
}

func (s ScopedIdentifier) Pos() source.Pos {
	if len(s.Scopes) == 0 {
		return s.Name.Pos
	}
	return s.Scopes[0].Pos
}

func (s ScopedIdentifier) String() string {
	if len(s.Scopes) == 0 && len(s.Name.Lit) == 0 {
		return "<nil>"
	}

	ss := make([]string, 0, len(s.Scopes) + 1)
	for _, name := range s.Scopes {
		ss = append(ss, name.String())
	}
	ss = append(ss, s.Name.String())
	return strings.Join(ss, "::")
}
