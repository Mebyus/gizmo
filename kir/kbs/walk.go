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
//
// Returns ordered list of regular include file paths.
func Walk(prefix string, start string) ([]string, error) {
	var w Walker
	w.prefix = prefix

	err := w.Walk(start)
	if err != nil {
		return nil, err
	}
	return w.includes, nil
}

type Walker struct {
	// Accumulated list of regular include file paths.
	includes []string

	prefix string
}

func (w *Walker) Walk(start string) error {
	includes, err := ParseFile(filepath.Join(w.prefix, start, "unit.build.kir"))
	if err != nil {
		return err
	}
	for _, include := range includes {
		switch filepath.Ext(include) {
		case ".kir", ".c":
			w.includes = append(w.includes, filepath.Join(w.prefix, start, include))
		default:
			err = w.Walk(include)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Gen(w io.Writer, files []string) error {
	var gen Generator

	for _, path := range files {
		var err error
		switch filepath.Ext(path) {
		case ".c":
			err = copyFile(w, path)
		case ".kir":
			err = genFile(&gen, w, path)
		default:
			panic(fmt.Sprintf("unexpected file \"%s\" extension", path))
		}
		if err != nil {
			return err
		}
	}
	return nil
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
