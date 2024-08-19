package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/tps"
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/enums/tpk"
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

// Result maps chunk elem type into chunk type itself.
func (x *TypeIndex) Chunks() map[*Type]*Type {
	if len(x.tm) == 0 {
		return nil
	}
	m := make(map[*Type]*Type)
	for _, t := range x.tm {
		if t.Kind == tpk.Chunk {
			e := t.Def.(ChunkTypeDef).ElemType
			m[e] = t
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

// Result maps array elem type into list of array types with this element.
func (x *TypeIndex) Arrays() map[*Type][]*Type {
	if len(x.tm) == 0 {
		return nil
	}
	m := make(map[*Type][]*Type)
	for _, t := range x.tm {
		if t.Kind == tpk.Array {
			e := t.Def.(ArrayTypeDef).ElemType
			m[e] = append(m[e], t)
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

func (x *TypeIndex) Lookup(ctx *Context, spec ast.TypeSpec) (*Type, error) {
	if ctx == nil {
		panic("nil context")
	}
	if spec == nil {
		// handle never and void return "types"
		return nil, nil
	}

	return x.lookup(ctx, spec)
}

func (x *TypeIndex) lookup(ctx *Context, spec ast.TypeSpec) (*Type, error) {
	switch spec.Kind() {
	case tps.Name:
		return x.lookupNamed(spec.(ast.TypeName).Name)
	case tps.Pointer:
		return x.lookupPointer(ctx, spec.(ast.PointerType).RefType)
	case tps.Struct:
		return x.lookupStruct(ctx, spec.(ast.StructType))
	case tps.Chunk:
		return x.lookupChunk(ctx, spec.(ast.ChunkType).ElemType)
	case tps.Enum:
		return x.lookupEnum(ctx, spec.(ast.EnumType))
	case tps.ArrayPointer:
		return x.lookupArrayPointer(ctx, spec.(ast.ArrayPointerType).ElemType)
	case tps.Array:
		s := spec.(ast.ArrayType)
		return x.lookupArray(ctx, s.ElemType, s.Size)
	case tps.ImportName:
		i := spec.(ast.ImportType)
		return x.lookupImportName(i.Unit, i.Name)
	case tps.RawMemoryPointer:
		return RawMemoryPointerType, nil
	default:
		panic(fmt.Sprintf("not implemented for %s", spec.Kind().String()))
	}
}

func (x *TypeIndex) lookupImportName(unit ast.Identifier, name ast.Identifier) (*Type, error) {
	uname := unit.Lit
	upos := unit.Pos
	usymbol := x.scope.Lookup(uname, upos.Num)
	if usymbol == nil {
		return nil, fmt.Errorf("%s: undefined symbol \"%s\"", upos, uname)
	}
	if usymbol.Kind != smk.Import {
		return nil, fmt.Errorf("%s: symbol \"%s\" is not an import", upos, uname)
	}
	u := usymbol.Def.(ImportSymDef).Unit
	n := name.Lit
	pos := name.Pos
	s := u.Scope.sym(n)
	if s == nil {
		return nil, fmt.Errorf("%s: unit %s has no \"%s\" symbol", pos, u.Name, name)
	}
	if s.Kind != smk.Type {
		return nil, fmt.Errorf("%s: imported symbol \"%s.%s\" is a %s, not a type", pos, u.Name, name, s.Kind)
	}
	if !s.Pub {
		return nil, fmt.Errorf("%s: imported symbol \"%s.%s\" is not public", pos, u.Name, name)
	}
	t := s.Def.(*Type)
	if t.Builtin() || t.Kind == tpk.Custom {
		return t, nil
	}
	panic(fmt.Sprintf("unexpected type kind: %s", t.Kind.String()))
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

func (x *TypeIndex) storeArray(elem *Type, length uint64) *Type {
	return x.store(newArrayType(elem, length))
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

func (x *TypeIndex) lookupArray(ctx *Context, spec ast.TypeSpec, length ast.Expression) (*Type, error) {
	elem, err := x.lookup(ctx, spec)
	if err != nil {
		return nil, err
	}
	expr, err := x.scope.scan(ctx, length)
	if err != nil {
		return nil, err
	}
	if !expr.Type().IsIntegerType() {
		return nil, fmt.Errorf("%s: only integer types can be used for array length", length.Pin())
	}
	r, err := x.scope.evalStaticExp(expr)
	if err != nil {
		return nil, err
	}

	// extract array length from evaluated integer
	i := r.(Integer)
	if i.Neg {
		return nil, fmt.Errorf("%s: array cannot have negative length", length.Pin())
	}
	l := i.Val
	if l == 0 {
		return Trivial, nil
	}

	return x.storeArray(elem, l), nil
}

func (x *TypeIndex) lookupArrayPointer(ctx *Context, spec ast.TypeSpec) (*Type, error) {
	elem, err := x.lookup(ctx, spec)
	if err != nil {
		return nil, err
	}
	return x.storeArrayPointer(elem), nil
}

func (x *TypeIndex) lookupEnum(ctx *Context, spec ast.EnumType) (*Type, error) {
	base, err := x.lookup(ctx, spec.Base)
	if err != nil {
		return nil, err
	}
	if !base.IsIntegerType() {
		return nil, fmt.Errorf("%s: \"%s\" is not an integer type", spec.Pos, spec.Base.Name)
	}

	var entries []EnumEntry
	if len(spec.Entries) != 0 {
		entries = make([]EnumEntry, 0, len(entries))
	}
	for _, entry := range spec.Entries {
		if entry.Expression != nil {
			panic("not implemented")
		}

		entries = append(entries, EnumEntry{
			Pos:  entry.Name.Pos,
			Name: entry.Name.Lit,
		})
	}

	return newEnumType(base, entries)
}

func (x *TypeIndex) lookupChunk(ctx *Context, spec ast.TypeSpec) (*Type, error) {
	elem, err := x.lookup(ctx, spec)
	if err != nil {
		return nil, err
	}
	return x.storeChunk(elem), nil
}

func (x *TypeIndex) lookupStruct(ctx *Context, spec ast.StructType) (*Type, error) {
	if len(spec.Fields) == 0 {
		return Trivial, nil
	}

	var members MembersList
	members.Init(len(spec.Fields))

	for _, f := range spec.Fields {
		t, err := x.lookup(ctx, f.Type)
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

func (x *TypeIndex) lookupPointer(ctx *Context, spec ast.TypeSpec) (*Type, error) {
	ref, err := x.lookup(ctx, spec)
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
		return nil, fmt.Errorf("%s: undefined symbol \"%s\"", pos, name)
	}
	if s.Kind != smk.Type {
		return nil, fmt.Errorf("%s: symbol \"%s\" is a %s, not a type", pos, name, s.Kind)
	}
	t := s.Def.(*Type)
	if t.Builtin() || t.Kind == tpk.Custom {
		return t, nil
	}
	panic(fmt.Sprintf("unexpected type kind: %s", t.Kind))
}
