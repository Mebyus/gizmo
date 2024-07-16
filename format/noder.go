package format

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

type Noder struct {
	tokens []token.Token
	nodes  []Node
}

func NewNoder(tokens []token.Token) *Noder {
	// TODO: prealloc space for nodes
	return &Noder{tokens: tokens}
}

func (g *Noder) Nodes(atom *ast.Atom) []Node {
	for _, top := range atom.Nodes {
		g.TopLevel(top)
		g.blank()
	}
	return g.nodes
}
