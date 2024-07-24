package genc

import (
	"fmt"

	"github.com/mebyus/gizmo/tt"
	"github.com/mebyus/gizmo/tt/scp"
	"github.com/mebyus/gizmo/tt/typ"
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

	// cached type specs
	specs map[*tt.Type]string
}

func (g *Builder) Bytes() []byte {
	return g.buf
}

func (g *Builder) Gen(u *tt.Unit) {
	g.prelude()

	for _, s := range u.Types {
		g.TypeDef(s)
		g.nl()
	}

	for _, s := range u.Cons {
		g.Con(s)
		g.nl()
	}

	for _, s := range u.Funs {
		g.FunDef(s)
		g.nl()
	}
}

func (g *Builder) TypeDef(s *tt.Symbol) {
	t := s.Def.(*tt.Type)

	g.puts("typedef")
	g.space()
	g.typeSpecForDef(t.Base)
	g.space()
	g.SymbolName(s)
	g.semi()
	g.nl()
}

func (g *Builder) typeSpecForDef(t *tt.Type) {
	if t == nil {
		panic("nil type")
	}
	if t.Kind == typ.Struct {
		g.StructType(t.Base.Def.(*tt.StructTypeDef))
		return
	}
	g.TypeSpec(t)
}

func (g *Builder) structFields(members []tt.Member) {
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

func (g *Builder) structField(member *tt.Member) {
	g.TypeSpec(member.Type)
	g.space()
	g.puts(member.Name)
	g.semi()
}

func (g *Builder) StructType(def *tt.StructTypeDef) {
	g.puts("struct")
	g.space()
	g.structFields(def.Members.Members)
}

func (g *Builder) Con(s *tt.Symbol) {
	def := s.Def.(*tt.ConstDef)

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

func (g *Builder) getSymbolName(s *tt.Symbol) string {
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
		return g.prefix + s.Name
	}
	return s.Name
}

func (g *Builder) SymbolName(s *tt.Symbol) {
	g.puts(g.getSymbolName(s))
}

func (g *Builder) getTypeSpec(t *tt.Type) string {
	if t == nil {
		return "void"
	}
	if t.Builtin {
		return g.getSymbolName(t.Symbol)
	}

	switch t.Kind {
	case typ.Trivial:
		return "struct {}"
	case typ.Named:
		return g.getSymbolName(t.Symbol)
	case typ.Pointer:
		return g.getTypeSpec(t.Def.(tt.PointerTypeDef).RefType) + "*"
	case typ.ArrayPointer:
		return g.getTypeSpec(t.Def.(tt.ArrayPointerTypeDef).RefType) + "*"
	default:
		panic(fmt.Sprintf("%s types not implemented", t.Base.Kind.String()))
	}
}

func (g *Builder) TypeSpec(t *tt.Type) {
	s, ok := g.specs[t]
	if !ok {
		s = g.getTypeSpec(t)
		g.specs[t] = s
	}
	g.puts(s)
}

func (g *Builder) FnParams(params []*tt.Symbol) {
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

func (g *Builder) FnParam(p *tt.Symbol) {
	g.TypeSpec(p.Type)
	g.space()
	g.SymbolName(p)
}

func (g *Builder) Block(block *tt.Block) {
	if len(block.Nodes) == 0 {
		g.puts("{}")
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

func (g *Builder) FunDef(s *tt.Symbol) {
	def := s.Def.(*tt.FunDef)

	g.TypeSpec(def.Result)
	g.nl()
	g.SymbolName(s)
	g.FnParams(def.Params)
	g.space()
	g.Block(&def.Body)
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
