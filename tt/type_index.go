package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/tps"
	"github.com/mebyus/gizmo/tt/sym"
	"github.com/mebyus/gizmo/tt/typ"
)

type TypeIndex struct {
	// Scope which is used for resolving symbols when performing
	// type specifier lookups.
	scope *Scope

	// Type map. Maps Type.Stable() to saved type.
	tm map[Stable]*Type
}

func (x *TypeIndex) Lookup(spec ast.TypeSpecifier) (*Type, error) {
	if spec == nil {
		return nil, nil
	}

	return x.lookup(spec)
}

func (x *TypeIndex) lookup(spec ast.TypeSpecifier) (*Type, error) {
	switch spec.Kind() {
	case tps.Name:
		return x.lookupNamed(spec.(ast.TypeName).Name.Name)
	case tps.Pointer:
		return x.lookupPointer(spec.(ast.PointerType).RefType)
	default:
		panic(fmt.Sprintf("not implemented for %s", spec.Kind().String()))
	}
}

func (x *TypeIndex) storePointer(ref *Type) *Type {
	return x.store(newPointerType(ref))
}

func (x *TypeIndex) store(a *Type) *Type {
	stable := a.Stable()
	t := x.tm[stable]
	if t == nil {
		x.tm[stable] = a
		t = a
	}
	return t
}

func (x *TypeIndex) lookupPointer(spec ast.TypeSpecifier) (*Type, error) {
	ref, err := x.lookup(spec)
	if err != nil {
		return nil, err
	}
	return x.storePointer(ref), nil 
}

func (x *TypeIndex) lookupNamed(idn ast.Identifier) (*Type, error) {
	name := idn.Lit
	pos := idn.Pos
	s := x.scope.Lookup(name, pos.Num)
	if s == nil {
		return nil, fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}
	if s.Kind != sym.Type {
		panic(fmt.Sprintf("unexpected symbol kind: %s", s.Kind.String()))
	}
	t := s.Def.(*Type)
	if t.Builtin || t.Kind == typ.Named {
		return t, nil
	}
	panic(fmt.Sprintf("unexpected type kind: %s", t.Kind.String()))
}
