package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/stg/scp"
	"github.com/mebyus/gizmo/stg/sfp"
)

type Block struct {
	NodeS

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
	case stm.Call:
		return b.addCall(ctx, statement.(ast.CallStatement))
	case stm.Assign:
		return b.addAssign(ctx, statement.(ast.AssignStatement))
	case stm.ShortInit:
		return b.addShortInit(ctx, statement.(ast.ShortInitStatement))
	case stm.For:
		return b.addFor(ctx, statement.(ast.For))
	case stm.While:
		return b.addForCond(ctx, statement.(ast.ForIf))
	case stm.ForRange:
		return b.addForRange(ctx, statement.(ast.ForRange))
	case stm.Match:
		return b.addMatch(ctx, statement.(ast.MatchStatement))
	case stm.Never:
		return b.addNever(ctx, statement.(ast.NeverStatement))
	// case stm.Jump:
	// g.JumpStatement(statement.(ast.JumpStatement))
	// case stm.ForEach:
	// g.ForEachStatement(statement.(ast.ForEachStatement))
	case stm.Let:
		return b.addLet(ctx, statement.(ast.LetStatement))
	case stm.Defer:
		return b.addDefer(ctx, statement.(ast.DeferStatement))
	default:
		panic(fmt.Sprintf("not implemented for %s statement", statement.Kind().String()))
	}
}

func (b *Block) addForRange(ctx *Context, stmt ast.ForRange) error {
	exp, err := b.Scope.scan(ctx, stmt.Range)
	if err != nil {
		return err
	}
	if !exp.Type().IsIntegerType() {
		return fmt.Errorf("%s: expression inside \"range\" must result in integer type", exp.Pin())
	}

	node := &ForRange{
		Range: exp,
		Body:  Block{Pos: stmt.Body.Pos},
		Sym: &Symbol{
			Pos:  stmt.Name.Pos,
			Name: stmt.Name.Lit,
			Type: UintType,
			Kind: smk.Let,
		},
	}
	node.Body.Scope = NewScope(scp.Loop, b.Scope, &node.Body.Pos)
	node.Body.Scope.Bind(node.Sym)

	err = node.Body.Fill(ctx, stmt.Body.Statements)
	if err != nil {
		return err
	}

	b.addNode(node)
	return nil
}

func (b *Block) addShortInit(ctx *Context, stmt ast.ShortInitStatement) error {
	name := stmt.Name.Lit
	pos := stmt.Name.Pos
	s := b.Scope.Sym(name, pos.Num)
	exp, err := b.Scope.scan(ctx, stmt.Init)
	if err != nil {
		return err
	}
	if s == nil {
		t := exp.Type()
		s = &Symbol{
			Pos:  pos,
			Name: name,
			Type: t,
			Kind: smk.Var,
		}

		// bind occurs after expression scan, because variable
		// that is being defined must not be visible in init expression
		b.Scope.Bind(s)

		b.addNode(&VarStatement{
			Sym: s,
			Exp: exp,
		})
		return nil
	}

	ctx.m.warn(pos, "short init statement used for assignment")
	panic("assign via short init not implemented")
}

func (b *Block) addDefer(ctx *Context, stmt ast.DeferStatement) error {
	pos := stmt.Pos
	if b.Scope.LoopLevel != 0 {
		return fmt.Errorf("%s: defer cannot be used inside a loop", pos)
	}

	exp, err := b.Scope.scanChainOperand(ctx, stmt.Call)
	if err != nil {
		return err
	}
	call := exp.(*CallExp)
	details, err := getCallDetails(call.Callee)
	if err != nil {
		return err
	}

	index := len(ctx.defers)
	uncertain := ctx.rets != 0 || b.Scope.BranchLevel != 0
	b.addNode(&DeferStatement{
		Pos:       pos,
		Args:      call.Arguments,
		Index:     uint32(index),
		Uncertain: uncertain,
	})
	ctx.defers = append(ctx.defers, Defer{
		Params:    details.Params,
		Symbol:    details.Symbol,
		Index:     uint32(index),
		Uncertain: uncertain,
	})
	return nil
}

func (b *Block) addNever(ctx *Context, stmt ast.NeverStatement) error {
	b.addNode(&NeverStatement{
		Pos: stmt.Pos,
	})
	return nil
}

func (b *Block) addMatch(ctx *Context, stmt ast.MatchStatement) error {
	exp, err := b.Scope.scan(ctx, stmt.Exp)
	if err != nil {
		return err
	}
	t := exp.Type()
	switch {
	case t.IsEnumType():
		old := ctx.pushEnum(t.Symbol())
		defer func() {
			ctx.enum = old
		}()
	case t.IsIntegerType():
		// continue execution
	default:
		return fmt.Errorf("%s: only integer or enum types can be matched", exp.Pin())
	}

	cases := make([]MatchCase, 0, len(stmt.Cases))
	for _, mc := range stmt.Cases {
		c, err := b.fillMatchCase(ctx, t, mc)
		if err != nil {
			return err
		}
		cases = append(cases, c)
	}

	elseCase, err := b.fillElseCase(ctx, stmt.Else)
	if err != nil {
		return err
	}

	b.addNode(&MatchStatement{
		Pos:   stmt.Pos,
		Exp:   exp,
		Cases: cases,
		Else:  elseCase,
	})
	return nil
}

