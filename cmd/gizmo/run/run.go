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
	result, err := interp.Interpret(unit)
	if err != nil {
		return err
	}
	fmt.Println("files      =", result.Files)
	fmt.Println("test_files =", result.TestFiles)
	return nil
}
