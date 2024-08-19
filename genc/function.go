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

func (g *Builder) FunDecl(s *stg.Symbol) {
	def := s.Def.(*stg.FunDef)

	g.TypeSpec(def.Result)
	g.space()
	g.SymbolName(s)
	g.FunParams(def.Params)
	g.semi()
}

func (g *Builder) FunDef(s *stg.Symbol) {
	def := s.Def.(*stg.FunDef)

	g.TypeSpec(def.Result)
	g.nl()
	g.SymbolName(s)
	g.FunParams(def.Params)
	g.space()
	g.Block(&def.Body)
}
