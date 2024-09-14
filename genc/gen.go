/*
package genc generates C code from a Ku program (in STG representation).
Almost all logic for generating output (C code) is bound to Builder.

Since C language does not have native namespacing mechanism for symbols
we need to implement mangling ourselves. To avoid name collisions with external
existing C code, all symbol names, generated by this package, have a prefix
"ku_" for global and unit level objects.

We also need to mangle symbol names from different units, and thus naming scheme
becomes:

	"ku_" - prefix for global objects and generated types (chunks for example)
	"ku_" + "<mangled_unit_name>" + "_" - prefix for unit level objects

The moniker string "<mangled_unit_name>" inside the latter must be unique for
each unit even if some unit names collide. To achieve this unit names are
mangled as follows. If two or more units have the same name, then we substitute
each "name" with:

	"name" + "<unit origin>" + "<unit graph index>"
*/
package genc

import (
	"io"
	"strconv"

	"github.com/mebyus/gizmo/stg"
	"github.com/mebyus/gizmo/uwalk"
)

const defaultPrefix = "ku_"

func GenUnit(w io.Writer, u *stg.Unit) error {
	var g Builder
	g.prefix = defaultPrefix
	g.specs = make(map[*stg.Type]string)

	g.prelude()
	g.builtinDerivativeTypes(u.Scope.Types)
	g.Gen(u)

	_, err := w.Write(g.Bytes())
	return err
}

func GenProgram(w io.Writer, p *uwalk.Program) error {
	var g Builder
	g.prefix = defaultPrefix
	g.specs = make(map[*stg.Type]string)

	g.MangleUnitNames(p.Units)

	g.prelude()
	g.builtinDerivativeTypes(p.Global.Types)
	for _, c := range p.Graph.Cohorts {
		for _, i := range c {
			u := p.Graph.Nodes[i].Unit
			g.Gen(u)
		}
	}

	if true { // TODO: make this as p.Main != nil
		g.entrypoint()
	}

	_, err := w.Write(g.Bytes())
	return err
}

func (g *Builder) getMangledUnitName(u *stg.Unit) string {
	return g.unames[u.Index]
}

func (g *Builder) MangleUnitNames(units []*stg.Unit) {
	g.unames = make([]string, len(units))

	m := make(map[ /* unit name */ string]*stg.Unit, len(units))
	var shouldMangle []uint32 // list of unit indices which should be mangled
	for i := range len(units) {
		u := units[i]
		old, ok := m[u.Name]
		if ok {
			m[u.Name] = nil
			if old != nil {
				shouldMangle = append(shouldMangle, old.Index)
			}
			shouldMangle = append(shouldMangle, u.Index)
		} else {
			m[u.Name] = u
			g.unames[i] = u.Name
		}
	}

	for _, i := range shouldMangle {
		u := units[i]
		g.unames[i] = u.Name + strconv.FormatUint(uint64(u.Path.Origin), 10) + strconv.FormatUint(uint64(i), 10)
	}
}
