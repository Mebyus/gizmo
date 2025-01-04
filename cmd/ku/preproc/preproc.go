package preproc

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/kir/kbs"
	"github.com/mebyus/gizmo/lexer"
)

var Preproc = &butler.Lackey{
	Name: "preproc",

	Short: "list token stream produced by a given source file (after preprocessor)",
	Usage: "ku preproc [options] <file>",

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

	var p kbs.MacroProcessor
	p.SetInput(lx)
	return lexer.List(os.Stdout, &p)
}
