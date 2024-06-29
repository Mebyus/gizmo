package lex

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/vm"
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

	return lex(files[0])
}

func lex(path string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lx := vm.NewLexer(src)
	return vm.ListTokens(os.Stdout, lx)
}
