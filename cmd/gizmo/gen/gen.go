package gen

import (
	"fmt"
	"io"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/gencpp"
	"github.com/mebyus/gizmo/parser"
)

var Gen = &butler.Lackey{
	Name:  "gen",
	Short: "generate C++ code from a given gizmo source file",
	Usage: "gizmo gen [options] <files>",

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
			Name: "output-file",
			Kind: butler.String,
			Desc: "path to file where output should be stored",
		},
	}
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}
	return gen(r.Params.(*Config), files[0])
}

func gen(config *Config, filename string) error {
	unit, err := parser.ParseFile(filename)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filename, err)
	}
	var out io.Writer
	if config.OutputFile == "" {
		out = os.Stdout
	} else {
		f, err := os.Create(config.OutputFile)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}
	return gencpp.Gen(out, nil, unit)
}
