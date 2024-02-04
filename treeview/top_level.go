package treeview

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
)

func ConvertTopLevel(top ast.TopLevel) Node {
	switch top.Kind() {
	case toplvl.Declare:
		return ConvertTopFunctionDeclaration(top.(ast.TopFunctionDeclaration))
	case toplvl.Fn:
		return ConvertTopFunctionDefinition(top.(ast.TopFunctionDefinition))
	default:
		return Node{Name: fmt.Sprintf("<top level %s not implemented>", top.Kind().String())}
	}
}

func ConvertTopFunctionDeclaration(top ast.TopFunctionDeclaration) Node {
	return Node{
		Name:  "declare fn",
		Nodes: ConvertFunctionDeclaration(top.Declaration),
	}
}

func ConvertTopFunctionDefinition(top ast.TopFunctionDefinition) Node {
	return Node{
		Name: "fn",
		Nodes: []Node{
			{
				Name:  "head",
				Nodes: ConvertFunctionDeclaration(top.Definition.Head),
			},
			ConvertFunctionBody(top.Definition.Body),
		},
	}
}

func ConvertFunctionBody(body ast.BlockStatement) Node {
	name := "body"
	if len(body.Statements) == 0 {
		name += ": <empty>"
	}
	return Node{
		Name:  name,
		Nodes: ConvertStatements(body.Statements),
	}
}

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
	return Node{Name: fmt.Sprintf("<%s statement not implemented>", statement.Kind().String())}
}

func ConvertFunctionDeclaration(declaration ast.FunctionDeclaration) []Node {
	nodes := make([]Node, 0, 3)

	name := "name"
	if declaration.Name.Lit == "" {
		name += ": <nil>"
	}
	nodes = append(nodes, Node{
		Name: "name: " + declaration.Name.Lit,
	})
	nodes = append(nodes, ConvertFunctionSignature(declaration.Signature)...)
	return nodes
}

func ConvertFunctionSignature(signature ast.FunctionSignature) []Node {
	argsTitle := "args"
	if len(signature.Params) == 0 {
		argsTitle += ": <void>"
	}

	resultTitle := "rest"
	if signature.Never {
		resultTitle += ": <never>"
	} else if signature.Result == nil {
		resultTitle += ": <void>"
	}

	return []Node{
		{
			Name: argsTitle,
		},
		{
			Name: resultTitle,
		},
	}
}
