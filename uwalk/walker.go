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
}

func (w *Walker) WalkFrom(path origin.Path) error {
	q := NewUnitQueue()
	q.AddPath(path)

	for {
		p := q.NextPath()
		if p.IsEmpty() {
			w.Units = q.Sorted()
			return nil
		}

		u, err := w.AnalyzeHeaders(p)
		if err != nil {
			return err
		}

		if u.Name == "main" {
			if w.Main != nil {
				return fmt.Errorf("main unit [%s] cannot be imported", u.Path)
			}
			w.Main = u
		}
		q.AddUnit(u)
	}
}

func (w *Walker) AnalyzeHeaders(path origin.Path) (*stg.Unit, error) {
	dir, err := w.Resolve(path)
	if err != nil {
		return nil, err
	}
	files, err := source.LoadUnitFiles(dir)
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
