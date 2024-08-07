package stg

import (
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/enums/tpk"
	"github.com/mebyus/gizmo/stg/scp"
)

func newTypeSymbol(name string, t *Type) *Symbol {
	s := &Symbol{
		Name: name,
		Kind: smk.Type,
		Def:  t,
		Type: &Type{}, // TODO: make type for other types
	}
	return s
}

func (s *Scope) addTypeSymbol(name string, t *Type) {
	sym := newTypeSymbol(name, t)
	t.Def = BuiltinTypeDef{Symbol: sym}
	s.BindTypeSymbol(sym)
}

func (s *Scope) addBooleanType() {
	s.addTypeSymbol("bool", BoolType)
}

func (s *Scope) addStringType() {
	s.addTypeSymbol("str", StrType)
}

func newStaticType(kind tpk.Kind) *Type {
	t := &Type{Kind: kind}
	return t
}

func newBooleanType() *Type {
	t := &Type{
		Kind:    tpk.Boolean,
		Size:    1,
		Builtin: true,
	}
	return t
}

func newStringType() *Type {
	t := &Type{
		Kind:    tpk.String,
		Size:    16,
		Builtin: true,
	}
	return t
}

func newUnsignedType(size uint32) *Type {
	t := &Type{
		Kind:    tpk.Unsigned,
		Size:    size,
		Builtin: true,
	}
	return t
}

func newSignedType(size uint32) *Type {
	t := &Type{
		Kind:    tpk.Signed,
		Size:    size,
		Builtin: true,
	}
	return t
}

var (
	Trivial = newStaticType(tpk.Trivial)

	StaticInteger = newStaticType(tpk.StaticInteger)
	StaticFloat   = newStaticType(tpk.StaticFloat)
	StaticString  = newStaticType(tpk.StaticString)
	StaticBoolean = newStaticType(tpk.StaticBoolean)
	StaticNil     = newStaticType(tpk.StaticNil)

	Uint8Type  = newUnsignedType(1)
	Uint16Type = newUnsignedType(2)
	Uint32Type = newUnsignedType(4)
	Uint64Type = newUnsignedType(8)
	UintType   = newUnsignedType(8) // TODO: adjust uint size based on target arch

	Sint8Type  = newSignedType(1)
	Sint16Type = newSignedType(2)
	Sint32Type = newSignedType(4)
	Sint64Type = newSignedType(8)
	SintType   = newSignedType(8) // TODO: adjust int size based on target arch

	BoolType = newBooleanType()
	StrType  = newStringType()
	RuneType = newUnsignedType(4)
)

func (s *Scope) addStaticTypes() {
	// TODO: this is probably redundant
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

	s.addTypeSymbol("s8", Sint8Type)
	s.addTypeSymbol("s16", Sint16Type)
	s.addTypeSymbol("s32", Sint32Type)
	s.addTypeSymbol("s64", Sint64Type)
	s.addTypeSymbol("sint", SintType)

	s.addBooleanType()
	s.addStringType()
	s.addTypeSymbol("rune", RuneType)

	return s
}
