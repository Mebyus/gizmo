package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/enums/tpk"
)

// perform second scan of named types present in unit,
// this operation uses type graph constructed during shallow scan.
//
// after this phase each symbol holding a named type has its Def field
// set to *Type
func (m *Merger) eval() error {
	// TODO: remove debug print
	fmt.Println("list of isolated symbols")
	for _, i := range m.graph.Isolated {
		n := &m.graph.Nodes[i]
		err := m.evalNode(n)
		if err != nil {
			return err
		}
	}

	fmt.Println()
	fmt.Println("list of components")
	for k := 0; k < len(m.graph.Comps); k += 1 {
		c := &m.graph.Comps[k]
		fmt.Printf("component %d\n", k)
		for rank, cohort := range c.Cohorts {
			fmt.Printf("cohort %d\n", rank)
			for _, i := range cohort {
				n := &m.graph.Nodes[c.V[i].Index]
				err := m.evalNode(n)
				if err != nil {
					return err
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}

	return nil
}

func (m *Merger) evalNode(n *GraphNode) error {
	s := n.Symbol
	fmt.Printf("%s %s\n", s.Kind, s.Name)

	var err error
	switch s.Kind {
	case smk.Let:
		if n.SelfLoop {
			panic("self loop is impossible for const symbol")
		}
		err = m.evalConstant(s)
	case smk.Type:
		err = m.evalType(s, n.SelfLoop)
	default:
		panic(fmt.Sprintf("unexpected %s symbol", s.Kind))
	}
	return err
}

func (m *Merger) evalConstant(s *Symbol) error {
	scope := m.unit.Scope
	node := m.nodes.Con(s.Index())

	t, err := scope.Types.Lookup(node.Type)
	if err != nil {
		return err
	}

	exp := node.Expr
	if exp == nil {
		panic("nil expression in constant definition")
	}
	ctx := m.newConstCtx()
	e, err := scope.Scan(ctx, exp)
	if err != nil {
		return err
	}
	e2, err := scope.evalStaticExp(e)
	if err != nil {
		return err
	}

	// TODO: remove debug print
	if e2.Kind() == exn.Integer {
		i := e2.(Integer)
		sign := ""
		if i.Neg {
			sign = "-"
		}
		fmt.Printf("integer const %s = %s%d\n", s.Name, sign, i.Val)
	}

	if t == nil {
		t = e2.Type()
	} else {
		panic("type check not implemented")
	}

	// TODO: eval constant value (reduce expression)

	def := &ConstDef{
		Exp:  e2,
		Type: t,
	}
	s.Def = def
	s.Type = t
	return nil
}

func (m *Merger) evalType(s *Symbol, selfLoop bool) error {
	if selfLoop {
		m.evalRecursiveType(s)
	} else {
		m.evalSimpleType(s)
	}
	return nil
}

func (m *Merger) evalRecursiveType(s *Symbol) {
	node := m.nodes.Type(s.Index())
	t := &Type{
		Flags: TypeFlagRecursive,
		Kind:  tpk.Custom,
	}
	def := CustomTypeDef{
		Base:   t,
		Symbol: s,
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

func (m *Merger) getReceiverMethods(name string) []*Symbol {
	ii := m.nodes.MedsByReceiver[name]
	if len(ii) == 0 {
		return nil
	}
	symbols := make([]*Symbol, 0, len(ii))
	for _, i := range ii {
		symbols = append(symbols, m.unit.Meds[i])
	}
	return symbols
}

func (m *Merger) evalSimpleType(s *Symbol) {
	node := m.nodes.Type(s.Index())
	base, err := m.unit.Scope.Types.lookup(node.Spec)
	if err != nil {
		// type graph structure must guarantee successful lookup
		panic(err)
	}
	fmt.Printf("%s: %T\n", s.Name, base.Def)
	t := &Type{
		Def: CustomTypeDef{
			Symbol:  s,
			Base:    base,
			Methods: m.getReceiverMethods(s.Name),
		},
		Kind: tpk.Custom,
	}
	s.Def = t
}
