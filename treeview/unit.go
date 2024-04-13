package treeview

import (
	"io"

	"github.com/mebyus/gizmo/ast"
)

func RenderUnitAtom(w io.Writer, atom ast.Atom) error {
	tree := ConvertUnitAtom(atom)
	Render(w, tree)
	return nil
}

func ConvertUnitAtom(atom ast.Atom) Node {
	nodes := make([]Node, 0, len(atom.Nodes)+1)

	if atom.Header.Unit != nil {
		nodes = append(nodes, ConvertUnitBlock(atom.Header.Unit))
	}

	for _, top := range atom.Nodes {
		node := ConvertTopLevel(top)
		nodes = append(nodes, node)
	}

	return Node{
		Text:  "atom",
		Nodes: nodes,
	}
}

func ConvertUnitBlock(block *ast.UnitBlock) Node {
	return Node{
		Text:  "unit",
		Nodes: nil,
	}
}
