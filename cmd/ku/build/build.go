package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/cmd/ku/env"
	"github.com/mebyus/gizmo/compiler/build"
	"github.com/mebyus/gizmo/compiler/cc"
	"github.com/mebyus/gizmo/genc"
	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/stg"
	"github.com/mebyus/gizmo/uwalk"
)

var Build = &butler.Lackey{
	Name: "build",

	Short: "build specified targets and their dependencies",
	Usage: "ku build [options] <targets>",

	Exec:   execute,
	Params: &Config{},
}

type Config struct {
	OutputFile string

	BuildKind string

	// Project build configuration file
	BuildFile string

	// Local build configuration file
	EnvFile string

	// Compiler root directory.
	RootDir string

	kind build.Kind
}

func (c *Config) Apply(p *butler.Param) error {
	switch p.Name {
	case "":
		panic("param with empty name")
	case "o":
		c.OutputFile = p.Str()
	case "k":
		c.BuildKind = p.Str()
	case "file":
		c.BuildFile = p.Str()
	case "env":
		c.EnvFile = p.Str()
	default:
		panic(fmt.Sprintf("unexpected param: {%s}", p.Name))
	}
	return nil
}

func (c *Config) Recipe() []butler.Param {
	return []butler.Param{
		{
			Name:        "k",
			Kind:        butler.String,
			Def:         "debug",
			ValidValues: []string{"debug", "test", "safe", "fast"},
			Desc:        "select build kind (optimizations, some defaults, etc.)",
		},
		{
			Name: "o",
			Kind: butler.String,
			Def:  "",
			Desc: "specify output file path",
		},
		{
			Name: "file",
			Kind: butler.String,
			Def:  filepath.Join("build", "build.gm"),
			Desc: "specify a file to use as a build script",
		},
		{
			Name: "env",
			Kind: butler.String,
			Def:  filepath.Join("build", "env.gm"),
			Desc: "specify a file to use for local build environment definitions",
		},
	}
}

func execute(r *butler.Lackey, targets []string) error {
	if len(targets) == 0 {
		return fmt.Errorf("at least one unit must be specified")
	}

	start := time.Now()

	root, err := env.RootDir()
	if err != nil {
		return err
	}

	config := r.Params.(*Config)
	kind, err := build.Parse(config.BuildKind)
	if err != nil {
		return err
	}

	config.kind = kind
	config.RootDir = root

	outexe, err := buildTarget(config, targets[0])
	if err != nil {
		return err
	}

	fmt.Printf("total: %s\n", time.Since(start))
	if outexe != "" {
		fmt.Println()
		fmt.Println("exe:", outexe)
	}
	return nil
}

func makeUnitsBundle(config *Config, path string) (*uwalk.Bundle, error) {
	path = filepath.Clean(path)

	// TODO: remove this hack
	// it is here only for convenient development with autocomplete
	path = strings.TrimPrefix(path, "src/")

	start := time.Now()
	bundle, err := uwalk.Walk(&uwalk.Config{
		StdDir:   filepath.Join(config.RootDir, "src/std"),
		LocalDir: "src",
	}, origin.Local(path))
	if err != nil {
		return nil, err
	}

	fmt.Printf("uwalk: %s\n", time.Since(start))
	return bundle, nil
}

func makeProgram(config *Config, bundle *uwalk.Bundle) (*uwalk.Program, error) {
	start := time.Now()
	resolver := stg.MapResolver(bundle.Map)
	for _, c := range bundle.Graph.Cohorts {
		for _, i := range c {
			u := bundle.Graph.Nodes[i].Unit

			parsers := bundle.GetUnitParsers(u)
			if len(parsers) == 0 {
				panic("no unit parsers")
			}

			atoms := make([]*ast.Atom, 0, len(parsers))
			for _, p := range parsers {
				atom, err := p.Parse()
				if err != nil {
					return nil, err
				}
				atoms = append(atoms, atom)
			}
			_, err := stg.Merge(stg.UnitContext{
				Resolver: resolver,
				Global:   bundle.Global,
				Unit:     u,
			}, atoms)
			if err != nil {
				return nil, err
			}
		}
	}

	fmt.Printf("stg: %s\n", time.Since(start))
	return bundle.Program()
}

// returns path to resulting executable if any
func buildTarget(config *Config, path string) (string, error) {
	bundle, err := makeUnitsBundle(config, path)
	if err != nil {
		return "", err
	}

	program, err := makeProgram(config, bundle)
	if err != nil {
		return "", err
	}

	base := filepath.Base(path)
	outc := filepath.Join("build/.cache", base+".gen.c")
	err = genCode(outc, program)
	if err != nil {
		return "", err
	}
	outobj := filepath.Join("build/.cache", base+".o")
	err = compile(outobj, outc)
	if err != nil {
		return "", err
	}
	if program.Main == nil {
		// skip making executable if there is no main unit
		return "", nil
	}

	outexe := filepath.Join("build/bin", base)
	err = link(outexe, outobj)
	if err != nil {
		return "", err
	}
	return outexe, nil
}

func genCode(out string, p *uwalk.Program) error {
	start := time.Now()

	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	err = genc.GenProgram(f, p)
	if err != nil {
		return err
	}

	fmt.Printf("genc: %s\n", time.Since(start))
	return nil
}

func compile(out string, cfile string) error {
	start := time.Now()
	err := cc.CompileObj(build.Debug, out, cfile)
	if err != nil {
		return err
	}

	fmt.Printf("cc (obj): %s\n", time.Since(start))
	return nil
}

func link(out string, obj string) error {
	start := time.Now()
	err := cc.Link(build.Debug, out, "_start", []string{obj})
	if err != nil {
		return err
	}

	fmt.Printf("cc (link): %s\n", time.Since(start))
	return nil
}
