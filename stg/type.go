package stg

// TODO: rename package to tgr - Type Graph Representation

import (
	"encoding/binary"
	"fmt"
	"hash"
	"hash/fnv"

	"github.com/mebyus/gizmo/enums/tpk"
	"github.com/mebyus/gizmo/source"
)

// Type represents a value type of symbol, field, expression or subexpression
// in a program. In other words when something is used in the program source code
// type describes what kind of result that usage yields.
type Type struct {
	nodeSymDef

	// Type definition.
	//
	// Meaning and contents of this field heavily depends on the kind of this type.
	// List of kinds for which this field is always nil:
	//
	//	- StaticInteger
	//	- StaticFloat
	//	- StaticString
	//	- StaticBoolean
	//	- StaticNil
	//	- Named
	//	- String
	//	- Boolean
	Def TypeDef

	// Pseudo-unique id (hash code) of the type. Depends only on type definition
	// of represented type and (maybe) compiler version. One should use this field
	// in conjunction with kind to reduce collision chance of different types.
	//
	// Does not depend upon:
	//	- position in source code
	//	- formatting
	//	- compiler run (stays the same across different builds)
	//
	// For named types Unit field should be used together with this field. In general
	// for robust unique type id one should use Stable() method.
	hash uint64

	// Not 0 only for user-defined (not builtin) named types. Contains hash id of a Unit
	// where the type is defined.
	Unit uint64

	// Bit flags with additional type properties. Actual meaning may differ
	// upon Kind.
	Flags TypeFlag

	// Byte size of this type's value. May be 0 for some types.
	// More specifically this field equals the stride between two
	// consecutive elements of this type inside an array.
	Size uint32

	// Discriminator for type definition category.
	Kind tpk.Kind
}

// TypeFlag bit flags for specifing additional type properties.
type TypeFlag uint64

const (
	// Static variant of the type.
	TypeFlagStatic TypeFlag = 1 << iota

	// Type is a builtin.
	TypeFlagBuiltin

	// Type has recursive definition.
	TypeFlagRecursive

	// Signed integer type.
	TypeFlagSigned
)

// Custom returns true for custom (user defined opposed to language builtin) types.
func (t *Type) Custom() bool {
	return t.Kind == tpk.Custom
}

func (t *Type) Static() bool {
	return t.Flags&TypeFlagStatic != 0
}

// Returns true for types that are language builtins.
func (t *Type) Builtin() bool {
	return t.Flags&TypeFlagBuiltin != 0
}

// Returns true for types which definition references itself.
func (t *Type) Recursive() bool {
	return t.Flags&TypeFlagRecursive != 0
}

// Returns true for signed integer type and false otherwise.
// In particular returns false for unsigned integer types.
func (t *Type) Signed() bool {
	return t.Flags&TypeFlagSigned != 0
}

// Returns true for static integer type.
func (t *Type) IsStaticInteger() bool {
	return t.Kind == tpk.Integer && t.Static() && t.Size == 0
}

// Returns true for static string type.
func (t *Type) IsStaticString() bool {
	return t.Kind == tpk.String && t.Static() && t.Size == 0
}

// Returns true for static boolean type.
func (t *Type) IsStaticBoolean() bool {
	return t.Kind == tpk.Boolean && t.Static() && t.Size == 0
}

func (t *Type) Symbol() *Symbol {
	if t.Builtin() {
		return t.Def.(BuiltinTypeDef).Symbol
	}

	if t.Kind == tpk.Custom {
		return t.Def.(CustomTypeDef).Symbol
	}

	panic(fmt.Sprintf("%s types cannot be bound to symbols", t.Kind))
}

func (t *Type) ElemType() *Type {
	switch t.Kind {
	case tpk.Chunk:
		return t.Def.(ChunkTypeDef).ElemType
	case tpk.Array:
		return t.Def.(ArrayTypeDef).ElemType
	default:
		panic(fmt.Sprintf("%s types do not have elements", t.Kind))
	}
}

func (t *Type) IsIntegerType() bool {
	switch t.Kind {
	case tpk.Integer:
		return true
	case tpk.Custom:
		return t.Def.(CustomTypeDef).Base.Kind == tpk.Integer
	default:
		return false
	}
}

func (t *Type) IsStringType() bool {
	switch t.Kind {
	case tpk.String:
		return true
	case tpk.Custom:
		return t.Def.(CustomTypeDef).Base.Kind == tpk.String
	default:
		return false
	}
}

func (t *Type) IsBooleanType() bool {
	switch t.Kind {
	case tpk.Boolean:
		return true
	case tpk.Custom:
		return t.Def.(CustomTypeDef).Base.Kind == tpk.Boolean
	default:
		return false
	}
}

func (t *Type) IsEnumType() bool {
	switch t.Kind {
	case tpk.Enum:
		return true
	case tpk.Custom:
		return t.Def.(CustomTypeDef).Base.Kind == tpk.Enum
	default:
		return false
	}
}

