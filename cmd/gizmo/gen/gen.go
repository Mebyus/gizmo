package gen

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/gencpp"
	"github.com/mebyus/gizmo/parser"
)

var Lackey = &butler.Lackey{
	Name:  "gen",
	Short: "generate C++ code from a given unit",
	Usage: "gizmo gen [options] <files>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}
	return Gen(files[0])
}

func Gen(filename string) error {
	unit, err := parser.ParseFile(filename)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filename, err)
	}
	return gencpp.Gen(os.Stdout, unit)
}
