package builder

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/mebyus/gizmo/builder/impgraph"
	"github.com/mebyus/gizmo/interp"
	"github.com/mebyus/gizmo/ir/origin"
	"github.com/mebyus/gizmo/parser"
)

const debug = true

type Builder struct {
	cfg Config
}

func New(config Config) *Builder {
	return &Builder{
		cfg: config,
	}
}

func Build(config Config, unit string) error {
	g := New(config)
	return g.Build(unit)
}

func (g *Builder) Build(unit string) error {
	cache, err := NewCache(&g.cfg)
	if err != nil {
		return err
	}
	cache.LookupBuild(unit)

	deps, err := g.WalkFrom(unit)
	if err != nil {
		return err
	}
	graph, err := g.MakeImportGraph(deps)
	if err != nil {
		return err
	}
	return g.Scribe(graph)
}

// Scribe takes import graph, gathers unit files, parses them and combines
// generated code into singular file build result
func (g *Builder) Scribe(graph *impgraph.Graph) error {

	return nil
}

func (g *Builder) WalkFrom(unit string) ([]*DepEntry, error) {
	walker := NewWalker()
	walker.AddPath(origin.Local(unit))

	for {
		p := walker.NextPath()
		if p.IsEmpty() {
			return walker.Sorted(), nil
		}

		entry, err := g.FindUnitBuildInfo(p)
		if err != nil {
			return nil, err
		}

		walker.AddEntry(entry)
	}
}

func (g *Builder) MakeImportGraph(deps []*DepEntry) (*impgraph.Graph, error) {
	graph := impgraph.New(len(deps))
	for _, d := range deps {
		err := graph.Add(d)
		if err != nil {
			return nil, err
		}
	}

	err := graph.Scan()
	if err != nil {
		return nil, err
	}

	if debug {
		for _, node := range graph.Nodes {
			g.debug("node %-3d => %s", node.Index, node.Bud.UID())
		}

		if len(graph.Nodes) > 1 && len(graph.Roots) == 0 {
			g.debug("warn: import graph does not have roots")
		}
		if len(graph.Nodes) > 1 && len(graph.Pinnacles) == 0 {
			g.debug("warn: import graph does not have pinnacles")
		}
	}

	cycle := graph.Chart()
	if cycle != nil {
		for _, node := range cycle.Nodes {
			fmt.Println(node.Bud.UID())
		}
		return nil, fmt.Errorf("import cycle detected")
	}

	graph.Rank()
	if debug {
		for rank, cohort := range graph.Cohorts {
			g.debug("rank %-3d => %v", rank, cohort)
		}
	}

	return graph, nil
}

func (g *Builder) FindUnitBuildInfo(p origin.Path) (*DepEntry, error) {
	if debug {
		g.debug(p.String())
	}

	switch p.Origin {
	case origin.Std:
		panic("not implemented for std")
	case origin.Pkg:
		panic("not implemented for pkg")
	case origin.Loc:
		filename := filepath.Join(g.cfg.BaseSourceDir, p.ImpStr, "unit.gzm")
		unit, err := parser.ParseFile(filename)
		if err != nil {
			return nil, err
		}
		if unit.Unit == nil {
			return nil, fmt.Errorf("file \"%s\" does not contain unit block", filename)
		}
		result, err := interp.Interpret(unit.Unit)
		if err != nil {
			return nil, err
		}
		return &DepEntry{
			BuildInfo: UnitBuildInfo{
				Files:     result.Files,
				TestFiles: result.TestFiles,

				DefaultNamespace: result.Name,
			},
			Imports: origin.Locals(result.Imports),
			Name:    result.Name,
			Path:    p,
		}, nil
	default:
		panic("unexpected import origin: " + strconv.FormatInt(int64(p.Origin), 10))
	}
}

func (g *Builder) debug(format string, args ...any) {
	fmt.Print("[debug] builder | ")
	fmt.Printf(format, args...)
	fmt.Println()
}

type BuildKind uint32

const (
	// debug-friendly optimizations + debug information in binaries + safety checks
	BuildDebug = iota + 1

	// moderate-level optimizations + debug information in binaries + safety checks
	BuildTest

	// most optimizations enabled + safety checks
	BuildSafe

	// all optimizations enabled + disabled safety checks
	BuildFast
)

var buildKindText = [...]string{
	0: "<nil>",

	BuildDebug: "debug",
	BuildTest:  "test",
	BuildSafe:  "safe",
	BuildFast:  "fast",
}

func (k BuildKind) String() string {
	return buildKindText[k]
}

func ParseKind(s string) (BuildKind, error) {
	switch s {
	case "":
		return 0, fmt.Errorf("empty build kind")
	case "debug":
		return BuildDebug, nil
	case "test":
		return BuildTest, nil
	case "safe":
		return BuildSafe, nil
	case "fast":
		return BuildFast, nil
	default:
		return 0, fmt.Errorf("unknown build kind: %s", s)
	}
}

type Config struct {
	BaseOutputDir string
	BaseCacheDir  string
	BaseSourceDir string

	BaseNamespace string

	BuildKind BuildKind
}
