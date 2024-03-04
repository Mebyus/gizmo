package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mebyus/gizmo/builder/impgraph"
	"github.com/mebyus/gizmo/interp"
	"github.com/mebyus/gizmo/ir/origin"
	"github.com/mebyus/gizmo/parser"
)

const debug = true

type Builder struct {
	cfg Config

	cache *Cache
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
	g.cache = cache

	cache.LookupBuild(unit)

	deps, err := g.WalkFrom(unit)
	if err != nil {
		return err
	}
	graph, err := g.MakeImportGraph(deps)
	if err != nil {
		return err
	}
	genout, err := g.Scribe(graph)
	if err != nil {
		return err
	}
	return g.SaveAndCompile(unit, genout)
}

func (g *Builder) SaveAndCompile(mod string, code []byte) error {
	path, err := g.cache.SaveModGenout(mod, code)
	if err != nil {
		return err
	}

	name := filepath.Base(mod) + ".o"
	dir := filepath.Join(g.cache.dir, "mod", "obj", filepath.Dir(mod))
	err = os.MkdirAll(dir, 0o775)
	if err != nil {
		return err
	}
	out := filepath.Join(dir, name)
	return g.Compile(path, out)
}

// Scribe takes import graph, gathers unit files, parses them and combines
// generated code into singular file build result
func (g *Builder) Scribe(graph *impgraph.Graph) ([]byte, error) {
	pool := NewPool(g.cache, len(graph.Nodes))
	for _, cohort := range graph.Cohorts {
		for _, node := range cohort {
			pool.AddTask(&BuildTask{dep: graph.Nodes[node].Bud.(*DepEntry)})
		}
	}

	pool.Start()
	output := pool.WaitOutput()
	if output.err != nil {
		return nil, output.err
	}
	return output.code, nil
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

	src, err := g.cache.LoadSourceFile(p, "unit.gzm")
	if err != nil {
		return nil, fmt.Errorf("load unit build source: %w", err)
	}
	unit, err := parser.ParseSource(src)
	if err != nil {
		return nil, err
	}
	if unit.Unit == nil {
		return nil, fmt.Errorf("file \"%s\" does not contain unit block", src.Path)
	}
	result, err := interp.Interpret(unit.Unit)
	if err != nil {
		return nil, err
	}
	for _, s := range result.Imports {
		cleaned := filepath.Clean(s)
		if cleaned != s {
			return nil, fmt.Errorf("import path \"%s\" is badly formatted", s)
		}
		if s == "" || s == "." {
			return nil, fmt.Errorf("empty import path")
		}
	}
	return &DepEntry{
		BuildInfo: UnitBuildInfo{
			Files:     result.Files,
			TestFiles: result.TestFiles,

			DefaultNamespace: result.DefaultNamespace,
		},
		Imports: origin.Locals(result.Imports),
		Name:    result.Name,
		Path:    p,
	}, nil
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
