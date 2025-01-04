package lex

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/lexer"
)

var Lex = &butler.Lackey{
	Name: "lex",

	Short: "list token stream produced by a given source file",
	Usage: "ku lex [options] <file>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}
	return lex(files[0])
}

func lex(filename string) error {
	lx, err := lexer.FromFile(filename)
	if err != nil {
		return err
	}

	fmt.Printf("***   tokens   ***\n")
	fmt.Printf("=============================================\n")
	err = lexer.List(os.Stdout, lx)
	if err != nil {
		return err
	}
	fmt.Printf("=============================================\n")

	stats := lx.Stats()
	fmt.Printf("\n***   stats   ***\n")
	fmt.Printf("=================\n")
	fmt.Printf("size:     %d\n", stats.Size)
	fmt.Printf("sloc:     %d\n", stats.HardLines)
	fmt.Printf("lines:    %d\n", stats.Lines)
	fmt.Printf("tokens:   %d\n", stats.Tokens)
	return nil
}
