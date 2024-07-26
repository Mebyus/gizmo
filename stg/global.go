package stg

import (
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/stg/scp"
	"github.com/mebyus/gizmo/stg/typ"
)

func newTypeSymbol(name string, t *Type) *Symbol {
	s := &Symbol{
		Name: name,
		Kind: smk.Type,
		Def:  t,
		Type: &Type{}, // TODO: make type for other types
	}
	t.Symbol = s
	return s
}

func (s *Scope) addTypeSymbol(name string, t *Type) {
	s.BindTypeSymbol(newTypeSymbol(name, t))
}

func (s *Scope) addBuiltinBoolType() {
	t := &Type{
		Kind:    typ.Boolean,
		Size:    1,
		Builtin: true,
	}
	t.Base = t
	s.BindTypeSymbol(newTypeSymbol("bool", t))
}

func (s *Scope) addBuiltinStringType() {
	t := &Type{
		Kind:    typ.String,
		Size:    16,
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

func newUnsignedType(size uint32) *Type {
	t := &Type{
		Kind:    typ.Unsigned,
		Size:    size,
		Builtin: true,
	}
	t.Base = t
	return t
}

func newSignedType(size uint32) *Type {
	t := &Type{
		Kind:    typ.Signed,
		Size:    size,
		Builtin: true,
	}
	t.Base = t
	return t
}

var (
	Trivial = newStaticType(typ.Trivial)

	StaticInteger = newStaticType(typ.StaticInteger)
	StaticFloat   = newStaticType(typ.StaticFloat)
	StaticString  = newStaticType(typ.StaticString)
	StaticBoolean = newStaticType(typ.StaticBoolean)
	StaticNil     = newStaticType(typ.StaticNil)

	Uint8Type  = newUnsignedType(1)
	Uint16Type = newUnsignedType(2)
	Uint32Type = newUnsignedType(4)
	Uint64Type = newUnsignedType(8)
	UintType   = newUnsignedType(8) // TODO: adjust uint size based on target arch

	Int8Type  = newSignedType(1)
	Int16Type = newSignedType(2)
	Int32Type = newSignedType(4)
	Int64Type = newSignedType(8)
	IntType   = newSignedType(8) // TODO: adjust int size based on target arch

	RuneType = newUnsignedType(4)
)

func (s *Scope) addStaticTypes() {
	s.Types.tm[StaticInteger.Stable()] = StaticInteger
	s.Types.tm[StaticFloat.Stable()] = StaticFloat
	s.Types.tm[StaticString.Stable()] = StaticString
	s.Types.tm[StaticBoolean.Stable()] = StaticBoolean
	s.Types.tm[StaticNil.Stable()] = StaticNil
}

func (s *Scope) addTrivialType() {
	// s.Types.tm[]
}

func NewGlobalScope() *Scope {
	s := NewScope(scp.Global, nil, nil)

	s.addTrivialType()

	s.addStaticTypes()

	s.addTypeSymbol("u8", Uint8Type)
	s.addTypeSymbol("u16", Uint16Type)
	s.addTypeSymbol("u32", Uint32Type)
	s.addTypeSymbol("u64", Uint64Type)
	s.addTypeSymbol("uint", UintType)

	s.addTypeSymbol("i8", Int8Type)
	s.addTypeSymbol("i16", Int16Type)
	s.addTypeSymbol("i32", Int32Type)
	s.addTypeSymbol("i64", Int64Type)
	s.addTypeSymbol("int", IntType)

	s.addBuiltinBoolType()
	s.addBuiltinStringType()
	s.addTypeSymbol("rune", RuneType)

	return s
}
