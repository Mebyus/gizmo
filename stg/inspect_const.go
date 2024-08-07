package stg

import "github.com/mebyus/gizmo/ast"

// performs preliminary unit level constant definition scan in order to obtain data
// necessary for constructing dependency graph between symbols
func (m *Merger) inspectConstant(ctx *SymbolContext) error {
	i := ctx.Symbol.Index()
	exp := m.nodes.Con(i).Expr
	return m.inspectExp(ctx, exp)
}

func (m *Merger) inspectExp(ctx *SymbolContext, exp ast.Expression) error {
	return nil
}
