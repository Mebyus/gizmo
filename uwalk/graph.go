package uwalk

import (
	"fmt"
	"slices"
	"sort"

	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/stg"
)

type Node struct {
	// List of ancestor nodes indices. For root nodes this list is always empty.
	// Sorted by node index.
	//
	// These nodes correspond to units imported by this node's unit.
	Anc []int

	// List of descendant nodes indices. For pinnacle nodes this list is always empty.
	// Sorted by link node index.
	//
	// These nodes correspond to units which import this node's unit.
	Des []int

	Rank int

	// Unit carried by this node.
	Unit *stg.Unit
}

type Graph struct {
	// Sorted by unit index.
	Nodes []Node

	// List of node indices.
	Roots []int

	// Each cohort is a list of node indices.
	Cohorts [][]int
}

type Cycle struct {
	// Always has at least two nodes.
	Nodes []Node
}

// Fills Unit.Imports.Units according to import paths.
func (b *Bundle) mapUnits() {
	m := make(map[origin.Path]*stg.Unit, len(b.Units))
	b.Map = m

	for _, unit := range b.Units {
		m[unit.Path] = unit
	}

	for _, unit := range b.Units {
		unit.Imports.Units = make([]*stg.Unit, 0, len(unit.Imports.Paths))
		for _, p := range unit.Imports.Paths {
			u, ok := m[p]
			if !ok {
				panic(fmt.Sprintf("imported unit [%s] not found", p))
			}
			if u == unit {
				panic("unit imported itself")
			}
			unit.Imports.Units = append(unit.Imports.Units, u)
		}
	}
}

func listIndices(units []*stg.Unit) []int {
	if len(units) == 0 {
		return nil
	}

	list := make([]int, 0, len(units))
	for _, u := range units {
		list = append(list, int(u.Index))
	}
	return list
}

func (b *Bundle) makeGraphNodes() {
	b.Graph.Nodes = make([]Node, len(b.Units))

	for i, u := range b.Units {
		anc := listIndices(u.Imports.Units)

		b.Graph.Nodes[i].Unit = u
		b.Graph.Nodes[i].Anc = anc

		for _, j := range anc {
			b.Graph.Nodes[j].Des = append(b.Graph.Nodes[j].Des, i)
		}

		if len(anc) == 0 {
			b.Graph.Roots = append(b.Graph.Roots, i)
		}
	}
}

func rankGraphOrFindCycle(g *Graph) *Cycle {
	if len(g.Roots) == 0 {
		// TODO: algorithm for finding nodes in cycle
		return &Cycle{}
	}

	// number of nodes successfully ranked
	var total int

	// how many ancestors are still unranked for a node with
	// particular index
	left := make([]int, len(g.Nodes))
	for i, n := range g.Nodes {
		left[i] = len(n.Anc)
	}

	// nodes to scan in this wave
	wave := slices.Clone(g.Roots)
	g.Cohorts = make([][]int, 0, 2)
	g.Cohorts = append(g.Cohorts, g.Roots)

	// buffer for next wave
	var next []int

	for len(wave) != 0 {
		for _, i := range wave {
			waiters := g.Nodes[i].Des
			if len(waiters) == 0 {
				continue
			}

			// rank that will be passed to waiters
			rank := g.Nodes[i].Rank + 1

			for _, j := range waiters {
				left[j] -= 1

				if rank > g.Nodes[j].Rank {
					// select highest rank from all nodes inside the wave
					g.Nodes[j].Rank = rank
				}

				// check if waiter node has finished ranking
				if left[j] == 0 {
					for k := len(g.Cohorts); k <= rank; k += 1 {
						g.Cohorts = append(g.Cohorts, nil)
					}
					g.Cohorts[rank] = append(g.Cohorts[rank], j)

					// next wave is constructed from nodes that finished
					// ranking during this wave
					next = append(next, j)
				}
			}
		}

		total += len(wave)
		wave, next = next, wave
		next = next[:0]
	}

	if total < len(g.Nodes) {
		// TODO: algorithm for finding nodes in cycle
		return &Cycle{}
	}

	for _, c := range g.Cohorts[1:] {
		sort.Ints(c)
	}

	return nil
}

func (b *Bundle) makeGraph() *Cycle {
	b.mapUnits()
	b.makeGraphNodes()
	return rankGraphOrFindCycle(&b.Graph)
}

func printGraph(g *Graph) {
	for i, n := range g.Nodes {
		fmt.Printf("%-3d =>  [%s]\n", i, n.Unit.Path)
	}
	fmt.Println()
	for i, c := range g.Cohorts {
		fmt.Printf("%-3d =>  %v\n", i, c)
	}
}
