package treeview

import (
	"io"

	"github.com/mebyus/gizmo/ast"
)

func RenderUnitAtom(w io.Writer, atom ast.UnitAtom) error {
	tree := ConvertUnitAtom(atom)
	Render(w, tree)
	return nil
}

func ConvertUnitAtom(atom ast.UnitAtom) Node {
	nodes := make([]Node, 0, len(atom.Blocks)+1)

	if atom.Unit != nil {
		nodes = append(nodes, ConvertUnitBlock(atom.Unit))
	}

	for _, block := range atom.Blocks {
		nodes = append(nodes, ConvertNamespaceBlock(block))
	}

	return Node{
		Text:  "atom",
		Nodes: nodes,
	}
}

func ConvertUnitBlock(block *ast.UnitBlock) Node {
	// uid := block.Block.UID()
	return Node{
		Text:  "unit",
		Nodes: nil,
	}
}

func ConvertNamespaceBlock(block ast.NamespaceBlock) Node {
	nodesTitle := "nodes"
	if len(block.Nodes) == 0 {
		nodesTitle += ": <empty>"
	}
	nodes := make([]Node, 0, len(block.Nodes))
	for _, top := range block.Nodes {
		node := ConvertTopLevel(top)
		nodes = append(nodes, node)
	}

	return Node{
		Text: "namespace",
		Nodes: []Node{
			{
				Text: "name: " + formatScopedIdentifier(block.Name),
			},
			{
				Text:  nodesTitle,
				Nodes: nodes,
			},
		},
	}
}
