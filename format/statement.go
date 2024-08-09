package format

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

func (g *Noder) Block(block ast.BlockStatement) {
	if len(block.Statements) == 0 {
		g.emptyBlock(block.Pos)
		return
	}

	g.startBlock(block.Pos)
	for _, node := range block.Statements {
		g.Statement(node)
	}
	g.endBlock()
}

func (g *Noder) BlockStatement(block ast.BlockStatement) {
	g.Block(block)
}

func (g *Noder) emptyBlock(pos source.Pos) {
	g.genpos(token.LeftCurly, pos)
	g.sep()
	g.gen(token.RightCurly)
}

func (g *Noder) Statement(node ast.Statement) {
	g.start()

	switch node.Kind() {
	case stm.Block:
		g.BlockStatement(node.(ast.BlockStatement))
	case stm.If:
		g.IfStatement(node.(ast.IfStatement))
	case stm.Let:
		g.LetStatement(node.(ast.LetStatement))
	case stm.Return:
		g.ReturnStatement(node.(ast.ReturnStatement))
	case stm.Var:
		g.VarStatement(node.(ast.VarStatement))
	case stm.Defer:
		g.DeferStatement(node.(ast.DeferStatement))
	default:
		panic(fmt.Sprintf("node %s statement not implemented", node.Kind().String()))
	}
}

func (g *Noder) DeferStatement(node ast.DeferStatement) {
	g.genpos(token.Defer, node.Pos)
	g.ss()
	g.Expression(node.Call)
	g.semi()
}

func (g *Noder) VarStatement(node ast.VarStatement) {
	g.genpos(token.Var, node.Pos)
	g.ss()
	g.idn(node.Name)
	g.gen(token.Colon)
	g.ss()
	g.TypeSpecifier(node.Type)
	g.ss()
	g.gen(token.Assign)
	g.ss()
	if node.Exp == nil {
		g.gen(token.Dirty)
	} else {
		g.Expression(node.Exp)
	}
	g.semi()
}

func (g *Noder) ReturnStatement(node ast.ReturnStatement) {
	g.genpos(token.Return, node.Pos)

	if node.Expression == nil {
		g.semi()
		return
	}

	g.ss()
	g.Expression(node.Expression)
	g.semi()
}

func (g *Noder) LetStatement(node ast.LetStatement) {
	g.gen(token.Let)
	g.ss()
	g.idn(node.Name)
	g.gen(token.Colon)
	g.ss()
	g.TypeSpecifier(node.Type)
	g.ss()
	g.gen(token.Assign)
	g.ss()
	g.Expression(node.Exp)
	g.semi()
}

func (g *Noder) IfStatement(node ast.IfStatement) {
	g.genpos(token.If, node.If.Pos)
	g.ss()
	g.Expression(node.If.Condition)
	g.ss()
	g.Block(node.If.Body)

	for _, clause := range node.ElseIf {
		g.space()
		g.gen(token.Else)
		g.space()
		g.genpos(token.If, clause.Pos)
		g.space()
		g.Expression(clause.Condition)
		g.space()
		g.Block(clause.Body)
	}

	if node.Else != nil {
		g.space()
		g.genpos(token.Else, node.Else.Pos)
		g.space()
		g.Block(node.Else.Body)
	}
}
