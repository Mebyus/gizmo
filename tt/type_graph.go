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

type TypeGraph struct {
	Nodes []TypeGraphNode

	// Stores all non-trivial (more than 1 node) components inside the graph.
	Comps []TypeGraphComponent

	// List of isolated node indices.
	Isolated []int
}

type TypeGraphBuilder struct {
	Nodes []TypeGraphNode

	// Stores all non-trivial (more than 1 node) components inside the graph.
	Comps []TypeGraphComponent

	// List of isolated node indices.
	Isolated []int

	// Symbol map. Maps symbol to its node index.
	sm map[*Symbol]int

	/* Internal state for components discovery */

	// Counter for assigning component number to nodes.
	comp int

	// Maximum number of vertices among components.
	maxCompSize int

	// List of node indices for current BFS scan.
	wave []int

	// List of node indices for next BFS scan.
	next []int

	// maps node index to component vertex index,
	// note that each component has its own indexing
	// for vertices
	remap []int
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

			if len(c.V) < 2 {
				panic("connected component must have at least 2 vertices")
			}

			for p := 0; p < len(c.V); p += 1 {
				i := c.V[p].Index

				c.V[p].Anc = remapLinks(g.remap, g.Nodes[i].Anc)
				c.V[p].Des = remapLinks(g.remap, g.Nodes[i].Des)
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

func (g *TypeGraphBuilder) Scan() *TypeGraph {
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
	g.discoverClusters()
	// TODO: scan clusters for direct-indirect correctness
	g.rank()

	return &TypeGraph{
		Nodes:    g.Nodes,
		Comps:    g.Comps,
		Isolated: g.Isolated,
	}
}

func (g *TypeGraphBuilder) rank() {
	var r TypeGraphRanker

	for k := 0; k < len(g.Comps); k += 1 {
		c := &g.Comps[k]
		if c.isTwoLevelNoCluster() {
			for _, i := range c.Pinnacles {
				c.V[i].Rank = 1
			}
			c.Cohorts = [][]int{c.Roots, c.Pinnacles}
			continue
		}

		if len(c.Clusters) != 0 {
			panic("recursive cluster types not implemented")
		}

		r.Rank(c)
	}
}

func (g *TypeGraphBuilder) discoverClusters() {
	if len(g.Comps) == 0 {
		return
	}

	if g.maxCompSize < 2 {
		panic("connected components are present, but max component size is less than 2 vertices, which is impossible")
	}

	var w TypeGraphComponentWalker

	var k int

	// search for first component with possible clusters
	// to properly initialize walker
	for k < len(g.Comps) {
		c := &g.Comps[k]
		if c.isTwoLevelNoCluster() {
			k += 1
			continue
		}

		w.Init(g.maxCompSize, c)
		w.Walk()
		k += 1
		break
	}

	for k < len(g.Comps) {
		c := &g.Comps[k]
		if !c.isTwoLevelNoCluster() {
			w.Reset(c)
			w.Walk()
		}
		k += 1
	}
}

func remapLinks(remap []int, links []TypeGraphLink) []int {
	if len(links) == 0 {
		return nil
	}
	s := make([]int, 0, len(links))
	for _, l := range links {
		s = append(s, remap[l.Index])
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

// Walk implements Tarjan’s algorithm for searching Strongly Connected Components
// inside directed graph.
func (w *TypeGraphComponentWalker) Walk() {
	for _, i := range w.c.Roots {
		w.walk(i)
	}

	if w.step >= len(w.c.V) {
		// all vertices have been discovered
		return
	}

	// we need additional scan due to root cluster(s)
	// inside the component
	for i := 0; i < len(w.c.V); i += 1 {
		if w.disc[i] == 0 {
			w.walk(i)
		}
		if w.step >= len(w.c.V) {
			// all vertices have been discovered
			return
		}
	}
}

// recursive depth-first walk
func (w *TypeGraphComponentWalker) walk(v int) {
	w.step += 1
	w.disc[v] = w.step
	w.low[v] = w.step
	w.stack.Push(v)

	for _, i := range w.c.V[v].Des {
		if w.disc[i] == 0 {
			// if vertex is not yet visited, traverse its subtree
			w.walk(i)

			// after subtree traversal current vertex
			// should have the lowest low discovery step
			// of all its descendant vertices
			w.low[v] = min(w.low[v], w.low[i])
		} else if w.stack.Has(i) {
			// this vertex is already present in stack,
			// thus forming a cycle, we must update
			// low discovery step of subtree start
			w.low[v] = min(w.low[v], w.disc[i])
		}
	}

	if w.low[v] == w.disc[v] {
		// we found head vertex of the cluster,
		// pop the stack until reaching head

		i := w.stack.Pop()
		if i == v {
			// do not keep track of trivial clusters (that contains one vertex)
			return
		}

		// cluster number of newly discovered cluster
		num := len(w.c.Clusters) + 1
		w.c.V[i].Cluster = num

		// cluster has at least 2 vertices by definition
		list := make([]int, 0, 2)
		list = append(list, i)
		for i != v {
			i = w.stack.Pop()
			w.c.V[i].Cluster = num
			list = append(list, i)
		}
		sort.Ints(list)
		w.c.Clusters = append(w.c.Clusters, list)
		fmt.Printf("cluster (%d): %v\n", w.low[v], list)
	}
}

type TypeGraphVertex struct {
	// list of ancestor indices inside V
	Anc []int

	// list of descendant indices inside V
	Des []int

	// Graph roots have rank of zero. Each descent step increases rank by one.
	// Thus all non-root nodes have positive rank value.
	Rank int

	// original node index
	Index int

	// Cluster number in which this vertex resides.
	// Equals 0 if vertex does not belong to cluster.
	Cluster int
}

type TypeGraphComponent struct {
	V []TypeGraphVertex

	// list of indices inside V
	Roots []int

	// list of indices inside V
	Pinnacles []int

	// each element in this slice is a list of non-trivial vertices
	// which belong to the same cluster
	Clusters [][]int

	Cohorts [][]int

	// component number
	Num int
}

func (c *TypeGraphComponent) isTwoLevelNoCluster() bool {
	// component has 2 levels and no cycles
	return len(c.V) == len(c.Roots)+len(c.Pinnacles)
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

type TypeGraphRanker struct {
	// Indicates how many ancestors are still unranked for a vertex with
	// particular index, directly corresponding to index in this slice.
	// If vertex was already ranked corresponding left value will be 0
	left []int

	// list of vertex indices which will be ranked
	// during current iteration
	wave []int

	// buffer for preparing next wave
	next []int

	c *TypeGraphComponent
}

func (r *TypeGraphRanker) Rank(c *TypeGraphComponent) {
	r.c = c
	r.left = r.left[:0]
	for i := 0; i < len(c.V); i += 1 {
		r.left = append(r.left, len(c.V[i].Anc))
	}
	r.wave = r.wave[:0]
	r.next = r.next[:0]

	r.wave = append(r.wave, c.Roots...)

	// components that are being ranked by this method
	// have at least 3 cohorts
	c.Cohorts = make([][]int, 0, 3)
	c.Cohorts = append(c.Cohorts, c.Roots)

	r.rank()
}

func (r *TypeGraphRanker) swap() {
	r.wave, r.next = r.next, r.wave
	r.next = r.next[:0]
}

// add node with specified graph index to cohort of the specified rank
func (r *TypeGraphRanker) add(node int, rank int) {
	for j := len(r.c.Cohorts); j <= rank; j++ {
		// allocate place for storing slices of graph indices
		// for cohorts with rank not initialized previously
		r.c.Cohorts = append(r.c.Cohorts, nil)
	}

	r.c.Cohorts[rank] = append(r.c.Cohorts[rank], node)
}

func (r *TypeGraphRanker) rank() {
	for len(r.wave) != 0 {
		for _, i := range r.wave {
			waiters := r.c.V[i].Des
			if len(waiters) == 0 {
				continue
			}

			// rank that will be passed to waiters
			rank := r.c.V[i].Rank + 1

			for _, j := range waiters {
				r.left[j] -= 1

				if rank > r.c.V[j].Rank {
					// select highest rank from all nodes inside the wave
					r.c.V[j].Rank = rank
				}

				// check if waiter node has finished ranking
				if r.left[j] == 0 {
					r.add(j, r.c.V[j].Rank)

					// next wave is constructed from nodes that finished
					// ranking during this wave
					r.next = append(r.next, j)
				}
			}
		}

		r.swap()
	}

	for _, cohort := range r.c.Cohorts[1:] {
		sort.Ints(cohort)
	}
}
