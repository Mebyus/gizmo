package builder

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/mebyus/gizmo/interp"
	"github.com/mebyus/gizmo/ir/origin"
	"github.com/mebyus/gizmo/parser"
)

const debug = true

type Builder struct {
	cfg Config
}

func New(config *Config) *Builder {
	return &Builder{
		cfg: *config,
	}
}

func (g *Builder) Build(unit string) error {
	return g.WalkFrom(unit)
}

func (g *Builder) WalkFrom(unit string) error {
	walker := NewWalker()
	walker.AddPath(origin.Local(unit))

	for {
		p := walker.NextPath()
		if p.IsEmpty() {
			return nil
		}

		entry, err := g.FindUnitBuildInfo(p)
		if err != nil {
			return err
		}

		walker.AddEntry(entry)
	}
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
