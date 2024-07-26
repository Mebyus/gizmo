package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/source/origin"
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
	// Unit that is currently being built by merger.
	unit Unit

	ctx UnitContext

	nodes NodesBox

	Warns []Warn

	graph *TypeGraph
}

type NodesBox struct {
	Funs  []ast.TopFun
	Cons  []ast.TopLet
	Vars  []ast.TopVar
	Meds  []ast.Method
	Types []ast.TopType

	// Maps custom type receiver name to a list of
	// its method indices inside Meds slice.
	MedsByReceiver map[string][]astIndexSymDef
}

func (n *NodesBox) addType(node ast.TopType) astIndexSymDef {
	i := len(n.Types)
	n.Types = append(n.Types, node)
	return astIndexSymDef(i)
}

func (n *NodesBox) addFun(node ast.TopFun) astIndexSymDef {
	i := len(n.Funs)
	n.Funs = append(n.Funs, node)
	return astIndexSymDef(i)
}

func (n *NodesBox) addCon(node ast.TopLet) astIndexSymDef {
	i := len(n.Cons)
	n.Cons = append(n.Cons, node)
	return astIndexSymDef(i)
}

func (n *NodesBox) addVar(node ast.TopVar) astIndexSymDef {
	i := len(n.Vars)
	n.Vars = append(n.Vars, node)
	return astIndexSymDef(i)
}

func (n *NodesBox) addMed(node ast.Method) astIndexSymDef {
	i := len(n.Meds)
	n.Meds = append(n.Meds, node)
	return astIndexSymDef(i)
}

func (n *NodesBox) bindMethod(rname string, i astIndexSymDef) {
	n.MedsByReceiver[rname] = append(n.MedsByReceiver[rname], i)
}

func (n *NodesBox) Type(i astIndexSymDef) ast.TopType {
	return n.Types[i]
}

func (n *NodesBox) Fun(i astIndexSymDef) ast.TopFun {
	return n.Funs[i]
}

func (n *NodesBox) Con(i astIndexSymDef) ast.TopLet {
	return n.Cons[i]
}

func (n *NodesBox) Var(i astIndexSymDef) ast.TopVar {
	return n.Vars[i]
}

func (n *NodesBox) Med(i astIndexSymDef) ast.Method {
	return n.Meds[i]
}

func New(ctx UnitContext) *Merger {
	return &Merger{
		ctx: ctx,
		unit: Unit{
			Name:  ctx.Name,
			Scope: NewUnitScope(ctx.Global),
		},
		nodes: NodesBox{
			MedsByReceiver: make(map[string][]astIndexSymDef),
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

	Resolver Resolver

	Global *Scope
}

// Resolver gives access to other units by their origin path.
type Resolver interface {
	// Returns nil if there is no unit with such path.
	Resolve(origin.Path) *Unit
}

// EmptyResolver implements Resolver interface by returning
// nil for all requested paths.
type EmptyResolver struct{}

func NewEmptyResolver() EmptyResolver {
	return EmptyResolver{}
}

// Explicit interface implementation check.
var _ Resolver = EmptyResolver{}

func (EmptyResolver) Resolve(origin.Path) *Unit {
	return nil
}

func (m *Merger) Add(atom *ast.Atom) error {
	err := m.addImports(atom.Header.Imports.Blocks)
	if err != nil {
		return err
	}
	err = m.addTypes(atom.Types)
	if err != nil {
		return err
	}
	err = m.addLets(atom.Lets)
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
	err = m.addMethods(atom.Meds)
	if err != nil {
		return err
	}

	return nil
}

func (m *Merger) addImports(blocks []ast.ImportBlock) error {
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

func (m *Merger) addLets(cons []ast.TopLet) error {
	for _, con := range cons {
		err := m.addLet(con)
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

func (m *Merger) addMethods(meds []ast.Method) error {
	for _, med := range meds {
		err := m.addMethod(med)
		if err != nil {
			return err
		}
	}
	return nil
}

// Merge is called after all unit atoms were added to merger to build type tree of the unit.
func (m *Merger) Merge() (*Unit, error) {
	err := m.merge()
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
	u := m.ctx.Resolver.Resolve(bind.Path)
	if u == nil {
		panic("impossible due to previous checks")
	}

	name := bind.Name.Lit
	pos := bind.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	// TODO: add context search for imported unit
	s := &Symbol{
		Kind:   smk.Import,
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
		Kind:   smk.Fun,
		Name:   name,
		Pos:    pos,
		Public: top.Pub,
		Def:    m.nodes.addFun(top),
	}
	m.add(s)
	m.unit.addFun(s)
	return nil
}

func (m *Merger) addType(top ast.TopType) error {
	name := top.Name.Lit
	pos := top.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind:   smk.Type,
		Name:   name,
		Pos:    pos,
		Public: top.Pub,
		Def:    m.nodes.addType(top),
	}
	m.add(s)
	m.unit.addType(s)
	return nil
}

func (m *Merger) addLet(top ast.TopLet) error {
	name := top.Name.Lit
	pos := top.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind:   smk.Let,
		Name:   name,
		Pos:    pos,
		Public: top.Pub,
		Def:    m.nodes.addCon(top),
	}
	m.add(s)
	m.unit.addLet(s)
	return nil
}

func (m *Merger) addVar(top ast.TopVar) error {
	name := top.Name.Lit
	pos := top.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind:   smk.Var,
		Name:   name,
		Pos:    pos,
		Public: top.Pub,
		Def:    m.nodes.addVar(top),
	}
	m.add(s)
	m.unit.addVar(s)
	return nil
}

func (m *Merger) addMethod(top ast.Method) error {
	rname := top.Receiver.Name.Lit
	mname := top.Name.Lit
	pos := top.Name.Pos
	name := rname + "." + mname

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind:   smk.Method,
		Name:   name,
		Pos:    pos,
		Public: top.Pub,
		Def:    m.nodes.addMed(top),
	}
	m.add(s)
	m.unit.addMed(s)
	return nil
}
