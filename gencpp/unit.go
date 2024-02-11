package gencpp

import "github.com/mebyus/gizmo/ast"

func (g *Builder) UnitAtom(atom ast.UnitAtom) {
	for _, block := range atom.Blocks {
		g.NamespaceBlock(block)
		g.nl()
	}
}

func (g *Builder) NamespaceBlock(block ast.NamespaceBlock) {
	if len(block.Nodes) == 0 {
		return
	}

	g.write("namespace ")
	g.ScopedIdentifier(block.Name)
	g.write(" {")
	g.nl()
	g.nl()

	for _, node := range block.Nodes {
		g.TopLevel(node)
		g.nl()
	}

	g.wb('}')
	g.nl()
}
