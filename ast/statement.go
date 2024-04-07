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
type nodeStatement struct{}

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

// <AssignStatement> = <Target> "=" <Expression> ";"
//
// <Target> = <ChainOperand>
type AssignStatement struct {
	nodeStatement

	// cannot be call expression
	Target ChainOperand

	Expression Expression
}

// interface implementation check
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

// <AddAssignStatement> = <Target> "+=" <Expression> ";"
//
// <Target> = <ChainOperand>
type AddAssignStatement struct {
	nodeStatement

	// cannot be call expression
	Target ChainOperand

	Expression Expression
}

// interface implementation check
var _ Statement = AddAssignStatement{}

func (AddAssignStatement) Kind() stm.Kind {
	return stm.AddAssign
}

func (s AddAssignStatement) Pin() source.Pos {
	return s.Target.Pin()
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

	If IfClause

	// Equals nil if there are no "else if" clauses in statement
	ElseIf []ElseIfClause

	// Equals nil if there is "else" clause in statement
	Else *ElseClause
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

	Body BlockStatement
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

// <ExpressionStatement> = <Expression> ";"
type ExpressionStatement struct {
	nodeStatement

	// must be call expression
	Expression Expression
}

var _ Statement = ExpressionStatement{}

func (ExpressionStatement) Kind() stm.Kind {
	return stm.Expr
}

func (s ExpressionStatement) Pin() source.Pos {
	return s.Expression.Pin()
}

// <ForStatement> = "for" <BlockStatement>
type ForStatement struct {
	nodeStatement

	Pos source.Pos

	// must contain at least one statement
	Body BlockStatement
}

var _ Statement = ForStatement{}

func (ForStatement) Kind() stm.Kind {
	return stm.For
}

func (s ForStatement) Pin() source.Pos {
	return s.Pos
}

// <ForConditionStatement> = "for" <Expression> <BlockStatement>
type ForConditionStatement struct {
	nodeStatement

	Pos source.Pos

	// Always not nil
	Condition Expression

	// must contain at least one statement
	Body BlockStatement
}

var _ Statement = ForConditionStatement{}

func (ForConditionStatement) Kind() stm.Kind {
	return stm.ForCond
}

func (s ForConditionStatement) Pin() source.Pos {
	return s.Pos
}

type MatchStatement struct {
	nodeStatement

	Pos source.Pos

	// Expression at the top of match statement that is being matched
	Expression Expression

	// Match cases in order they appear
	Cases []MatchCase

	// Always present, but can have zero statements in block
	//
	// Else case is always the last in match statement
	Else BlockStatement
}

type MatchCase struct {
	// Always not nil
	Expression Expression

	Body BlockStatement
}

var _ Statement = MatchStatement{}

func (MatchStatement) Kind() stm.Kind {
	return stm.Match
}

func (s MatchStatement) Pin() source.Pos {
	return s.Pos
}

type JumpStatement struct {
	nodeStatement

	Label Label
}

// Explicit interface implementation check
var _ Statement = JumpStatement{}

func (JumpStatement) Kind() stm.Kind {
	return stm.Jump
}

func (s JumpStatement) Pin() source.Pos {
	return s.Label.Pin()
}

// <ForEachStatement> = "for" <Name> "in" <Expression> <BlockStatement>
type ForEachStatement struct {
	nodeStatement

	Pos source.Pos

	Name Identifier

	// Always not nil
	Iterator Expression

	// must contain at least one statement
	Body BlockStatement
}

var _ Statement = ForEachStatement{}

func (ForEachStatement) Kind() stm.Kind {
	return stm.ForEach
}

func (s ForEachStatement) Pin() source.Pos {
	return s.Pos
}

type LetStatement struct {
	nodeStatement

	Pos source.Pos

	Name Identifier

	Type TypeSpecifier

	// Always not nil
	Expression Expression
}

var _ Statement = LetStatement{}

func (LetStatement) Kind() stm.Kind {
	return stm.Let
}

func (s LetStatement) Pin() source.Pos {
	return s.Pos
}
