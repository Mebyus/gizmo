package impgraph

import "fmt"

func (g *Graph) Add(bud Bud) error {
	uid := bud.UID()
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

// Chart traverses graph to fill rank information and detect cycles.
// Upon detecting first cycle charting stops and method returns
// detected cycle. Returns nil if chart was complete and no cycles
// were detected
func (g *Graph) Chart() *Cycle {
	return nil
}
