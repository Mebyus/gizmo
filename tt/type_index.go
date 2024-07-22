package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/tps"
	"github.com/mebyus/gizmo/tt/sym"
	"github.com/mebyus/gizmo/tt/typ"
)

// TypeIndex is a helper managment object which is used to lookup and create
// a type by its type specifier.
//
// If type with given type specifier already exists, then existed type is returned,
// without creating a new one. If type index fails to find an existing type, then
// a new one will be created, stored for later lookups and returned.
type TypeIndex struct {
	// Scope which is used for resolving symbols when performing
	// type specifier lookups.
	scope *Scope

	// Type map. Maps Type.Stable() to saved type.
	tm map[Stable]*Type
}

func (x *TypeIndex) Lookup(spec ast.TypeSpecifier) (*Type, error) {
	if spec == nil {
		// handle never and void return "types"
		return nil, nil
	}

	return x.lookup(spec)
}

func (x *TypeIndex) lookup(spec ast.TypeSpecifier) (*Type, error) {
	switch spec.Kind() {
	case tps.Name:
		return x.lookupNamed(spec.(ast.TypeName).Name)
	case tps.Pointer:
		return x.lookupPointer(spec.(ast.PointerType).RefType)
	case tps.Struct:
		return x.lookupStruct(spec.(ast.StructType))
	case tps.Chunk:
		return x.lookupChunk(spec.(ast.ChunkType).ElemType)
	case tps.Enum:
		return x.lookupEnum(spec.(ast.EnumType))
	case tps.ArrayPointer:
		return x.lookupArrayPointer(spec.(ast.ArrayPointerType).ElemType)
	default:
		panic(fmt.Sprintf("not implemented for %s", spec.Kind().String()))
	}
}

func (x *TypeIndex) storePointer(ref *Type) *Type {
	return x.store(newPointerType(ref))
}

func (x *TypeIndex) storeArrayPointer(ref *Type) *Type {
	return x.store(newArrayPointerType(ref))
}

func (x *TypeIndex) storeChunk(elem *Type) *Type {
	return x.store(newChunkType(elem))
}

// store saves a given type if it is not yet known to index
// or returns and old one with the same stable hash.
func (x *TypeIndex) store(a *Type) *Type {
	stable := a.Stable()
	t := x.tm[stable]
	if t == nil {
		x.tm[stable] = a
		t = a
	}
	return t
}

func (x *TypeIndex) lookupArrayPointer(spec ast.TypeSpecifier) (*Type, error) {
	elem, err := x.lookup(spec)
	if err != nil {
		return nil, err
	}
	return x.storeArrayPointer(elem), nil
}

func (x *TypeIndex) lookupEnum(spec ast.EnumType) (*Type, error) {
	panic("not implemented")
	return nil, nil
}

func (x *TypeIndex) lookupChunk(spec ast.TypeSpecifier) (*Type, error) {
	elem, err := x.lookup(spec)
	if err != nil {
		return nil, err
	}
	return x.storeChunk(elem), nil
}

func (x *TypeIndex) lookupStruct(spec ast.StructType) (*Type, error) {
	if len(spec.Fields) == 0 {
		return Trivial, nil
	}

	var members MembersList
	members.Init(len(spec.Fields))

	for _, f := range spec.Fields {
		t, err := x.lookup(f.Type)
		if err != nil {
			return nil, err
		}

		members.Add(Member{
			Pos:  f.Name.Pos,
			Name: f.Name.Lit,
			Kind: MemberField,
			Type: t,
		})
	}

	return x.store(newStructType(members)), nil
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
