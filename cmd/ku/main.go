package main

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/cmd/ku/atom"
	"github.com/mebyus/gizmo/cmd/ku/build"
	"github.com/mebyus/gizmo/cmd/ku/clean"
	"github.com/mebyus/gizmo/cmd/ku/format"
	"github.com/mebyus/gizmo/cmd/ku/gen"
	"github.com/mebyus/gizmo/cmd/ku/lex"
	"github.com/mebyus/gizmo/cmd/ku/unit"
	"github.com/mebyus/gizmo/cmd/ku/vm"
)

func main() {
	if len(os.Args) == 0 {
		panic("os args are empty")
	}
	args := os.Args[1:]

	err := root.Run(args)
	if err != nil {
		fatal(err)
	}
}

var root = &butler.Lackey{
	Name: "ku",

	Short: "Ku is a command line tool for managing Ku source code.",
	Usage: "ku <command> [arguments]",

	Sub: []*butler.Lackey{
		lex.Lex,
		gen.Gen,
		build.Build,
		clean.Clean,
		atom.Atom,
		unit.Unit,
		format.Format,
		vm.VM,
	},
}

func fatal(v any) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}