package ast

import (
	"github.com/mebyus/gizmo/source"
)

// <UnitBlock> = "unit" <UnitName> "{" { <Statement> ";" } "}"
type UnitBlock struct {
	Name Identifier

	Statements []Statement
}

// <Statement> = <BlockStatement> | <AssignStatement> | <DefineStatement> | <IfStatement> |
//
//	<DeferStatement> | <ExpressionStatement> | <ReturnStatement> | <MatchStatement> |
//	<LoopStatement> | <VarStatement>
type Statement any

// <ReturnStatement> = "return" <Expression> ";"
// type ReturnStatement struct {
// 	Expression Expression
// }

// <SimpleAssignStatement> = <Identifier> "=" <Expression>
type AssignStatement struct {
	Target     Identifier
	Expression Expression
}

// <AddAssign> = <Identifier> "+=" <Expression>
type AddAssign struct {
	Target     Identifier
	Expression Expression
}

type BlockStatement struct {
	// starting position of block
	Pos source.Pos

	Statements []Statement
}

// <IfStatement> = <IfClause> [ <ElseClause> ]
type IfStatement struct {
	If   IfClause
	Else *ElseClause
}

// <IfClause> = "if" <Expression> <BlockStatement>
type IfClause struct {
	Condition Expression
	Body      BlockStatement
}

// <ElseClause> = "else" <BlockStatement>
type ElseClause struct {
	Body BlockStatement
}
