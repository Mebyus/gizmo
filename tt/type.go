package tt

// TODO: rename package to tgr - Type Graph Representation

import (
	"encoding/binary"
	"fmt"
	"hash"
	"hash/fnv"

	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/tt/typ"
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

	// Not nil only for custom and builtin types. Contains symbol which names
	// this type.
	Symbol *Symbol

	// Each type has a base (even builtin types). Base is a type representing
	// the raw definition (struct, enum, union, etc.) and memory layout of a type.
	// Has nothing to do with inheritance.
	//
	// Types which has the same base are called flavors of that base or type (base is always
	// a flavor of itself by definition). Flavors can be cast among each other via
	// explicit type casts and under the hood such casts do not alter memory or
	// perform any operation. One and only one base exists among all of its flavors.
	// From formal mathematical point of view all flavors of specific type form
	// equivalence class in set of all types in a program.
	Base *Type

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

	// Byte size of this type's value. May be 0 for some types.
	// More specifically this field equals the stride between two
	// consecutive elements of this type inside an array.
	Size uint32

	// Discriminator for type definition category.
	Kind typ.Kind

	// True for types that are language builtins.
	Builtin bool

	// True for types which definition references itself.
	Recursive bool
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

	// equal to hash field in Type.Base struct
	BaseHash uint64

	// hash of unit in which type is defined (if it's a named type)
	Unit uint64

	// kind of type itself
	Kind typ.Kind

	// kind of Type.Base
	BaseKind typ.Kind
}

