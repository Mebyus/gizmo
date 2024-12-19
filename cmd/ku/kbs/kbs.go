package kbs

import (
	"fmt"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/kir/kbs"
)

var Kbs = &butler.Lackey{
	Name: "kbs",

	Short: "list includes produced by a given build script file",
	Usage: "gizmo kbs [options] <file>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}
	return lex(files[0])
}

func lex(filename string) error {
	includes, err := kbs.ParseFile(filename)
	if err != nil {
		return err
	}
	for _, include := range includes {
		fmt.Println(include)
	}
	return nil
}
