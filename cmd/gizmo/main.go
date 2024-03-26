package main

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/cmd/gizmo/atom"
	"github.com/mebyus/gizmo/cmd/gizmo/build"
	"github.com/mebyus/gizmo/cmd/gizmo/clean"
	"github.com/mebyus/gizmo/cmd/gizmo/gen"
	"github.com/mebyus/gizmo/cmd/gizmo/lex"
	"github.com/mebyus/gizmo/cmd/gizmo/tree"
	"github.com/mebyus/gizmo/cmd/gizmo/ubev"
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
	Name: "gizmo",

	Short: "Gizmo is a command line tool for managing Gizmo source code.",
	Usage: "gizmo <command> [arguments]",

	Sub: []*butler.Lackey{
		lex.Lex,
		tree.Tree,
		gen.Gen,
		build.Build,
		ubev.Ubev,
		clean.Clean,
		atom.Atom,
	},
}

func fatal(v any) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
