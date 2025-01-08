package kir

import (
	"fmt"
	"path/filepath"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/compiler/build"
	"github.com/mebyus/gizmo/compiler/cc"
	"github.com/mebyus/gizmo/kir/kbs"
)

var Kir = &butler.Lackey{
	Name: "kir",

	Short: "list includes produced by a given build script file",
	Usage: "ku kir [options] <file>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}
	return kir(files[0])
}

func kir(filename string) error {
	script, err := kbs.Walk("src", filename)
	if err != nil {
		return err
	}

	outc := filepath.Join("build/.cache", script.Name+".gen.c")
	err = kbs.GenIntoFile(outc, script.Includes)
	if err != nil {
		return err
	}
	outobj := filepath.Join("build/.cache", script.Name+".o")
	err = cc.CompileObj(build.Fast, outobj, outc)
	if err != nil {
		return err
	}
	outexe := filepath.Join("build/bin", script.Name)
	err = cc.LinkDynamic(build.Fast, outexe, "_start", []string{outobj}, "build/bin/lib", script.Links)
	if err != nil {
		return err
	}
	return nil
}
