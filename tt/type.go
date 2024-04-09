package tt

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/tt/typ"
)

// Type represents a value type of symbol, field, expression or subexpression
// in a program. In other words when something is used in the program source code
// type describes what kind of result that usage yields.
type Type struct {
	nodeSymDef

	// Only named and builtin types have a name.
	// Example of creating a type with name:
	//
	//	// Create a named type from builtin "int" type
	//	type MyType int
	//
	// On the other hand []i32 is just a chunk of sized
	// i32 numbers with no specific name for that type.
	Name string

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

	// Discriminator for type definition category.
	Kind typ.Kind

	// True for types that are language builtins.
	Builtin bool
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

func (m *Merger) lookupType(spec ast.TypeSpecifier) *Type {
	if spec == nil {
		return nil
	}
	return &Type{}
}

func HashName(name string) uint64 {
	h := fnv.New64a()
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
		if t.Name == "" {
			panic("empty name of builtin type")
		}
		return HashName(t.Name)
	}

	switch t.Kind {
	case typ.Named:
		return HashName(t.Name)
	case typ.Pointer:
		return HashPointerType(t.Def.(PtrTypeDef).RefType)
	default:
		panic(fmt.Sprintf("not implemented for %s", t.Kind.String()))
	}
}

func putUint64(buf []byte, x uint64) {
	binary.LittleEndian.PutUint64(buf, x)
}

func HashPointerType(ref *Type) uint64 {
	var buf [10]byte

	h := fnv.New64a()
	buf[0] = byte(typ.Pointer)
	putUint64(buf[2:], ref.Hash())
	h.Write(buf[:])
	return h.Sum64()
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

type PtrTypeDef struct {
	nodeTypeDef

	RefType *Type
}

func newPointerType(ref *Type) *Type {
	t := &Type{
		Kind: typ.Pointer,
		Def:  PtrTypeDef{RefType: ref},
	}
	t.Base = t
	return t
}
