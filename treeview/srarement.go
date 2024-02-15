package treeview

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/stm"
)

func ConvertStatements(statements []ast.Statement) []Node {
	if len(statements) == 0 {
		return nil
	}

	nodes := make([]Node, 0, len(statements))
	for _, s := range statements {
		nodes = append(nodes, ConvertStatement(s))
	}

	return nodes
}

func ConvertStatement(statement ast.Statement) Node {
	switch statement.Kind() {
	case stm.Return:
		return ConvertReturnStatement(statement.(ast.ReturnStatement))
	case stm.Const:
		return ConvertConstStatement(statement.(ast.ConstStatement))
	case stm.If:
		return ConvertIfStatement(statement.(ast.IfStatement))
	case stm.Assign:
		return ConvertAssignStatement(statement.(ast.AssignStatement))
	default:
		return Node{Text: fmt.Sprintf("<%s statement not implemented>", statement.Kind().String())}
	}
}

func ConvertAssignStatement(statement ast.AssignStatement) Node {
	return Node{
		Text: "assign",
		Nodes: []Node{
			{
				Text:  "target",
				Nodes: []Node{ConvertChainOperand(statement.Target)},
			},
			ConvertExpression(statement.Expression),
		},
	}
}

func ConvertIfStatement(statement ast.IfStatement) Node {
	var clauseNodes []Node
	clauseNodes = append(clauseNodes, ConvertIfClause(statement.If))

	for _, clause := range statement.ElseIf {
		clauseNodes = append(clauseNodes, ConvertElseIfClause(clause))
	}

	if statement.Else != nil {
		clauseNodes = append(clauseNodes, ConvertElseClause(*statement.Else))
	}

	return Node{
		Text:  "if",
		Nodes: clauseNodes,
	}
}

func ConvertIfClause(clause ast.IfClause) Node {
	return Node{
		Text: "main",
		Nodes: []Node{
			ConvertIfClauseCondition(clause.Condition),
			ConvertIfClauseBody(clause.Body),
		},
	}
}

func ConvertIfClauseCondition(expr ast.Expression) Node {
	return Node{
		Text:  "cond",
		Nodes: []Node{ConvertExpression(expr)},
	}
}

func ConvertIfClauseBody(body ast.BlockStatement) Node {
	name := "body"
	if len(body.Statements) == 0 {
		name += ": <empty>"
	}
	return Node{
		Text:  name,
		Nodes: ConvertStatements(body.Statements),
	}
}

func ConvertElseIfClause(clause ast.ElseIfClause) Node {
	return Node{
		Text: "elif",
		Nodes: []Node{
			ConvertIfClauseCondition(clause.Condition),
			ConvertIfClauseBody(clause.Body),
		},
	}
}

func ConvertElseClause(clause ast.ElseClause) Node {
	return Node{
		Text: "else",
		Nodes: []Node{
			ConvertIfClauseBody(clause.Body),
		},
	}
}

func ConvertReturnStatement(statement ast.ReturnStatement) Node {
	title := "return"
	if statement.Expression == nil {
		title += ": <void>"
	}

	var exprNodes []Node
	if statement.Expression != nil {
		exprNodes = []Node{ConvertExpression(statement.Expression)}
	}

	return Node{
		Text:  title,
		Nodes: exprNodes,
	}
}

func ConvertConstStatement(statement ast.ConstStatement) Node {
	nameTitle := "name: "
	if len(statement.Name.Lit) == 0 {
		nameTitle += "<nil>"
	} else {
		nameTitle += statement.Name.Lit
	}
	return Node{
		Text: "const",
		Nodes: []Node{
			{
				Text: nameTitle,
			},
			ConvertTypeSpecifier(statement.Type),
			ConvertExpression(statement.Expression),
		},
	}
}

func ConvertVarInit(v ast.VarInit) Node {
	return Node{
		Text: "var",
		Nodes: []Node{
			{
				Text: "name: " + formatIdentifier(v.Name),
			},
			ConvertTypeSpecifier(v.Type),
			ConvertExpression(v.Expression),
		},
	}
}
