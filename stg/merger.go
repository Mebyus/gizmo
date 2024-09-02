package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/source/origin"
	"github.com/mebyus/gizmo/stg/scp"
)

// Merger is a high-level algorithm driver that gathers multiple ASTs of unit's atoms
// to produce that unit's type tree.
//
// Unit merging is done in several separate phases:
//
//   - 1 (gather) - gather all atoms of the unit
//   - 2 (index) - index unit level symbol
//   - 3 (method bind) - bind methods to corresponding receivers
//   - 4 (inspect) - determine dependency relations between unit level symbols
//   - 5 (graph) - construct, map and rank symbol dependency graph
//   - 6 (static eval) - eval and finalize all properties of unit level types and constants
//   - 7 (block scan) - recursively scan statements and expressions inside functions
type Merger struct {
	nodes NodesBox

	Warns []Warn

	resolver Resolver

	// Unit that is currently being built by merger.
	unit *Unit

	graph *Graph
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

	// Could be nil when constructing merger.
	Unit *Unit
}

func New(ctx UnitContext) *Merger {
	if ctx.Global == nil {
		panic("nil global scope")
	}
	u := ctx.Unit
	if u == nil {
		u = &Unit{Name: ctx.Name}
	}
	u.Scope = NewUnitScope(u, ctx.Global)
	u.TestsScope = NewScope(scp.UnitTests, u.Scope, nil)

	return &Merger{
		resolver: ctx.Resolver,
		unit:     u,
		nodes: NodesBox{
			MethodsByReceiver: make(map[string][]astIndexSymDef),
		},
	}
}

func Merge(ctx UnitContext, atoms []*ast.Atom) (*Unit, error) {
	if len(atoms) == 0 {
		panic("no atoms")
	}
	m := New(ctx)

	m.nodes.prealloc(atoms)
	for _, a := range atoms {
		err := m.Add(a)
		if err != nil {
			return nil, err
		}
	}

	return m.Merge()
}

type NodesBox struct {
	Funs      []ast.TopFun
	Vars      []ast.TopVar
	Tests     []ast.TopFun
	Types     []ast.TopType
	Methods   []ast.Method
	Constants []ast.TopLet

	// Maps custom type receiver name to a list of
	// its method indices inside Meds slice.
	MethodsByReceiver map[ /* receiver type name */ string][]astIndexSymDef
}

