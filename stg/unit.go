package stg

import (
	"fmt"
	"os"
	"sort"

	"github.com/mebyus/gizmo/parser"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/source/origin"
)

type UnitImports struct {
	// Sorted.
	Paths []origin.Path

	// Has order corresponding to Paths field.
	// That means expression Units[i].Path == Paths[i] is true.
	Units []*Unit
}

type Unit struct {
	Imports UnitImports

	// List of all unit level function symbols defined in unit.
	Funs []*Symbol

	// List of all unit level constant symbols defined in unit.
	Lets []*Symbol

	// List of all unit level custom type symbols defined in unit.
	Types []*Symbol

	// List of all unit level variable symbols defined in unit.
	Vars []*Symbol

	// List of all method symbols defined in unit.
	// Methods can only be defined at unit level.
	Meds []*Symbol

	// Import path of this unit.
	Path origin.Path

	Name string

	// Scope that holds all unit level symbols from all unit atoms.
	//
	// This field is always not nil and Scope.Kind is always equal to scp.Unit.
	Scope *Scope

	DiscoveryIndex uint32
}

func SortUnits(units []*Unit) {
	u := units
	sort.Slice(u, func(i, j int) bool {
		a := u[i]
		b := u[j]
		return origin.Less(a.Path, b.Path)
	})
}

func (u *Unit) addFun(s *Symbol) {
	u.Funs = append(u.Funs, s)
}

func (u *Unit) addLet(s *Symbol) {
	u.Lets = append(u.Lets, s)
}

func (u *Unit) addType(s *Symbol) {
	u.Types = append(u.Types, s)
}

func (u *Unit) addVar(s *Symbol) {
	u.Vars = append(u.Vars, s)
}

func (u *Unit) addMed(s *Symbol) {
	u.Meds = append(u.Meds, s)
}

// UnitFromDir scans given directory for source files, processes them as
// single unit and constructs graph of unit's symbols.
//
// Given directory path should be cleaned by the client.
func UnitFromDir(resolver Resolver, dir string) (*Unit, error) {
	files, err := source.LoadUnitFiles(dir)
	if err != nil {
		return nil, err
	}

	var name string
	parsers := make([]*parser.Parser, 0, len(files))
	iset := origin.NewSet()
	for _, file := range files {
		p := parser.FromSource(file)
		h, err := p.Header()
		if err != nil {
			return nil, err
		}
		if h.Unit != nil {
			// if file has unit clause, use it to
			// determine unit name

			n := h.Unit.Name.Lit
			if name == "" {
				name = n
			} else if n != name {
				return nil, fmt.Errorf("%s: inconsistent unit name \"%s\" (previous was \"%s\")",
					h.Unit.Name.Pos.String(), n, name)
			}
		}

		for _, p := range h.Imports.Paths {
			u := resolver.Resolve(p)
			if u == nil {
				return nil, fmt.Errorf("unable to resolve imported unit \"%s\"", p.String())
			}
			if iset.Has(p) {
				return nil, fmt.Errorf("multiple imports of the same unit \"%s\"", p.String())
			}
			iset.Add(p)
		}

		parsers = append(parsers, p)
	}
	if name == "" {
		// unit does not have files with unit clause
		// determine unit name by its directory

		stat, err := os.Lstat(dir)
		if err != nil {
			return nil, err
		}
		if !stat.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", dir)
		}
		name = stat.Name()
		if name == "" {
			return nil, fmt.Errorf("os: no directory name")
		}
	}

	m := New(UnitContext{
		Name:     name,
		Global:   NewGlobalScope(),
		Resolver: resolver,
	})
	for _, p := range parsers {
		atom, err := p.Parse()
		if err != nil {
			return nil, err
		}
		err = m.Add(atom)
		if err != nil {
			return nil, err
		}
	}
	return m.Merge()
}
