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
	Nodes []Node
}
