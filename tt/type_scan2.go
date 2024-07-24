package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/tt/typ"
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
	name := node.Name.Lit
	t := &Type{
		Recursive: true,

		Name: name,
		Kind: typ.Named,
		Def:  NamedTypeDef{Symbol: s},
	}
	s.Def = t

	base, err := m.unit.Scope.Types.lookup(node.Spec)
	if err != nil {
		// type graph structure must guarantee successful lookup
		panic(err)
	}
	t.Base = base.Base
}

func (m *Merger) bindType(s *Symbol) *Type {
	node := m.nodes.Type(s.Def.(astIndexSymDef))
	name := node.Name.Lit
	base, err := m.unit.Scope.Types.lookup(node.Spec)
	if err != nil {
		// type graph structure must guarantee successful lookup
		panic(err)
	}

	return &Type{
		Name: name,
		Base: base.Base,
		Kind: typ.Named,
		Def:  NamedTypeDef{Symbol: s},
	}
}
