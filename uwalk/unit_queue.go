package uwalk

import (
	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/stg"
)

type QueueItem struct {
	Path origin.Path

	IncludeTestFiles bool
}

// UnitQueue keeps track which unit were already visited and
// which unit should be visited next during unit discovery phase.
type UnitQueue struct {
	// List of paths waiting in queue.
	backlog []QueueItem

	// List of all collected units.
	units []*stg.Unit

	// Set which contains paths of already visited units.
	visited map[origin.Path]struct{}

	// If true then every item placed into queue
	// will have its IncludeTestFiles flag raised regardless
	// of its initial value.
	RaiseIncludeTestFilesFlagsForAllItems bool
}

func NewUnitQueue() *UnitQueue {
	return &UnitQueue{
		visited: make(map[origin.Path]struct{}),
	}
}

// Sorted returns stored units sorted by import path.
func (q *UnitQueue) Sorted() []*stg.Unit {
	stg.SortAndOrderUnits(q.units)
	return q.units
}

// Add tries to add item to backlog. If a given path was
// already visited then this call will be no-op.
func (q *UnitQueue) Add(item QueueItem) {
	_, ok := q.visited[item.Path]
	if ok {
		// given path is already known to walker
		return
	}

	item.IncludeTestFiles = item.IncludeTestFiles || q.RaiseIncludeTestFilesFlagsForAllItems
	q.visited[item.Path] = struct{}{}
	q.backlog = append(q.backlog, item)
}

// Next get next item from backlog and write it to a given pointer.
// Returns true if operation was successfull and false when there are no items left.
func (q *UnitQueue) Next(item *QueueItem) bool {
	if len(q.backlog) == 0 {
		return false
	}

	last := len(q.backlog) - 1
	*item = q.backlog[last]

	// shrink slice, but keep its capacity
	q.backlog = q.backlog[:last]
	return true
}

func (q *UnitQueue) AddUnit(unit *stg.Unit) {
	unit.DiscoveryIndex = uint32(len(q.units))
	q.units = append(q.units, unit)

	for _, p := range unit.Imports.Paths {
		q.Add(QueueItem{Path: p})
	}
}
