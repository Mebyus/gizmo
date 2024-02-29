package impgraph

import (
	"fmt"
	"sort"
)

func (g *Graph) Add(bud Bud) error {
	uid := bud.UID()
	if uid == "" {
		return fmt.Errorf("empty node uid (index=%d)", len(g.Nodes))
	}
	_, ok := g.bm[uid]
	if ok {
		return fmt.Errorf("duplicate node uid: \"%s\"", uid)
	}

	index := len(g.Nodes)
	g.bm[uid] = index
	g.Nodes = append(g.Nodes, Node{
		Index: index,
		Bud:   bud,
	})
	return nil
}

// Scan performs preliminary validation of all stored nodes, fill slices with
// node connections (ancestors and descendants) and find roots
//
// Must be called after all graph nodes were added via Add method
func (g *Graph) Scan() error {
	if len(g.Nodes) == 0 {
		return fmt.Errorf("graph does not contain nodes")
	}

	for i := 0; i < len(g.Nodes); i++ {
		node := g.Nodes[i]

		anc, err := g.mapAncestors(i, node.Bud.Ancestors())
		if err != nil {
			return fmt.Errorf("scan node (uid=\"%s\"): %w", node.Bud.UID(), err)
		}

		// refer to node by index because we want to mutate it
		g.Nodes[i].Anc = anc

		if len(anc) == 0 {
			// found root node
			g.Roots = append(g.Roots, i)
		}

		// Fill descendants based on the list of ancestors of current node
		//
		// If node A is in the list of ancestors of node B than by definition it means
		// that A has node B as one of its descendants
		for _, j := range anc {
			// refer to node by index because we want to mutate it
			g.Nodes[j].Des = append(g.Nodes[j].Des, i)
		}
	}

	// By the nature of loop above we do not need to sort descendants slice (Node.Des)
	// in each node, because descendants were added exactly in ascending order to
	// each such slice

	// Find pinnacles and make quick check for stray nodes
	//
	// Special case when a stray node is ok
	if len(g.Nodes) == 1 {
		g.Pinnacles = append(g.Pinnacles, 0)
		return nil
	}

	// Sweep nodes in general case
	for i := 0; i < len(g.Nodes); i++ {
		node := g.Nodes[i]

		if len(node.Anc) == 0 && len(node.Des) == 0 {
			return fmt.Errorf("stray node (uid=\"%s\")", node.Bud.UID())
		}

		if len(node.Des) == 0 {
			// found pinnacle node
			g.Pinnacles = append(g.Pinnacles, i)
		}
	}

	return nil
}

func (g *Graph) mapAncestors(node int, uids []string) ([]int, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	s := make([]int, 0, len(uids))
	for _, uid := range uids {
		if uid == "" {
			return nil, fmt.Errorf("empty ancestor uid")
		}
		i, ok := g.bm[uid]
		if !ok {
			return nil, fmt.Errorf("unknown ancestor uid: \"%s\"", uid)
		}
		if i == node {
			return nil, fmt.Errorf("ancestor points to self")
		}
		s = append(s, i)
	}
	if len(s) == 1 {
		return s, nil
	}

	sort.Ints(s)
	for j := 1; j < len(s); j++ {
		if s[j-1] == s[j] {
			i := s[j]
			uid := g.Nodes[i].Bud.UID()
			return nil, fmt.Errorf("duplicate ancestor uid: \"%s\"", uid)
		}
	}
	return s, nil
}

type ScoutPos struct {
	// node node in Graph.Nodes slice
	node int

	// index of next descendant in Node.Des slice,
	// to clarify: this value is not a Node index
	next int
}

// Scout helper object for graph charting
type Scout struct {
	// stores scout's path inside Graph
	stack []ScoutPos

	// Stack map
	//
	// maps Node index to its index in stack (only if Node is present in stack)
	sm map[int]int

	// indicates whether a Node was already visited or not
	//
	// element at specific index always corresponds to Node with
	// the same index
	visited []bool

	g *Graph
}

