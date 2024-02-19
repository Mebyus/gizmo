package ir

import (
	"strings"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
)

type StructsMap struct {
	// maps full struct name (including namespace) to a list of its methods
	//
	// for structs in default namespace only struct name is used, without namespace
	Methods map[string][]ast.FunctionDeclaration
}

func NamespaceScopes(block ast.NamespaceBlock) []string {
	if block.Default {
		return nil
	}
	s := make([]string, 0, len(block.Name.Scopes)+1)
	for _, scope := range block.Name.Scopes {
		s = append(s, scope.Lit)
	}
	s = append(s, block.Name.Name.Lit)
	return s
}

func IndexMethods(atom ast.UnitAtom) StructsMap {
	m := make(map[string][]ast.FunctionDeclaration)

	for _, block := range atom.Blocks {
		scopes := NamespaceScopes(block)

		for _, node := range block.Nodes {
			if node.Kind() != toplvl.Method {
				continue
			}

			method := node.(ast.Method)
			key := strings.Join(append(scopes, method.Receiver.Lit), "::")
			m[key] = append(m[key], ast.FunctionDeclaration{
				Signature: method.Signature,
				Name:      method.Name,
			})
		}
	}

	if len(m) == 0 {
		m = nil
	}
	return StructsMap{Methods: m}
}
