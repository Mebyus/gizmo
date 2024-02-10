package gen

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/gencpp"
	"github.com/mebyus/gizmo/parser"
)

func Gen(filename string) error {
	unit, err := parser.ParseFile(filename)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filename, err)
	}
	return gencpp.Gen(os.Stdout, unit)
}
