package tree

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/parser"
	"github.com/mebyus/gizmo/treeview"
)

var Lackey = &butler.Lackey{
	Name:  "tree",
	Short: "display AST of a given unit or source file",
	Usage: "gizmo tree [options] <files>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}
	return Tree(files[0])
}

func Tree(filename string) error {
	unit, err := parser.ParseFile(filename)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filename, err)
	}
	return treeview.RenderUnitAtom(os.Stdout, unit)
}
