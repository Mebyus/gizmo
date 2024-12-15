package ast

import (
	"github.com/mebyus/gizmo/ast/aop"
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
)

// <AssignStatement> = <Target> <AssignOperator> <Expression> ";"
//
// <Target> = <ChainOperand>
type AssignStatement struct {
	NodeS

	Chain ChainOperand

	Exp Exp

	Operator aop.Kind
}

// Explicit interface implementation check.
var _ Statement = AssignStatement{}

func (AssignStatement) Kind() stm.Kind {
	return stm.Assign
}

func (s AssignStatement) Pin() source.Pos {
	return s.Chain.Start.Pos
}
