package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mebyus/gizmo/cmd/gizmo/highlight"
	"github.com/mebyus/gizmo/cmd/gizmo/lex"
	"github.com/mebyus/gizmo/cmd/gizmo/parse"
	"github.com/mebyus/gizmo/cmd/gizmo/run"
	"github.com/mebyus/gizmo/cmd/gizmo/tree"
)

func main() {
	var cmd string
	var args []string

	if len(os.Args) < 2 {
		cmd = "help"
	} else {
		cmd = os.Args[1]
		args = os.Args[2:]
	}

	var err error
	switch cmd {
	case "lex":
		if len(args) == 0 {
			fatal("filename must be specified to invoke lex command")
		}
		filename := args[0]
		err = lex.Lex(filename)
	case "highlight":
		if len(args) == 0 {
			fatal("filename must be specified to invoke highlight command")
		}
		filename := args[0]
		if len(args) < 2 {
			fatal("token index must be specified to invoke highlight command")
		}
		s := args[1]
		var idx uint64
		idx, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			fatal(err)
		}
		err = highlight.Highlight(filename, int(idx))
	case "parse":
		if len(args) == 0 {
			fatal("filename must be specified to invoke parse command")
		}
		filename := args[0]
		err = parse.Parse(filename)
	case "run":
		if len(args) == 0 {
			fatal("filename must be specified to invoke run command")
		}
		filename := args[0]
		err = run.Run(filename)
	case "tree":
		if len(args) == 0 {
			fatal("filename must be specified to invoke tree command")
		}
		filename := args[0]
		err = tree.Tree(filename)
	case "help":
		fmt.Print(usage)
	default:
		fatal("unknown command: " + cmd)
	}

	if err != nil {
		fatal(err)
	}
}

const usage = `
gizmo is a command line tool for managing Gizmo source code

Usage:

	gizmo <command> [arguments]

Available commands:

	lex
	parse
	tree
	help

`

func fatal(v any) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
