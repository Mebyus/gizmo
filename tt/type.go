package tt

import (
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

	// Discriminator for type definition category.
	Kind typ.Kind

	// True for types that are language builtins.
	Builtin bool
}

func (m *Merger) lookupType(spec ast.TypeSpecifier) *Type {
	if spec == nil {
		return nil
	}
	return &Type{}
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
