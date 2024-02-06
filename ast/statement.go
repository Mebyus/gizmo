package ast

import (
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
)

// <Statement> = <BlockStatement> | <AssignStatement> | <DefineStatement> | <IfStatement> |
//
//	<DeferStatement> | <ExpressionStatement> | <ReturnStatement> | <MatchStatement> |
//	<LoopStatement> | <VarStatement>
type Statement interface {
	Node

	// dummy discriminator method
	Statement()

	Kind() stm.Kind
}

// Dummy provides quick, easy to use implementation of discriminator Statement() method
//
// Used for embedding into other (non-dummy) statement nodes
type nodeStatement struct{ uidHolder }

func (nodeStatement) Statement() {}

type BlockStatement struct {
	nodeStatement

	// starting position of block
	Pos source.Pos

	Statements []Statement
}

// interface implementation check
var _ Statement = BlockStatement{}

func (BlockStatement) Kind() stm.Kind {
	return stm.Block
}

func (s BlockStatement) Pin() source.Pos {
	return s.Pos
}

// <AssignStatement> = <Identifier> "=" <Expression> ";"
type AssignStatement struct {
	nodeStatement

	Target     Identifier
	Expression Expression
}

// interface implementation check
var _ Statement = AssignStatement{}

func (AssignStatement) Kind() stm.Kind {
	return stm.Assign
}

func (s AssignStatement) Pin() source.Pos {
	return s.Target.Pos
}

// <AddAssignStatement> = <Identifier> "+=" <Expression> ";"
type AddAssignStatement struct {
	nodeStatement

	Target     Identifier
	Expression Expression
}

// interface implementation check
var _ Statement = AddAssignStatement{}

func (AddAssignStatement) Kind() stm.Kind {
	return stm.AddAssign
}

func (s AddAssignStatement) Pin() source.Pos {
	return s.Target.Pos
}

// <ReturnStatement> = "return" [ <Expression> ] ";"
type ReturnStatement struct {
	nodeStatement

	Pos source.Pos

	// Equals nil if return does not have expression
	Expression Expression
}

var _ Statement = ReturnStatement{}

func (ReturnStatement) Kind() stm.Kind {
	return stm.Return
}

func (s ReturnStatement) Pin() source.Pos {
	return s.Pos
}

type ConstStatement struct {
	nodeStatement

	Pos source.Pos

	Name Identifier

	Type TypeSpecifier

	// Always not nil
	Expression Expression
}

var _ Statement = ConstStatement{}

func (ConstStatement) Kind() stm.Kind {
	return stm.Const
}

func (s ConstStatement) Pin() source.Pos {
	return s.Pos
}
