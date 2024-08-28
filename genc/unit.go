package genc

import "github.com/mebyus/gizmo/stg"

func (g *Builder) Gen(u *stg.Unit) {
	g.BlockTitle(u.Name, "unit "+u.Path.String())

	for _, s := range u.Types {
		g.SymbolSourceComment(s)
		g.TypeDef(s)
		g.nl()

		t := s.Def.(*stg.Type)
		c, ok := g.chunks[t]
		if ok {
			g.ChunkTypeDef(c, t)
			g.nl()
			g.ChunkTypeMethods(c, t)
			g.nl()
		}
	}

	g.Constants(u)
	g.FunDecs(u)
	g.MethodDecs(u)
	g.FunDefs(u)
	g.MethodDefs(u)
}

func (g *Builder) MethodDefs(u *stg.Unit) {
	if len(u.Meds) == 0 {
		return
	}

	g.nl()
	g.BlockTitle(u.Name, "method definitions")
	g.nl()
	for _, s := range u.Meds {
		g.SymbolSourceComment(s)
		g.MethodDef(s)
		g.nl()
	}
}

func (g *Builder) FunDefs(u *stg.Unit) {
	if len(u.Funs) == 0 {
		return
	}

	g.nl()
	g.BlockTitle(u.Name, "function definitions")
	g.nl()
	for _, s := range u.Funs {
		g.SymbolSourceComment(s)
		g.FunDef(s)
		g.nl()
	}
}

func (g *Builder) MethodDecs(u *stg.Unit) {
	if len(u.Meds) == 0 {
		return
	}

	g.nl()
	g.BlockTitle(u.Name, "method declarations")
	g.nl()
	for _, s := range u.Meds {
		g.MethodDec(s)
		g.nl()
	}
}

func (g *Builder) FunDecs(u *stg.Unit) {
	if len(u.Funs) == 0 {
		return
	}

	g.BlockTitle(u.Name, "function declaraions")
	g.nl()
	for _, s := range u.Funs {
		g.FunDec(s)
		g.nl()
	}
}

func (g *Builder) Constants(u *stg.Unit) {
	if len(u.Lets) == 0 {
		return
	}

	g.BlockTitle(u.Name, "constants")
	g.nl()
	for _, s := range u.Lets {
		g.Constant(s)
		g.nl()
	}
}

func (g *Builder) Constant(s *stg.Symbol) {
	def := s.Def.(*stg.ConstDef)

	g.puts("const")
	g.space()
	g.TypeSpec(s.Type)
	g.nl()
	g.SymbolName(s)
	g.puts(" = ")
	g.Exp(def.Exp)
	g.semi()
	g.nl()
}
