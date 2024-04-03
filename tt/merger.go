package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/tt/sym"
)

// Merger is a high-level algorithm driver that gathers multiple ASTs of unit's atoms
// to produce that unit's type tree.
type Merger struct {
	ctx Context

	// Unit that is currently built by merger
	unit Unit

	// Processing of this nodes is deferred until phase 2, because they do not produce
	// a top-level symbols but need to be bound with other top-level symbols. May contain
	// one of the following
	//
	//	- method
	//	- prototype method blueprint
	noSymbolNodes []ast.TopLevel
}

func New(ctx Context) *Merger {
	return &Merger{
		ctx: ctx,
		unit: Unit{
			Scope: NewUnitScope(ctx.Global),
		},
	}
}

// Context is a reference data structure that contains type and symbol information
// about imported units. It also holds build conditions under which unit compilation
// is performed.
type Context struct {
	Global *Scope
}

func (m *Merger) Add(atom ast.UnitAtom) error {
	for _, block := range atom.Header.Imports.ImportBlocks {
		for _, spec := range block.Specs {
			err := m.addImport(ast.ImportBind{
				Name:   spec.Name,
				Public: block.Public,
				Path: origin.Path{
					Origin: block.Origin,
					ImpStr: spec.String.Lit,
				},
			})
			if err != nil {
				return err
			}
		}
	}

	for _, block := range atom.Blocks {
		// TODO: remove namespace blocks support
		if !block.Default {
			continue
		}

		for _, top := range block.Nodes {
			var err error

			switch top.Kind() {
			case toplvl.Fn:
				err = m.addFn(top.(ast.TopFunctionDefinition))
			case toplvl.Type:
				err = m.addType(top.(ast.TopType))
			case toplvl.Const:
				err = m.addConst(top.(ast.TopConst))
			case toplvl.Method, toplvl.Pmb:
				// defer processing until phase 2
				m.addPhaseTwoNode(top)
			default:
				panic(fmt.Sprintf("<top-level %s node not implemented>", top.Kind().String()))
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Merger) addPhaseTwoNode(top ast.TopLevel) {
	m.noSymbolNodes = append(m.noSymbolNodes, top)
}

// Merge is called after all unit atoms were added to merger to build type tree of the unit.
func (m *Merger) Merge() (*Unit, error) {
	return &m.unit, nil
}

// add top-level symbol to unit
func (m *Merger) add(s *Symbol) {
	m.unit.Scope.Bind(s)
}

// check if top-level symbol with a given name already exists in unit
func (m *Merger) has(name string) bool {
	s := m.unit.Scope.sym(name)
	return s != nil
}

func (m *Merger) errMultDef(name string, pos source.Pos) error {
	return fmt.Errorf("%s: multiple definitions of symbol \"%s\"", pos.String(), name)
}

func (m *Merger) addImport(bind ast.ImportBind) error {
	name := bind.Name.Lit
	pos := bind.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	// TODO: add context search for imported unit
	s := &Symbol{
		Kind:   sym.Import,
		Name:   name,
		Pos:    pos,
		Public: bind.Public,
	}
	m.add(s)
	return nil
}

func (m *Merger) addFn(top ast.TopFunctionDefinition) error {
	name := top.Definition.Head.Name.Lit
	pos := top.Pin()

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind:   sym.Fn,
		Name:   name,
		Pos:    pos,
		Public: top.Public,
	}
	m.add(s)
	return nil
}

func (m *Merger) addType(top ast.TopType) error {
	name := top.Name.Lit
	pos := top.Pin()

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind:   sym.Const,
		Name:   name,
		Pos:    pos,
		Public: top.Public,
	}
	m.add(s)
	return nil
}

func (m *Merger) addConst(top ast.TopConst) error {
	name := top.Name.Lit
	pos := top.Pin()

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind:   sym.Const,
		Name:   name,
		Pos:    pos,
		Public: top.Public,
	}
	m.add(s)
	return nil
}
