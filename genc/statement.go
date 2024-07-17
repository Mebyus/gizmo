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
	case stm.SymbolAssign:
		g.symbolAssignStatement(node.(*tt.SymbolAssignStatement))
	case stm.SymbolCall:
		g.symbolCallStatement(node.(*tt.SymbolCallStatement))
	case stm.AddAssign:
		// TODO: refactor assignment statements into single type
		// with operation (assignment) type specified
		g.addAssignStatement(node.(*tt.AddAssignStatement))
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

func (g *Builder) addAssignStatement(node *tt.AddAssignStatement) {
	g.Expression(node.Target)
	g.puts(" += ")
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

func (g *Builder) symbolCallStatement(node *tt.SymbolCallStatement) {
	g.SymbolName(node.Callee)
	g.CallArgs(node.Arguments)
}

func (g *Builder) symbolAssignStatement(node *tt.SymbolAssignStatement) {
	g.SymbolName(node.Target)
	g.puts(" = ")
	g.Expression(node.Expr)
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
