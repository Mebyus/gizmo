package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/enums/tpk"
)

// perform second scan of named types present in unit,
// this operation uses type graph constructed during shallow scan.
//
// after this phase each symbol holding a named type has its Def field
// set to *Type
func (m *Merger) bindTypes() error {
	// TODO: remove debug print
	fmt.Println("list of isolated types")
	for _, i := range m.graph.Isolated {
		n := m.graph.Nodes[i]
		s := n.Sym
		fmt.Printf("%s\n", s.Name)

		if n.SelfLoop {
			m.bindRecursiveType(s)
		} else {
			s.Def = m.bindType(s)
		}
	}

	fmt.Println()
	fmt.Println("list of component types")
	for k := 0; k < len(m.graph.Comps); k += 1 {
		c := &m.graph.Comps[k]
		fmt.Printf("component %d\n", k)
		for rank, cohort := range c.Cohorts {
			fmt.Printf("cohort %d\n", rank)
			for _, i := range cohort {
				n := m.graph.Nodes[c.V[i].Index]
				fmt.Printf("%s\n", n.Sym.Name)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	return nil
}

func (m *Merger) bindRecursiveType(s *Symbol) {
	node := m.nodes.Type(s.Def.(astIndexSymDef))
	t := &Type{
		Recursive: true,

		Kind: tpk.Custom,
	}
	def := CustomTypeDef{
		Base: t,
		Sym:  s,
	}
	t.Def = def
	s.Def = t

	base, err := m.unit.Scope.Types.lookup(node.Spec)
	if err != nil {
		// type graph structure must guarantee successful lookup
		panic(err)
	}
	if base != t {
		panic("different recursive type after lookup")
	}
}

func (m *Merger) bindType(s *Symbol) *Type {
	node := m.nodes.Type(s.Def.(astIndexSymDef))
	base, err := m.unit.Scope.Types.lookup(node.Spec)
	if err != nil {
		// type graph structure must guarantee successful lookup
		panic(err)
	}
	fmt.Printf("%s: %T\n", s.Name, base.Def)
	return &Type{
		Def: CustomTypeDef{
			Sym:  s,
			Base: base,
		},
		Kind: tpk.Custom,
	}
}
