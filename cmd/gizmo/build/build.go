package build

import (
	"fmt"
	"path/filepath"

	"github.com/mebyus/gizmo/builder"
	"github.com/mebyus/gizmo/butler"
)

var Build = &butler.Lackey{
	Name:  "build",
	Short: "build specified targets and their dependencies",
	Usage: "gizmo build [options] <targets>",

	Exec:   execute,
	Params: &Config{},
}

type Config struct {
	BuildKind string

	// Project build configuration file
	BuildFile string

	// Local build configuration file
	EnvFile string
}

func (c *Config) Apply(p *butler.Param) error {
	switch p.Name {
	case "":
		panic("param with empty name")
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
	return build(r.Params.(*Config), targets[0])
}

func build(config *Config, path string) error {
	kind, err := builder.ParseKind(config.BuildKind)
	if err != nil {
		return err
	}

	path = filepath.Clean(path)
	cfg := builder.Config{
		BaseOutputDir: filepath.Join("build", "target"),
		BaseCacheDir:  filepath.Join("build", ".cache"),
		BaseSourceDir: "src",
		BaseNamespace: "coven",

		BuildKind: kind,
	}

	return builder.Build(cfg, path)
}
