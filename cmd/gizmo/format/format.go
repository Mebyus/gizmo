package format

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/lexer"
)

var Format = &butler.Lackey{
	Name: "fmt",

	Short: "apply standard formatting to specified gizmo source file",
	Usage: "gizmo fmt [options] <file>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}
	return format(files[0])
}

func format(filename string) error {
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
