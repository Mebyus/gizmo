package gen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/genc"
	"github.com/mebyus/gizmo/stg"
)

var Gen = &butler.Lackey{
	Name: "gen",

	Short: "generate C code from a given ku unit",
	Usage: "ku unit gen [options] <units>",

	Exec:   execute,
	Params: &Config{},
}

type Config struct {
	OutputFile string
}

func (c *Config) Apply(p *butler.Param) error {
	switch p.Name {
	case "output-file":
		c.OutputFile = p.Str()
	default:
		panic(fmt.Sprintf("unexpected param: {%s}", p.Name))
	}
	return nil
}

func (c *Config) Recipe() []butler.Param {
	return []butler.Param{
		{
			Name:     "output-file",
			Kind:     butler.String,
			Desc:     "path to file where output should be stored",
			Required: true,
		},
	}
}

func execute(r *butler.Lackey, units []string) error {
	if len(units) == 0 {
		return fmt.Errorf("at least one unit must be specified")
	}
	return gen(r.Params.(*Config), filepath.Clean(units[0]))
}

func gen(config *Config, dir string) error {
	u, err := stg.UnitFromDir(stg.NewEmptyResolver(), dir)
	if err != nil {
		return err
	}
	fmt.Printf("unit name: %s\n", u.Name)

	f, err := os.Create(config.OutputFile)
	if err != nil {
		return err
	}

	return genc.GenUnit(f, u)
}
