package build

import (
	"fmt"
	"path/filepath"

	"github.com/mebyus/gizmo/builder"
	"github.com/mebyus/gizmo/butler"
)

var Build = &butler.Lackey{
	Name:  "build",
	Short: "build gizmo unit and its dependencies",
	Usage: "gizmo build [options] <unit>",

	Exec:   execute,
	Params: &Config{},
}

type Config struct {
	BuildKind string
}

func (c *Config) Apply(p *butler.Param) error {
	switch p.Name {
	case "":
		panic("param with empty name")
	case "kind":
		c.BuildKind = p.String()
	default:
		panic(fmt.Sprintf("unexpected param: {%s}", p.Name))
	}
	return nil
}

func (c *Config) Recipe() []butler.Param {
	return []butler.Param{
		{
			Name: "kind",
			Kind: butler.String,
			Def:  "debug",
		},
	}
}

func execute(r *butler.Lackey, units []string) error {
	if len(units) == 0 {
		return fmt.Errorf("at least one unit must be specified")
	}
	return build(r.Params.(*Config), units[0])
}

func build(config *Config, path string) error {
	kind, err := builder.ParseKind(config.BuildKind)
	if err != nil {
		return err
	}

	path = filepath.Clean(path)
	cfg := builder.Config{
		BaseOutputDir: "build",
		BaseCacheDir:  "build",

		BuildKind: kind,
	}

	return builder.Build(cfg, path)
}
