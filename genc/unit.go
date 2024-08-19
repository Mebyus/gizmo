package genc

import "github.com/mebyus/gizmo/stg"

func (g *Builder) Gen(u *stg.Unit) {
	g.prelude()

	chunks := u.Scope.Types.Chunks()
	g.genBuiltinChunkTypes(chunks)
	arrays := u.Scope.Types.Arrays()
	g.genBuiltinArrayTypes(arrays)

	for _, s := range u.Types {
		g.SymbolSourceComment(s)
		g.TypeDef(s)
		g.nl()

		t := s.Def.(*stg.Type)
		c, ok := chunks[t]
		if ok {
			g.ChunkTypeDef(c, t)
			g.nl()
			g.ChunkTypeMethods(c, t)
			g.nl()
		}
	}

	for _, s := range u.Lets {
		g.Con(s)
		g.nl()
	}

	g.BlockTitle(u.Name, "function declaraions")
	g.nl()
	for _, s := range u.Funs {
		g.FunDecl(s)
		g.nl()
	}

	g.nl()
	g.BlockTitle(u.Name, "method declarations")
	g.nl()
	for _, s := range u.Meds {
		g.MethodDecl(s)
		g.nl()
	}

	g.nl()
	g.BlockTitle(u.Name, "function implementations")
	g.nl()
	for _, s := range u.Funs {
		g.SymbolSourceComment(s)
		g.FunDef(s)
		g.nl()
	}

	g.nl()
	g.BlockTitle(u.Name, "method implementations")
	g.nl()
	for _, s := range u.Meds {
		g.SymbolSourceComment(s)
		g.MethodDef(s)
		g.nl()
	}
}

func (g *Builder) Con(s *stg.Symbol) {
	def := s.Def.(*stg.ConstDef)

	g.puts("const")
	g.space()
	g.TypeSpec(s.Type)
	g.nl()
	g.SymbolName(s)
	g.puts(" = ")
	g.Expression(def.Exp)
	g.semi()
	g.nl()
}
