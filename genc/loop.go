package genc

import "github.com/mebyus/gizmo/stg"

func (g *Builder) loopStatement(node *stg.Loop) {
	g.indent()
	g.puts("while (true) ")
	g.Block(&node.Body)
}

func (g *Builder) whileStatement(node *stg.While) {
	g.indent()
	g.puts("while (")
	g.Exp(node.Exp)
	g.puts(") ")
	g.Block(&node.Body)
}

func (g *Builder) forRange(node *stg.ForRange) {
	g.indent()
	g.puts("for (")
	g.TypeSpec(node.Sym.Type)
	g.space()
	g.SymbolName(node.Sym)
	g.puts(" = 0; ")
	g.SymbolName(node.Sym)
	g.puts(" < ")
	g.exp(node.Range)
	g.puts("; ")
	g.SymbolName(node.Sym)
	g.puts(" += 1")
	g.puts(") ")
	g.Block(&node.Body)
}
