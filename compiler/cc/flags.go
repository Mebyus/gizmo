package cc

import (
	"fmt"
	"strings"

	"github.com/mebyus/gizmo/compiler/build"
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

var codegenFlags = []string{
	"-fwrapv",
	"-fno-asynchronous-unwind-tables",
	"-funsigned-char",
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
const linkTimeOptimization = "-flto"
const wholeProgramOptimizations = "-fwhole-program"

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

func compilerEntryFlag(entry string) string {
	return "--entry=" + entry
}

func linkerFlagsViaCC(flags []string) string {
	if len(flags) == 0 {
		return ""
	}
	return "-Wl," + strings.Join(flags, ",")
}

func linkerEntryPointFlag(entry string) string {
	return "-e" + entry
}

func buildKindFlags(k build.Kind) []string {
	switch k {
	case 0:
		panic("empty build kind")
	case build.Debug:
		return []string{optzFlag(debugCompilerOptimizations), debugInfoFlag}
	case build.Test:
		return []string{optzFlag(testCompilerOptimizations), debugInfoFlag}
	case build.Safe:
		return []string{
			optzFlag(safeCompilerOptimizations),

			// TODO: enabling optimizations below breaks entrypoint linkage
			// Research how to fix this.

			// wholeProgramOptimizations,
			// linkTimeOptimization,
		}
	case build.Fast:
		return []string{optzFlag(fastCompilerOptimizations)} // wholeProgramOptimizations, linkTimeOptimization,

	default:
		panic(fmt.Sprintf("%s (%d) build not implemented", k, k))
	}
}
