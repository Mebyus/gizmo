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

	Roots     []int
	Pinnacles []int

	// Symbol map. Maps symbol to its node index.
	sm map[*Symbol]int

	// keeps track on number of steps happened during traversal
	step int

	stack TypeGraphStack

	// Stores discovery step of visited nodes.
	//
	// Equals 0 for nodes which are not yet visited.
	disc []int

	// Earliest visited node (the node with minimum
	// discovery time) that can be reached from
	// subtree rooted with current node.
	low []int
}

type TypeGraphStack struct {
	// used to check whether or not node is present the stack (by node index)
	m []bool

	// stack elements, each element is a node index
	s []int
}

func (s *TypeGraphStack) Init(size int) {
	s.m = make([]bool, size)
	s.s = make([]int, 0, size/2) // prealloc some space, to reduce number of reallocations
}

func (s *TypeGraphStack) Has(i int) bool {
	return s.m[i]
}

func (s *TypeGraphStack) Push(i int) {
	s.s = append(s.s, i)

	// mark added element as present in stack
	s.m[i] = true
}

// Pop removes element stored on top of the stack and returns it
// to the caller.
func (s *TypeGraphStack) Pop() int {
	i := s.Top()

	// shrink stack by one element, but keep underlying space
	// for future use
	s.s = s.s[:s.tip()]

	// mark removed element as no longer present in stack
	s.m[i] = false

	return i
}

// Top returns element stored on top of the stack. Does not alter
// stack state.
func (s *TypeGraphStack) Top() int {
	return s.s[s.tip()]
}

func (s *TypeGraphStack) tip() int {
	return len(s.s) - 1
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

	if c == 0 {
		// all nodes are stray nodes
		return
	}
	if c == 1 {
		panic("graph with connected nodes should have at least 2 nodes")
	}

	nodes := make([]TypeGraphNode, 0, c)
	var roots []int
	var pinnacles []int
	// sm := make(map[*Symbol]int, c)

	for i := 0; i < len(g.Nodes); i += 1 {
		if stray[i] {
			continue
		}

		node := g.Nodes[i]
		j := remap[i] // node's new index
		anc := remapLinks(stray, remap, node.Anc)
		des := remapLinks(stray, remap, node.Des)

		if len(anc) == 0 {
			roots = append(roots, j)
		}
		if len(des) == 0 {
			pinnacles = append(pinnacles, j)
		}

		nodes = append(nodes, TypeGraphNode{
			Anc:      anc,
			Des:      des,
			Index:    j,
			Sym:      g.Nodes[i].Sym,
			SelfLoop: g.Nodes[i].SelfLoop,
		})
	}

	if len(roots) == 0 || len(pinnacles) == 0 {
		fmt.Println("detected cluster inside graph")
	}

	// TODO: remove debug print
	for _, n := range nodes {
		var anc []string
		for _, l := range n.Anc {
			s := nodes[l.Index].Sym.Name
			if l.Kind == linkIndirect {
				s = "*" + s
			}

			anc = append(anc, s)
		}

		fmt.Printf("%s: %v\n", n.Sym.Name, anc)
	}

	g.Nodes = nodes
	g.Roots = roots
	g.Pinnacles = pinnacles
	g.sm = nil

	g.Traverse()
}

func remapLinks(filter []bool, remap []int, links []TypeGraphLink) []TypeGraphLink {
	var s []TypeGraphLink
	for _, l := range links {
		if filter[l.Index] {
			continue
		}
		s = append(s, TypeGraphLink{
			Index: remap[l.Index],
			Kind:  l.Kind,
		})
	}
	return s
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

// Traverse implements Tarjan’s algorithm for searching Strongly Connected Components
// inside directed graph.
func (g *TypeGraphBuilder) Traverse() {
	g.disc = make([]int, len(g.Nodes))
	g.low = make([]int, len(g.Nodes))

	g.stack.Init(len(g.Nodes))

	for _, i := range g.Roots {
		// TODO: find a way to separate
		// isolated components from each other
		// while traversing the graph
		g.traverse(i)
	}

	// TODO: rank nodes
	// In order to do that we need to start from
	// roots and rank-by-rank eliminate connections
	// from descendant's ancestors
	//
	// after each such step there will be a new set of
	// "roots" with no ancestors, if that is not the case,
	// then we encountered one of the clusters
	//
	// we can check list of nodes inside this cluster
	// from what we obtained in the previous step,
	// clusters are ranked as a whole
}

// recursive depth-first traverse
func (g *TypeGraphBuilder) traverse(n int) {
	g.step += 1
	g.disc[n] = g.step
	g.low[n] = g.step
	g.stack.Push(n)

	for _, l := range g.Nodes[n].Des {
		i := l.Index

		if g.disc[i] == 0 {
			// if node is not yet visited, traverse its subtree
			g.traverse(i)

			// after subtree traversal current node
			// should have the lowest low discovery step
			// of all its descendant nodes
			g.low[n] = min(g.low[n], g.low[i])
		} else if g.stack.Has(i) {
			// this node is already present in stack,
			// thus forming a cycle, we must update
			// low discovery step of subtree start
			g.low[n] = min(g.low[n], g.disc[i])
		}
	}

	if g.low[n] == g.disc[n] {
		// we found head node of the cluster,
		// pop the stack until reaching head

		var list []int
		i := g.stack.Pop()
		list = append(list, i)
		for i != n {
			i = g.stack.Pop()
			list = append(list, i)
		}
		fmt.Printf("cluster (%d): %v\n", g.low[n], list)
	}
}
