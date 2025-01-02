package kbs

import (
	"fmt"
	"io"
	"strconv"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/bop"
	"github.com/mebyus/gizmo/ast/lbl"
	"github.com/mebyus/gizmo/ast/uop"
	"github.com/mebyus/gizmo/token"
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
			g.TopType(atom.Types[node.Index])
		case ast.NodeLet:
			g.TopLet(atom.Constants[node.Index])
		case ast.NodeVar:
			g.TopVar(atom.Vars[node.Index])
		case ast.NodeFun:
			g.Fun(atom.Funs[node.Index])
		case ast.NodeStub:
			g.Stub(atom.Decs[node.Index])
		default:
			panic(fmt.Sprintf("unknown node (%d)", node.Kind))
		}
	}
}

func (g *Generator) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(g.buf)
	return int64(n), err
}

func (g *Generator) TopType(node ast.TopType) {
	g.puts("typedef ")
	g.TypeSpec(node.Spec)
	g.space()
	g.puts(node.Name.Lit)
	g.semi()
	g.nl()
	g.nl()
}

func (g *Generator) funHeader(name ast.Identifier, signature ast.Signature) {
	if signature.Never {
		g.puts("_Noreturn ")
	}
	// if 
	g.puts("static ")
	if signature.Result == nil {
		g.puts("void")
	} else {
		g.TypeSpec(signature.Result)
	}
	g.nl()

	g.puts(name.Lit)
	g.funParams(signature.Params)
}

func (g *Generator) Stub(node ast.TopDec) {
	g.funHeader(node.Name, node.Signature)
	g.semi()
	g.nl()
	g.nl()
}

func (g *Generator) Fun(node ast.TopFun) {
	g.funHeader(node.Name, node.Signature)
	g.space()
	g.Block(node.Body)
	g.nl()
	g.nl()
}

func (g *Generator) TopVar(node ast.TopVar) {
	g.puts("static ")
	g.Var(ast.VarStatement{Var: node.Var})
	g.nl()
}

func (g *Generator) TopLet(node ast.TopLet) {
	g.puts("static ")
	g.Let(ast.LetStatement{Let: node.Let})
	g.nl()
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
}

func (g *Generator) TypeSpec(spec ast.TypeSpec) {
	switch s := spec.(type) {
	case ast.TypeName:
		g.puts(s.Name.Lit)
	case ast.StructType:
		g.Struct(s)
	case ast.PointerType:
		g.PointerType(s)
	case ast.ArrayPointerType:
		g.PointerArray(s)
	case ast.ChunkType:
		g.Chunk(s)
	default:
		panic(fmt.Sprintf("unexpected %s type", spec.Kind()))
	}
}

func (g *Generator) Chunk(c ast.ChunkType) {
	name := c.ElemType.(ast.TypeName).Name.Lit
	switch name {
	case "u8":
		g.puts("bx")
	default:
		panic(fmt.Sprintf("%s chunks not implemented", name))
	}
}

func (g *Generator) PointerType(p ast.PointerType) {
	g.TypeSpec(p.RefType)
	g.puts("*")
}

func (g *Generator) PointerArray(a ast.ArrayPointerType) {
	g.TypeSpec(a.ElemType)
	g.puts("*")
}

func (g *Generator) Struct(s ast.StructType) {
	g.puts("struct {")
	g.nl()
	g.inc()

	for _, field := range s.Fields {
		g.indent()
		g.field(field)
		g.semi()
		g.nl()
	}

	g.dec()
	g.puts("}")
}

func (g *Generator) field(field ast.FieldDefinition) {
	g.TypeSpec(field.Type)
	g.space()
	g.puts(field.Name.Lit)
}

func (g *Generator) Statement(s ast.Statement) {
	switch s := s.(type) {
	case ast.ReturnStatement:
		g.Return(s)
	case ast.VarStatement:
		g.Var(s)
	case ast.LetStatement:
		g.Let(s)
	case ast.AssignStatement:
		g.Assign(s)
	case ast.IfStatement:
		g.If(s)
	case ast.For:
		g.For(s)
	case ast.CallStatement:
		g.Call(s)
	case ast.ForIf:
		g.While(s)
	case ast.JumpStatement:
		g.Jump(s)
	case ast.NeverStatement:
		g.Never(s)
	case ast.ForRange:
		g.ForRange(s)
	case ast.MatchStatement:
		g.Match(s)
	default:
		panic(fmt.Sprintf("unexpected %s statement", s.Kind()))
	}
}

