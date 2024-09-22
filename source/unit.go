package source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const MaxFileSize = 1 << 26

type UnitParams struct {
	// System (not unit origin) path to directory where unit files are located.
	Dir string

	// Default value will be used if this field equals 0.
	MaxFileSize uint64

	IncludeTestFiles bool
}

func LoadUnitFiles(params *UnitParams) ([]*File, error) {
	if params.MaxFileSize == 0 {
		params.MaxFileSize = MaxFileSize
	}
	dir := params.Dir
	if dir == "" {
		panic("empty unit directory path")
	}

	entries, err := os.ReadDir(params.Dir)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("directory \"%s\" is empty", params.Dir)
	}

	var files []*File
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".ku" {
			continue
		}
		if !params.IncludeTestFiles && strings.HasSuffix(name, ".test.ku") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		if !info.Mode().IsRegular() {
			continue
		}

		size := uint64(info.Size())
		if size == 0 {
			continue
		}
		if size > params.MaxFileSize {
			return nil, fmt.Errorf("file \"%s\" is larger (%d bytes) than max allowed size", name, size)
		}

		path := filepath.Join(dir, name)
		file, err := Load(path)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("directory \"%s\" does not contain ku source files", dir)
	}
	SortAndOrder(files)

	return files, nil
}
