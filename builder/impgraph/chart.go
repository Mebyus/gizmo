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

		anc, err := g.mapAncestors(node.Bud.Ancestors())
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

func (g *Graph) mapAncestors(uids []string) ([]int, error) {
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

// Chart traverses graph to fill rank information and detect cycles.
// Upon detecting first cycle charting stops and method returns
// detected cycle. Returns nil if chart was complete and no cycles
// were detected
func (g *Graph) Chart() *Cycle {
	return nil
}
