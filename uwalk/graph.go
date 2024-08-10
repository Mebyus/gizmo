package uwalk

import (
	"fmt"

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

	// Unit carried by this node.
	Unit *stg.Unit
}

type Graph struct {
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

func (b *Bundle) makeGraph() *Cycle {
	b.mapUnits()
	return nil
}
