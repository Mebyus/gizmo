package tt

import (
	"github.com/mebyus/gizmo/tt/scp"
	"github.com/mebyus/gizmo/tt/sym"
	"github.com/mebyus/gizmo/tt/typ"
)

func newTypeSymbol(name string, typ *Type) *Symbol {
	return &Symbol{
		Name: name,
		Kind: sym.Type,
		Def:  typ,
		Type: &Type{}, // TODO: make type for other types
	}
}

func (s *Scope) addBuiltinUnsignedType(name string, size uint32) {
	t := &Type{
		Name:    name,
		Kind:    typ.Unsigned,
		Def:     IntTypeDef{Size: size},
		Builtin: true,
	}
	t.Base = t
	s.BindTypeSymbol(newTypeSymbol(name, t))
}

func (s *Scope) addBuiltinSignedType(name string, size uint32) {
	t := &Type{
		Name:    name,
		Kind:    typ.Signed,
		Def:     IntTypeDef{Size: size},
		Builtin: true,
	}
	t.Base = t
	s.BindTypeSymbol(newTypeSymbol(name, t))
}

func (s *Scope) addBuiltinBoolType() {
	t := &Type{
		Name:    "bool",
		Kind:    typ.Boolean,
		Builtin: true,
	}
	t.Base = t
	s.BindTypeSymbol(newTypeSymbol("bool", t))
}

func (s *Scope) addBuiltinStringType() {
	t := &Type{
		Name:    "str",
		Kind:    typ.String,
		Builtin: true,
	}
	t.Base = t
	s.BindTypeSymbol(newTypeSymbol("str", t))
}

func newStaticType(kind typ.Kind) *Type {
	t := &Type{Kind: kind}
	t.Base = t
	return t
}

var (
	StaticInteger = newStaticType(typ.StaticInteger)
	StaticFloat   = newStaticType(typ.StaticFloat)
	StaticString  = newStaticType(typ.StaticString)
	StaticBoolean = newStaticType(typ.StaticBoolean)
	StaticNil     = newStaticType(typ.StaticNil)
)

func (s *Scope) addStaticTypes() {
	s.Types.tm[StaticInteger.Stable()] = StaticInteger
	s.Types.tm[StaticFloat.Stable()] = StaticFloat
	s.Types.tm[StaticString.Stable()] = StaticString
	s.Types.tm[StaticBoolean.Stable()] = StaticBoolean
	s.Types.tm[StaticNil.Stable()] = StaticNil
}

func NewGlobalScope() *Scope {
	s := NewScope(scp.Global, nil, nil)

	s.addStaticTypes()

	s.addBuiltinUnsignedType("u8", 1)
	s.addBuiltinUnsignedType("u16", 2)
	s.addBuiltinUnsignedType("u32", 4)
	s.addBuiltinUnsignedType("u64", 8)
	s.addBuiltinUnsignedType("uint", 8) // TODO: adjust uint size based on target arch

	s.addBuiltinSignedType("i8", 1)
	s.addBuiltinSignedType("i16", 2)
	s.addBuiltinSignedType("i32", 4)
	s.addBuiltinSignedType("i64", 8)
	s.addBuiltinSignedType("int", 8) // TODO: adjust int size based on target arch

	s.addBuiltinBoolType()
	s.addBuiltinStringType()
	s.addBuiltinUnsignedType("rune", 4)

	return s
}
