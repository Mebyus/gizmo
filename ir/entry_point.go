package ir

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
)

// CheckEntryPoint returns true if atom contains a function with specified name
// and signature appropriate to be program entrypoint. Entrypoint can only be found
// in default namespace block
func CheckEntryPoint(atom ast.Atom, entry string) bool {
	for _, node := range atom.Nodes {
		if node.Kind() != toplvl.Fn {
			continue
		}

		fn := node.(ast.TopFunctionDefinition)
		if !fn.Public {
			continue
		}
		if fn.Definition.Head.Signature.Result != nil {
			continue
		}
		if len(fn.Definition.Head.Signature.Params) != 0 {
			continue
		}

		if fn.Definition.Head.Name.Lit == entry {
			return true
		}
	}
	return false
}
