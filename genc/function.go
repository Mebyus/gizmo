package genc

import "github.com/mebyus/gizmo/stg"

func (g *Builder) FunParams(params []*stg.Symbol) {
	if len(params) == 0 {
		g.puts("()")
		return
	}

	g.puts("(")
	g.FunParam(params[0])
	for _, p := range params[1:] {
		g.puts(", ")
		g.FunParam(p)
	}
	g.puts(")")
}

func (g *Builder) FunParam(p *stg.Symbol) {
	g.TypeSpec(p.Type)
	g.space()
	g.SymbolName(p)
}

func (g *Builder) FunDec(s *stg.Symbol) {
	def := s.Def.(*stg.FunDef)

	g.puts("static ")
	g.TypeSpec(def.Result)
	g.space()
	g.SymbolName(s)
	g.FunParams(def.Params)
	g.semi()
}

func (g *Builder) FunDef(s *stg.Symbol) {
	def := s.Def.(*stg.FunDef)
	if len(def.Defers) == 0 {
		g.funDefNoDefers(s)
	} else {
		g.funDefWithDefers(s)
	}
}

func (g *Builder) funDefWithDefers(s *stg.Symbol) {
	def := s.Def.(*stg.FunDef)

	g.puts("static ")
	g.TypeSpec(def.Result)
	g.nl()
	g.SymbolName(s)
	g.FunParams(def.Params)
	g.space()
	g.funDefersWrapperBody(def)
}

func (g *Builder) funDefersWrapperBody(def *stg.FunDef) {
	g.puts("{")
	g.nl()
	g.inc()

	g.indent()
	g.puts("<defer args type name>")
	g.puts("ku_defer_args;")
	g.nl()

	g.indent()
	g.puts("<function with defer args setup>")
	g.puts("(")
	for _, p := range def.Params {
		g.SymbolName(p)
		g.puts(", ")
	}
	g.puts("&ku_defer_args")
	g.puts(")")
	g.semi()
	g.nl()

	for i := range len(def.Defers) {
		d := &def.Defers[len(def.Defers)-i-1] // reverse order when defer is actually performed
		g.callDefer(d)
		g.nl()
	}

	g.dec()
	g.indent()
	g.puts("}")
	g.nl()
}

func (g *Builder) callDefer(d *stg.Defer) {
	if d.Uncertain {
		g.callUncertainDefer(d)
	} else {
		g.callCertainDefer(d)
	}
}

func (g *Builder) callCertainDefer(d *stg.Defer) {
	g.indent()
	g.puts("<defer call name>")
	g.deferCallArgs(d)
	g.semi()
}

func (g *Builder) deferCallArgs(d *stg.Defer) {
	if len(d.Params) == 0 {
		g.puts("()")
		return
	}
	g.puts("(")
	g.deferCallArg(d.Index, 0)
	for i := 1; i < len(d.Params); i += 1 {
		g.puts(", ")
		g.deferCallArg(d.Index, i)
	}
	g.puts(")")
}

func (g *Builder) deferCallFlag(i uint32) {
	g.puts("ku_defer_args.defer")
	g.putn(uint64(i))
	g.puts(".call")
}

func (g *Builder) deferCallArg(i uint32, n int) {
	g.puts("ku_defer_args.defer")
	g.putn(uint64(i))
	g.puts(".arg")
	g.putn(uint64(n))
}

func (g *Builder) callUncertainDefer(d *stg.Defer) {
	g.indent()
	g.puts("if (")
	g.deferCallFlag(d.Index)
	g.puts(") {")
	g.nl()
	g.inc()
	g.callCertainDefer(d)
	g.nl()
	g.dec()
	g.indent()
	g.puts("}")
}

func (g *Builder) funDefNoDefers(s *stg.Symbol) {
	def := s.Def.(*stg.FunDef)

	g.puts("static ")
	g.TypeSpec(def.Result)
	g.nl()
	g.SymbolName(s)
	g.FunParams(def.Params)
	g.space()
	g.Block(&def.Body)
}
