package ast

import (
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/source/origin"
)

// <ImportBlock> = [ "pub" ] "import" [ <ImportOrigin> ] "{" { <ImportSpec> } "}"
//
// <ImportOrigin> = "std" | "pkg" | "loc"
//
// If <ImportOrigin> is absent in block, then it is equivalent to <ImportOrigin> = "loc".
// Canonical gizmo code style omits import origin in such cases, instead of specifying
// it to "loc" explicitly
type ImportBlock struct {
	Pos source.Pos

	Specs []ImportSpec

	Origin origin.Origin

	Public bool
}

// <ImportSpec> = <Name> "=>" <ImportString>
//
// <ImportString> = <String> (cannot be empty)
type ImportSpec struct {
	Name   Identifier
	String ImportString
}

type ImportString struct {
	Pos source.Pos
	Lit string
}

// ImportBind is not produced by parser. It is convenience struct for passing around
// preprocessed data related to a single import spec
type ImportBind struct {
	Name   Identifier
	Path   origin.Path
	Public bool
}
