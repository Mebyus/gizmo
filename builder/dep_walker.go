package builder

import (
	"sort"

	"github.com/mebyus/gizmo/ir/origin"
)

// DepWalker finds reachable (by imports) units, checks existence of
// unit files (specified in build block), constructs units dependecy
// graph, checks it for import cycles and sorts units by dependency rank
type DepWalker struct {
	// list of paths waiting to be walked to
	backlog []origin.Path

	// list of all collected entries
	entries []*DepEntry

	// maps unit origin path to its entry
	index map[origin.Path]*DepEntry
}

func NewWalker() *DepWalker {
	return &DepWalker{
		index: make(map[origin.Path]*DepEntry),
	}
}

// Sorted returns stored entries sorted by origin path
func (w *DepWalker) Sorted() []*DepEntry {
	e := w.entries
	sort.Slice(e, func(i, j int) bool {
		a := e[i]
		b := e[j]
		return origin.Less(a.Path, b.Path)
	})
	return e
}

// AddPath tries to add origin path to backlog. If a given path is
// already known to walker then this call will be no-op
func (w *DepWalker) AddPath(p origin.Path) {
	_, ok := w.index[p]
	if ok {
		// given path is already known to walker
		return
	}

	// mark path as already known to walker,
	// later, when unit is actually parsed, SaveEntry must be called
	// to supply actual value to this slot in map
	w.index[p] = nil

	// and save it into backlog for later NextPath() call
	w.backlog = append(w.backlog, p)
}

// NextPath get next origin path from backlog. When there are none left
// returns empty origin path
func (w *DepWalker) NextPath() origin.Path {
	if len(w.backlog) == 0 {
		return origin.Empty
	}

	last := len(w.backlog) - 1
	p := w.backlog[last]

	// shrink slice, but keep its capacity
	w.backlog = w.backlog[:last]
	return p
}

func (w *DepWalker) SaveEntry(e *DepEntry) {
	w.entries = append(w.entries, e)
	w.index[e.Path] = e
}

// DepEntry describes dependency relation between unit and
// imports which appear in this unit
type DepEntry struct {
	// Parsed unit info
	BuildInfo UnitBuildInfo

	// Paths imported by this unit
	Imports []origin.Path

	// Unit import path
	Path origin.Path

	// Unit name
	Name string
}

type UnitBuildInfo struct {
	Files     []string
	TestFiles []string

	DefaultNamespace string
}
