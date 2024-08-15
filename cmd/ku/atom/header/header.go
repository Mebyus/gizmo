package header

import (
	"fmt"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/parser"
	"github.com/mebyus/gizmo/source"
)

var Header = &butler.Lackey{
	Name: "header",

	Short: "output header of a given atom (gizmo source file)",
	Usage: "gizmo atom header <file>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}
	return header(files[0])
}

func header(filename string) error {
	src, err := source.Load(filename)
	if err != nil {
		return err
	}
	p := parser.FromSource(src)
	h, err := p.Header()
	if err != nil {
		return err
	}

	for _, path := range h.Imports.Paths {
		fmt.Println(path.String())
	}

	return nil
}
