package interp

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

type Result struct {
	Files     []string
	TestFiles []string
}

func Interpret(unit ast.UnitBlock) (Result, error) {
	var files []string
	var testFiles []string

	for _, statement := range unit.Block.Statements {
		switch s := statement.(type) {
		case ast.AssignStatement:
			switch s.Target.Lit {
			case "files":
				list, err := getStringsFromExpression(s.Expression)
				if err != nil {
					return Result{}, err
				}
				files = list
			case "test_files":
				list, err := getStringsFromExpression(s.Expression)
				if err != nil {
					return Result{}, err
				}
				testFiles = list
			default:
				return Result{}, fmt.Errorf("reference to udefined symbol: %s (at %s)", s.Target.Lit, s.Target.Pos.String())
			}
		default:
			return Result{}, fmt.Errorf("unexpected statement: %v (%T)", s, s)
		}
	}

	return Result{
		Files:     files,
		TestFiles: testFiles,
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
