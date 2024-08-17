package genc

import (
	"fmt"

	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/stg"
)

func (g *Builder) Statement(node stg.Statement) {
	g.indent()

	switch node.Kind() {
	case stm.Return:
		g.returnStatement(node.(*stg.ReturnStatement))
	case stm.Let:
		g.letStatement(node.(*stg.LetStatement))
	case stm.Var:
		g.varStatement(node.(*stg.VarStatement))
	case stm.Call:
		g.callStatement(node.(*stg.CallStatement))
	case stm.Assign:
		g.assignStatement(node.(*stg.AssignStatement))
	case stm.SimpleIf:
		g.simpleIfStatement(node.(*stg.SimpleIfStatement))
		return
	case stm.For:
		g.loopStatement(node.(*stg.LoopStatement))
		return
	case stm.ForCond:
		g.whileStatement(node.(*stg.WhileStatement))
		return
	case stm.Match:
		g.matchStatement(node.(*stg.MatchStatement))
		return
	default:
		panic(fmt.Sprintf("%s statement not implemented", node.Kind().String()))
	}

	g.semi()
	g.nl()
}

func (g *Builder) matchStatement(node *stg.MatchStatement) {
	g.puts("switch (")
	g.exp(node.Exp)
	g.puts(") {")
	g.nl()

	for i := range len(node.Cases) {
		c := &node.Cases[i]
		g.matchCase(c)
		g.nl()
	}
	g.matchElseCase(node.Else)

	g.indent()
	g.puts("}")
	g.nl()
}

func (g *Builder) matchCase(c *stg.MatchCase) {
	exp := c.ExpList[0]
	g.indent()
	g.puts("case ")
	g.exp(exp)
	g.puts(":")
	for _, exp := range c.ExpList[1:] {
		g.nl()
		g.indent()
		g.puts("case ")
		g.exp(exp)
		g.puts(":")
	}

	g.space()
	g.Block(&c.Body)
	g.indent()
	g.puts("break;")
	g.nl()
}

func (g *Builder) matchElseCase(c *stg.Block) {
	if c == nil {
		return
	}

	g.indent()
	g.puts("default: ")
	g.Block(c)
}

func (g *Builder) loopStatement(node *stg.LoopStatement) {
	g.puts("while (true) ")
	g.Block(&node.Body)
}

func (g *Builder) assignStatement(node *stg.AssignStatement) {
	g.ChainOperandTarget(node.Target)
	g.space()
	g.puts(node.Operation.String())
	g.space()
	g.Expression(node.Expr)
}

// generate target expression for assignment
func (g *Builder) ChainOperandTarget(node stg.ChainOperand) {
	switch node.Kind() {
	case exn.Chain:
		g.ChainSymbol(node.(*stg.ChainSymbol))
	case exn.Member:
		g.MemberExp(node.(*stg.MemberExpression))
	case exn.Indirect:
		g.IndirectExp(node.(*stg.IndirectExpression))
	case exn.IndirectIndex:
		g.IndirectIndexExp(node.(*stg.IndirectIndexExpression))
	case exn.IndirectMember:
		g.IndirectMemberExp(node.(*stg.IndirectMemberExpression))
	case exn.ChunkIndex:
		g.ChunkIndirectElemExp(node.(*stg.ChunkIndexExpression))
	case exn.ArrayIndex:
		g.ArrayIndirectElemExp(node.(*stg.ArrayIndexExp))
	default:
		panic(fmt.Sprintf("unexpected %s operand", node.Kind()))
	}
}

func (g *Builder) whileStatement(node *stg.WhileStatement) {
	g.puts("while (")
	g.Expression(node.Condition)
	g.puts(") ")
	g.Block(&node.Body)
}

func (g *Builder) simpleIfStatement(node *stg.SimpleIfStatement) {
	g.puts("if (")
	g.Expression(node.Condition)
	g.puts(") ")
	g.Block(&node.Body)
}

func (g *Builder) callStatement(node *stg.CallStatement) {
	g.CallExp(node.Call)
}

func (g *Builder) varStatement(node *stg.VarStatement) {
	g.TypeSpec(node.Sym.Type)
	g.space()
	g.SymbolName(node.Sym)
	if node.Expr == nil {
		return
	}
	g.puts(" = ")
	g.Expression(node.Expr)
}

func (g *Builder) letStatement(node *stg.LetStatement) {
	g.TypeSpec(node.Sym.Type)
	g.space()
	g.SymbolName(node.Sym)
	g.puts(" = ")
	g.Expression(node.Expr)
}

func (g *Builder) returnStatement(node *stg.ReturnStatement) {
	g.puts("return")
	if node.Expr == nil {
		return
	}
	g.space()
	g.Expression(node.Expr)
}
