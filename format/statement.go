package format

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

func (g *Builder) Block(block ast.BlockStatement) {
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

func (g *Builder) BlockStatement(block ast.BlockStatement) {
	g.Block(block)
}

func (g *Builder) startBlock(pos source.Pos) {
	g.genpos(token.LeftCurly, pos)
}

func (g *Builder) endBlock() {
	g.gen(token.RightCurly)
}

func (g *Builder) emptyBlock(pos source.Pos) {
	g.genpos(token.LeftCurly, pos)
	g.sep()
	g.gen(token.RightCurly)
}

func (g *Builder) Statement(node ast.Statement) {
	g.start()

	switch node.Kind() {
	case stm.SymbolAssign:
		g.SymbolAssignStatement(node.(ast.SymbolAssignStatement))
	case stm.Block:
		g.BlockStatement(node.(ast.BlockStatement))
	case stm.If:
		g.IfStatement(node.(ast.IfStatement))
	case stm.Let:
		g.LetStatement(node.(ast.LetStatement))
	default:
		panic(fmt.Sprintf("node %s statement not implemented", node.Kind().String()))
	}
}

func (g *Builder) SymbolAssignStatement(node ast.SymbolAssignStatement) {
	g.idn(node.Target)
	g.space()
	g.gen(token.Assign)
	g.space()
	g.Expression(node.Expression)
	g.semi()
}

func (g *Builder) LetStatement(node ast.LetStatement) {
	g.genpos(token.Let, node.Pos)
	g.space()
	g.idn(node.Name)
	g.gen(token.Colon)
	g.space()
	g.TypeSpecifier(node.Type)
	g.space()
	g.gen(token.Assign)
	g.space()
	g.Expression(node.Expression)
	g.semi()
}

func (g *Builder) IfStatement(node ast.IfStatement) {
	g.genpos(token.If, node.If.Pos)
	g.space()
	g.Expression(node.If.Condition)
	g.space()
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
