package treeview

import (
	"fmt"
	"strconv"

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
		return ConvertSubsExpression(expr.(ast.SubsExpression))
	case exn.Basic:
		return ConvertBasicLiteral(expr.(ast.BasicLiteral))
	case exn.Call:
		return ConvertCallExpression(expr.(ast.CallExpression))
	case exn.Paren:
		return ConvertParenthesizedExpression(expr.(ast.ParenthesizedExpression))
	default:
		return Node{Text: fmt.Sprintf("<%s expression not implemented>", expr.Kind().String())}
	}
}

func ConvertCallExpression(expr ast.CallExpression) Node {
	return Node{
		Text: "call",
		Nodes: []Node{
			ConvertCallTarget(expr.Callee),
			ConvertCallArguments(expr.Arguments),
		},
	}
}

func ConvertCallTarget(target ast.ChainOperand) Node {
	return Node{
		Text: "target",
	}
}

func ConvertCallArguments(args []ast.Expression) Node {
	argsTitle := "args"
	if len(args) == 0 {
		argsTitle += ": <void>"
	}

	nodes := make([]Node, 0, len(args))
	for i, arg := range args {
		nodes = append(nodes, Node{
			Text:  strconv.FormatInt(int64(i), 10),
			Nodes: []Node{ConvertExpression(arg)},
		})
	}

	return Node{
		Text:  argsTitle,
		Nodes: nodes,
	}
}

func ConvertSubsExpression(expr ast.SubsExpression) Node {
	return Node{
		Text: "subs: " + formatScopedIdentifier(expr.Identifier),
	}
}

func ConvertBasicLiteral(lit ast.BasicLiteral) Node {
	return Node{
		Text: "lit: " + lit.Token.Literal(),
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

func ConvertChainOperand(operand ast.ChainOperand) Node {
	switch operand.Kind() {
	case exn.Start:
		return ConvertChainStart(operand.(ast.ChainStart))
	case exn.Call:
		return ConvertCallExpression(operand.(ast.CallExpression))
	default:
		return Node{Text: fmt.Sprintf("<%s operand not implemented>", operand.Kind().String())}
	}
}

func ConvertChainStart(start ast.ChainStart) Node {
	return Node{Text: "idn: " + formatScopedIdentifier(start.Identifier)}
}

func ConvertParenthesizedExpression(expr ast.ParenthesizedExpression) Node {
	return Node{
		Text:  "paren",
		Nodes: []Node{ConvertInnerExpression(expr.Inner)},
	}
}
