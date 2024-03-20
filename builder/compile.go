package builder

import (
	"fmt"
	"os"
	"os/exec"
)

var warningFlags = []string{
	"-Wall",
	"-Wextra",
	"-Wconversion",
	"-Wunreachable-code",
	"-Wshadow",
	"-Wundef",
	"-Wfloat-equal",
	"-Wformat=0",
	"-Wpointer-arith",
	"-Winit-self",
	"-Wduplicated-branches",
	"-Wduplicated-cond",
	"-Wnull-dereference",
	"-Wvla",
	"-Wnoexcept",
	"-Wswitch-default",
	"-Wno-main",
	"-Wno-shadow",
	"-Wshadow=local",
}

var genFlags = []string{
	"-fwrapv",
	"-fno-exceptions",
	"-fno-rtti",
}

var customEntryLinkFlags = []string{
	"-nodefaultlibs",
	"-nolibc",
	"-nostdlib",
	"-nostartfiles",
}

const compiler = "g++"
const compilerStdVersion = "20"
const debugCompilerOptimizations = "g"
const testCompilerOptimizations = "1"
const safeCompilerOptimizations = "2"
const fastCompilerOptimizations = "fast"
const debugInfoFlag = "-ggdb"
const maxCompilerErrorsFlag = "-fmax-errors=1"

var otherFlags = []string{
	"-Werror",
	"-pipe",
}

func stdFlag(v string) string {
	return "-std=c++" + v
}

func optzFlag(v string) string {
	return "-O" + v
}

func (g *Builder) Link(objs []string, entry string, out string) error {
	if len(objs) == 0 {
		panic("empty object files list")
	}
	if len(entry) == 0 {
		panic("entry point name is not specified")
	}

	args := make([]string, 0, 10+len(objs))

	args = append(args, g.buildKindFlags()...)
	args = append(args, customEntryLinkFlags...)
	args = append(args, "--entry="+entry)
	args = append(args, "-o", out)
	args = append(args, objs...)

	// TODO: link external libraries

	cmd := exec.Command(compiler, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (g *Builder) Compile(src string, out string) error {
	args := make([]string, 0, 10+len(genFlags)+len(warningFlags)+len(otherFlags))
	args = append(args, genFlags...)
	args = append(args, maxCompilerErrorsFlag)
	args = append(args, warningFlags...)
	args = append(args, otherFlags...)
	args = append(args, stdFlag(compilerStdVersion))

	args = append(args, g.buildKindFlags()...)

	args = append(args, "-o", out)
	args = append(args, "-c", src)

	cmd := exec.Command(compiler, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (g *Builder) buildKindFlags() []string {
	switch g.cfg.BuildKind {
	case BuildDebug:
		return []string{optzFlag(debugCompilerOptimizations), debugInfoFlag}
	case BuildTest:
		return []string{optzFlag(testCompilerOptimizations), debugInfoFlag}
	case BuildSafe:
		return []string{optzFlag(safeCompilerOptimizations)}
	case BuildFast:
		return []string{optzFlag(fastCompilerOptimizations)}
	default:
		panic(fmt.Sprintf("unexpected build kind: %d", g.cfg.BuildKind))
	}
}
