package origin

// Origin denotes context in which import string should be resolved.
// Resolution mechanisms and rules may significantly vary between
// different origins
//
// Origin is specified at the start of import block, like this:
//
//	import std {
//		// import entries
//	}
type Origin uint8

const (
	// Zero value of Origin. Valid only as intermidiate value (e.g. check if struct was
	// filled or just created with zero values). Typically should not be used directly
	empty Origin = iota

	// Std denotes import context of standard library
	Std

	// Pkg denotes import context of units from other (third-party) projects.
	// Can be managed via "pkg.ku" file in project's root directory
	Pkg

	// Loc denotes import context of units "local" to current project.
	// In contrast with other origins it is specified by omitting origin name
	// in import block:
	//
	//	import {
	//		// import entries
	//	}
	Loc
)

func (o Origin) IsEmpty() bool {
	return o == empty
}

var text = [...]string{
	empty: "<nil>",

	Std: "std",
	Pkg: "pkg",
	Loc: "loc",
}

func (o Origin) String() string {
	return text[o]
}

// Parse attempts to parse origin from a string. Second return value indicates
// whether this operation was successful or not
//
// Empty string yields local origin
func Parse(s string) (Origin, bool) {
	switch s {
	case "":
		return Loc, true
	case "std":
		return Std, true
	case "pkg":
		return Pkg, true
	default:
		return empty, false
	}
}
