package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/tt/sym"
)

// Phase 2 merger indexing coagulates info from AST nodes that were skipped during
// phase 1. It also performs additional scans on the tree to determine:
//
//   - top-level symbol resolution
//   - symbol hoisting information inside the unit
//   - used and defined types indexing
//   - type checking of various operations
func (m *Merger) merge() error {
	err := m.bindMethods()
	if err != nil {
		return err
	}
	err = m.shallowScanTypes()
	if err != nil {
		return err
	}
	err = m.bindTypes()
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

// Checks that each method has a custom type receiver defined in unit.
// Binds each method to its own receiver (based on receiver name).
func (m *Merger) bindMethods() error {
	if len(m.nodes.Meds) != len(m.unit.Meds) {
		panic(fmt.Sprintf("inconsistent number of method nodes (%d) and symbols (%d)",
			len(m.nodes.Meds), len(m.unit.Meds)))
	}

	for j := 0; j < len(m.unit.Meds); j += 1 {
		i := astIndexSymDef(j)
		med := m.nodes.Med(i)

		rname := med.Receiver.Name.Lit
		pos := med.Receiver.Name.Pos

		r := m.unit.Scope.sym(rname)
		if r == nil {
			return fmt.Errorf("%s: unresolved symbol \"%s\" in method receiver", pos.String(), rname)
		}
		if r.Kind != sym.Type {
			return fmt.Errorf("%s: only custom types can have methods, but \"%s\" is not a type",
				pos.String(), rname)
		}

		// method names are guaranteed to be unique for specific receiver due
		// to earlier scope symbol name check
		m.nodes.bindMethod(rname, i)
	}
	return nil
}

func (m *Merger) scanConstants() error {
	for _, s := range m.unit.Cons {
		c, err := m.scanCon(s)
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

func (m *Merger) scanCon(s *Symbol) (*ConstDef, error) {
	scope := m.unit.Scope
	node := m.nodes.Con(s.Def.(astIndexSymDef))

	t, err := scope.Types.Lookup(node.Type)
	if err != nil {
		return nil, err
	}

	expr := node.Expression
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
	for _, s := range m.unit.Funs {
		fn, err := m.scanFun(s)
		if err != nil {
			return err
		}
		s.Def = fn
	}
	return nil
}

func (m *Merger) scanFun(s *Symbol) (*FunDef, error) {
	scope := m.unit.Scope
	node := m.nodes.Fun(s.Def.(astIndexSymDef))

	var params []*Symbol
	for _, p := range node.Signature.Params {
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

	pos := node.Body.Pos
	result, err := scope.Types.Lookup(node.Signature.Result)
	if err != nil {
		return nil, err
	}
	fn := &FunDef{
		Params: params,
		Result: result,
		Never:  node.Signature.Never,
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

	err = m.scanFunBody(fn, node.Body.Statements)
	if err != nil {
		return nil, err
	}

	return fn, nil
}

func (m *Merger) scanFunBody(def *FunDef, statements []ast.Statement) error {
	if len(statements) == 0 {
		if def.Result != nil {
			return fmt.Errorf("%s: function with return type cannot have empty body", def.Body.Pos.String())
		}
		if def.Never {
			return fmt.Errorf("%s: function marked as never returning cannot have empty body", def.Body.Pos.String())
		}
		return nil
	}

	ctx := m.newFunCtx(def)
	err := def.Body.Fill(ctx, statements)
	if err != nil {
		return err
	}
	def.Refs = ctx.ref.Elems()
	return nil
}
