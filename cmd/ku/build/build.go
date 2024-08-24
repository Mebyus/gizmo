package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mebyus/gizmo/butler"
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

	RootDir string
}

func (c *Config) Apply(p *butler.Param) error {
	switch p.Name {
	case "":
		panic("param with empty name")
	case "output":
		c.OutputFile = p.Str()
	case "kind":
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
			Name:        "kind",
			Kind:        butler.String,
			Def:         "debug",
			ValidValues: []string{"debug", "test", "safe", "fast"},
			Desc:        "select build kind (optimizations, some defaults, etc.)",
		},
		{
			Name: "output",
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
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	execDir := filepath.Dir(execPath)
	rootDir := filepath.Join(execDir, "./../..")
	fmt.Println(execPath)
	fmt.Println(execDir)
	fmt.Println(rootDir)
	return build(r.Params.(*Config), targets[0])
}

func build(config *Config, path string) error {
	path = filepath.Clean(path)

	// TODO: remove this hack
	// it is here only for convenient development with autocomplete
	path = strings.TrimPrefix(path, "src/")

	bundle, err := uwalk.Walk(origin.Local(path))
	if err != nil {
		return err
	}

	global := stg.NewGlobalScope()
	resolver := stg.MapResolver(bundle.Map)
	for _, c := range bundle.Graph.Cohorts {
		for _, i := range c {
			u := bundle.Graph.Nodes[i].Unit

			parsers := bundle.GetUnitParsers(u)
			if len(parsers) == 0 {
				panic("no unit parsers")
			}

			m := stg.New(stg.UnitContext{
				Resolver: resolver,
				Global:   global,
				Unit:     u,
			})
			for _, p := range parsers {
				atom, err := p.Parse()
				if err != nil {
					return err
				}
				// TODO: pass all units to merger at once
				// as a slice
				err = m.Add(atom)
				if err != nil {
					return err
				}
			}
			_, err := m.Merge()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
