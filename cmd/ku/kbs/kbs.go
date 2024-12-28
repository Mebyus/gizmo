package kbs

import (
	"fmt"
	"os"

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
	includes, err := kbs.Walk("src", filename)
	if err != nil {
		return err
	}
	return kbs.Gen(os.Stdout, includes)
}
