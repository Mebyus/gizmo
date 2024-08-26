package build

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
	"-Wswitch-default",
	"-Wshadow=local",
	"-Wno-main",
	"-Wno-shadow",
	"-Wno-unused-parameter",
	"-Wno-unused-function",
}

var genFlags = []string{
	"-fwrapv",
	"-fno-asynchronous-unwind-tables",
}

var customEntryLinkFlags = []string{
	"-nodefaultlibs",
	"-nolibc",
	"-nostdlib",
	"-nostartfiles",
}

const compiler = "cc"
const compilerStdVersion = "11"
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
	return "-std=c" + v
}

func optzFlag(v string) string {
	return "-O" + v
}

// func (g *Builder) Link(objs []string, entry string, out string) error {
// 	if len(objs) == 0 {
// 		panic("empty object files list")
// 	}
// 	if entry == "" {
// 		panic("entry point name is not specified")
// 	}

// 	args := make([]string, 0, 10+len(objs))

// 	args = append(args, g.buildKindFlags()...)
// 	args = append(args, customEntryLinkFlags...)
// 	args = append(args, "--entry="+entry)
// 	args = append(args, "-o", out)
// 	args = append(args, objs...)

// 	// TODO: link external libraries

// 	cmd := exec.Command(compiler, args...)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	return cmd.Run()
// }

func Compile(src string, out string) error {
	args := make([]string, 0, 10+len(genFlags)+len(warningFlags)+len(otherFlags))
	args = append(args, genFlags...)
	args = append(args, maxCompilerErrorsFlag)
	args = append(args, warningFlags...)
	args = append(args, otherFlags...)
	args = append(args, stdFlag(compilerStdVersion))

	args = append(args, buildKindFlags(BuildDebug)...)

	args = append(args, "-o", out)
	args = append(args, "-c", src)

	cmd := exec.Command(compiler, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type BuildKind uint32

const (
	// debug-friendly optimizations + debug information in binaries + safety checks
	BuildDebug = iota + 1

	// moderate-level optimizations + debug information in binaries + safety checks
	BuildTest

	// most optimizations enabled + safety checks
	BuildSafe

	// all optimizations enabled + disabled safety checks
	BuildFast
)

var buildKindText = [...]string{
	0: "<nil>",

	BuildDebug: "debug",
	BuildTest:  "test",
	BuildSafe:  "safe",
	BuildFast:  "fast",
}

func (k BuildKind) String() string {
	return buildKindText[k]
}

func buildKindFlags(k BuildKind) []string {
	switch k {
	case BuildDebug:
		return []string{optzFlag(debugCompilerOptimizations), debugInfoFlag}
	case BuildTest:
		return []string{optzFlag(testCompilerOptimizations), debugInfoFlag}
	case BuildSafe:
		return []string{optzFlag(safeCompilerOptimizations)}
	case BuildFast:
		return []string{optzFlag(fastCompilerOptimizations)}
	default:
		panic(fmt.Sprintf("unexpected build kind: %d", k))
	}
}
