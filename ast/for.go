package ast

import (
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
)

// <For> = "for" <BlockStatement>
type For struct {
	NodeS

	Pos source.Pos

	Body Block
}

var _ Statement = For{}

func (For) Kind() stm.Kind {
	return stm.For
}

func (s For) Pin() source.Pos {
	return s.Pos
}

// <ForIf> = "for" <Expression> <Block>
type ForIf struct {
	NodeS

	// Always not nil.
	If Exp

	Body Block
}

var _ Statement = ForIf{}

func (ForIf) Kind() stm.Kind {
	return stm.While
}

func (s ForIf) Pin() source.Pos {
	return s.If.Pin()
}

// <ForIn> => "for" <Name> "in" <Exp> <Block>
type ForIn struct {
	NodeS

	Name Identifier

	// Always not nil
	Iterator Exp

	// must contain at least one statement
	Body Block
}

var _ Statement = ForIn{}

func (ForIn) Kind() stm.Kind {
	return stm.ForIn
}

func (s ForIn) Pin() source.Pos {
	return s.Name.Pos
}

// <ForRangeStatement> = "for" <Name> "in" "range" "(" <Exp> ")" <Block>
type ForRange struct {
	NodeS

	Name Identifier

	// Always not nil
	Range Exp

	Body Block
}

var _ Statement = ForRange{}

func (ForRange) Kind() stm.Kind {
	return stm.ForRange
}

func (s ForRange) Pin() source.Pos {
	return s.Name.Pos
}
