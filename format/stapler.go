package format

import "github.com/mebyus/gizmo/token"

type Stapler struct {
	// read-only input field
	tokens []token.Token

	// accumulator output field
	nodes []Node
}

func NewStapler(tokens []token.Token) *Stapler {
	return &Stapler{tokens: tokens}
}

func (s *Stapler) Nodes() []Node {
	return s.nodes
}

func Staple(tokens []token.Token, nodes []Node) []Node {
	s := NewStapler(tokens)

	for _, node := range nodes {
		switch node.Kind {
		default:
			s.add(node)
		}
	}

	return s.Nodes()
}

func (s *Stapler) add(node Node) {
	s.nodes = append(s.nodes, node)
}
