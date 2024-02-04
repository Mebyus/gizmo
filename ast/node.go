package ast

import "github.com/mebyus/gizmo/source"

// UID identifies a node inside a program AST
//
// It consists of two parts: atom id and node local id (unique only inside
// particular atom)
type UID struct {
	atom  uint32
	local uint32
}

type Node interface {
	source.Pin

	UID() UID
}

type uidHolder UID

func (h uidHolder) UID() UID {
	return UID(h)
}
