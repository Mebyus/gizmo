package ast

import (
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
)

// <AssignStatement> = <Target> "=" <Expression> ";"
//
// <Target> = <ChainOperand>
type AssignStatement struct {
	nodeStatement

	// cannot be call expression
	Target ChainOperand

	Expression Expression
}

// Explicit interface implementation check
var _ Statement = AssignStatement{}

func (AssignStatement) Kind() stm.Kind {
	return stm.Assign
}

func (s AssignStatement) Pin() source.Pos {
	return s.Target.Pin()
}

// <SymbolAssignStatement> = <Identifier> "=" <Expression> ";"
type SymbolAssignStatement struct {
	nodeStatement

	Target Identifier

	Expression Expression
}

// Explicit interface implementation check
var _ Statement = SymbolAssignStatement{}

func (SymbolAssignStatement) Kind() stm.Kind {
	return stm.SymbolAssign
}

func (s SymbolAssignStatement) Pin() source.Pos {
	return s.Target.Pos
}

// <IndirectAssignStatement> = <Identifier> ".@" "=" <Expression> ";"
type IndirectAssignStatement struct {
	nodeStatement

	Target Identifier

	Expression Expression
}

// Explicit interface implementation check
var _ Statement = IndirectAssignStatement{}

func (IndirectAssignStatement) Kind() stm.Kind {
	return stm.IndirectAssign
}

func (s IndirectAssignStatement) Pin() source.Pos {
	return s.Target.Pos
}

// <AddAssignStatement> = <Target> "+=" <Expression> ";"
//
// <Target> = <ChainOperand>
type AddAssignStatement struct {
	nodeStatement

	// cannot be call expression
	Target ChainOperand

	Expression Expression
}

// Explicit interface implementation check
var _ Statement = AddAssignStatement{}

func (AddAssignStatement) Kind() stm.Kind {
	return stm.AddAssign
}

func (s AddAssignStatement) Pin() source.Pos {
	return s.Target.Pin()
}
