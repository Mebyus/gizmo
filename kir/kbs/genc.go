package kbs

import (
	"fmt"
	"io"
	"strconv"

	"github.com/mebyus/gizmo/ast"
)

type Generator struct {
	buf []byte

	// Indentation buffer.
	//
	// Stores sequence of bytes which is used for indenting current line
	// in output. When a new line starts this buffer is used to add indentation.
	ib []byte
}

func (g *Generator) Reset() {
	// reset buffer position, but keep underlying memory
	g.buf = g.buf[:0]
}

func (g *Generator) Atom(atom *ast.Atom) {
	g.Reset()

	for _, node := range atom.Nodes {
		switch node.Kind {
		case ast.NodeType:
			g.Type(atom.Types[node.Index])
		case ast.NodeLet:
			g.Let(atom.Constants[node.Index])
		case ast.NodeVar:
			g.Var(atom.Vars[node.Index])
		case ast.NodeFun:
			g.Fun(atom.Funs[node.Index])
		default:
			panic(fmt.Sprintf("unknown node (%d)", node.Kind))
		}
	}
}

func (g *Generator) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(g.buf)
	return int64(n), err
}

func (g *Generator) Type(node ast.TopType) {

}

func (g *Generator) Fun(node ast.TopFun) {
	if node.Signature.Never {
		g.puts("_Noreturn ")
	}
	g.puts("static ")
	if node.Signature.Result == nil {
		g.puts("void")
	} else {
		g.TypeSpec(node.Signature.Result)
	}
	g.nl()

	g.puts(node.Name.Lit)
	g.funParams(node.Signature.Params)
	g.space()
	g.Block(node.Body)
	g.nl()
}

func (g *Generator) Var(node ast.TopVar) {

}

func (g *Generator) Let(node ast.TopLet) {

}

func (g *Generator) funParams(params []ast.FieldDefinition) {
	if len(params) == 0 {
		g.puts("(void)")
		return
	}

	g.puts("(")
	g.funParam(params[0])
	for _, param := range params[1:] {
		g.puts(", ")
		g.funParam(param)
	}
	g.puts(")")
}

func (g *Generator) funParam(param ast.FieldDefinition) {
	g.TypeSpec(param.Type)
	g.space()
	g.puts(param.Name.Lit)
}

func (g *Generator) Block(block ast.Block) {
	if len(block.Statements) == 0 {
		g.puts("{}")
		g.nl()
		return
	}

	g.puts("{")
	g.nl()
	g.inc()

	for _, s := range block.Statements {
		g.Statement(s)
	}

	g.dec()
	g.indent()
	g.puts("}")
	g.nl()
}

func (g *Generator) TypeSpec(spec ast.TypeSpec) {
	switch s := spec.(type) {
	case ast.TypeName:
		g.puts(s.Name.Lit)
	default:
		panic(fmt.Sprintf("unexpected %s type", spec.Kind()))
	}
}

func (g *Generator) Statement(s ast.Statement) {
	switch s := s.(type) {
	case ast.ReturnStatement:
		g.Return(s)
	default:
		panic(fmt.Sprintf("unexpected %s statement", s.Kind()))
	}
}

func (g *Generator) Return(s ast.ReturnStatement) {
	g.indent()
	if s.Exp == nil {
		g.puts("return;")
		return
	}

	g.puts("return ")
	g.Exp(s.Exp)
	g.semi()
	g.nl()
}

func (g *Generator) Exp(exp ast.Exp) {

}

// put decimal formatted integer into output buffer
func (g *Generator) putn(n uint64) {
	g.puts(strconv.FormatUint(n, 10))
}

// put string into output buffer
func (g *Generator) puts(s string) {
	g.buf = append(g.buf, s...)
}

// put single byte into output buffer
func (g *Generator) putb(b byte) {
	g.buf = append(g.buf, b)
}

func (g *Generator) put(b []byte) {
	g.buf = append(g.buf, b...)
}

func (g *Generator) nl() {
	g.putb('\n')
}

func (g *Generator) space() {
	g.putb(' ')
}

func (g *Generator) semi() {
	g.putb(';')
}

// increment indentation by one level.
func (g *Generator) inc() {
	g.ib = append(g.ib, '\t')
}

// decrement indentation by one level.
func (g *Generator) dec() {
	g.ib = g.ib[:len(g.ib)-1]
}

// add indentation to current line.
func (g *Generator) indent() {
	g.put(g.ib)
}
