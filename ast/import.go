package ast

import "github.com/mebyus/gizmo/token"

// <ImportBlock> = [ "pub" ] "import" [ <ImportOrigin> ] "{" { <Import> } "}"
type ImportBlock struct {
	Specs []Import

	Origin Identifier

	Public bool
}

// <Import> = <Name> "=>" <ImportString>
type Import struct {
	Name   Identifier
	String ImportString
}

// Token.Kind is String
type ImportString token.Token
