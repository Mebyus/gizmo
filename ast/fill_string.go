package ast

import "github.com/mebyus/gizmo/source"

type FillString struct {
	NodeO

	// Always has at least one element.
	Parts []FillPart
}

type FillPart interface {
	FillPart()
}

type FillPartString struct {
	Pos source.Pos
	Lit string
}

func (FillPartString) FillPart() {}

type FillPartExp struct {
	Exp Exp
}

func (FillPartExp) FillPart() {}