// NewScout creates new scout with enough space to handle graph with
// specified number of nodes
func NewScout(graph *Graph) *Scout {
	return &Scout{
		g: graph,

		sm:      make(map[int]int),
		visited: make([]bool, len(graph.Nodes)),
	}
}

func (s *Scout) isEmpty() bool {
	return len(s.stack) == 0
}

// length of stored path
func (s *Scout) plen() int {
	return len(s.stack)
}

// index of the element at the top of the stack
func (s *Scout) top() int {
	return s.plen() - 1
}

// push node's index onto the stack
func (s *Scout) push(i int) {
	s.visited[i] = true
	s.sm[i] = s.plen()
	s.stack = append(s.stack, ScoutPos{node: i})
}

func (s *Scout) pop() {
	tip := s.tip()
	delete(s.sm, tip.node)

	// shrink stack, but keep underlying memory
	s.stack = s.stack[:s.top()]
}

func (s *Scout) tip() ScoutPos {
	return s.stack[s.top()]
}

// shorthand for retrieving node by its graph index
func (s *Scout) node(i int) Node {
	return s.g.Nodes[i]
}

func (s *Scout) next() int {
	tip := s.tip()
	next := tip.next
	tip.next++
	s.stack[s.top()] = tip
	return next
}

type StepKind uint8

const (
	// scout descends, its path length increases
	descend StepKind = iota

	// scout ascends, its path length decreases
	ascend

	// scout found cycle in its path
	cycle
)

type ScoutStep struct {
	// meaning depends on Kind
	//
	//	descend: next Node index
	//	ascend: ignored
	//	cycle: stack index of cycle start in scout's path
	Val int

	Kind StepKind
}

func (s *Scout) step() ScoutStep {
	tip := s.tip()
	des := s.node(tip.node).Des

	if len(des) == 0 {
		return ScoutStep{Kind: ascend}
	}

	j := s.next()
	for j < len(des) {
		// next node's index in path chosen among current node's descendants
		next := des[j]
		if !s.visited[next] {
			return ScoutStep{
				Val:  next,
				Kind: descend,
			}
		}

		k, ok := s.sm[next]
		if ok {
			// cycle found
			return ScoutStep{
				Val:  k,
				Kind: cycle,
			}
		}

		j = s.next()
	}

	return ScoutStep{Kind: ascend}
}

// traverse Graph starting from Node with specified index
func (s *Scout) traverse(node int) *Cycle {
	s.push(node)

	for !s.isEmpty() {
		step := s.step()

		switch step.Kind {
		case descend:
			s.push(step.Val)
		case ascend:
			s.pop()
		case cycle:
			return s.cycle(step.Val)
		default:
			panic(step.Kind)
		}
	}

	return nil
}

// gathers nodes in found cycle from scout's path
//
// argument is a stack index of cycle start, it is assumed that
// nodes up until after top of the stack form the cycle
func (s *Scout) cycle(p int) *Cycle {
	num := s.top() - p
	if num <= 1 {
		panic("not enough nodes to form a cycle")
	}
	nodes := make([]Node, 0, num)

	// iterate over scout path elements
	for i := p; i < s.plen(); i++ {
		// graph index of node in cycle
		n := s.stack[i].node

		nodes = append(nodes, s.node(n))
	}

	c := &Cycle{Nodes: nodes}
	c.Sort()
	return c
}

func (s *Scout) Traverse() *Cycle {
	for _, root := range s.g.Roots {
		c := s.traverse(root)
		if c != nil {
			return c
		}
	}

	for i := 0; i < len(s.visited); i++ {
		if !s.visited[i] {
			// This can only happen if there are isolated cycles
			// in graph. Such cycles are unreachable from roots
			c := s.traverse(i)
			if c == nil {
				panic(fmt.Sprintf("isolated nodes must yield a cycle (node=%d)", i))
			}
			return c
		}
	}

	return nil
}

// Chart traverses graph to fill rank information and detect cycles.
// Upon detecting first cycle charting stops and method returns
// detected cycle. Returns nil if chart was complete and no cycles
// were detected
func (g *Graph) Chart() *Cycle {
	return nil
}
