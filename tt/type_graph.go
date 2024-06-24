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

	// List of adjacent node indices. This list is obtained by turning
	// digraph (directed graph) into ugraph (undirected graph). In other
	// words this is a list of node indices merged from ancestors and
	// descendants.
	Adj []int

	// Graph roots have rank of zero. Each descent step increases rank by one.
	// Thus all non-root nodes have positive rank value.
	Rank int

	// Component number. Each distinct number marks connected isolated component
	// within a graph. By definition nodes from different components do not
	// have connections between them. More formally:
	//
	//	if node A belongs to c1 (component 1) and node B belongs to c2 then
	//	by definition there is no edge A -> B or B -> A in this graph
	//
	// Components separate all graph nodes into equivalence classes with no
	// intersections between them.
	Comp int

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

	// Stores all non-trivial (more than 1 node) components inside the graph.
	Comps []TypeGraphComponent

	Roots     []int
	Pinnacles []int

	// List of isolated node indices.
	Isolated []int

	// Symbol map. Maps symbol to its node index.
	sm map[*Symbol]int

	/* Internal state for components discovery */

	// Counter for assigning component number to nodes.
	comp int

	// Maximum number of vertices among components.
	maxCompSize int

	// True for already visited node during BFS.
	// vis []bool

	// List of node indices for current BFS scan.
	wave []int

	// List of node indices for next BFS scan.
	next []int

	// maps node index to component vertex index,
	// note that each component has its own indexing
	// for vertices
	remap []int

	/* Internal state for clusters discovery */

	// keeps track on number of steps happened during traversal
	step int

	stack TypeGraphStack

	// Stores discovery step number of visited nodes.
	//
	// Equals 0 for nodes which are not yet visited.
	disc []int

	// Earliest visited node (the node with minimum
	// discovery step number) that can be reached from
	// subtree rooted with current node.
	low []int
}

type TypeGraphStack struct {
	// stack elements, each element is a node index
	s []int

	// used to check whether or not node is present the stack (by node index)
	m []bool
}

func (s *TypeGraphStack) Init(size int) {
	s.m = make([]bool, size)
	s.s = make([]int, 0, size/2) // prealloc some space, to reduce number of reallocations
}

func (s *TypeGraphStack) Reset() {
	s.s = s.s[:0]
	clear(s.m)
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
		// vis: make([]bool, size),
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

// MergeInts merges two sorted (in ascending order) slices into
// a single sorted slice with unique values. Each given slice must
// contain only unique values, but each value could be present in
// both slices. If one of the given slices is empty then other is returned
// as a result, thus avoiding unnecessary copying. The result is intended
// to be read-only.
//
//	[1, 2, 3, 4] + [2, 3, 5]  => [1, 2, 3, 4, 5]
//	[1, 1, 3] + [2]           => incorrect input: first slice contains duplicates
func MergeInts(a, b []int) []int {
	if len(a) == 0 && len(b) == 0 {
		return nil
	}
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}

	s := make([]int, 0, max(len(a), len(b)))

	i := 0
	j := 0

	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			s = append(s, a[i])
			i += 1
		} else if a[i] > b[j] {
			s = append(s, b[j])
			j += 1
		} else {
			// a[i] == b[j]
			s = append(s, a[i])
			i += 1
			j += 1
		}
	}

	for i < len(a) {
		s = append(s, a[i])
		i += 1
	}

	for j < len(b) {
		s = append(s, b[j])
		j += 1
	}

	return s
}

func linksToIndices(links []TypeGraphLink) []int {
	if len(links) == 0 {
		return nil
	}
	s := make([]int, 0, len(links))
	for _, l := range links {
		s = append(s, l.Index)
	}
	return s
}

func mergeLinks(a, b []TypeGraphLink) []int {
	return MergeInts(linksToIndices(a), linksToIndices(b))
}

