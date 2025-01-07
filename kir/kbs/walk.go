package kbs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mebyus/gizmo/parser"
)

// Walk traverses includes from build files.
// Start must point to project root directory.
func Walk(prefix string, start string) (*ScriptOutput, error) {
	var w Walker
	w.prefix = prefix
	w.eval.env = nil

	err := w.Walk(start)
	if err != nil {
		return nil, err
	}
	return &w.out, nil
}

func WalkWithTests(prefix string, start string) (*ScriptOutput, error) {
	var w Walker
	w.prefix = prefix
	w.eval.env = map[string]string{
		"BUILD_KIND": "test",
	}

	err := w.Walk(start)
	if err != nil {
		return nil, err
	}
	return &w.out, nil
}

type Walker struct {
	// Accumulated output of scripts (includes, links, etc.).
	out ScriptOutput

	eval Evaluator

	prefix string
}

func (w *Walker) Walk(start string) error {
	directives, err := ParseFile(filepath.Join(w.prefix, start, "unit.kbs"))
	if err != nil {
		return err
	}
	script, err := w.eval.Eval(directives)
	if err != nil {
		return err
	}

	if w.out.Name == "" {
		// assign script name to the first non-empty name in the tree
		w.out.Name = script.Name
	}
	w.out.Links = append(w.out.Links, script.Links...)

	for _, include := range script.Includes {
		switch filepath.Ext(include) {
		case ".kir", ".c":
			w.out.Includes = append(w.out.Includes, filepath.Join(w.prefix, start, include))
		default:
			err = w.Walk(include)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GenWithTests(w io.Writer, files []string) error {
	var gen Generator
	gen.IncludeTests = true

	err := driveGenerator(w, &gen, files)
	if err != nil {
		return err
	}

	err = genTests(&gen, w)
	if err != nil {
		return err
	}

	return nil
}

func Gen(w io.Writer, files []string) error {
	var gen Generator
	return driveGenerator(w, &gen, files)
}

func driveGenerator(w io.Writer, gen *Generator, files []string) error {
	for _, path := range files {
		var err error
		switch filepath.Ext(path) {
		case ".c":
			err = copyFile(w, path)
		case ".kir":
			err = genFile(gen, w, path)
		default:
			panic(fmt.Sprintf("unexpected file \"%s\" extension", path))
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func GenIntoFile(path string, files []string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	return Gen(out, files)
}

func GenWithTestsIntoFile(path string, files []string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	return GenWithTests(out, files)
}

func genTests(gen *Generator, w io.Writer) error {
	gen.TestsAndDriver(gen.Tests)
	_, err := gen.WriteTo(w)
	return err
}

func copyFile(w io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	_, err = file.WriteTo(w)
	return err
}

func genFile(gen *Generator, w io.Writer, path string) error {
	atom, err := parser.ParseFile(path)
	if err != nil {
		return err
	}
	gen.Atom(atom)
	_, err = gen.WriteTo(w)
	return err
}
