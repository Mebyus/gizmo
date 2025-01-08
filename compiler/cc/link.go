package cc

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mebyus/gizmo/compiler/build"
)

func LinkStatic(build build.Kind, out string, entry string, objs []string) error {
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
	args = append(args, staticLinkFlag)
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

func LinkDynamic(build build.Kind, out string, entry string, objs []string, libdir string, libs []string) error {
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
	args = append(args, linkLibSearchDirFlag(libdir))
	// args = append(args, linkLibSearchDirFlag("/usr/lib/x86_64-linux-gnu"))
	args = append(args, "-o", out)
	args = append(args, objs...)

	for _, l := range libs {
		args = append(args, linkLibFlag(l))
	}
	fmt.Println(args)
	// TODO: link external libraries

	cmd := exec.Command(compiler, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
