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

type Block struct {
	nodeStatement

	// starting position of block
	Pos source.Pos

	Statements []Statement
}

// interface implementation check
var _ Statement = Block{}

func (Block) Kind() stm.Kind {
	return stm.Block
}

func (s Block) Pin() source.Pos {
	return s.Pos
}

// <ReturnStatement> = "return" [ <Expression> ] ";"
type ReturnStatement struct {
	nodeStatement

	Pos source.Pos

	// Equals nil if return does not have expression.
	Expression Expression
}

var _ Statement = ReturnStatement{}

func (ReturnStatement) Kind() stm.Kind {
	return stm.Return
}

func (s ReturnStatement) Pin() source.Pos {
	return s.Pos
}

// <NeverStatement> = "never"
type NeverStatement struct {
	nodeStatement

	Pos source.Pos
}

var _ Statement = NeverStatement{}

func (NeverStatement) Kind() stm.Kind {
	return stm.Never
}

func (s NeverStatement) Pin() source.Pos {
	return s.Pos
}

type Var struct {
	Pos source.Pos

	Name Identifier

	Type TypeSpec

	// Equals nil if init expression is dirty
	Exp Expression
}

type VarStatement struct {
	nodeStatement

	Var
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
	Body      Block
}

// <IfClause> = "else" "if" <Expression> <BlockStatement>
type ElseIfClause struct {
	Pos source.Pos

	// Always not nil
	Condition Expression

	Body Block
}

// <ElseClause> = "else" <BlockStatement>
type ElseClause struct {
	Pos source.Pos

	Body Block
}

var _ Statement = IfStatement{}

func (IfStatement) Kind() stm.Kind {
	return stm.If
}

func (s IfStatement) Pin() source.Pos {
	return s.If.Pos
}

// <CallStatement> = <ChainOperand> ";"
type CallStatement struct {
	nodeStatement

	// Must be call operand.
	Call ChainOperand
}

var _ Statement = CallStatement{}

func (CallStatement) Kind() stm.Kind {
	return stm.Call
}

func (s CallStatement) Pin() source.Pos {
	return s.Call.Pin()
}

// <ForStatement> = "for" <BlockStatement>
type ForStatement struct {
	nodeStatement

	Pos source.Pos

	Body Block
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

	// Always not nil.
	Condition Expression

	Body Block
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

	// Exp at the top of match statement that is being matched.
	Exp Expression

	// Match cases in order they appear.
	Cases []MatchCase

	// Else is always last in match statement.
	Else *Block
}

type MatchCase struct {
	Pos source.Pos

	// Always has at least one element.
	ExpList []Expression

	Body Block
}

var _ Statement = MatchStatement{}

func (MatchStatement) Kind() stm.Kind {
	return stm.Match
}

func (s MatchStatement) Pin() source.Pos {
	return s.Pos
}

type MatchBoolStatement struct {
	nodeStatement

	Pos source.Pos

	// Match cases in order they appear.
	Cases []MatchBoolCase

	// Else is always last in match statement.
	Else *Block
}

type MatchBoolCase struct {
	Pos source.Pos

	// Always not nil.
	Exp Expression

	Body Block
}

var _ Statement = MatchBoolStatement{}

func (MatchBoolStatement) Kind() stm.Kind {
	return stm.MatchBool
}

func (s MatchBoolStatement) Pin() source.Pos {
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
	Body Block
}

var _ Statement = ForEachStatement{}

func (ForEachStatement) Kind() stm.Kind {
	return stm.ForEach
}

func (s ForEachStatement) Pin() source.Pos {
	return s.Pos
}

type Let struct {
	Name Identifier

	// Equals nil if constant type is not specified explicitly.
	Type TypeSpec

	// Expression that defines constant value. Always not nil.
	Exp Expression
}

type LetStatement struct {
	nodeStatement

	Let
}

var _ Statement = LetStatement{}

func (LetStatement) Kind() stm.Kind {
	return stm.Let
}

func (s LetStatement) Pin() source.Pos {
	return s.Name.Pos
}

// <DeferStatement> = "defer" <CallExpression> ";"
type DeferStatement struct {
	nodeStatement

	Pos source.Pos

	// must be call expression
	Call ChainOperand
}

var _ Statement = DeferStatement{}

func (DeferStatement) Kind() stm.Kind {
	return stm.Defer
}

func (s DeferStatement) Pin() source.Pos {
	return s.Pos
}
