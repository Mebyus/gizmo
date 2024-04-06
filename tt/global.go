package tt

import (
	"github.com/mebyus/gizmo/tt/scp"
	"github.com/mebyus/gizmo/tt/sym"
)

func newNamedTypeSymbol(name string, typ *Type) *Symbol {
	return &Symbol{
		Name: name,
		Kind: sym.Type,
		Def:  typ,
		Type: &Type{}, // TODO: make type for other types
	}
}

func NewGlobalScope() *Scope {
	s := NewScope(scp.Global, nil, nil)

	s.Bind(newNamedTypeSymbol("u8", &Type{}))
	s.Bind(newNamedTypeSymbol("u16", &Type{}))
	s.Bind(newNamedTypeSymbol("u32", &Type{}))
	s.Bind(newNamedTypeSymbol("u64", &Type{}))

	s.Bind(newNamedTypeSymbol("i8", &Type{}))
	s.Bind(newNamedTypeSymbol("i16", &Type{}))
	s.Bind(newNamedTypeSymbol("i32", &Type{}))
	s.Bind(newNamedTypeSymbol("i64", &Type{}))

	s.Bind(newNamedTypeSymbol("uint", &Type{}))
	s.Bind(newNamedTypeSymbol("int", &Type{}))

	s.Bind(newNamedTypeSymbol("bool", &Type{}))
	s.Bind(newNamedTypeSymbol("rune", &Type{}))
	s.Bind(newNamedTypeSymbol("str", &Type{}))

	return s
}

type True struct {
}

type False struct {
}
