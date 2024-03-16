package ubev

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/interp"
	"github.com/mebyus/gizmo/parser"
)

var Ubev = &butler.Lackey{
	Name:  "ubev",
	Short: "evaluate unit build block",
	Usage: "gizmo ubev [options] <unit|file>",

	Exec: execute,
}

func execute(r *butler.Lackey, paths []string) error {
	if len(paths) == 0 {
		return fmt.Errorf("at least one path must be specified")
	}
	return ubev(paths[0])
}

func determineUnitFile(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return filepath.Join(path, "unit.gm"), nil
	}
	return path, nil
}

func ubev(path string) error {
	filename, err := determineUnitFile(path)
	if err != nil {
		return err
	}

	unit, err := parser.ParseFile(filename)
	if err != nil {
		return err
	}
	if unit.Unit == nil {
		return fmt.Errorf("file does not contain unit block")
	}
	result, err := interp.Interpret(unit.Unit)
	if err != nil {
		return err
	}
	fmt.Println("imports    =", result.Imports)
	fmt.Println("files      =", result.Files)
	fmt.Println("test_files =", result.TestFiles)
	return nil
}
