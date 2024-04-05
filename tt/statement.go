package tt

import (
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
)

type Node interface {
	source.Pin
}

type Statement interface {
	Node

	// dummy discriminator method
	Statement()

	Kind() stm.Kind
}

// Dummy provides quick, easy to use implementation of discriminator Statement() method
//
// Used for embedding into other (non-dummy) statement nodes
type nodeStatement struct{}

func (nodeStatement) Statement() {}
