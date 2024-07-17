package genc

import (
	"fmt"

	"github.com/mebyus/gizmo/tt"
	"github.com/mebyus/gizmo/tt/scp"
	"github.com/mebyus/gizmo/tt/sym"
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

	for _, s := range u.Scope.Symbols {
		switch s.Kind {
		case sym.Fn:
			g.FnDef(s, s.Def.(*tt.FnDef))
		case sym.Const:
			g.Const(s, s.Def.(*tt.ConstDef))
		case sym.Type:
			g.Type(s, s.Def.(*tt.Type))
		default:
			panic(fmt.Sprintf("not implemented for %s symbols", s.Kind.String()))
		}

		g.nl()
	}
}

func (g *Builder) Type(s *tt.Symbol, t *tt.Type) {
	g.puts("typedef")
	g.space()
	switch t.Base.Kind {
	case typ.Trivial:
		g.puts("struct {}")
	case typ.Struct:
		g.StructType(t.Base.Def.(*tt.StructTypeDef))
	default:
		panic(fmt.Sprintf("%s types not implemented", t.Base.Kind.String()))
	}
	g.space()
	g.SymbolName(s)
	g.semi()
	g.nl()
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

func (g *Builder) Const(s *tt.Symbol, def *tt.ConstDef) {
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
	if t.Kind == typ.Boolean {
		return "bool"
	}
	if t.Kind == typ.Unsigned {
		size := t.Def.(tt.IntTypeDef).Size
		switch size {
		case 1:
			return "u8"
		case 2:
			return "u16"
		case 4:
			return "u32"
		case 8:
			return "u64"
		case 16:
			return "u128"
		default:
			panic(fmt.Sprintf("unxpected size %d", size))
		}
	}
	if t.Kind == typ.Signed {
		size := t.Def.(tt.IntTypeDef).Size
		switch size {
		case 1:
			return "i8"
		case 2:
			return "i16"
		case 4:
			return "i32"
		case 8:
			return "i64"
		case 16:
			return "i128"
		default:
			panic(fmt.Sprintf("unxpected size %d", size))
		}

	}
	if t.Kind == typ.Named {
		return g.getSymbolName(t.Def.(tt.NamedTypeDef).Symbol)
	}
	if t.Kind == typ.Pointer {
		return g.getTypeSpec(t.Def.(tt.PtrTypeDef).RefType) + "*"
	}

	panic(fmt.Sprintf("%s types not implemented", t.Kind))
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

func (g *Builder) FnDef(s *tt.Symbol, def *tt.FnDef) {
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
