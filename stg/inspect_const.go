package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/ast/tps"
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/stg/scp"
)

// performs preliminary unit level constant definition scan in order to obtain data
// necessary for constructing dependency graph between symbols
func (m *Merger) inspectConstant(ctx *SymbolContext) error {
	i := ctx.Symbol.Index()
	node := m.nodes.Con(i)
	err := m.inspectConstTypeSpec(ctx, node.Type)
	if err != nil {
		return err
	}
	if node.Expr == nil {
		panic("nil expression in constant definition")
	}
	return m.inspectExp(ctx, node.Expr)
}

// determines which symbols are used in constant definition type specifier.
func (m *Merger) inspectConstTypeSpec(ctx *SymbolContext, spec ast.TypeSpec) error {
	if spec == nil {
		// constant have static or inferred type
		return nil
	}

	switch spec.Kind() {
	case tps.Name:
		return m.inspectConstTypeName(ctx, spec.(ast.TypeName))
	case tps.Pointer, tps.Array, tps.Struct, tps.ArrayPointer,
		tps.Bag, tps.Chunk, tps.Function, tps.Union:

		return fmt.Errorf("%s: %s type is not allowed in constant definition", spec.Pin(), spec.Kind())
	default:
		panic(fmt.Sprintf("%s types not implemented", spec.Kind()))
	}
}

func (m *Merger) inspectConstTypeName(ctx *SymbolContext, spec ast.TypeName) error {
	name := spec.Name.Lit
	pos := spec.Name.Pos
	s := m.unit.Scope.lookup(name, 0)
	if s == nil {
		return fmt.Errorf("%s: undefined symbol \"%s\"", pos, name)
	}
	if s.Kind != smk.Type {
		return fmt.Errorf("%s: %s symbol \"%s\" cannot be used in constant type specifier",
			pos, s.Kind, name)
	}
	if s.Scope.Kind == scp.Unit {
		ctx.Links.AddDirect(s)
	}
	return nil
}

func (m *Merger) inspectExp(ctx *SymbolContext, exp ast.Expression) error {
	switch exp.Kind() {
	case exn.Basic:
		return nil
	case exn.Symbol:
		return m.inspectSymbolExp(ctx, exp.(ast.SymbolExpression))
	case exn.Unary:
		return m.inspectUnaryExp(ctx, exp.(ast.UnaryExpression))
	case exn.Binary:
		return m.inspectBinaryExp(ctx, exp.(ast.BinaryExpression))
	case exn.Paren:
		return m.inspectParenExp(ctx, exp.(ast.ParenthesizedExpression))
	case exn.Address, exn.Chain, exn.Indirect, exn.Receiver:
		return fmt.Errorf("%s: %s is not allowed in constant expression", exp.Pin(), exp.Kind())
	default:
		panic(fmt.Sprintf("%s expressions not implemented", exp.Kind()))
	}
}

func (m *Merger) inspectSymbolExp(ctx *SymbolContext, exp ast.SymbolExpression) error {
	name := exp.Identifier.Lit
	pos := exp.Identifier.Pos
	s := m.unit.Scope.lookup(name, 0)
	if s == nil {
		return fmt.Errorf("%s: undefined symbol \"%s\"", pos, name)
	}
	if s.Kind != smk.Let {
		return fmt.Errorf("%s: %s symbol \"%s\" cannot be used in constant definition",
			pos, s.Kind, name)
	}
	ctx.Links.AddDirect(s)
	return nil
}

func (m *Merger) inspectUnaryExp(ctx *SymbolContext, exp ast.UnaryExpression) error {
	return m.inspectExp(ctx, exp.Inner)
}

func (m *Merger) inspectParenExp(ctx *SymbolContext, exp ast.ParenthesizedExpression) error {
	return m.inspectExp(ctx, exp.Inner)
}

func (m *Merger) inspectBinaryExp(ctx *SymbolContext, exp ast.BinaryExpression) error {
	err := m.inspectExp(ctx, exp.Left)
	if err != nil {
		return err
	}
	return m.inspectExp(ctx, exp.Right)
}
