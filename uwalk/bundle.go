package uwalk

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/stg"
)

type Bundle struct {
	Graph Graph

	// List of all program units sorted by import path.
	Units []*stg.Unit

	// Index in this slice corresponds to Unit.DiscoveryIndex.
	Source []ParserSet

	Map map[origin.Path]*stg.Unit

	// Not nil if bundle has main unit inside.
	Main *stg.Unit
}

type Program struct {
	Graph Graph

	// List of all program units sorted by import path.
	Units []*stg.Unit

	// Not nil if program has main unit.
	Main *stg.Unit
}

func (b *Bundle) Program() *Program {
	return &Program{
		Graph: b.Graph,
		Units: b.Units,
		Main:  b.Main,
	}
}

func (b *Bundle) GetUnitParsers(unit *stg.Unit) ParserSet {
	return b.Source[unit.DiscoveryIndex]
}

type Config struct {
	// Root directory for searching units from standard library.
	StdDir string

	// Root directory for searching local units (from project being built).
	LocalDir string
}

func Walk(cfg *Config, path origin.Path) (*Bundle, error) {
	w := Walker{
		StdDir:   cfg.StdDir,
		LocalDir: cfg.LocalDir,
	}
	err := w.WalkFrom(path)
	if err != nil {
		return nil, err
	}

	b := Bundle{
		Source: w.Source,
		Units:  w.Units,
		Main:   w.Main,
	}
	cycle := b.makeGraph()
	if cycle != nil {
		fmt.Fprintln(os.Stderr, "found import cycle")
		for _, n := range cycle.Nodes {
			fmt.Fprintf(os.Stderr, "[%s]\n", n.Unit.Path)
		}
		return nil, fmt.Errorf("import cycle")
	}

	// TODO: remove debug prints
	printGraph(&b.Graph)

	return &b, nil
}
