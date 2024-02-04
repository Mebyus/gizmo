package run

import (
	"fmt"

	"github.com/mebyus/gizmo/interp"
	"github.com/mebyus/gizmo/parser"
)

func Run(filename string) error {
	unit, err := parser.ParseFile(filename)
	if err != nil {
		return err
	}
	if unit.Unit == nil {
		return fmt.Errorf("file does not contain unit block")
	}
	result, err := interp.Interpret(*unit.Unit)
	if err != nil {
		return err
	}
	fmt.Println("files      =", result.Files)
	fmt.Println("test_files =", result.TestFiles)
	return nil
}
