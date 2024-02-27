package interp

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

type Result struct {
	Files     []string
	TestFiles []string
	Imports   []string

	Name string
}

func Interpret(unit *ast.UnitBlock) (*Result, error) {
	var files []string
	var testFiles []string
	var imports []string

	for _, statement := range unit.Block.Statements {
		switch s := statement.(type) {
		case ast.AssignStatement:
			target := s.Target.(ast.ChainStart).Identifier.Name
			switch target.Lit {
			case "files":
				list, err := getStringsFromExpression(s.Expression)
				if err != nil {
					return nil, err
				}
				files = list
			case "test_files":
				list, err := getStringsFromExpression(s.Expression)
				if err != nil {
					return nil, err
				}
				testFiles = list
			case "imports":
				list, err := getStringsFromExpression(s.Expression)
				if err != nil {
					return nil, err
				}
				imports = list
			default:
				return nil, fmt.Errorf("reference to undefined symbol: %s (at %s)", target.Lit, target.Pos.String())
			}
		default:
			return nil, fmt.Errorf("unexpected statement: %v (%T)", s, s)
		}
	}

	return &Result{
		Files:     files,
		TestFiles: testFiles,
		Imports:   imports,

		Name: unit.Name.Lit,
	}, nil
}

func getStringsFromExpression(expr ast.Expression) ([]string, error) {
	list, ok := expr.(ast.ListLiteral)
	if !ok {
		return nil, fmt.Errorf("unexpected expression: %v (%T)", expr, expr)
	}
	var ss []string
	for _, elem := range list.Elems {
		lit, ok := elem.(ast.BasicLiteral)
		if !ok {
			return nil, fmt.Errorf("unexpected list element: %v (%T)", elem, elem)
		}
		if lit.Token.Kind != token.String {
			return nil, fmt.Errorf("unexpected literal inside list: %v (%T)", lit, lit)
		}
		ss = append(ss, lit.Token.Lit)
	}
	return ss, nil
}
