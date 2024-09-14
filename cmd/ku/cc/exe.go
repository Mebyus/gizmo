package cc

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/compiler/build"
	"github.com/mebyus/gizmo/compiler/cc"
)

var Exe = &butler.Lackey{
	Name: "exe",

	Short: "invoke c compiler to compile executable",
	Usage: "ku cc exe [options]",

	Exec:   executeExe,
	Params: &Config{},
}

func executeExe(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one c source file must be specified")
	}
	return compileExe(r.Params.(*Config), files[0])

}

func compileExe(config *Config, src string) error {
	if src == "" {
		return fmt.Errorf("empty source file path")
	}
	kind, err := build.Parse(config.BuildKind)
	if err != nil {
		return err
	}

	out := config.OutputFile
	if out == "" {
		base := filepath.Base(src)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext) + ".exe"
		out = filepath.Join("build/bin", name)
	}
	return cc.CompileExe(kind, out, src)
}
