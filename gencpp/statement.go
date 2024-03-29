package gencpp

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/lbl"
	"github.com/mebyus/gizmo/ast/stm"
)

// Same as BlockStatement method, but indentation formatting is different
// to start block on the same line
func (g *Builder) Block(block ast.BlockStatement) {
	if len(block.Statements) == 0 {
		g.write("{}")
		return
	}

	g.write("{")
	g.nl()
	g.inc()

	for _, statement := range block.Statements {
		g.Statement(statement)
	}

	g.dec()
	g.indent()
	g.write("}")
}

func (g *Builder) BlockStatement(block ast.BlockStatement) {
	g.indent()
	g.Block(block)
	g.nl()
}

func (g *Builder) Statement(statement ast.Statement) {
	switch statement.Kind() {
	case stm.Block:
		g.BlockStatement(statement.(ast.BlockStatement))
	case stm.Return:
		g.ReturnStatement(statement.(ast.ReturnStatement))
	case stm.Const:
		g.ConstStatement(statement.(ast.ConstStatement))
	case stm.Var:
		g.VarStatement(statement.(ast.VarStatement))
	case stm.If:
		g.IfStatement(statement.(ast.IfStatement))
	case stm.Expr:
		g.ExpressionStatement(statement.(ast.ExpressionStatement))
	case stm.Assign:
		g.AssignStatement(statement.(ast.AssignStatement))
	case stm.AddAssign:
		g.AddAssignStatement(statement.(ast.AddAssignStatement))
	case stm.For:
		g.ForStatement(statement.(ast.ForStatement))
	case stm.ForCond:
		g.ForConditionStatement(statement.(ast.ForConditionStatement))
	case stm.Match:
		g.MatchStatement(statement.(ast.MatchStatement))
	case stm.Jump:
		g.JumpStatement(statement.(ast.JumpStatement))
	case stm.ForEach:
		g.ForEachStatement(statement.(ast.ForEachStatement))
	case stm.Let:
		g.LetStatement(statement.(ast.LetStatement))
	default:
		g.indent()
		g.write(fmt.Sprintf("<%s statement not implemented>", statement.Kind().String()))
		g.nl()
	}
}

func (g *Builder) ForEachStatement(statement ast.ForEachStatement) {
	start, end, ok := analyzeRange(statement.Iterator)

	g.indent()
	g.write("for (")
	g.write("iarch")
	g.space()
	g.Identifier(statement.Name)
	g.write(" = ")

	if ok {
		if start == nil {
			g.write("0")
		} else {
			g.Expression(start)
		}
	} else {
		g.write("<nil>")
	}

	g.write("; ")
	g.Identifier(statement.Name)
	g.write(" < ")

	if ok {
		g.Expression(end)
	} else {
		g.write("<nil>")
	}

	g.write("; ")
	g.Identifier(statement.Name)
	g.write(" += 1")

	g.write(") ")
	g.Block(statement.Body)
	g.nl()
}

func analyzeRange(expr ast.Expression) (start, end ast.Expression, ok bool) {
	call, ok := expr.(ast.CallExpression)
	if !ok {
		return nil, nil, false
	}
	rangeFuncCall, ok := call.Callee.(ast.ChainStart)
	if !ok {
		return nil, nil, false
	}
	if len(rangeFuncCall.Identifier.Scopes) != 0 || rangeFuncCall.Identifier.Name.Lit != "range" {
		return nil, nil, false
	}

	args := call.Arguments
	if len(args) == 0 {
		return nil, nil, false
	}
	if len(args) == 1 {
		return nil, args[0], true
	}
	if len(args) == 2 {
		return args[0], args[1], true
	}
	return nil, nil, false
}

func (g *Builder) JumpStatement(statement ast.JumpStatement) {
	label := statement.Label

	g.indent()
	switch label.Kind() {
	case lbl.Next:
		g.write("continue")
	case lbl.End:
		g.write("break")
	default:
		g.write(fmt.Sprintf("jump <%s>", label.Kind().String()))
	}
	g.semi()
	g.nl()
}

