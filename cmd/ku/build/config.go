package build

import (
	"path/filepath"
	"strings"

	"github.com/mebyus/gizmo/cmd/ku/env"
	"github.com/mebyus/gizmo/compiler/build"
)

type Config struct {
	// Build or test initial path that was specified for the command.
	//
	// Initial path can point to unit or file where program discovery starts.
	InitPath string

	// Output file path.
	//
	// For "test" command it stores path to temporary executable file
	// that does the testing.
	//
	// If this field is empty, then compiler will place output file
	// into default location and use generated name. Default location
	// is based on command being executed and unit name.
	//
	// For "build" command executables are placed into "build/bin/"
	// directory by default.
	//
	// For "test" command executables are placed into "build/.test/"
	// directory by default.
	OutFile string

	// Compiler root directory.
	RootDir string

	Kind build.Kind

	// If true then "test" command is being executed.
	Test bool
}

func NewConfigFromParams(params *Params, path string) (*Config, error) {
	root, err := env.RootDir()
	if err != nil {
		return nil, err
	}

	kind, err := build.Parse(params.BuildKind)
	if err != nil {
		return nil, err
	}

	path = filepath.Clean(path)

	// TODO: remove this hack
	// it is here only for convenient development with autocomplete
	path = strings.TrimPrefix(path, "src/")

	return &Config{
		InitPath: path,
		OutFile:  params.OutputFile,
		RootDir:  root,
		Kind:     kind,
	}, nil
}
