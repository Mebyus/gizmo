package cc

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/compiler/build"
	"github.com/mebyus/gizmo/compiler/cc"
)

var Obj = &butler.Lackey{
	Name: "obj",

	Short: "invoke c compiler to compile object file",
	Usage: "ku cc obj [options]",

	Exec:   executeObj,
	Params: &Config{},
}

func executeObj(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one c source file must be specified")
	}
	return compileObj(r.Params.(*Config), files[0])

}

func compileObj(config *Config, src string) error {
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
		name := strings.TrimSuffix(base, ext) + ".o"
		out = filepath.Join("build/.cache", name)
	}
	return cc.CompileObj(kind, out, src)
}