func (g *Generator) Match(m ast.MatchStatement) {
	g.indent()
	g.puts("switch (")
	g.Exp(m.Exp)
	g.puts(") {")
	g.nl()

	for i := range len(m.Cases) {
		c := m.Cases[i]
		g.matchCase(c)
		g.nl()
	}
	g.matchElseCase(m.Else)

	g.indent()
	g.puts("}")
	g.nl()
}

func (g *Generator) matchCase(c ast.MatchCase) {
	exp := c.ExpList[0]
	g.indent()
	g.puts("case ")
	g.Exp(exp)
	g.puts(":")
	for _, exp := range c.ExpList[1:] {
		g.nl()
		g.indent()
		g.puts("case ")
		g.Exp(exp)
		g.puts(":")
	}

	g.space()
	g.Block(c.Body)
	g.nl()
	g.indent()
	g.puts("break;")
	g.nl()
}

func (g *Generator) matchElseCase(c *ast.Block) {
	if c == nil {
		return
	}

	g.indent()
	g.puts("default: ")
	g.Block(*c)
	g.nl()
	g.indent()
	g.puts("break;")
	g.nl()
}

func (g *Generator) ForRange(r ast.ForRange) {
	g.indent()
	g.puts("for (")
	g.puts("uint")
	g.space()
	g.puts(r.Name.Lit)
	g.puts(" = 0; ")
	g.puts(r.Name.Lit)
	g.puts(" < ")
	g.Exp(r.Range)
	g.puts("; ")
	g.puts(r.Name.Lit)
	g.puts(" += 1")
	g.puts(") ")
	g.Block(r.Body)
	g.nl()
	g.nl()
}

func (g *Generator) Never(s ast.NeverStatement) {
	g.indent()
	g.puts("panic_never();")
	g.nl()
}

func (g *Generator) Call(s ast.CallStatement) {
	g.indent()
	g.CallExp(s.Call)
	g.semi()
	g.nl()
}

func (g *Generator) Jump(s ast.JumpStatement) {
	g.indent()
	switch s.Label.(ast.ReservedLabel).ResKind {
	case lbl.Next:
		g.puts("continue;")
	case lbl.Out:
		g.puts("break;")
	default:
		panic(s.Label)
	}
	g.nl()
}

func (g *Generator) For(s ast.For) {
	g.indent()
	g.puts("while (true) ")
	g.Block(s.Body)
	g.nl()
}

func (g *Generator) While(s ast.ForIf) {
	g.indent()
	g.puts("while (")
	g.Exp(s.If)
	g.puts(") ")
	g.Block(s.Body)
	g.nl()
}

func (g *Generator) If(s ast.IfStatement) {
	g.indent()
	g.ifClause(s.If)
	for _, c := range s.ElseIf {
		g.elseIfClause(c)
	}
	if s.Else != nil {
		g.elseClause(s.Else)
	}
	g.nl()
}

func (g *Generator) ifClause(c ast.IfClause) {
	g.puts("if (")
	g.Exp(c.Condition)
	g.puts(") ")
	g.Block(c.Body)
}

func (g *Generator) elseIfClause(c ast.ElseIfClause) {
	g.puts("else ")
	g.ifClause(ast.IfClause(c))
}

func (g *Generator) elseClause(c *ast.ElseClause) {
	g.puts("else ")
	g.Block(c.Body)
}

func (g *Generator) Assign(s ast.AssignStatement) {
	g.indent()
	g.ChainOperand(s.Chain)
	g.space()
	g.puts(s.Operator.String())
	g.space()
	g.Exp(s.Exp)
	g.semi()
	g.nl()
}

func (g *Generator) Let(s ast.LetStatement) {
	g.indent()
	g.puts("const ")
	g.TypeSpec(s.Type)
	g.space()
	g.puts(s.Name.Lit)
	g.puts(" = ")
	g.Exp(s.Exp)
	g.semi()
	g.nl()
}

func (g *Generator) Var(s ast.VarStatement) {
	g.indent()
	array, ok := s.Type.(ast.ArrayType)
	if ok {
		g.VarArray(s, array)
	} else {
		g.TypeSpec(s.Type)
		g.space()
		g.puts(s.Name.Lit)
	}

	_, dirty := s.Exp.(ast.Dirty)
	if dirty {
		g.semi()
		g.nl()
		return
	}

	g.puts(" = ")
	g.Exp(s.Exp)
	g.semi()
	g.nl()
}

func (g *Generator) VarArray(s ast.VarStatement, a ast.ArrayType) {
	g.TypeSpec(a.ElemType)
	g.space()
	g.puts(s.Name.Lit)
	g.puts("[")
	g.Exp(a.Size)
	g.puts("]")
}

