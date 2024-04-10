package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/stm"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/tt/scp"
	"github.com/mebyus/gizmo/tt/sym"
)

type Block struct {
	nodeStatement

	// Position where block starts.
	Pos source.Pos

	Nodes []Statement

	Scope *Scope
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
	return b.Scope.CheckUsage()
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
	case stm.Block:
		// g.BlockStatement(statement.(ast.BlockStatement))
	case stm.Return:
		return b.addReturn(ctx, statement.(ast.ReturnStatement))
	case stm.Const:
		// g.ConstStatement(statement.(ast.ConstStatement))
	case stm.Var:
		return b.addVar(ctx, statement.(ast.VarStatement))
	case stm.If:
		// g.IfStatement(statement.(ast.IfStatement))
	case stm.Expr:
		// g.ExpressionStatement(statement.(ast.ExpressionStatement))
	case stm.SymbolAssign:
		return b.addSymbolAssign(ctx, statement.(ast.SymbolAssignStatement))
	case stm.IndirectAssign:
		return b.addIndirectAssign(ctx, statement.(ast.IndirectAssignStatement))
	case stm.Assign:
		// return b.addAssign(ctx, statement.(ast.AssignStatement))
	case stm.AddAssign:
		// g.AddAssignStatement(statement.(ast.AddAssignStatement))
	case stm.For:
		// g.ForStatement(statement.(ast.ForStatement))
	case stm.ForCond:
		// g.ForConditionStatement(statement.(ast.ForConditionStatement))
	case stm.Match:
		// g.MatchStatement(statement.(ast.MatchStatement))
	case stm.Jump:
		// g.JumpStatement(statement.(ast.JumpStatement))
	case stm.ForEach:
		// g.ForEachStatement(statement.(ast.ForEachStatement))
	case stm.Let:
		// g.LetStatement(statement.(ast.LetStatement))
	case stm.Defer:
		//
	default:
		panic(fmt.Sprintf("%s statement not implemented", statement.Kind().String()))
	}
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

	b.addNode(&ReturnStatement{
		Pos: pos,
		// TODO: fill expression
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

	b.addNode(&IndirectAssignStatement{
		Pos:    pos,
		Target: s,
		// TODO: fill expression
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
	if s.Scope.Kind == scp.Unit {
		ctx.ref.Add(s)
	}

	b.addNode(&SymbolAssignStatement{
		Pos:    pos,
		Target: s,
		// TODO: fill expression
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

	s = &Symbol{
		Pos:  pos,
		Name: name,
		Type: ctx.m.lookupType(stmt.Type),
		Kind: sym.Var,
	}
	b.Scope.Bind(s)

	b.addNode(&VarStatement{
		Sym: s,
		// TODO: fill expression
	})
	return nil
}
