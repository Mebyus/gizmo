package format

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

type Builder struct {
	tokens []token.Token
	nodes  []Node
}

func New(tokens []token.Token) *Builder {
	// TODO: prealloc space for nodes
	return &Builder{tokens: tokens}
}

func (g *Builder) Nodes(atom ast.Atom) []Node {
	for _, top := range atom.Nodes {
		g.TopLevel(top)
	}
	return g.nodes
}
