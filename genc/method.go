package genc

import "github.com/mebyus/gizmo/stg"

func (g *Builder) MethodDef(s *stg.Symbol) {
	def := s.Def.(*stg.MethodDef)

	g.puts("static ")
	g.TypeSpec(def.Result)
	g.nl()
	g.SymbolName(s)
	g.MethodParams(def.Receiver, def.Params)
	g.space()
	g.Block(&def.Body)
}

func (g *Builder) MethodParams(r *stg.Type, params []*stg.Symbol) {
	g.puts("(")
	g.TypeSpec(r)
	g.space()
	g.puts("g")
	for _, p := range params {
		g.puts(", ")
		g.FunParam(p)
	}
	g.puts(")")
}

func (g *Builder) MethodDec(s *stg.Symbol) {
	def := s.Def.(*stg.MethodDef)

	g.puts("static ")
	g.TypeSpec(def.Result)
	g.space()
	g.SymbolName(s)
	g.MethodParams(def.Receiver, def.Params)
	g.semi()
}

func (g *Builder) ChunkTypeMethods(c, elem *stg.Type) {
	g.ChunkTypeIndexMethod(c, elem)
	g.nl()
	g.ChunkTypeElemMethod(c, elem)
}

func (g *Builder) ArrayTypeMethods(a, elem *stg.Type) {
	g.ArrayTypeElemMethod(a, elem)
	g.nl()
	g.ArrayTypeHeadSliceMethod(a, elem)
}

func (g *Builder) ChunkTypeIndexMethodName(elem *stg.Type) {
	g.puts("ku_chunk_")
	g.TypeSpec(elem)
	g.puts("_index")
}

func (g *Builder) ChunkTypeElemMethodName(elem *stg.Type) {
	g.puts("ku_chunk_")
	g.TypeSpec(elem)
	g.puts("_elem")
}

func (g *Builder) ArrayTypeElemMethodName(t *stg.Type) {
	g.puts("ku_array")
	g.putn(t.Def.(stg.ArrayTypeDef).Len)
	g.puts("_")
	g.TypeSpec(t.ElemType())
	g.puts("_elem")
}

func (g *Builder) ArrayTypeHeadSliceMethodName(t *stg.Type) {
	g.puts("ku_array")
	g.putn(t.Def.(stg.ArrayTypeDef).Len)
	g.puts("_")
	g.TypeSpec(t.ElemType())
	g.puts("_head_slice")
}

func (g *Builder) ChunkTypeIndexMethod(c, elem *stg.Type) {
	g.puts("static ")
	g.TypeSpec(elem)
	g.nl()
	g.ChunkTypeIndexMethodName(elem)

	g.puts("(")

	g.puts(g.getChunkTypeName(c))
	g.puts(" c, u64 i) {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("ku_must(c.ptr != nil);")
	g.nl()

	g.indent()
	g.puts("ku_must(i < c.len);")
	g.nl()

	g.indent()
	g.puts("return c.ptr[i];")
	g.nl()

	g.dec()
	g.puts("}")
	g.nl()
}

func (g *Builder) ChunkTypeElemMethod(c, elem *stg.Type) {
	g.puts("static ")
	g.TypeSpec(elem)
	g.puts("*")
	g.nl()
	g.ChunkTypeElemMethodName(elem)

	g.puts("(")

	g.puts(g.getChunkTypeName(c))
	g.puts(" c, uint i) {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("ku_must(c.ptr != nil);")
	g.nl()

	g.indent()
	g.puts("ku_must(i < c.len);")
	g.nl()

	g.indent()
	g.puts("return c.ptr + i;")
	g.nl()

	g.dec()
	g.puts("}")
	g.nl()
}

func (g *Builder) ArrayTypeElemMethod(a, elem *stg.Type) {
	g.puts("static ")
	g.TypeSpec(elem)
	g.puts("*")
	g.nl()
	g.ArrayTypeElemMethodName(a)

	g.puts("(")

	g.puts(g.getArrayTypeName(a))
	g.puts("* a, uint i) {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("ku_must(i < ")
	g.putn(a.Def.(stg.ArrayTypeDef).Len)
	g.puts(");")
	g.nl()

	g.indent()
	g.puts("return a->arr + i;")
	g.nl()

	g.dec()
	g.puts("}")
	g.nl()
}

func (g *Builder) ArrayTypeHeadSliceMethod(a, elem *stg.Type) {
	g.puts("static ")
	g.puts(g.getChunkTypeNameByElem(elem))
	g.nl()
	g.ArrayTypeHeadSliceMethodName(a)

	g.puts("(")

	g.puts(g.getArrayTypeName(a))
	g.puts("* a, uint i) {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("ku_must(i < ")
	g.putn(a.Def.(stg.ArrayTypeDef).Len)
	g.puts(");")
	g.nl()

	g.indent()
	g.puts(g.getChunkTypeNameByElem(elem))
	g.puts(" s;")
	g.nl()

	g.indent()
	g.puts("s.ptr = a->arr;")
	g.nl()

	g.indent()
	g.puts("s.len = i;")
	g.nl()

	g.indent()
	g.puts("return s;")
	g.nl()

	g.dec()
	g.puts("}")
	g.nl()
}
