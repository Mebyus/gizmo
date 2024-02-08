package treeview

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
)

func ConvertExpression(expr ast.Expression) Node {
	return Node{
		Text:  "expr",
		Nodes: []Node{ConvertInnerExpression(expr)},
	}
}

func ConvertInnerExpression(expr ast.Expression) Node {
	switch expr.Kind() {
	case exn.Binary:
		return ConvertBinaryExpression(expr.(ast.BinaryExpression))
	case exn.Subs:
		return CovnertSubsExpression(expr.(ast.SubsExpression))
	default:
		return Node{Text: fmt.Sprintf("<%s expression not implemented>", expr.Kind().String())}
	}
}

func CovnertSubsExpression(expr ast.SubsExpression) Node {
	return Node{
		Text: "subs: " + formatScopedIdentifier(expr.Identifier),
	}
}

func ConvertBinaryExpression(expr ast.BinaryExpression) Node {
	return Node{
		Text: "binary",
		Nodes: []Node{
			{
				Text: "oper: " + expr.Operator.Kind.String(),
			},
			ConvertInnerExpression(expr.Left),
			ConvertInnerExpression(expr.Right),
		},
	}
}
