package uwalk

import (
	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/stg"
)

// UnitQueue keeps track which unit were already visited and
// which unit should be visited next during unit discovery phase.
type UnitQueue struct {
	// List of paths waiting in queue.
	backlog []origin.Path

	// List of all collected units.
	units []*stg.Unit

	// Set which contains paths of already visited units.
	visited map[origin.Path]struct{}
}

func NewUnitQueue() *UnitQueue {
	return &UnitQueue{
		visited: make(map[origin.Path]struct{}),
	}
}

// Sorted returns stored units sorted by import path.
func (w *UnitQueue) Sorted() []*stg.Unit {
	stg.SortAndOrderUnits(w.units)
	return w.units
}

// AddPath tries to add unit path to backlog. If a given path was
// already visited then this call will be no-op.
func (w *UnitQueue) AddPath(p origin.Path) {
	_, ok := w.visited[p]
	if ok {
		// given path is already known to walker
		return
	}

	w.visited[p] = struct{}{}
	w.backlog = append(w.backlog, p)
}

// NextPath get next unit path from backlog. When there are none left
// returns empty path.
func (w *UnitQueue) NextPath() origin.Path {
	if len(w.backlog) == 0 {
		return origin.Empty
	}

	last := len(w.backlog) - 1
	p := w.backlog[last]

	// shrink slice, but keep its capacity
	w.backlog = w.backlog[:last]
	return p
}

func (w *UnitQueue) AddUnit(unit *stg.Unit) {
	unit.DiscoveryIndex = uint32(len(w.units))
	w.units = append(w.units, unit)

	for _, p := range unit.Imports.Paths {
		w.AddPath(p)
	}
}
