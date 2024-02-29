package impgraph

type Node struct {
	// List of ancestor nodes indices. For root nodes this list is always empty
	//
	// These nodes correspond to units imported by this node's unit
	Anc []int

	// List of descendant nodes indices. For pinnacle nodes this list is always empty
	//
	// These nodes correspond to units which import this node's unit
	Des []int

	// Graph roots have rank of zero. Each descent step increases rank by one.
	// Thus all non-root nodes have positive rank value
	Rank int

	// Node index inside graph's list of nodes
	Index int

	// Data attached to this node. Contains payload information which does not
	// affect graph structure
	Bud Bud
}

type Bud interface {
	// Gives unique (among all possible data values) string identifier
	UID() string

	// List all ancestor identifiers for this data value
	Ancestors() []string
}

type Graph struct {
	Nodes []Node

	// Each cohort contains list of node indices with the same rank.
	// Cohort index in this slice directly corresponds to that rank
	Cohorts [][]int

	// List of root nodes indices
	Roots []int

	// List of pinnacle nodes indices
	Pinnacles []int

	// Bud map. Maps bud uid to node index
	bm map[string]int
}

// Create a new blank graph with enough preallocated space to hold
// specified number of nodes
func New(size int) *Graph {
	return &Graph{
		Nodes: make([]Node, 0, size),

		bm: make(map[string]int, size),
	}
}

// Cycle contains information about node cycle inside graph
type Cycle struct {
	// Always has at least two nodes
	Nodes []Node
}

// Sort shift nodes (by cycling them inside the list) in such a way
// that node with the smallest index becomes first in the list of cycle's nodes
func (c *Cycle) Sort() {
	// find node with the smallest graph idnex
	var j int // index of the smallest node inside cycle's nodes slice
	m := c.Nodes[0].Index
	for i, node := range c.Nodes {
		if node.Index < m {
			m = node.Index
			j = i
		}
	}

	c.Shift(-j)
}

// Shift nodes in cyclic manner by specified number of positions.
// If argument is positive shift to the right, if it's negative
// shift to the left, if it's zero of a multiple of number of nodes
// do nothing
func (c *Cycle) Shift(n int) {
	if n % len(c.Nodes) == 0 {
		return
	}

	if n < 0 {
		
	}
}
