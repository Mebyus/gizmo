package tt

import (
	"fmt"
	"sort"
)

// TypeGraphLink represents a link between two graph nodes.
type TypeGraphLink struct {
	// Index of connected node inside graph's list of nodes.
	Index int

	// Describes how link was formed: through direct or indirect inclusion.
	Kind TypeLinkKind
}

type TypeGraphNode struct {
	// Ancestor links between symbols. This slice is created during
	// initial graph construction. Graph links between nodes are created
	// based on this information.
	Links []TypeLink

	// List of ancestor nodes indices. For root nodes this list is always empty.
	//
	// These nodes correspond to symbols used by this node's symbol.
	Anc []TypeGraphLink

	// List of descendant nodes indices. For pinnacle nodes this list is always empty.
	//
	// These nodes correspond to symbols which use this node's symbol.
	Des []TypeGraphLink

	// Graph roots have rank of zero. Each descent step increases rank by one.
	// Thus all non-root nodes have positive rank value.
	Rank int

	// Node index inside graph's list of nodes.
	Index int

	// Symbol which defines named type attached to this node.
	// Contains payload information which does not affect graph structure.
	Sym *Symbol

	// If true means that symbol has indirect link to itself among
	// its ancestors.
	SelfLoop bool
}

type TypeGraphStrayNode struct {
	Sym      *Symbol
	SelfLoop bool
}

type TypeGraphBuilder struct {
	Stray []TypeGraphStrayNode
	Nodes []TypeGraphNode

	Roots []int

	// Symbol map. Maps symbol to its node index.
	sm map[*Symbol]int
}

func NewTypeGraphBuilder(size int) *TypeGraphBuilder {
	return &TypeGraphBuilder{
		Nodes: make([]TypeGraphNode, 0, size),

		sm: make(map[*Symbol]int, size),
	}
}

// Add a symbol with the list of its ancestor links.
func (g *TypeGraphBuilder) Add(s *Symbol, links []TypeLink) {
	if s == nil {
		panic("nil symbol")
	}
	_, ok := g.sm[s]
	if ok {
		panic(fmt.Sprintf("duplicate symbol: %s", s.Name))
	}

	index := len(g.Nodes)
	g.sm[s] = index
	g.Nodes = append(g.Nodes, TypeGraphNode{
		Links: links,
		Index: index,
		Sym:   s,
	})
}

func (g *TypeGraphBuilder) Scan() {
	for i := 0; i < len(g.Nodes); i += 1 {
		anc := g.mapAncestors(i)
		g.Nodes[i].Anc = anc

		if len(anc) == 0 {
			// found root node
			// TODO: maybe it's not needed here due to later remapping
			g.Roots = append(g.Roots, i)
		}

		for _, l := range anc {
			g.Nodes[l.Index].Des = append(g.Nodes[l.Index].Des, TypeGraphLink{
				Index: i,
				Kind:  l.Kind,
			})
		}
	}

	// find stray nodes and remap the graph
	stray := make([]bool, len(g.Nodes)) // true for stray node's index
	remap := make([]int, len(g.Nodes))  // maps old node indices to new ones (after prunning is done)
	c := 0                              // how many non-stray nodes were discovered
	for i := 0; i < len(g.Nodes); i += 1 {
		node := g.Nodes[i]

		if len(node.Anc) == 0 && len(node.Des) == 0 {
			stray[i] = true
			g.Stray = append(g.Stray, TypeGraphStrayNode{
				Sym:      node.Sym,
				SelfLoop: node.SelfLoop,
			})
			// TODO: remove this debug print
			fmt.Printf("%s: stray type symbol \"%s\"\n", node.Sym.Pos.String(), node.Sym.Name)
		} else {
			remap[i] = c
			c += 1
		}
	}

}

func (g *TypeGraphBuilder) mapAncestors(node int) []TypeGraphLink {
	links := g.Nodes[node].Links
	if len(links) == 0 {
		return nil
	}

	anc := make([]TypeGraphLink, 0, len(links))
	for _, l := range links {
		if l.Symbol == nil {
			panic("nil symbol")
		}
		if l.Kind == linkEmpty {
			panic("empty link")
		}

		i, ok := g.sm[l.Symbol]
		if !ok {
			panic(fmt.Sprintf("unknown ancestor symbol: %s", l.Symbol.Name))
		}

		if i == node {
			g.Nodes[node].SelfLoop = true
		} else {
			anc = append(anc, TypeGraphLink{
				Index: i,
				Kind:  l.Kind,
			})
		}
	}

	if len(anc) == 0 {
		return nil
	}
	if len(anc) == 1 {
		return anc
	}

	sort.Slice(anc, func(i, j int) bool {
		return anc[i].Index < anc[j].Index
	})
	for j := 1; j < len(anc); j += 1 {
		if anc[j-1].Index == anc[j].Index {
			i := anc[j].Index
			name := g.Nodes[i].Sym.Name
			panic(fmt.Sprintf("duplicate ancestor link: %s (i=%d)", name, i))
		}
	}
	return anc
}
