package ast

import (
	"github.com/mebyus/gizmo/ast/lbl"
	"github.com/mebyus/gizmo/source"
)

type Label interface {
	Node

	// dummy discriminator method
	Label()

	Kind() lbl.Kind
}

// Dummy provides quick, easy to use implementation of discriminator Label() method
//
// Used for embedding into other (non-dummy) label nodes
type nodeLabel struct{}

func (nodeLabel) Label() {}

type ReservedLabel struct {
	nodeLabel

	Pos     source.Pos
	ResKind lbl.Kind
}

// Explicit interface implementation check
var _ Label = ReservedLabel{}

func (l ReservedLabel) Kind() lbl.Kind {
	return l.ResKind
}

func (l ReservedLabel) Pin() source.Pos {
	return l.Pos
}
