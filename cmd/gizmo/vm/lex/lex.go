package lex

import (
	"fmt"

	"github.com/mebyus/gizmo/butler"
)

var Lex = &butler.Lackey{
	Name: "lex",

	Short: "list token stream produced by a given kuvm asm file",
	Usage: "gizmo vm lex [options] <file>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}

	return nil
}
