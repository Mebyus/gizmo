package genc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mebyus/gizmo/enums/tpk"
	"github.com/mebyus/gizmo/stg"
)

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
