package tt

type TypeGraphNode struct {
	// List of ancestor nodes indices. For root nodes this list is always empty.
	//
	// These nodes correspond to symbols used by this node's symbol.
	Anc []int

	// List of descendant nodes indices. For pinnacle nodes this list is always empty.
	//
	// These nodes correspond to symbols which use this node's symbol.
	Des []int

	// Graph roots have rank of zero. Each descent step increases rank by one.
	// Thus all non-root nodes have positive rank value.
	Rank int

	// Node index inside graph's list of nodes.
	Index int

	// Symbol which defines named type attached to this node.
	// Contains payload information which does not affect graph structure.
	Sym *Symbol
}

type TypeGraphBuilder struct {
}
