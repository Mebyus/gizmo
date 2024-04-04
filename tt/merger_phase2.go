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
	err = m.scanFns()
	if err != nil {
		return err
	}
	return nil
}

func (m *Merger) bindNodes() error {
	for _, top := range m.symBindNodes {
		var err error

		switch top.Kind() {
		case toplvl.Method:
			err = m.bindMethod(top.(ast.Method))
		case toplvl.Pmb:
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

	err := rvSym.Def.(*TempTypeDef).addMethod(method)
	return err
}

func (m *Merger) scanTypes() error {
	// TODO: graph based scanning
	for _, s := range m.unit.Scope.Symbols {
		if s.Kind == sym.Type {
			def := s.Def.(*TempTypeDef)
			t, err := m.scanType(def)
			if err != nil {
				return err
			}
			s.Def = t
		}
	}
	return nil
}

func (m *Merger) scanType(def *TempTypeDef) (*Type, error) {
	return &Type{}, nil
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
	var params []*Symbol
	for _, p := range def.top.Definition.Head.Signature.Params {
		params = append(params, &Symbol{
			Pos:  p.Name.Pos,
			Name: p.Name.Lit,
			Type: m.lookupType(p.Type),
			Kind: sym.Param,
		})
	}
	return &FnDef{
		Params: params,
		Result: m.lookupType(def.top.Definition.Head.Signature.Result),
		Never:  def.top.Definition.Head.Signature.Never,
	}, nil
}