// Stable is a unique identifier of a type inside a program. Via this identifier
// different *Type instances (with different pointers) can be compared to
// establish identity of represented types.
//
// Main scenario for Stable usage is as follows:
//
//  1. Encounter TypeSpecifier during AST traversal
//  2. Create fresh &Type{} struct and populate it with data according to TypeSpecifier
//  3. Compute Type.Stable()
//  4. Search Stable lookup table for already created identical type
//  5. Add new Stable entry to lookup table if no existing type was found
type Stable struct {
	// equal to hash field in type struct
	Hash uint64

	// hash of unit in which type is defined (if it's a named type)
	Unit uint64

	Flags TypeFlag

	// kind of type itself
	Kind tpk.Kind
}

func (t *Type) Stable() Stable {
	return Stable{
		Hash:  t.Hash(),
		Unit:  t.Unit,
		Flags: t.Flags,
		Kind:  t.Kind,
	}
}

func HashName(name string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(name))
	return h.Sum64()
}

func HashRecursiveName(name string) uint64 {
	h := fnv.New64a()
	h.Write([]byte{0})
	h.Write([]byte(name))
	return h.Sum64()
}

func (t *Type) Hash() uint64 {
	if t.hash != 0 {
		return t.hash
	}
	t.hash = t.computeHash()
	if t.hash == 0 {
		panic("hashing a type should not produce 0")
	}
	return t.hash
}

func (t *Type) computeHash() uint64 {
	if t.IsStaticInteger() {
		return uint64(tpk.Integer)
	}
	if t.IsStaticString() {
		return uint64(tpk.String)
	}
	if t.IsStaticBoolean() {
		return uint64(tpk.Boolean)
	}
	if t.Builtin() {
		return HashName(t.Symbol().Name)
	}

	switch t.Kind {
	case tpk.Custom:
		if t.Recursive() {
			return HashRecursiveName(t.Def.(CustomTypeDef).Symbol.Name)
		}
		// TODO: we probably should panic here, since
		// named types are not meant to be looked up by hash
		return HashName(t.Def.(CustomTypeDef).Symbol.Name)
	case tpk.Pointer:
		return HashPointerType(t.Def.(PointerTypeDef).RefType)
	case tpk.ArrayPointer:
		return HashArrayPointerType(t.Def.(ArrayPointerTypeDef).RefType)
	case tpk.Chunk:
		return HashChunkType(t.Def.(ChunkTypeDef).ElemType)
	case tpk.Struct:
		return HashStructType(t)
	case tpk.Array:
		return HashArrayType(t)
	case tpk.StaticFloat:
		return uint64(tpk.StaticFloat)
	case tpk.StaticNil:
		return uint64(tpk.StaticNil)
	default:
		panic(fmt.Sprintf("not implemented for %s", t.Kind.String()))
	}
}

func putUint64(buf []byte, x uint64) {
	binary.LittleEndian.PutUint64(buf, x)
}

func HashArrayType(t *Type) uint64 {
	def := t.Def.(ArrayTypeDef)
	elem := def.ElemType

	var buf [19]byte
	h := fnv.New64a()
	buf[0] = byte(tpk.Array)
	putUint64(buf[2:10], elem.Hash())
	putUint64(buf[11:], def.Len)
	return h.Sum64()
}

func HashChunkType(elem *Type) uint64 {
	var buf [10]byte

	h := fnv.New64a()
	buf[0] = byte(tpk.Chunk)
	putUint64(buf[2:], elem.Hash())
	h.Write(buf[:])
	return h.Sum64()
}

func HashPointerType(ref *Type) uint64 {
	var buf [10]byte

	h := fnv.New64a()
	buf[0] = byte(tpk.Pointer)
	putUint64(buf[2:], ref.Hash())
	h.Write(buf[:])
	return h.Sum64()
}

func HashArrayPointerType(ref *Type) uint64 {
	var buf [10]byte

	h := fnv.New64a()
	buf[0] = byte(tpk.ArrayPointer)
	putUint64(buf[2:], ref.Hash())
	h.Write(buf[:])
	return h.Sum64()
}

func HashStructType(t *Type) uint64 {
	fields := t.Def.(*StructTypeDef).Fields

	h := fnv.New64a()
	var buf [1]byte
	buf[0] = byte(tpk.Struct)
	h.Write(buf[:])
	for i := range len(fields) {
		m := &fields[i]
		hashStructField(h, m)
	}
	return h.Sum64()
}

func hashStructField(h hash.Hash64, field *Field) {
	var buf [10]byte
	putUint64(buf[1:], field.Type.Hash())
	h.Write(buf[:])
	h.Write([]byte(field.Name))
}

type TypeDef interface {
	TypeDef()
}

// This is dummy implementation of TypeDef interface.
//
// Used for embedding into other (non-dummy) type definition nodes.
type nodeTypeDef struct{}

func (nodeTypeDef) TypeDef() {}

