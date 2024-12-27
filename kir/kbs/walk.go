package kbs

import "path/filepath"

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
