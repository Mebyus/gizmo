package main

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/cmd/gizmo/gen"
	"github.com/mebyus/gizmo/cmd/gizmo/lex"
	"github.com/mebyus/gizmo/cmd/gizmo/tree"
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

	// var err error
	// switch cmd {
	// case "lex":
	// 	if len(args) == 0 {
	// 		fatal("filename must be specified to invoke lex command")
	// 	}
	// 	filename := args[0]
	// 	err = lex.Lex(filename)
	// case "highlight":
	// 	if len(args) == 0 {
	// 		fatal("filename must be specified to invoke highlight command")
	// 	}
	// 	filename := args[0]
	// 	if len(args) < 2 {
	// 		fatal("token index must be specified to invoke highlight command")
	// 	}
	// 	s := args[1]
	// 	var idx uint64
	// 	idx, err = strconv.ParseUint(s, 10, 64)
	// 	if err != nil {
	// 		fatal(err)
	// 	}
	// 	err = highlight.Highlight(filename, int(idx))
	// case "parse":
	// 	if len(args) == 0 {
	// 		fatal("filename must be specified to invoke parse command")
	// 	}
	// 	filename := args[0]
	// 	err = parse.Parse(filename)
	// case "run":
	// 	if len(args) == 0 {
	// 		fatal("filename must be specified to invoke run command")
	// 	}
	// 	filename := args[0]
	// 	err = run.Run(filename)
	// case "tree":
	// 	if len(args) == 0 {
	// 		fatal("filename must be specified to invoke tree command")
	// 	}
	// 	filename := args[0]
	// 	err = tree.Tree(filename)
	// case "gen":
	// 	if len(args) == 0 {
	// 		fatal("filename must be specified to invoke gen command")
	// 	}
	// 	filename := args[0]
	// 	err = gen.Gen(filename)
	// case "help":
	// 	fmt.Print(usage)
	// default:
	// 	fatal("unknown command: " + cmd)
	// }

	// if err != nil {
	// 	fatal(err)
	// }
}

var root = &butler.Lackey{
	Name: "gizmo",

	Short: "gizmo is a command line tool for managing Gizmo source code",
	Usage: "gizmo <command> [arguments]",

	Sub: []*butler.Lackey{
		lex.Lackey,
		tree.Lackey,
		gen.Gen,
	},
}

func fatal(v any) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
