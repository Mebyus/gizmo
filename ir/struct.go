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

func Index(atom ast.UnitAtom) Meta {
	sm := make(map[string][]ast.FunctionDeclaration)
	sp := make(map[string]PropsRef)

	for _, block := range atom.Blocks {
		scopes := NamespaceScopes(block)

		for _, node := range block.Nodes {
			if node.Kind() == toplvl.Method {
				method := node.(ast.Method)
				key := strings.Join(append(scopes, method.Receiver.Lit), "::")
				sm[key] = append(sm[key], ast.FunctionDeclaration{
					Signature: method.Signature,
					Name:      method.Name,
				})
			}

			if node.Kind() == toplvl.Declare {
				decl := node.(ast.TopFunctionDeclaration)
				if len(decl.Props) == 0 {
					continue
				}

				ref := NewPropsRef(decl.Props)
				key := strings.Join(append(scopes, decl.Declaration.Name.Lit), "::")
				sp[key] = ref
			}

			if node.Kind() == toplvl.Fn {
				fn := node.(ast.TopFunctionDefinition)
				if len(fn.Props) == 0 {
					continue
				}

				ref := NewPropsRef(fn.Props)
				key := strings.Join(append(scopes, fn.Definition.Head.Name.Lit), "::")
				sp[key] = ref
			}
		}
	}

	if len(sm) == 0 {
		sm = nil
	}
	if len(sp) == 0 {
		sp = nil
	}
	return Meta{
		Structs: StructsMap{Methods: sm},
		Symbols: SymbolsMap{Props: sp},
	}
}
