package tree

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/parser"
	"github.com/mebyus/gizmo/treeview"
)

func Tree(filename string) error {
	unit, err := parser.ParseFile(filename)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filename, err)
	}
	return treeview.RenderUnitAtom(os.Stdout, unit)
}
