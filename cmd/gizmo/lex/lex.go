package lex

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/lexer"
)

var Lackey = &butler.Lackey{
	Name:  "lex",
	Short: "list token stream produced by a given source file",
	Usage: "gizmo lex [options] <file>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}
	return Lex(files[0])
}

func Lex(filename string) error {
	lx, err := lexer.FromFile(filename)
	if err != nil {
		return err
	}
	err = lexer.List(os.Stdout, lx)
	if err != nil {
		return err
	}
	return nil
}
