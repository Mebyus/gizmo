package cc

import (
	"os"
	"os/exec"

	"github.com/mebyus/gizmo/compiler/build"
)

func Link(build build.Kind, out string, entry string, objs []string) error {
	if len(objs) == 0 {
		panic("no object files")
	}
	if entry == "" {
		panic("entry point name is not specified")
	}
	if out == "" {
		panic("output file is not specified")
	}

	args := make([]string, 0, 10+len(objs))

	args = append(args, buildKindFlags(build)...)
	args = append(args, customEntryLinkFlags...)
	args = append(args, compilerEntryFlag(entry))
	args = append(args, "-o", out)
	args = append(args, objs...)

	// TODO: link external libraries

	cmd := exec.Command(compiler, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