func (b *Block) fillMatchCase(ctx *Context, want *Type, mc ast.MatchCase) (MatchCase, error) {
	if len(mc.ExpList) == 0 {
		panic("empty list")
	}

	list := make([]Exp, 0, len(mc.ExpList))
	for _, exp := range mc.ExpList {
		e, err := b.Scope.scan(ctx, exp)
		if err != nil {
			return MatchCase{}, err
		}
		err = typeCheckExp(want, e)
		if err != nil {
			return MatchCase{}, err
		}
		list = append(list, e)
	}

	c := MatchCase{
		Pos:     mc.Pos,
		ExpList: list,
		Body:    Block{Pos: mc.Body.Pos},
	}
	c.Body.Scope = NewScope(scp.Case, b.Scope, &c.Body.Pos)

	err := c.Body.Fill(ctx, mc.Body.Statements)
	if err != nil {
		return MatchCase{}, nil
	}

	return c, nil
}

func (b *Block) fillElseCase(ctx *Context, c *ast.Block) (*Block, error) {
	if c == nil {
		return nil, nil
	}

	block := Block{Pos: c.Pos}
	block.Scope = NewScope(scp.Else, b.Scope, &block.Pos)

	err := block.Fill(ctx, c.Statements)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (b *Block) addFor(ctx *Context, stmt ast.For) error {
	node := &Loop{
		Pos:  stmt.Pos,
		Body: Block{Pos: stmt.Body.Pos},
	}
	node.Body.Scope = NewScope(scp.Loop, b.Scope, &node.Body.Pos)

	err := node.Body.Fill(ctx, stmt.Body.Statements)
	if err != nil {
		return err
	}

	b.addNode(node)
	return nil
}

func (b *Block) addCall(ctx *Context, stmt ast.CallStatement) error {
	o, err := b.Scope.scanChainOperand(ctx, stmt.Call)
	if err != nil {
		return err
	}

	pos := stmt.Call.Identifier.Pos
	b.addNode(&CallStatement{
		Pos:  pos,
		Call: o.(*CallExp),
	})
	return nil
}

func (b *Block) addForCond(ctx *Context, stmt ast.ForIf) error {
	if stmt.If == nil {
		panic("nil condition in for statement")
	}
	condition, err := b.Scope.Scan(ctx, stmt.If)
	if err != nil {
		return err
	}

	node := &While{
		Exp:  condition,
		Body: Block{Pos: stmt.Body.Pos},
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
			return fmt.Errorf("%s: empty return in function which return value is not empty", pos)
		}
		if ctx.never {
			return fmt.Errorf("%s: return used in function which is marked as never returning", pos)
		}

		ctx.rets += 1
		b.addNode(&ReturnStatement{Pos: pos})
		return nil
	}

	if ctx.ret == nil {
		return fmt.Errorf("%s: return with expression in function which does not return a value", pos)
	}
	if ctx.never {
		panic("unreachable: impossible condition")
	}

	exp, err := b.Scope.Scan(ctx, stmt.Expression)
	if err != nil {
		return err
	}
	t := exp.Type()
	if t == nil {
		panic(fmt.Sprintf("%s: %s expression has no type", exp.Pin(), exp.Kind()))
	}
	err = typeCheckExp(ctx.ret, exp)
	if err != nil {
		return err
	}

	ctx.rets += 1
	b.addNode(&ReturnStatement{
		Pos:  pos,
		Expr: exp,
	})
	return nil
}

func (b *Block) addAssign(ctx *Context, stmt ast.AssignStatement) error {
	o, err := b.Scope.scanChainOperand(ctx, stmt.Target)
	if err != nil {
		return err
	}

	expr, err := b.Scope.Scan(ctx, stmt.Expression)
	if err != nil {
		return err
	}

	// TODO: type check for target + operation + expression

	b.addNode(&AssignStatement{
		Target:    o,
		Expr:      expr,
		Operation: stmt.Operator,
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

	t, err := b.Scope.Types.Lookup(ctx, stmt.Type)
	if err != nil {
		return err
	}
	exp, err := b.Scope.Scan(ctx, stmt.Exp)
	if err != nil {
		return err
	}
	if t == nil {
		// infer type from expression if it is specified explicitly
		t = exp.Type()
	}

	s = &Symbol{
		Pos:  pos,
		Name: name,
		Type: t,
		Kind: smk.Var,
	}

	// bind occurs after expression scan, because variable
	// that is being defined must not be visible in init expression
	b.Scope.Bind(s)

	b.addNode(&VarStatement{
		Sym: s,
		Exp: exp,
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

	t, err := b.Scope.Types.Lookup(ctx, stmt.Type)
	if err != nil {
		return err
	}
	exp, err := b.Scope.Scan(ctx, stmt.Exp)
	if err != nil {
		return err
	}
	if t == nil {
		// infer type from expression if it is specified explicitly
		t = exp.Type()
	}

	s = &Symbol{
		Pos:  pos,
		Name: name,
		Type: t,
		Kind: smk.Let,
	}

	if stmt.Exp == nil {
		panic("nil init expression in let statement")
	}

	// bind occurs after expression scan, because variable
	// that is being defined must not be visible in init expression
	b.Scope.Bind(s)

	b.addNode(&LetStatement{
		Sym:  s,
		Expr: exp,
	})
	return nil
}
