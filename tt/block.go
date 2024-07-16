package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/tt/scp"
	"github.com/mebyus/gizmo/tt/sfp"
	"github.com/mebyus/gizmo/tt/sym"
	"github.com/mebyus/gizmo/tt/typ"
)

type Block struct {
	nodeStatement

	// Position where block starts.
	Pos source.Pos

	Nodes []Statement

	Flow []FlowPoint

	Scope *Scope

	FlowKind sfp.Kind
}

// Explicit interface implementation check
var _ Statement = &Block{}

func (b *Block) Pin() source.Pos {
	return b.Pos
}

func (b *Block) Kind() stm.Kind {
	return stm.Block
}

func (b *Block) addNode(node Statement) {
	b.Nodes = append(b.Nodes, node)
}

func (b *Block) Fill(ctx *Context, statements []ast.Statement) error {
	err := b.fill(ctx, statements)
	if err != nil {
		return err
	}
	return b.Scope.CheckUsage(ctx)
}

func (b *Block) fill(ctx *Context, statements []ast.Statement) error {
	i := 0
	for ; i < len(statements); i++ {
		s := statements[i]

		err := b.add(ctx, s)
		if err != nil {
			return err
		}

		if s.Kind() == stm.Return {
			if i+1 < len(statements) {
				pos := statements[i+1].Pin()
				ctx.m.warn(pos, "dead code after return")
			}
			return nil
		}
	}
	return nil
}

func (b *Block) add(ctx *Context, statement ast.Statement) error {
	switch statement.Kind() {
	// case stm.Block:
	// g.BlockStatement(statement.(ast.BlockStatement))
	case stm.Return:
		return b.addReturn(ctx, statement.(ast.ReturnStatement))
	// case stm.Const:
	// g.ConstStatement(statement.(ast.ConstStatement))
	case stm.Var:
		return b.addVar(ctx, statement.(ast.VarStatement))
	case stm.If:
		return b.addIf(ctx, statement.(ast.IfStatement))
	// case stm.Expr:
	// g.ExpressionStatement(statement.(ast.ExpressionStatement))
	case stm.SymbolCall:
		return b.addSymbolCall(ctx, statement.(ast.SymbolCallStatement))
	case stm.SymbolAssign:
		return b.addSymbolAssign(ctx, statement.(ast.SymbolAssignStatement))
	case stm.IndirectAssign:
		return b.addIndirectAssign(ctx, statement.(ast.IndirectAssignStatement))
	// case stm.Assign:
	// return b.addAssign(ctx, statement.(ast.AssignStatement))
	// case stm.AddAssign:
	// g.AddAssignStatement(statement.(ast.AddAssignStatement))
	// case stm.For:
	// g.ForStatement(statement.(ast.ForStatement))
	case stm.ForCond:
		return b.addForCond(ctx, statement.(ast.ForConditionStatement))
	// case stm.Match:
	// g.MatchStatement(statement.(ast.MatchStatement))
	// case stm.Jump:
	// g.JumpStatement(statement.(ast.JumpStatement))
	// case stm.ForEach:
	// g.ForEachStatement(statement.(ast.ForEachStatement))
	case stm.Let:
		return b.addLet(ctx, statement.(ast.LetStatement))
	// case stm.Defer:
	//
	default:
		panic(fmt.Sprintf("not implemented for %s statement", statement.Kind().String()))
	}
}

func (b *Block) addSymbolCall(ctx *Context, stmt ast.SymbolCallStatement) error {
	name := stmt.Callee.Lit
	pos := stmt.Callee.Pos
	s := b.Scope.Lookup(name, pos.Num)
	if s == nil {
		return fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}
	if s.Scope.Kind == scp.Unit {
		ctx.ref.Add(s)
	}
	if s.Kind != sym.Fn {
		return fmt.Errorf("%s: call to symbol \"%s\", which is not a function", pos.String(), name)
	}

	def := s.Def.(*FnDef)
	if len(stmt.Arguments) < len(def.Params) {
		return fmt.Errorf("%s: not enough arguments (got %d) to call \"%s\" function (want %d)",
			pos.String(), len(stmt.Arguments), name, len(def.Params))
	}
	if len(stmt.Arguments) > len(def.Params) {
		return fmt.Errorf("%s: too many arguments (got %d) in function \"%s\" call (want %d)",
			pos.String(), len(stmt.Arguments), name, len(def.Params))
	}
	args, err := b.Scope.scanCallArgs(ctx, def.Params, stmt.Arguments)
	if err != nil {
		return err
	}

	b.addNode(&SymbolCallStatement{
		Pos:       pos,
		Callee:    s,
		Arguments: args,
	})
	return nil
}

func (b *Block) addForCond(ctx *Context, stmt ast.ForConditionStatement) error {
	if stmt.Condition == nil {
		panic("nil condition in for statement")
	}
	condition, err := b.Scope.Scan(ctx, stmt.Condition)
	if err != nil {
		return err
	}

	node := &WhileStatement{
		Pos:       stmt.Pos,
		Condition: condition,
		Body:      Block{Pos: stmt.Body.Pos},
	}
	node.Body.Scope = NewScope(scp.Loop, b.Scope, &node.Body.Pos)

	err = node.Body.Fill(ctx, stmt.Body.Statements)
	if err != nil {
		return err
	}

	b.addNode(node)
	return nil
}