type BuiltinTypeDef struct {
	nodeTypeDef

	Symbol *Symbol
}

type CustomTypeDef struct {
	nodeTypeDef

	// Symbol which creates custom type.
	Symbol *Symbol

	// Type which was given a custom name by the symbol.
	Base *Type

	// List of methods which are bound to this custom type.
	Methods []*Symbol
}

type ArrayPointerTypeDef struct {
	nodeTypeDef

	RefType *Type
}

type PointerTypeDef struct {
	nodeTypeDef

	RefType *Type
}

type ChunkTypeDef struct {
	nodeTypeDef

	ElemType *Type
}

type ArrayTypeDef struct {
	nodeTypeDef

	ElemType *Type

	// Number of elements in array.
	Len uint64
}

func newArrayType(elem *Type, len uint64) *Type {
	t := &Type{
		Kind: tpk.Array,
		Def: ArrayTypeDef{
			ElemType: elem,
			Len:      len,
		},
	}
	return t
}

func newArrayPointerType(ref *Type) *Type {
	t := &Type{
		Kind: tpk.ArrayPointer,
		Def:  ArrayPointerTypeDef{RefType: ref},
	}
	return t
}

func newPointerType(ref *Type) *Type {
	t := &Type{
		Kind: tpk.Pointer,
		Def:  PointerTypeDef{RefType: ref},
	}
	return t
}

func newChunkType(elem *Type) *Type {
	t := &Type{
		Kind: tpk.Chunk,
		Def:  ChunkTypeDef{ElemType: elem},
	}
	return t
}

func newStructType(fields []Field) *Type {
	t := &Type{
		Kind: tpk.Struct,
		Def:  &StructTypeDef{Fields: fields},
	}
	return t
}

func newEnumType(base *Type, entries []EnumEntry) (*Type, error) {
	var index map[string]*EnumEntry
	if len(entries) != 0 {
		index = make(map[string]*EnumEntry, len(entries))
	}
	for i := range len(entries) {
		e := &entries[i]
		if e.Exp != nil {
			panic("not implemented")
		}
		name := e.Name
		_, ok := index[name]
		if ok {
			return nil, fmt.Errorf("%s: duplicate enum entry \"%s\"", e.Pos, name)
		}

		index[name] = e
		e.Exp = Integer{
			Val: uint64(i),
			typ: StaticIntegerType,
		}
	}

	return &Type{
		Kind: tpk.Enum,
		Size: base.Size,
		Def: &EnumTypeDef{
			Base:    base,
			Entries: entries,

			index: index,
		},
	}, nil
}

// returns an error if argument type does not match parameter type
func typeCheckExp(want *Type, exp Exp) error {
	t := exp.Type()
	if t == want {
		return nil
	}

	switch want.Kind {
	case tpk.Integer:
		if t.Kind != tpk.Integer {
			return fmt.Errorf("%s: param expects integer type, but argument has %s type",
				exp.Pin(), t.Kind)
		}

		if t.IsStaticInteger() {
			// TODO: check that static integer fits into parameter type max value
			// and there is no "signedness conflict" between type and value
			return nil
		}

		if want.Signed() {
			if t.Signed() {
				if t.Size < want.Size {
					// argument fits into parameter type
					// silently cast signed integer
					// into bigger one
					return nil
				}
			} else {
				if t.Size < want.Size {
					// argument fits into parameter type
					// silently cast unsigned integer
					// into bigger signed one
					return nil
				}
			}
		} else {
			if t.Signed() {
				return fmt.Errorf("%s: cannot silently convert signed to unsigned integer", exp.Pin())
			} else {
				if t.Size < want.Size {
					// argument fits into parameter type
					// silently cast unsigned integer
					// into bigger one
					return nil
				}
			}
		}
	case tpk.Boolean:
		if t.IsStaticBoolean() {
			return nil
		}
	case tpk.String:
		if t.IsStaticString() {
			return nil
		}
	case tpk.Custom:
		panic("not implemented")
	default:
		fmt.Println("want ", want)
		fmt.Println("got t", t)
		panic(fmt.Errorf("%s param types not implemented", want.Kind.String()))
	}

	return fmt.Errorf("%s: mismatched types of call argument (%s) and parameter (%s)",
		exp.Pin().String(), t.Kind, want.Kind)
}

type StructTypeDef struct {
	nodeTypeDef

	Fields []Field
}

type EnumTypeDef struct {
	nodeTypeDef

	// Entries are ordered by their position in source code.
	Entries []EnumEntry

	// Enum base integer (storage) type.
	Base *Type

	// maps entry name to its index inside Entries slice
	index map[string]*EnumEntry
}

func (d *EnumTypeDef) Entry(name string) *EnumEntry {
	return d.index[name]
}

type EnumEntry struct {
	// Where this entry was defined.
	Pos source.Pos

	Name string

	// Can be nil before complete type is constructed
	// if entry does not have explicit assigned value.
	Exp Exp
}