func (t *Type) Stable() Stable {
	return Stable{
		Hash: t.Hash(),
		Unit: t.Unit,
		Kind: t.Kind,

		BaseHash: t.Base.Hash(),
		BaseKind: t.Base.Kind,
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
	if t.Builtin {
		if t.Symbol.Name == "" {
			panic("empty name of builtin type")
		}
		return HashName(t.Symbol.Name)
	}

	switch t.Kind {
	case typ.Named:
		if t.Recursive {
			return HashRecursiveName(t.Symbol.Name)
		}
		// TODO: we probably should panic here, since
		// named types are not meant to be looked up by hash
		return HashName(t.Symbol.Name)
	case typ.Pointer:
		return HashPointerType(t.Def.(PointerTypeDef).RefType)
	case typ.ArrayPointer:
		return HashArrayPointerType(t.Def.(ArrayPointerTypeDef).RefType)
	case typ.Chunk:
		return HashChunkType(t.Def.(ChunkTypeDef).ElemType)
	case typ.Struct:
		return HashStructType(t)
	case typ.StaticInteger:
		return uint64(typ.StaticInteger)
	case typ.StaticBoolean:
		return uint64(typ.StaticBoolean)
	case typ.StaticFloat:
		return uint64(typ.StaticFloat)
	case typ.StaticString:
		return uint64(typ.StaticString)
	case typ.StaticNil:
		return uint64(typ.StaticNil)
	default:
		panic(fmt.Sprintf("not implemented for %s", t.Kind.String()))
	}
}

func putUint64(buf []byte, x uint64) {
	binary.LittleEndian.PutUint64(buf, x)
}

func HashChunkType(elem *Type) uint64 {
	var buf [10]byte

	h := fnv.New64a()
	buf[0] = byte(typ.Chunk)
	putUint64(buf[2:], elem.Hash())
	h.Write(buf[:])
	return h.Sum64()
}

func HashPointerType(ref *Type) uint64 {
	var buf [10]byte

	h := fnv.New64a()
	buf[0] = byte(typ.Pointer)
	putUint64(buf[2:], ref.Hash())
	h.Write(buf[:])
	return h.Sum64()
}

func HashArrayPointerType(ref *Type) uint64 {
	var buf [10]byte

	h := fnv.New64a()
	buf[0] = byte(typ.ArrayPointer)
	putUint64(buf[2:], ref.Hash())
	h.Write(buf[:])
	return h.Sum64()
}

func HashStructType(t *Type) uint64 {
	members := t.Def.(*StructTypeDef).Members.Members

	h := fnv.New64a()
	var buf [1]byte
	buf[0] = byte(typ.Struct)
	h.Write(buf[:])
	for i := 0; i < len(members); i += 1 {
		m := &members[i]
		hashStructField(h, m)
	}
	return h.Sum64()
}

func hashStructField(h hash.Hash64, m *Member) {
	var buf [10]byte
	putUint64(buf[1:], m.Type.Hash())
	h.Write(buf[:])
	h.Write([]byte(m.Name))
}

type TypeDef interface {
	TypeDef()
}

// This is dummy implementation of TypeDef interface.
//
// Used for embedding into other (non-dummy) type definition nodes.
type nodeTypeDef struct{}

func (nodeTypeDef) TypeDef() {}

type IntTypeDef struct {
	nodeTypeDef

	// Integer size in bytes.
	Size uint32
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

func newArrayPointerType(ref *Type) *Type {
	t := &Type{
		Kind: typ.ArrayPointer,
		Def:  ArrayPointerTypeDef{RefType: ref},
	}
	t.Base = t
	return t
}

func newPointerType(ref *Type) *Type {
	t := &Type{
		Kind: typ.Pointer,
		Def:  PointerTypeDef{RefType: ref},
	}
	t.Base = t
	return t
}

func newChunkType(elem *Type) *Type {
	t := &Type{
		Kind: typ.Chunk,
		Def:  ChunkTypeDef{ElemType: elem},
	}
	t.Base = t
	return t
}

func newStructType(members MembersList) *Type {
	t := &Type{
		Kind: typ.Struct,
		Def:  &StructTypeDef{Members: members},
	}
	t.Base = t
	return t
}

// returns an error if argument type does not match parameter type
func checkCallArgType(param *Symbol, arg Expression) error {
	t := arg.Type()
	pt := param.Type
	if t == pt {
		return nil
	}

	if (pt.Base.Kind == typ.Unsigned || pt.Base.Kind == typ.Signed) && t.Kind == typ.StaticInteger {
		// TODO: check that static integer fits into parameter type max value
		// and there is no "signedness conflict" between type and value
		return nil
	}

	// panic("unhandled case")
	return fmt.Errorf("%s: mismatched types of call argument (%s) and parameter (%s)",
		arg.Pin(), t.Kind, pt.Kind)
}

type MemberKind uint8

const (
	memberEmpty MemberKind = iota

	MemberField
	MemberMethod
)

var memberText = [...]string{
	memberEmpty: "<nil>",

	MemberField:  "field",
	MemberMethod: "method",
}

func (k MemberKind) String() string {
	return memberText[k]
}

type Member struct {
	// position where this member is defined.
	Pos source.Pos

	// Field or method name. Cannot be empty.
	Name string

	Type *Type

	// Always nil for field members.
	Def *MethodDef

	// Index of member inside the list of members.
	Index int

	Kind MemberKind
}

type MembersList struct {
	Members []Member

	// maps member name to its index inside Members slice
	index map[string]int
}

func (l *MembersList) Init(size int) {
	if size == 0 {
		return
	}

	l.Members = make([]Member, 0, size)
	l.index = make(map[string]int, size)
}

func (l *MembersList) Find(name string) *Member {
	i, ok := l.index[name]
	if !ok {
		return nil
	}
	return &l.Members[i]
}

func (l *MembersList) Add(member Member) {
	l.index[member.Name] = len(l.Members)
	member.Index = len(l.Members)
	l.Members = append(l.Members, member)
}

type StructTypeDef struct {
	nodeTypeDef

	Members MembersList
}

type NamedTypeDef struct {
	nodeTypeDef

	// Symbol which was used to define a named type.
	Symbol *Symbol
}

type EnumTypeDef struct {
	nodeTypeDef

	Entries []EnumEntry

	// maps entry name to its index inside Entries slice
	index map[string]int
}

type EnumEntry struct {
	// Where this entry was defined.
	Pos source.Pos

	Name string

	// Can be nil if entry does not have explicit assigned value
	Expression Expression
}
