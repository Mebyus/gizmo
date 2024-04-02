package utyp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/er"
	"github.com/mebyus/gizmo/parser"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/tt"
)

var Utyp = &butler.Lackey{
	Name: "tt",

	Short: "construct type tree of a given unit",
	Usage: "gizmo unit tt <unit>",

	Exec: execute,
}

func execute(r *butler.Lackey, units []string) error {
	if len(units) == 0 {
		return fmt.Errorf("at least one unit must be specified")
	}
	err := utyp(units[0])
	if err != nil {
		e, ok := err.(er.Error)
		if !ok {
			return err
		}
		e.Format(os.Stdout, "text")
	}
	return err
}

func utyp(unit string) error {
	files, err := loadFiles(unit)
	if err != nil {
		return err
	}

	m := tt.New(tt.Context{})
	for _, file := range files {
		p := parser.FromSource(file)
		_, err := p.Header()
		if err != nil {
			return err
		}
		atom, err := p.Parse()
		if err != nil {
			return err
		}
		err = m.Add(atom)
		if err != nil {
			return err
		}
	}
	_, err = m.Merge()
	return err
}

const maxFileSize = 1 << 26

func loadFiles(dir string) ([]*source.File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("directory \"%s\" is empty", dir)
	}

	var files []*source.File
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".gm" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		if !info.Mode().IsRegular() {
			continue
		}

		size := info.Size()
		if size > maxFileSize {
			return nil, fmt.Errorf("file \"%s\" is larger than max allowed size", name)
		}

		path := filepath.Join(dir, name)
		file, err := source.Load(path)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("directory \"%s\" does not contain gizmo source files", dir)
	}

	return files, nil
}
