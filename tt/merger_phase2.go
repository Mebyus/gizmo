package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
	"github.com/mebyus/gizmo/tt/sym"
)

// Phase 2 merger indexing coagulates info from AST nodes that were skipped during
// phase 1. It also performs additional scans on the tree to determine:
//
//   - top-level symbol resolution
//   - symbol hoisting information inside the unit
//   - used and defined types indexing
//   - type checking of various operations
func (m *Merger) runPhaseTwo() error {
	err := m.bindNodes()
	if err != nil {
		return err
	}
	err = m.scanTypes()
	if err != nil {
		return err
	}
	err = m.scanConstants()
	if err != nil {
		return err
	}
	err = m.scanFns()
	if err != nil {
		return err
	}
	m.unit.Scope.WarnUnused(&Context{m: m})
	return nil
}

func (m *Merger) bindNodes() error {
	for _, top := range m.symBindNodes {
		var err error

		switch top.Kind() {
		case toplvl.Method:
			err = m.bindMethod(top.(ast.Method))
		// case toplvl.Pmb:
		default:
			panic(fmt.Sprintf("unexpected top-level %s node", top.Kind().String()))
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Merger) bindMethod(method ast.Method) error {
	rvName := method.Receiver.Lit
	rvSym := m.unit.Scope.sym(rvName)
	if rvSym == nil {
		return fmt.Errorf("%s: unresolved symbol \"%s\" in method receiver", method.Receiver.Pos.String(), rvName)
	}
	if rvSym.Kind != sym.Type {
		return fmt.Errorf("%s: only types can have methods, but \"%s\" is not a type", method.Receiver.Pos.String(), rvName)
	}

	err := rvSym.Def.(*TempTypeDef).addMethod(method)
	return err
}

func (m *Merger) scanConstants() error {
	for _, s := range m.constants {
		def := s.Def.(TempConstDef)
		c, err := m.scanConst(s, def)
		if err != nil {
			return err
		}
		s.Def = c

		if c.Type != nil {
			s.Type = c.Type
		} else {
			// TODO: this won't work if expression has other constants
			// which are not yet processed
			s.Type = c.Expr.Type()
		}
	}
	return nil
}

func (m *Merger) scanConst(s *Symbol, def TempConstDef) (*ConstDef, error) {
	scope := m.unit.Scope

	t, err := scope.Types.Lookup(def.top.Type)
	if err != nil {
		return nil, err
	}

	expr := def.top.Expression
	if expr == nil {
		panic("nil expression in constant definition")
	}
	ctx := m.newConstCtx()
	e, err := scope.Scan(ctx, expr)
	if err != nil {
		return nil, err
	}

	if ctx.ref.Has(s) {
		return nil, fmt.Errorf("%s: init cycle (constant \"%s\" definition references itself)", s.Pos.String(), s.Name)
	}

	// TODO: check that expression type and constant type match
	refs := ctx.ref.Elems()
	if len(refs) != 0 {
		panic("referencing other symbols in constant definition is not implemented")
	}

	return &ConstDef{
		Expr: e,
		Refs: refs,
		Type: t,
	}, nil
}

func (m *Merger) scanFns() error {
	for _, s := range m.fns {
		def := s.Def.(TempFnDef)
		fn, err := m.scanFn(def)
		if err != nil {
			return err
		}
		s.Def = fn
	}
	return nil
}

func (m *Merger) scanFn(def TempFnDef) (*FnDef, error) {
	scope := m.unit.Scope

	var params []*Symbol
	for _, p := range def.top.Definition.Head.Signature.Params {
		t, err := scope.Types.Lookup(p.Type)
		if err != nil {
			return nil, err
		}

		params = append(params, &Symbol{
			Pos:  p.Name.Pos,
			Name: p.Name.Lit,
			Type: t,
			Kind: sym.Param,
		})
	}

	pos := def.top.Definition.Body.Pos
	result, err := scope.Types.Lookup(def.top.Definition.Head.Signature.Result)
	if err != nil {
		return nil, err
	}
	fn := &FnDef{
		Params: params,
		Result: result,
		Never:  def.top.Definition.Head.Signature.Never,
		Body:   Block{Pos: pos},
	}
	fn.Body.Scope = NewTopScope(m.unit.Scope, &fn.Body.Pos)

	for _, param := range params {
		name := param.Name
		s := fn.Body.Scope.sym(name)
		if s != nil {
			return nil, fmt.Errorf("%s: parameter \"%s\" redeclared in this function", pos.String(), name)
		}
		fn.Body.Scope.Bind(param)
	}

	err = m.scanFnBody(fn, def.top.Definition.Body.Statements)
	if err != nil {
		return nil, err
	}

	return fn, nil
}

func (m *Merger) scanFnBody(def *FnDef, statements []ast.Statement) error {
	if len(statements) == 0 {
		if def.Result != nil {
			return fmt.Errorf("%s: function with return type cannot have empty body", def.Body.Pos.String())
		}
		if def.Never {
			return fmt.Errorf("%s: function marked as never returning cannot have empty body", def.Body.Pos.String())
		}
		return nil
	}

	ctx := m.newFnCtx(def)
	err := def.Body.Fill(ctx, statements)
	if err != nil {
		return err
	}
	def.Refs = ctx.ref.Elems()
	return nil
}
