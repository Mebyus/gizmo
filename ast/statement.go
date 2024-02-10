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

type ConstInit struct {
	Pos source.Pos

	Name Identifier

	Type TypeSpecifier

	// Always not nil
	Expression Expression
}

type ConstStatement struct {
	nodeStatement

	ConstInit
}

var _ Statement = ConstStatement{}

func (ConstStatement) Kind() stm.Kind {
	return stm.Const
}

func (s ConstStatement) Pin() source.Pos {
	return s.Pos
}

type VarInit struct {
	Pos source.Pos

	Name Identifier

	Type TypeSpecifier

	// Equals nil if init expression is dirty
	Expression Expression
}

type VarStatement struct {
	nodeStatement

	VarInit
}

var _ Statement = VarStatement{}

func (VarStatement) Kind() stm.Kind {
	return stm.Var
}

func (s VarStatement) Pin() source.Pos {
	return s.Pos
}

// <IfStatement> = <IfClause> { <ElseIfClause> } [ <ElseClause> ]
type IfStatement struct {
	nodeStatement

	If     IfClause

	// Equals nil if there are no "else if" clauses in statement
	ElseIf []ElseIfClause

	// Equals nil if there is "else" clause in statement
	Else   *ElseClause
}

// <IfClause> = "if" <Expression> <BlockStatement>
type IfClause struct {
	Pos source.Pos

	// Always not nil
	Condition Expression
	Body      BlockStatement
}

// <IfClause> = "else" "if" <Expression> <BlockStatement>
type ElseIfClause struct {
	Pos source.Pos

	// Always not nil
	Condition Expression
	Body      BlockStatement
}

// <ElseClause> = "else" <BlockStatement>
type ElseClause struct {
	Pos source.Pos

	Body BlockStatement
}

var _ Statement = IfStatement{}

func (IfStatement) Kind() stm.Kind {
	return stm.If
}

func (s IfStatement) Pin() source.Pos {
	return s.If.Pos
}
