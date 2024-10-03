package uwalk

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/stg"
)

type Bundle struct {
	Graph Graph

	// List of all program units sorted by import path.
	Units []*stg.Unit

	// Index in this slice corresponds to Unit.DiscoveryIndex.
	// Every parser in this slice has only its header parsed.
	Source []ParserSet

	Map map[origin.Path]*stg.Unit

	// Not nil if bundle has main unit inside.
	Main *stg.Unit

	Global *stg.Scope
}

type Program struct {
	Graph Graph

	// List of all program units sorted by import path.
	Units []*stg.Unit

	// Total number of test symbols in all units.
	TestCount uint64

	Map map[origin.Path]*stg.Unit

	// Not nil if program has main unit.
	Main *stg.Unit

	Global *stg.Scope
}

func (b *Bundle) Program() (*Program, error) {
	var tests uint64
	for _, u := range b.Units {
		tests += uint64(len(u.Tests))
	}

	p := &Program{
		Graph:  b.Graph,
		Units:  b.Units,
		Main:   b.Main,
		Global: b.Global,
		Map:    b.Map,

		TestCount: tests,
	}

	if p.Main == nil {
		return p, nil
	}

	mfun := p.Main.Scope.Lookup("main", 0)
	if mfun == nil {
		return nil, fmt.Errorf("main unit must define \"main\" function")
	}
	if mfun.Kind != smk.Fun {
		return nil, fmt.Errorf("%s: main unit contains symbol \"main\" which is a %s not a function",
			mfun.Pos, mfun.Kind)
	}
	return p, nil
}

func (b *Bundle) GetUnitParsers(unit *stg.Unit) ParserSet {
	return b.Source[unit.DiscoveryIndex]
}

type Config struct {
	// Root directory for searching units from standard library.
	StdDir string

	// Root directory for searching local units (from project being built).
	LocalDir string

	// When false files with suffix ".test.ku" are excluded
	// from unit when loading files.
	IncludeTestFiles bool
}

func Walk(cfg *Config, init ...QueueItem) (*Bundle, error) {
	w := Walker{
		StdDir:   cfg.StdDir,
		LocalDir: cfg.LocalDir,
	}
	err := w.WalkFrom(init...)
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

	b.Global = stg.NewGlobalScope()
	return &b, nil
}
