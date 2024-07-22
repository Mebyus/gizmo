package genc

import (
	"fmt"

	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/tt"
)

func (g *Builder) Statement(node tt.Statement) {
	g.indent()

	switch node.Kind() {
	case stm.Return:
		g.returnStatement(node.(*tt.ReturnStatement))
	case stm.Let:
		g.letStatement(node.(*tt.LetStatement))
	case stm.Var:
		g.varStatement(node.(*tt.VarStatement))
	case stm.Call:
		g.callStatement(node.(*tt.CallStatement))
	case stm.Assign:
		g.assignStatement(node.(*tt.AssignStatement))
	case stm.SimpleIf:
		g.simpleIfStatement(node.(*tt.SimpleIfStatement))
		return
	case stm.For:
		g.loopStatement(node.(*tt.LoopStatement))
		return
	case stm.ForCond:
		g.whileStatement(node.(*tt.WhileStatement))
		return
	default:
		panic(fmt.Sprintf("%s statement not implemented", node.Kind().String()))
	}

	g.semi()
	g.nl()
}

func (g *Builder) loopStatement(node *tt.LoopStatement) {
	g.puts("while (true) ")
	g.Block(&node.Body)
}

func (g *Builder) assignStatement(node *tt.AssignStatement) {
	g.ChainOperand(node.Target)
	g.space()
	g.puts(node.Operation.String())
	g.space()
	g.Expression(node.Expr)
}

func (g *Builder) whileStatement(node *tt.WhileStatement) {
	g.puts("while (")
	g.Expression(node.Condition)
	g.puts(") ")
	g.Block(&node.Body)
}

func (g *Builder) simpleIfStatement(node *tt.SimpleIfStatement) {
	g.puts("if (")
	g.Expression(node.Condition)
	g.puts(") ")
	g.Block(&node.Body)
}

func (g *Builder) callStatement(node *tt.CallStatement) {
	g.CallExpression(node.Call)
}

func (g *Builder) varStatement(node *tt.VarStatement) {
	g.TypeSpec(node.Sym.Type)
	g.space()
	g.SymbolName(node.Sym)
	if node.Expr == nil {
		return
	}
	g.puts(" = ")
	g.Expression(node.Expr)
}

func (g *Builder) letStatement(node *tt.LetStatement) {
	g.TypeSpec(node.Sym.Type)
	g.space()
	g.SymbolName(node.Sym)
	g.puts(" = ")
	g.Expression(node.Expr)
}

func (g *Builder) returnStatement(node *tt.ReturnStatement) {
	g.puts("return")
	if node.Expr == nil {
		return
	}
	g.space()
	g.Expression(node.Expr)
}