func (g *Generator) Return(s ast.ReturnStatement) {
	g.indent()
	if s.Exp == nil {
		g.puts("return;")
		g.nl()
		return
	}

	g.puts("return ")
	g.Exp(s.Exp)
	g.semi()
	g.nl()
}

func (g *Generator) Exp(exp ast.Exp) {
	switch e := exp.(type) {
	case ast.SymbolExp:
		g.SymbolExp(e)
	case ast.BinExp:
		g.BinExp(e)
	case *ast.UnaryExp:
		g.UnaryExp(e)
	case ast.ParenExp:
		g.ParenExp(e)
	case ast.ChainOperand:
		g.ChainOperand(e)
	case ast.CallExp:
		g.CallExp(e)
	case ast.AddressExp:
		g.AddressExp(e)
	case ast.BasicLiteral:
		g.BasicLiteral(e)
	case ast.TintExp:
		g.TintExp(e)
	case ast.CastExp:
		g.CastExp(e)
	default:
		panic(fmt.Sprintf("unexpected %s expression", exp.Kind()))
	}
}

func (g *Generator) AddressExp(exp ast.AddressExp) {
	g.puts("&")
	g.ChainOperand(exp.Chain)
}

func (g *Generator) CastExp(exp ast.CastExp) {
	g.puts("cast(")
	g.TypeSpec(exp.Type)
	g.puts(", ")
	g.Exp(exp.Target)
	g.puts(")")
}

func (g *Generator) TintExp(exp ast.TintExp) {
	g.puts("cast(")
	g.TypeSpec(exp.Type)
	g.puts(", ")
	g.Exp(exp.Target)
	g.puts(")")
}

func (g *Generator) CallExp(exp ast.CallExp) {
	g.ChainOperand(exp.Callee)
	g.callArgs(exp.Args)
}

func (g *Generator) callArgs(args []ast.Exp) {
	if len(args) == 0 {
		g.puts("()")
		return
	}

	g.puts("(")
	g.Exp(args[0])
	for _, arg := range args[1:] {
		g.puts(", ")
		g.Exp(arg)
	}
	g.puts(")")
}

func (g *Generator) ChainOperand(op ast.ChainOperand) {
	g.chainParts(op.Start, op.Parts)
}

func (g *Generator) chainParts(start ast.Identifier, parts []ast.ChainPart) {
	if len(parts) == 0 {
		g.puts(start.Lit)
		return
	}

	part := parts[len(parts)-1]
	rest := parts[:len(parts)-1]
	switch p := part.(type) {
	case ast.SelectPart:
		g.chainParts(start, rest)
		g.puts(".")
		g.puts(p.Name.Lit)
	case ast.IndirectPart:
		g.puts("*")
		g.chainParts(start, rest)
	case ast.IndirectIndexPart:
		g.chainParts(start, rest)
		g.puts("[")
		g.Exp(p.Index)
		g.puts("]")
	case ast.IndirectFieldPart:
		g.chainParts(start, rest)
		g.puts("->")
		g.puts(p.Name.Lit)
	default:
		panic(fmt.Sprintf("%s: unexpected %s chain part", part.Pin(), part.Kind()))
	}
}

func (g *Generator) BasicLiteral(lit ast.BasicLiteral) {
	if lit.Token.Kind == token.String {
		g.puts("make_ss(\"")
		g.puts(lit.Token.Lit)
		g.puts("\", ")
		g.putn(lit.Token.Val)
		g.puts(")")
		return
	}

	g.puts(lit.Token.Literal())
}

func (g *Generator) SymbolExp(exp ast.SymbolExp) {
	g.puts(exp.Identifier.Lit)
}

func (g *Generator) ParenExp(exp ast.ParenExp) {
	g.puts("(")
	g.Exp(exp.Inner)
	g.puts(")")
}

func (g *Generator) UnaryExp(exp *ast.UnaryExp) {
	if exp.Operator.Kind == uop.BitwiseNot {
		g.puts("~")
	} else {
		g.puts(exp.Operator.Kind.String())
	}
	g.Exp(exp.Inner)
}

func (g *Generator) BinExp(exp ast.BinExp) {
	kind := exp.Operator.Kind
	paren := kind == bop.LeftShift || kind == bop.RightShift

	if paren {
		g.puts("(")
	}

	g.Exp(exp.Left)
	g.space()
	g.puts(exp.Operator.Kind.String())
	g.space()
	g.Exp(exp.Right)

	if paren {
		g.puts(")")
	}
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
