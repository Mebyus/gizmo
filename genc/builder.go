package genc

import (
	"fmt"
	"strings"

	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/enums/tpk"
	"github.com/mebyus/gizmo/stg"
	"github.com/mebyus/gizmo/stg/scp"
)

type Builder struct {
	buf []byte

	// Indentation buffer.
	//
	// Stores sequence of bytes which is used for indenting current line
	// in output. When a new line starts this buffer is used to add identation.
	ib []byte

	// name prefix for generated symbols
	prefix string

	// name prefix for generated type symbols
	tprefix string

	// cached type specs
	specs map[*stg.Type]string
}

func (g *Builder) Bytes() []byte {
	return g.buf
}

func (g *Builder) Gen(u *stg.Unit) {
	g.prelude()

	chunks := u.Scope.Types.Chunks()
	g.genBuiltinChunkTypes(chunks)

	for _, s := range u.Types {
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
		g.FunDef(s)
		g.nl()
	}

	g.nl()
	g.BlockTitle(u.Name, "method implementations")
	g.nl()
	for _, s := range u.Meds {
		g.MethodDef(s)
		g.nl()
	}
}

func (g *Builder) ChunkTypeMethods(c, elem *stg.Type) {
	g.ChunkTypeIndexMethod(c, elem)
	g.nl()
	g.ChunkTypeElemMethod(c, elem)
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

func (g *Builder) ChunkTypeIndexMethod(c, elem *stg.Type) {
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
	g.TypeSpec(elem)
	g.puts("*")
	g.nl()
	g.ChunkTypeElemMethodName(elem)

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
	g.puts("return c.ptr + i;")
	g.nl()

	g.dec()
	g.puts("}")
	g.nl()
}

func (g *Builder) genBuiltinChunkTypes(chunks map[*stg.Type]*stg.Type) {
	list := []*stg.Type{
		stg.Uint8Type,
		stg.Uint16Type,
		stg.Uint32Type,
		stg.Uint64Type,
		stg.UintType,

		stg.Int8Type,
		stg.Int16Type,
		stg.Int32Type,
		stg.Int64Type,
		stg.IntType,

		stg.StrType,
		stg.BoolType,
		stg.RuneType,
	}

	for _, t := range list {
		c, ok := chunks[t]
		if ok {
			g.ChunkTypeDef(c, t)
			g.nl()
			g.ChunkTypeMethods(c, t)
			g.nl()
		}
	}
}

func (g *Builder) TypeDef(s *stg.Symbol) {
	t := s.Def.(*stg.Type)

	g.puts("typedef")
	g.space()
	g.typeSpecForDef(t.Def.(stg.CustomTypeDef).Base)
	g.space()
	g.SymbolName(s)
	g.semi()
	g.nl()
}

func (g *Builder) ChunkTypeDef(c, elem *stg.Type) {
	g.puts("typedef")
	g.puts(" struct {")
	g.nl()
	g.inc()

	g.indent()
	g.TypeSpec(elem)
	g.puts("* ptr;")
	g.nl()

	g.indent()
	g.puts("u64 len;") // TODO: make this integer arch dependant
	g.nl()

	g.dec()
	g.puts("} ")
	g.puts(g.getChunkTypeName(c))
	g.semi()
	g.nl()
}

func (g *Builder) typeSpecForDef(t *stg.Type) {
	if t == nil {
		panic("nil type")
	}
	if t.Kind == tpk.Struct {
		g.StructType(t.Def.(*stg.StructTypeDef))
		return
	}
	g.TypeSpec(t)
}

func (g *Builder) structFields(members []stg.Member) {
	if len(members) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.nl()
	g.inc()
	for i := 0; i < len(members); i += 1 {
		g.indent()
		g.structField(&members[i])
		g.nl()
	}
	g.dec()
	g.puts("}")
}

func (g *Builder) structField(member *stg.Member) {
	g.TypeSpec(member.Type)
	g.space()
	g.puts(member.Name)
	g.semi()
}

func (g *Builder) StructType(def *stg.StructTypeDef) {
	g.puts("struct")
	g.space()
	g.structFields(def.Members.Members)
}

func (g *Builder) Con(s *stg.Symbol) {
	def := s.Def.(*stg.ConstDef)

	g.puts("const")
	g.space()
	g.TypeSpec(s.Type)
	g.nl()
	g.SymbolName(s)
	g.puts(" = ")
	g.Expression(def.Expr)
	g.semi()
	g.nl()
}

func (g *Builder) getSymbolName(s *stg.Symbol) string {
	if s.Scope.Kind == scp.Global {
		switch s.Name {
		case "int":
			return "i64" // TODO: determine name based on type size
		case "uint":
			return "u64"
		default:
			return s.Name
		}
	}
	if s.Scope.Kind == scp.Unit {
		if s.Kind == smk.Type {
			return g.tprefix + s.Name
		}
		if s.Kind == smk.Method {
			return g.prefix + strings.Replace(s.Name, ".", "_", 1)
		}
		return g.prefix + s.Name
	}
	return s.Name
}

func (g *Builder) SymbolName(s *stg.Symbol) {
	g.puts(g.getSymbolName(s))
}

func (g *Builder) getChunkTypeName(t *stg.Type) string {
	return g.tprefix + "Chunk_" + g.getSymbolName(t.Def.(stg.ChunkTypeDef).ElemType.Symbol())
}

func (g *Builder) getTypeSpec(t *stg.Type) string {
	if t == nil {
		return "void"
	}
	if t.Builtin {
		return g.getSymbolName(t.Symbol())
	}

	switch t.Kind {
	case tpk.Trivial:
		return "struct {}"
	case tpk.Custom:
		return g.getSymbolName(t.Def.(stg.CustomTypeDef).Sym)
	case tpk.Pointer:
		return g.getTypeSpec(t.Def.(stg.PointerTypeDef).RefType) + "*"
	case tpk.ArrayPointer:
		return g.getTypeSpec(t.Def.(stg.ArrayPointerTypeDef).RefType) + "*"
	case tpk.Chunk:
		return g.getChunkTypeName(t)
	default:
		panic(fmt.Sprintf("%s types not implemented", t.Kind.String()))
	}
}

func (g *Builder) TypeSpec(t *stg.Type) {
	s, ok := g.specs[t]
	if !ok {
		s = g.getTypeSpec(t)
		g.specs[t] = s
	}
	g.puts(s)
}

func (g *Builder) FunParams(params []*stg.Symbol) {
	if len(params) == 0 {
		g.puts("()")
		return
	}

	g.puts("(")
	g.FnParam(params[0])
	for _, p := range params[1:] {
		g.puts(", ")
		g.FnParam(p)
	}
	g.puts(")")
}

func (g *Builder) FnParam(p *stg.Symbol) {
	g.TypeSpec(p.Type)
	g.space()
	g.SymbolName(p)
}

func (g *Builder) Block(block *stg.Block) {
	if len(block.Nodes) == 0 {
		g.puts("{}")
		g.nl()
		return
	}

	g.puts("{")
	g.nl()
	g.inc()
	for _, node := range block.Nodes {
		g.Statement(node)
	}
	g.dec()
	g.indent()
	g.puts("}")
	g.nl()
}

func (g *Builder) MethodParams(r *stg.Type, params []*stg.Symbol) {
	g.puts("(")
	g.TypeSpec(r)
	g.space()
	g.puts("g")
	for _, p := range params {
		g.puts(", ")
		g.FnParam(p)
	}
	g.puts(")")
}

func (g *Builder) MethodDecl(s *stg.Symbol) {
	def := s.Def.(*stg.MethodDef)

	g.TypeSpec(def.Result)
	g.space()
	g.SymbolName(s)
	g.MethodParams(def.Receiver, def.Params)
	g.semi()
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

func (g *Builder) MethodDef(s *stg.Symbol) {
	def := s.Def.(*stg.MethodDef)

	g.TypeSpec(def.Result)
	g.nl()
	g.SymbolName(s)
	g.MethodParams(def.Receiver, def.Params)
	g.space()
	g.Block(&def.Body)
}

func (g *Builder) BlockTitle(unit string, s string) {
	g.puts("/* ===== ")
	g.puts(unit)
	g.puts(": ")
	g.puts(s)
	g.puts(" ===== */")
	g.nl()
}

// put string into output buffer
func (g *Builder) puts(s string) {
	g.buf = append(g.buf, s...)
}

// put single byte into output buffer
func (g *Builder) putb(b byte) {
	g.buf = append(g.buf, b)
}

func (g *Builder) put(b []byte) {
	g.buf = append(g.buf, b...)
}

func (g *Builder) nl() {
	g.putb('\n')
}

func (g *Builder) space() {
	g.putb(' ')
}

func (g *Builder) semi() {
	g.putb(';')
}

// increment indentation by one level.
func (g *Builder) inc() {
	g.ib = append(g.ib, '\t')
}

// decrement indentation by one level.
func (g *Builder) dec() {
	g.ib = g.ib[:len(g.ib)-1]
}

// add indentation to current line.
func (g *Builder) indent() {
	g.put(g.ib)
}
