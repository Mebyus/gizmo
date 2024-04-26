package format

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

type Builder struct {
	s Stapler
}

func New() *Builder {
	return &Builder{}
}

func (g *Builder) Format(atom ast.Atom, tokens []token.Token) []byte {
	for _, top := range atom.Nodes {
		g.TopLevel(top)
	}
	return g.s.staple()
}
