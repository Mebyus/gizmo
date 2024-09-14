package cc

import (
	"fmt"

	"github.com/mebyus/gizmo/butler"
)

var CC = &butler.Lackey{
	Name: "cc",

	Short: "invoke c compiler with ku defaults",
	Usage: "ku cc <command> [options]",

	Sub: []*butler.Lackey{
		Exe,
		Obj,
	},
}

type Config struct {
	OutputFile string

	BuildKind string
}

func (c *Config) Apply(p *butler.Param) error {
	switch p.Name {
	case "":
		panic("param with empty name")
	case "o":
		c.OutputFile = p.Str()
	case "k":
		c.BuildKind = p.Str()
	default:
		panic(fmt.Sprintf("unexpected param \"%s\"", p.Name))
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
	}
}
