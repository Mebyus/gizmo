package uwalk

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/stg"
)

type Bundle struct {
	Graph Graph

	// Sorted.
	Units []*stg.Unit

	// Index in this slice corresponds to Unit.DiscoveryIndex.
	Source []ParserSet

	Map map[origin.Path]*stg.Unit

	// Not nil if bundle has main unit inside.
	Main *stg.Unit
}

func (b *Bundle) GetUnitParsers(unit *stg.Unit) ParserSet {
	return b.Source[unit.DiscoveryIndex]
}

func Walk(path origin.Path) (*Bundle, error) {
	w := Walker{}
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
	return &b, nil
}
