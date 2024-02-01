package lex

import (
	"os"

	"github.com/mebyus/gizmo/lexer"
)

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
