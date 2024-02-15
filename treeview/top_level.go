package treeview

import (
	"fmt"
	"strconv"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
)

func ConvertTopLevel(top ast.TopLevel) Node {
	switch top.Kind() {
	case toplvl.Declare:
		return ConvertTopFunctionDeclaration(top.(ast.TopFunctionDeclaration))
	case toplvl.Fn:
		return ConvertTopFunctionDefinition(top.(ast.TopFunctionDefinition))
	case toplvl.Var:
		return ConvertTopVar(top.(ast.TopVar))
	default:
		return Node{Text: fmt.Sprintf("<top level %s not implemented>", top.Kind().String())}
	}
}

func ConvertTopFunctionDeclaration(top ast.TopFunctionDeclaration) Node {
	return Node{
		Text:  "declare fn",
		Nodes: ConvertFunctionDeclaration(top.Declaration),
	}
}

func ConvertTopFunctionDefinition(top ast.TopFunctionDefinition) Node {
	return Node{
		Text: "fn",
		Nodes: []Node{
			{
				Text:  "head",
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
		Text:  name,
		Nodes: ConvertStatements(body.Statements),
	}
}

func ConvertFunctionDeclaration(declaration ast.FunctionDeclaration) []Node {
	nodes := make([]Node, 0, 3)

	name := "name"
	if declaration.Name.Lit == "" {
		name += ": <nil>"
	}
	nodes = append(nodes, Node{
		Text: "name: " + declaration.Name.Lit,
	})
	nodes = append(nodes, ConvertFunctionSignature(declaration.Signature)...)
	return nodes
}

func ConvertFunctionSignature(signature ast.FunctionSignature) []Node {
	argsTitle := "args"
	if len(signature.Params) == 0 {
		argsTitle += ": <void>"
	}

	var resultNodes []Node
	resultTitle := "rest"
	if signature.Never {
		resultTitle += ": <never>"
	} else if signature.Result == nil {
		resultTitle += ": <void>"
	} else {
		resultNodes = []Node{ConvertTypeSpecifier(signature.Result)}
	}

	return []Node{
		{
			Text:  argsTitle,
			Nodes: ConvertFunctionParams(signature.Params),
		},
		{
			Text:  resultTitle,
			Nodes: resultNodes,
		},
	}
}

func ConvertFunctionParams(params []ast.FieldDefinition) []Node {
	if len(params) == 0 {
		return nil
	}

	nodes := make([]Node, 0, len(params))
	for i, p := range params {
		nodes = append(nodes, Node{
			Text:  strconv.FormatInt(int64(i), 10),
			Nodes: ConvertFieldDefinition(p),
		})
	}
	return nodes
}

func ConvertTopVar(top ast.TopVar) Node {
	return ConvertVarInit(top.VarInit)
}