func (g *Builder) MatchStatement(statement ast.MatchStatement) {
	g.indent()

	g.write("switch (")
	g.Expression(statement.Expression)
	g.write(") {")
	g.nl()
	g.nl()

	g.inc()
	for _, c := range statement.Cases {
		g.indent()
		g.write("case ")
		g.Expression(c.Expression)
		g.write(": ")
		g.Block(c.Body)
		g.write(" break;")
		g.nl()
		g.nl()
	}

	g.indent()
	g.write("default: ")
	g.Block(statement.Else)
	g.write(" break;")
	g.dec()
	g.nl()
	g.nl()

	g.indent()
	g.write("}")
	g.nl()
}

func (g *Builder) ForStatement(statement ast.ForStatement) {
	g.indent()
	g.write("while (true) ")
	g.Block(statement.Body)
	g.nl()
}

func (g *Builder) ForConditionStatement(statement ast.ForConditionStatement) {
	g.indent()
	g.write("while (")
	g.Expression(statement.Condition)
	g.write(") ")
	g.Block(statement.Body)
	g.nl()
}

func (g *Builder) AddAssignStatement(statement ast.AddAssignStatement) {
	g.indent()

	g.Expression(statement.Target)
	g.write(" += ")
	g.Expression(statement.Expression)

	g.semi()
	g.nl()
}

func (g *Builder) ExpressionStatement(statement ast.ExpressionStatement) {
	g.indent()

	g.Expression(statement.Expression)

	g.semi()
	g.nl()
}

func (g *Builder) AssignStatement(statement ast.AssignStatement) {
	g.indent()

	g.Expression(statement.Target)
	g.write(" = ")
	g.Expression(statement.Expression)

	g.semi()
	g.nl()
}

func (g *Builder) ReturnStatement(statement ast.ReturnStatement) {
	g.indent()
	g.write("return")

	if statement.Expression != nil {
		g.space()
		g.Expression(statement.Expression)
	}

	g.semi()
	g.nl()
}

func (g *Builder) LetStatement(statement ast.LetStatement) {
	g.indent()
	g.write("const")

	g.space()
	g.TypeSpecifier(statement.Type)

	g.space()
	g.Identifier(statement.Name)

	g.write(" = ")
	g.Expression(statement.Expression)
	g.semi()
	g.nl()
}

func (g *Builder) ConstStatement(statement ast.ConstStatement) {
	g.indent()
	g.ConstInit(statement.ConstInit)
	g.semi()
	g.nl()
}

func (g *Builder) ConstInit(c ast.ConstInit) {
	g.write("constexpr")

	g.space()
	g.TypeSpecifier(c.Type)

	g.space()
	g.Identifier(c.Name)

	g.write(" = ")
	g.Expression(c.Expression)
}

func (g *Builder) VarStatement(statement ast.VarStatement) {
	g.indent()
	g.VarInit(statement.VarInit)
	g.semi()
	g.nl()
}

func (g *Builder) VarInit(v ast.VarInit) {
	g.TypeSpecifier(v.Type)

	g.space()
	g.Identifier(v.Name)

	if v.Expression != nil {
		g.write(" = ")
		g.Expression(v.Expression)
	}
}

func (g *Builder) IfStatement(statement ast.IfStatement) {
	g.indent()
	g.write("if (")
	g.Expression(statement.If.Condition)
	g.write(") ")
	g.Block(statement.If.Body)

	for _, clause := range statement.ElseIf {
		g.write(" else if (")
		g.Expression(clause.Condition)
		g.write(") ")
		g.Block(clause.Body)
	}

	if statement.Else != nil {
		g.write(" else ")
		g.Block(statement.Else.Body)
	}

	g.nl()
}
