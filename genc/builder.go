package genc

import (
	"fmt"
	"strconv"
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
	arrays := u.Scope.Types.Arrays()
	g.genBuiltinArrayTypes(arrays)

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

var builtinTypes = []*stg.Type{
	stg.Uint8Type,
	stg.Uint16Type,
	stg.Uint32Type,
	stg.Uint64Type,
	stg.UintType,

	stg.Sint8Type,
	stg.Sint16Type,
	stg.Sint32Type,
	stg.Sint64Type,
	stg.SintType,

	stg.StrType,
	stg.BoolType,
	stg.RuneType,
}

func (g *Builder) genBuiltinArrayTypes(arrays map[*stg.Type][]*stg.Type) {
	for _, t := range builtinTypes {
		list := arrays[t]
		for _, a := range list {
			g.ArrayTypeDef(a, t)
			g.nl()
			g.ArrayTypeMethods(a, t)
			g.nl()
		}
	}
}

func (g *Builder) genBuiltinChunkTypes(chunks map[*stg.Type]*stg.Type) {
	for _, t := range builtinTypes {
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

	def := t.Def.(stg.CustomTypeDef)
	if def.Base.Kind == tpk.Enum {
		g.TypeDefEnum(s, def.Base)
		g.nl()
		return
	}

	g.puts("typedef")
	g.space()

	g.typeSpecForDef(def.Base)

	g.space()
	g.SymbolName(s)
	g.semi()
	g.nl()
}

func (g *Builder) TypeDefEnum(s *stg.Symbol, t *stg.Type) {
	g.EnumType(s, t)
	g.semi()
	g.nl()
	g.nl()

	g.puts("typedef")
	g.space()
	g.SymbolName(t.Def.(*stg.EnumTypeDef).Base.Symbol())
	g.space()
	g.puts(g.enumTypeName(s))
	g.semi()
}

func (g *Builder) enumTypeName(s *stg.Symbol) string {
	return g.tprefix + s.Name
}

func (g *Builder) genEnumTypeName(s *stg.Symbol) string {
	return g.tprefix + "GenEnum" + s.Name
}

func getEnumEntryName(sname, name string) string {
	return "KU_ENUM_" + strings.ToUpper(sname) + "_" + name
}

func (g *Builder) ArrayTypeDef(a, elem *stg.Type) {
	g.puts("typedef")
	g.puts(" struct {")
	g.nl()
	g.inc()

	g.indent()
	g.TypeSpec(elem)
	g.puts(" arr[")
	g.putn(a.Def.(stg.ArrayTypeDef).Len)
	g.puts("];")
	g.nl()

	g.dec()
	g.puts("} ")
	g.puts(g.getArrayTypeName(a))
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
	g.puts("uint len;")
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

func (g *Builder) EnumType(s *stg.Symbol, t *stg.Type) {
	g.puts("enum")
	g.space()
	g.puts(g.genEnumTypeName(s))
	g.space()
	g.enumFields(s.Name, t.Def.(*stg.EnumTypeDef).Entries)
}

func (g *Builder) enumFields(name string, entries []stg.EnumEntry) {
	if len(entries) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.nl()
	g.inc()
	for i := 0; i < len(entries); i += 1 {
		g.indent()
		g.enumField(name, &entries[i])
		g.nl()
	}
	g.dec()
	g.puts("}")
}

func (g *Builder) enumField(name string, entry *stg.EnumEntry) {
	g.puts(getEnumEntryName(name, entry.Name))
	g.puts(" = ")
	g.exp(entry.Exp)
	g.puts(",")
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
	g.Expression(def.Exp)
	g.semi()
	g.nl()
}

func (g *Builder) getSymbolName(s *stg.Symbol) string {
	if s.Scope.Kind == scp.Global {
		return s.Name
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
	def := t.Def.(stg.ChunkTypeDef)
	return g.getChunkTypeNameByElem(def.ElemType)
}

func (g *Builder) getChunkTypeNameByElem(elem *stg.Type) string {
	return g.tprefix + "Chunk_" + g.getSymbolName(elem.Symbol())
}

func (g *Builder) getArrayTypeName(t *stg.Type) string {
	def := t.Def.(stg.ArrayTypeDef)
	return g.tprefix + "Array" + strconv.FormatUint(def.Len, 10) + "_" + g.getSymbolName(def.ElemType.Symbol())
}

func (g *Builder) getTypeSpec(t *stg.Type) string {
	if t == nil {
		return "void"
	}
	if t.PerfectInteger() {
		return "s64"
	}
	if t.Builtin() {
		return g.getSymbolName(t.Symbol())
	}

	switch t.Kind {
	case tpk.Trivial:
		return "struct {}"
	case tpk.Custom:
		return g.getSymbolName(t.Def.(stg.CustomTypeDef).Symbol)
	case tpk.Pointer:
		return g.getTypeSpec(t.Def.(stg.PointerTypeDef).RefType) + "*"
	case tpk.ArrayPointer:
		return g.getTypeSpec(t.Def.(stg.ArrayPointerTypeDef).RefType) + "*"
	case tpk.Chunk:
		return g.getChunkTypeName(t)
	case tpk.Array:
		return g.getArrayTypeName(t)
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

// put decimal formatted integer into output buffer
func (g *Builder) putn(n uint64) {
	g.puts(strconv.FormatUint(n, 10))
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
