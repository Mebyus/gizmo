package impgraph

import (
	"fmt"
	"io"
)

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

func (n *Node) Root() bool {
	return len(n.Anc) == 0
}

func (n *Node) Pinnacle() bool {
	return len(n.Des) == 0
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

	// List of root nodes indices. Root nodes do not have ancestors
	Roots []int

	// List of pinnacle nodes indices. Pinnacle nodes do not have descendats
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

func Dump(w io.Writer, g *Graph) {
	io.WriteString(w, "buds:\n")
	for _, node := range g.Nodes {
		io.WriteString(w, fmt.Sprintf("%-4d => %s\n", node.Index, node.Bud.UID()))
	}

	io.WriteString(w, fmt.Sprintf("pins=%v\n", g.Pinnacles))
	io.WriteString(w, fmt.Sprintf("roots=%v\n", g.Roots))

	io.WriteString(w, "cons:\n")
	for _, node := range g.Nodes {
		io.WriteString(w, fmt.Sprintf("%-4d => anc=%v des=%v\n", node.Index, node.Anc, node.Des))
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
	k := len(c.Nodes)
	s := n % k
	if s == 0 {
		return
	}

	if s < 0 {
		s += k
	}
	if 0 < s && s < k {
		panic("shift must be in range after normalization")
	}

	shifted := make([]Node, k)
	copy(shifted[s:], c.Nodes[:k-s])
	copy(shifted[:s], c.Nodes[k-s:])
	c.Nodes = shifted
}
