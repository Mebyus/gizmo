package gencpp

import "github.com/mebyus/gizmo/ast"

func (g *Builder) UnitAtom(atom ast.UnitAtom) {
	for _, block := range atom.Blocks {
		g.NamespaceBlock(block)
		g.nl()
	}
}

func (g *Builder) NamespaceBlock(block ast.NamespaceBlock) {
	g.write("hello, world!")
}