func (b *Block) addIf(ctx *Context, stmt ast.IfStatement) error {
	if len(stmt.ElseIf) != 0 || stmt.Else != nil {
		panic("not implemented")
	}

	if stmt.If.Condition == nil {
		panic("nil condition expression in if statement")
	}
	condition, err := b.Scope.Scan(ctx, stmt.If.Condition)
	if err != nil {
		return err
	}
	if len(stmt.If.Body.Statements) == 0 {
		ctx.m.warn(stmt.If.Body.Pos, "empty if branch")
	}

	node := &SimpleIfStatement{
		Pos:       stmt.If.Pos,
		Condition: condition,
		Body:      Block{Pos: stmt.If.Body.Pos},
	}
	node.Body.Scope = NewScope(scp.If, b.Scope, &node.Body.Pos)

	err = node.Body.Fill(ctx, stmt.If.Body.Statements)
	if err != nil {
		return err
	}

	b.addNode(node)
	return nil
}

func (b *Block) addReturn(ctx *Context, stmt ast.ReturnStatement) error {
	pos := stmt.Pos

	if stmt.Expression == nil {
		if ctx.ret != nil {
			return fmt.Errorf("%s: empty return in function which return value is not empty", pos.String())
		}
		if ctx.never {
			return fmt.Errorf("%s: return used in function which is marked as never returning", pos.String())
		}

		b.addNode(&ReturnStatement{Pos: pos})
		return nil
	}

	if ctx.ret == nil {
		return fmt.Errorf("%s: return with expression in function which does not return a value", pos.String())
	}
	if ctx.never {
		panic("unreachable: impossible condition")
	}

	expr, err := b.Scope.Scan(ctx, stmt.Expression)
	if err != nil {
		return err
	}
	t := expr.Type()
	if t == nil {
		panic(fmt.Sprintf("%s expression has no type", expr.Kind()))
	}

	b.addNode(&ReturnStatement{
		Pos:  pos,
		Expr: expr,
	})
	return nil
}

func (b *Block) addIndirectAssign(ctx *Context, stmt ast.IndirectAssignStatement) error {
	name := stmt.Target.Lit
	pos := stmt.Target.Pos
	s := b.Scope.Lookup(name, pos.Num)
	if s == nil {
		return fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}
	if s.Scope.Kind == scp.Unit {
		ctx.ref.Add(s)
	}

	if s.Type.Kind != typ.Pointer {
		return fmt.Errorf("%s: invalid operation (indirect on non-pointer type)", pos)
	}

	expr, err := b.Scope.Scan(ctx, stmt.Expression)
	if err != nil {
		return err
	}

	b.addNode(&IndirectAssignStatement{
		Pos:    pos,
		Target: s,
		Expr:   expr,
	})
	return nil
}

func (b *Block) addSymbolAssign(ctx *Context, stmt ast.SymbolAssignStatement) error {
	name := stmt.Target.Lit
	pos := stmt.Target.Pos
	s := b.Scope.Lookup(name, pos.Num)
	if s == nil {
		return fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}
	if s.Kind != sym.Var {
		return fmt.Errorf("%s: cannot assign to %s symbol \"%s\"",
			pos.String(), s.Kind.String(), name)
	}
	if s.Scope.Kind == scp.Unit {
		ctx.ref.Add(s)
	}

	expr, err := b.Scope.Scan(ctx, stmt.Expression)
	if err != nil {
		return err
	}

	b.addNode(&SymbolAssignStatement{
		Pos:    pos,
		Target: s,
		Expr:   expr,
	})
	return nil
}

func (b *Block) addVar(ctx *Context, stmt ast.VarStatement) error {
	name := stmt.Name.Lit
	pos := stmt.Name.Pos
	s := b.Scope.Sym(name, pos.Num)
	if s != nil {
		return fmt.Errorf("%s: symbol \"%s\" redeclared in this block", pos.String(), name)
	}

	t, err := b.Scope.Types.Lookup(stmt.Type)
	if err != nil {
		return err
	}

	s = &Symbol{
		Pos:  pos,
		Name: name,
		Type: t,
		Kind: sym.Var,
	}

	expr, err := b.Scope.Scan(ctx, stmt.Expression)
	if err != nil {
		return err
	}

	// bind occurs after expression scan, because variable
	// that is being defined must not be visible in init expression
	b.Scope.Bind(s)

	b.addNode(&VarStatement{
		Sym:  s,
		Expr: expr,
	})
	return nil
}

func (b *Block) addLet(ctx *Context, stmt ast.LetStatement) error {
	name := stmt.Name.Lit
	pos := stmt.Name.Pos
	s := b.Scope.Sym(name, pos.Num)
	if s != nil {
		return fmt.Errorf("%s: symbol \"%s\" redeclared in this block", pos.String(), name)
	}

	t, err := b.Scope.Types.Lookup(stmt.Type)
	if err != nil {
		return err
	}

	s = &Symbol{
		Pos:  pos,
		Name: name,
		Type: t,
		Kind: sym.Let,
	}

	if stmt.Expression == nil {
		panic("nil init expression in let statement")
	}
	expr, err := b.Scope.Scan(ctx, stmt.Expression)
	if err != nil {
		return err
	}

	// bind occurs after expression scan, because variable
	// that is being defined must not be visible in init expression
	b.Scope.Bind(s)

	b.addNode(&LetStatement{
		Sym:  s,
		Expr: expr,
	})
	return nil
}