func (n *NodesBox) prealloc(atoms []*ast.Atom) {
	var (
		funs      int
		vars      int
		tests     int
		types     int
		methods   int
		constants int
	)
	for _, a := range atoms {
		funs += len(a.Funs)
		vars += len(a.Vars)
		tests += len(a.Tests)
		types += len(a.Types)
		methods += len(a.Methods)
		constants += len(a.Constants)
	}

	if funs != 0 {
		n.Funs = make([]ast.TopFun, 0, funs)
	}
	if vars != 0 {
		n.Vars = make([]ast.TopVar, 0, vars)
	}
	if tests != 0 {
		n.Tests = make([]ast.TopFun, 0, tests)
	}
	if types != 0 {
		n.Types = make([]ast.TopType, 0, types)
	}
	if methods != 0 {
		n.Methods = make([]ast.Method, 0, methods)
	}
	if constants != 0 {
		n.Constants = make([]ast.TopLet, 0, constants)
	}
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

func (n *NodesBox) addConstant(node ast.TopLet) astIndexSymDef {
	i := len(n.Constants)
	n.Constants = append(n.Constants, node)
	return astIndexSymDef(i)
}

func (n *NodesBox) addVar(node ast.TopVar) astIndexSymDef {
	i := len(n.Vars)
	n.Vars = append(n.Vars, node)
	return astIndexSymDef(i)
}

func (n *NodesBox) addMethod(node ast.Method) astIndexSymDef {
	i := len(n.Methods)
	n.Methods = append(n.Methods, node)
	return astIndexSymDef(i)
}

func (n *NodesBox) addTest(node ast.TopFun) astIndexSymDef {
	i := len(n.Tests)
	n.Tests = append(n.Tests, node)
	return astIndexSymDef(i)
}

func (n *NodesBox) bindMethod(rname string, i astIndexSymDef) {
	n.MethodsByReceiver[rname] = append(n.MethodsByReceiver[rname], i)
}

func (n *NodesBox) Type(i astIndexSymDef) ast.TopType {
	return n.Types[i]
}

func (n *NodesBox) Fun(i astIndexSymDef) ast.TopFun {
	return n.Funs[i]
}

func (n *NodesBox) Constant(i astIndexSymDef) ast.TopLet {
	return n.Constants[i]
}

func (n *NodesBox) Var(i astIndexSymDef) ast.TopVar {
	return n.Vars[i]
}

func (n *NodesBox) Method(i astIndexSymDef) ast.Method {
	return n.Methods[i]
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

type MapResolver map[origin.Path]*Unit

// Explicit interface implementation check.
var _ Resolver = MapResolver(nil)

func (r MapResolver) Resolve(p origin.Path) *Unit {
	return r[p]
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
	err = m.addConstants(atom.Constants)
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
	err = m.addMethods(atom.Methods)
	if err != nil {
		return err
	}
	err = m.addTests(atom.Tests)
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

func (m *Merger) addConstants(constants []ast.TopLet) error {
	for _, c := range constants {
		err := m.addConstant(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Merger) addFuns(funs []ast.TopFun) error {
	for _, f := range funs {
		err := m.addFun(f)
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

func (m *Merger) addMethods(methods []ast.Method) error {
	for _, md := range methods {
		err := m.addMethod(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Merger) addTests(tests []ast.TopFun) error {
	for _, t := range tests {
		err := m.addTest(t)
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
	return m.unit, nil
}

// add unit level symbol to scope
func (m *Merger) addSymbol(s *Symbol) {
	m.unit.Scope.Bind(s)
}

func (m *Merger) addTestSymbol(s *Symbol) {
	m.unit.TestsScope.Bind(s)
}

// check if top-level symbol with a given name already exists in unit
func (m *Merger) has(name string) bool {
	s := m.unit.Scope.sym(name)
	return s != nil
}

func (m *Merger) hasTest(name string) bool {
	s := m.unit.TestsScope.sym(name)
	return s != nil
}

func (m *Merger) errMultDef(name string, pos source.Pos) error {
	return fmt.Errorf("%s: multiple definitions of symbol \"%s\"", pos.String(), name)
}

func (m *Merger) addImport(bind ast.ImportBind) error {
	u := m.resolver.Resolve(bind.Path)
	if u == nil {
		panic("impossible due to previous checks")
	}

	name := bind.Name.Lit
	pos := bind.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind: smk.Import,
		Name: name,
		Pos:  pos,
		Pub:  bind.Pub,
		Def:  ImportSymDef{Unit: u},
	}
	m.addSymbol(s)
	return nil
}

func (m *Merger) addFun(top ast.TopFun) error {
	name := top.Name.Lit
	pos := top.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind: smk.Fun,
		Name: name,
		Pos:  pos,
		Pub:  top.Pub,
		Def:  m.nodes.addFun(top),
	}
	m.addSymbol(s)
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
		Kind: smk.Type,
		Name: name,
		Pos:  pos,
		Pub:  top.Pub,
		Def:  m.nodes.addType(top),
	}
	m.addSymbol(s)
	m.unit.addType(s)
	return nil
}

func (m *Merger) addConstant(top ast.TopLet) error {
	name := top.Name.Lit
	pos := top.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind: smk.Let,
		Name: name,
		Pos:  pos,
		Pub:  top.Pub,
		Def:  m.nodes.addConstant(top),
	}
	m.addSymbol(s)
	m.unit.addConstant(s)
	return nil
}

func (m *Merger) addVar(top ast.TopVar) error {
	name := top.Name.Lit
	pos := top.Name.Pos

	if m.has(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind: smk.Var,
		Name: name,
		Pos:  pos,
		Pub:  top.Pub,
		Def:  m.nodes.addVar(top),
	}
	m.addSymbol(s)
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
		Kind: smk.Method,
		Name: name,
		Pos:  pos,
		Pub:  top.Pub,
		Def:  m.nodes.addMethod(top),
	}
	m.addSymbol(s)
	m.unit.addMethod(s)
	return nil
}

func (m *Merger) addTest(top ast.TopFun) error {
	name := top.Name.Lit
	pos := top.Name.Pos

	if m.hasTest(name) {
		return m.errMultDef(name, pos)
	}

	s := &Symbol{
		Kind: smk.Test,
		Name: name,
		Pos:  pos,
		Def:  m.nodes.addTest(top),
	}
	m.addTestSymbol(s)
	m.unit.addTest(s)
	return nil
}
