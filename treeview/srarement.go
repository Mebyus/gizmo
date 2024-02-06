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
	default:
		return Node{Text: fmt.Sprintf("<%s statement not implemented>", statement.Kind().String())}
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
