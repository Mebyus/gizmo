package cc

import (
	"os"
	"os/exec"

	"github.com/mebyus/gizmo/compiler/build"
)

func CompileObj(build build.Kind, out string, src string) error {
	return compile(build, out, src, "")
}

func CompileExe(build build.Kind, out string, src string) error {
	return compile(build, out, src, "_start")
}

func compile(build build.Kind, out string, src string, entry string) error {
	if src == "" {
		panic("source file is not specified")
	}
	if out == "" {
		panic("output file is not specified")
	}

	// true for compiling executable file
	exe := entry != ""

	args := make([]string, 0, 10+len(codegenFlags)+len(warningFlags)+len(otherFlags))
	args = append(args, codegenFlags...)
	args = append(args, maxCompilerErrorsFlag)
	args = append(args, warningFlags...)
	args = append(args, otherFlags...)
	args = append(args, stdFlag(compilerStdVersion))

	args = append(args, buildKindFlags(build)...)

	if exe {
		args = append(args, customEntryLinkFlags...)
		args = append(args, compilerEntryFlag(entry))
	}

	args = append(args, "-o", out)

	if !exe {
		args = append(args, "-c")
	}
	args = append(args, src)

	cmd := exec.Command(compiler, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