// split graph into connected components
func (g *TypeGraphBuilder) split() {
	for i := 0; i < len(g.Nodes); i += 1 {
		// during first pass we construct adjacent links
		// for all nodes

		n := g.Nodes[i]
		adj := mergeLinks(n.Anc, n.Des)
		g.Nodes[i].Adj = adj

		// mark isolated nodes as separate component
		if len(adj) == 0 {
			g.comp += 1
			g.Nodes[i].Comp = g.comp
			g.Isolated = append(g.Isolated, i)

			// later we will use Node.Comp == 0 for checking
			// if a node was not yet visited during components BFS
		}
	}

	// g.comp denotes total number of created components,
	// hence here it represents number of isolated nodes
	if g.comp >= len(g.Nodes) {
		// all nodes are isolated
		// we do not need second pass with BFS
		return
	}

	g.remap = make([]int, len(g.Nodes))
	for i := 0; i < len(g.Nodes); i += 1 {
		if g.Nodes[i].Comp == 0 {
			// this node was not yet visited
			g.comp += 1
			g.bfs(i)
		}

		// TODO: remove debug print
		n := g.Nodes[i]
		fmt.Printf("%s (%d): comp=%v\n", n.Sym.Name, i, n.Comp)
	}
}

func (g *TypeGraphBuilder) bfs(n int) {
	// create new component and store its index
	k := len(g.Comps)
	g.Comps = append(g.Comps, TypeGraphComponent{Num: g.comp})
	c := &g.Comps[k]

	// reset the slice from previous BFS, but keep
	// underlying memory
	g.wave = g.wave[:0]
	g.wave = append(g.wave, n)
	g.Nodes[n].Comp = g.comp
	for {
		for _, i := range g.wave {
			l := len(c.V)
			if len(g.Nodes[i].Anc) == 0 {
				c.Roots = append(c.Roots, l)
			}
			if len(g.Nodes[i].Des) == 0 {
				c.Pinnacles = append(c.Pinnacles, l)
			}
			g.remap[i] = l
			c.V = append(c.V, TypeGraphVertex{Index: i})

			adj := g.Nodes[i].Adj
			for _, j := range adj {
				if g.Nodes[j].Comp == 0 {
					g.Nodes[j].Comp = g.comp
					g.next = append(g.next, j)
				}
			}
		}

		if len(g.next) == 0 {
			// we gathered all vertices that belong to current component
			// do vertex descendants remap before exiting

			for p := 0; p < len(c.V); p += 1 {
				i := c.V[p].Index
				var des []int
				if len(g.Nodes[i].Des) != 0 {
					des = make([]int, 0, len(g.Nodes[i].Des))
					for _, j := range g.Nodes[i].Des {
						des = append(des, g.remap[j.Index])
					}
				}
				c.V[p].Des = des

				// TODO: remove debug print
				fmt.Printf("vertex %d-%d (%d): %v\n", p, c.Num, i, des)
			}

			if len(c.V) > g.maxCompSize {
				g.maxCompSize = len(c.V)
			}

			return
		}

		g.wave, g.next = g.next, g.wave
		g.next = g.next[:0]

	}
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

	g.split()

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
		//
		// possibly we can return component's "generation"
		// number from traverse method
		//
		// if returned number is 0 then that means we did not
		// find connections with previous generaions and
		// should assign a new one on method exit
		//
		// at recursion start (in this scope) we should increment
		// next assigned generation number if returned value is 0
		// meaning we discovered new generation
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

type TypeGraphVertex struct {
	// list of indices inside V
	Des []int

	Rank int

	// original node index
	Index int
}

type TypeGraphComponent struct {
	V []TypeGraphVertex

	// list of indices inside V
	Roots []int

	// list of indices inside V
	Pinnacles []int

	// component number
	Num int
}

// Keeps track of internal state for clusters discovery
// inside graph component.
type TypeGraphComponentWalker struct {
	stack TypeGraphStack

	// Stores discovery step number of visited vertices.
	//
	// Equals 0 for vertices which are not yet visited.
	disc []int

	// Earliest visited vertex (the vertex with minimum
	// discovery step number) that can be reached from
	// subtree rooted with current vertex.
	low []int

	// component currently being processed
	c *TypeGraphComponent

	// keeps track on number of steps happened during traversal
	step int
}

func (w *TypeGraphComponentWalker) Init(size int, c *TypeGraphComponent) {
	w.disc = make([]int, size)
	w.low = make([]int, size)

	w.stack.Init(size)

	w.c = c
}

func (w *TypeGraphComponentWalker) Reset(c *TypeGraphComponent) {
	w.step = 0
	clear(w.disc)
	clear(w.low)
	w.stack.Reset()

	w.c = c
}
