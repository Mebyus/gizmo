package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/tt/sym"
)

// Merger is a high-level algorithm driver that gathers multiple ASTs of unit's atoms
// to produce that unit's type tree.
//
// Unit merging is done in two separate phases:
//
//   - phase 1 - atoms gathering (via Add method)
//   - phase 2 - symbol indexing, type checking, etc. for the whole unit (via Merge method)
//
// Semantic checks and tree construction is split between this two phases.
type Merger struct {
	nodes   NodesBox
	symbols SymbolsBox

	ctx UnitContext

	// Unit that is currently being built by merger.
	unit Unit

	// list of all top-level function symbols defined in unit
	funs []*Symbol

	// list of all top-level constant symbols defined in unit
	cons []*Symbol

	// list of all top-level type symbols defined in unit
	types []*Symbol

	Warns []Warn

	graph *TypeGraph
}

type NodesBox struct {
}

type SymbolsBox struct {
	// list of all top-level function symbols defined in unit
	funs []*Symbol

	// list of all top-level constant symbols defined in unit
	cons []*Symbol

	// list of all top-level type symbols defined in unit
	types []*Symbol
}

func New(ctx UnitContext) *Merger {
	return &Merger{
		ctx: ctx,
		unit: Unit{
			Name:  ctx.Name,
			Scope: NewUnitScope(ctx.Global),
		},
	}
}

// UnitContext is a reference data structure that contains type and symbol information
// about imported units. It also holds build conditions under which unit compilation
// is performed.
type UnitContext struct {
	// Possible name for a unit that is under construction.
	//
	// Unit name could be determined by outside means if
	// all unit atoms do not have unit clause with specified name.
	Name string

	Global *Scope
}

func (m *Merger) Add(atom *ast.Atom) error {
	err := m.addImportBlocks(atom.Header.Imports.ImportBlocks)
	if err != nil {
		return err
	}
	err = m.addTypes(atom.Types)
	if err != nil {
		return err
	}
	err = m.addCons(atom.Cons)
	if err != nil {
		return err
	}
	err = m.addFuns(atom.Funs)
	if err != nil {
		return err
	}
	err = m.addVars(atom.Vars)
	if err != nil {
		return err
	}
	err = m.addMeds(atom.Meds)
	if err != nil {
		return err
	}

	return nil
}

func (m *Merger) addImportBlocks(blocks []ast.ImportBlock) error {
	for _, block := range blocks {
		for _, spec := range block.Specs {
			err := m.addImport(ast.ImportBind{
				Name: spec.Name,
				Pub:  block.Pub,
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
	return nil
}

func (m *Merger) addTypes(types []ast.TopType) error {
	for _, t := range types {
		err := m.addType(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Merger) addCons(cons []ast.TopCon) error {
	for _, con := range cons {
		err := m.addCon(con)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Merger) addFuns(funs []ast.TopFun) error {
	for _, fun := range funs {
		err := m.addFun(fun)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Merger) addVars(vars []ast.TopVar) error {
	for _, v := range vars {
		err := m.addVar(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Merger) addMeds(meds []ast.Method) error {
	for _, med := range meds {
		err := m.addMed(med)
		if err != nil {
			return err
		}
	}
	return nil
}

// Merge is called after all unit atoms were added to merger to build type tree of the unit.
func (m *Merger) Merge() (*Unit, error) {
	err := m.runPhaseTwo()
	if err != nil {
		return nil, err
	}
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
		Public: bind.Pub,
	}
	m.add(s)
	return nil
}

func (m *Merger) addFun(top ast.TopFun) error {
	name := top.Name.Lit
	pos := top.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind:   sym.Fn,
		Name:   name,
		Pos:    pos,
		Public: top.Pub,
		Def:    NewTempFnDef(top),
	}
	m.add(s)
	m.funs = append(m.funs, s)
	return nil
}

func (m *Merger) addType(top ast.TopType) error {
	name := top.Name.Lit
	pos := top.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind:   sym.Type,
		Name:   name,
		Pos:    pos,
		Public: top.Pub,
		Def:    NewTempTypeDef(top),
	}
	m.add(s)
	m.types = append(m.types, s)
	return nil
}

func (m *Merger) addCon(top ast.TopCon) error {
	name := top.Name.Lit
	pos := top.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind:   sym.Const,
		Name:   name,
		Pos:    pos,
		Public: top.Pub,
		Def:    NewTempConstDef(top),
	}
	m.add(s)
	m.cons = append(m.cons, s)
	return nil
}

func (m *Merger) addVar(v ast.TopVar) error {
	return nil
}

func (m *Merger) addMed(med ast.Method) error {
	return nil
}
