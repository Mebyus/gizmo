package gencpp

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ir"
)

func (g *Builder) UnitAtom(atom ast.UnitAtom) {
	for _, block := range atom.Blocks {
		g.currentScopes = ir.NamespaceScopes(block)

		g.NamespaceBlock(block)
		g.nl()
	}
}

// use default namespace if argument is nil
func (g *Builder) namespaceWithPrefix(block *ast.NamespaceBlock) {
	g.write(g.cfg.GlobalNamespacePrefix)
	if block.Default {
		g.write(g.cfg.DefaultNamespace)
	} else {
		g.ScopedIdentifier(block.Name)
	}
}

func (g *Builder) NamespaceBlock(block ast.NamespaceBlock) {
	if len(block.Nodes) == 0 {
		return
	}

	g.write("namespace ")
	g.namespaceWithPrefix(&block)
	g.write(" {")
	g.nl()
	g.nl()

	for _, node := range block.Nodes {
		g.TopLevel(node)
		g.nl()
	}

	g.write("} ")
	if block.Default {
		g.comment("namespace " + g.cfg.DefaultNamespace)
	} else {
		g.comment("namespace " + block.Name.String())
	}
}
