package uwalk

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/parser"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/stg"
)

type ParserSet []*parser.Parser

type Walker struct {
	Units  []*stg.Unit
	Source []ParserSet

	StdDir   string
	LocalDir string

	// Not nil if main unit is found.
	Main *stg.Unit

	IncludeTestFilesForAllItems bool
}

func (w *Walker) WalkFrom(init ...QueueItem) error {
	if len(init) == 0 {
		panic("no init items")
	}

	q := NewUnitQueue()
	q.RaiseIncludeTestFilesFlagsForAllItems = w.IncludeTestFilesForAllItems

	for _, item := range init {
		q.Add(item)
	}

	for {
		var item QueueItem
		if !q.Next(&item) {
			w.Units = q.Sorted()
			return nil
		}

		u, err := w.AnalyzeHeaders(item)
		if err != nil {
			return err
		}

		q.AddUnit(u)
		if u.Name == "main" {
			if u.DiscoveryIndex != 0 {
				return fmt.Errorf("main unit [%s] cannot be imported", u.Path)
			}
			if w.Main != nil {
				panic("multiple main units in uwalk graph")
			}
			w.Main = u
		}
	}
}

func (w *Walker) AnalyzeHeaders(item QueueItem) (*stg.Unit, error) {
	path := item.Path

	dir, err := w.Resolve(path)
	if err != nil {
		return nil, err
	}
	files, err := source.LoadUnitFiles(&source.UnitParams{
		Dir:              dir,
		IncludeTestFiles: item.IncludeTestFiles,
	})
	if err != nil {
		return nil, err
	}

	var name string
	parsers := make([]*parser.Parser, 0, len(files))
	pset := origin.NewSet()
	for _, file := range files {
		p := parser.FromSource(file)
		h, err := p.Header()
		if err != nil {
			return nil, err
		}
		if h.Unit != nil {
			// if file has unit clause, use it to
			// determine unit name

			n := h.Unit.Name.Lit
			if name == "" {
				name = n
			} else if n != name {
				return nil, fmt.Errorf("%s: inconsistent unit name \"%s\" (previous was \"%s\")",
					h.Unit.Name.Pos.String(), n, name)
			}
		}

		for _, p := range h.Imports.Paths {
			if p == path {
				return nil, fmt.Errorf("unit [%s] imports itself", path)
			}
			if pset.Has(p) {
				return nil, fmt.Errorf("multiple imports of the same unit [%s]", p)
			}
			pset.Add(p)
		}

		parsers = append(parsers, p)
	}
	if name == "" {
		// unit does not have files with unit clause
		// determine unit name by its directory

		stat, err := os.Lstat(dir)
		if err != nil {
			return nil, err
		}
		if !stat.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", dir)
		}
		name = stat.Name()
		if name == "" {
			return nil, fmt.Errorf("os: no directory name")
		}
	}

	if name == "" {
		panic("empty unit name")
	}

	w.Source = append(w.Source, parsers)
	return &stg.Unit{
		Name: name,
		Path: path,
		Imports: stg.UnitImports{
			Paths: pset.Sorted(),
		},
	}, nil
}

func (w *Walker) Resolve(path origin.Path) (string, error) {
	o := path.Origin
	switch path.Origin {
	case 0:
		panic("empty path")
	case origin.Std:
		return w.StdDir + "/" + path.ImpStr, nil
	case origin.Pkg:
		panic("not implemented")
	case origin.Loc:
		return w.LocalDir + "/" + path.ImpStr, nil
	default:
		panic(fmt.Sprintf("unexpected %s (%d) origin", o, o))
	}
}
