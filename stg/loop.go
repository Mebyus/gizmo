package stg

import (
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
)

// Loop is a loop statement without condition.
// It represents language construct:
//
//	for {
//		<...NodeS>
//	}
type Loop struct {
	NodeS

	Pos source.Pos

	Body Block
}

var _ Statement = &Loop{}

func (*Loop) Kind() stm.Kind {
	return stm.For
}

func (s *Loop) Pin() source.Pos {
	return s.Pos
}

// While is a loop statement with condition.
// It represents language construct:
//
//	for <Exp> {
//		<...NodeS>
//	}
type While struct {
	NodeS

	// Loop condition.
	//
	// Always not nil.
	Exp Exp

	Body Block
}

var _ Statement = &While{}

func (*While) Kind() stm.Kind {
	return stm.While
}

func (s *While) Pin() source.Pos {
	return s.Exp.Pin()
}

// ForRange is a loop statement that uses special "range" iterator.
// It represents language construct:
//
//	for <Name> in range(<Exp>) {
//		<...NodeS>
//	}
//
// This loop is the language analogue to traditional C-style "for" loop.
type ForRange struct {
	NodeS

	Body Block

	// Expression inside "range".
	//
	// Always not nil.
	// Must be of interger type.
	Range Exp

	// Symbol that represents iterating variable.
	Sym *Symbol
}

var _ Statement = &ForRange{}

func (*ForRange) Kind() stm.Kind {
	return stm.ForRange
}

func (s *ForRange) Pin() source.Pos {
	return s.Sym.Pos
}
