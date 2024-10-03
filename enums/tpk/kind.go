package tpk

// Kind indicates type kind.
type Kind uint8

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left unspecified
	empty Kind = iota

	// Trivial type has no size or properties. It is formed by language
	// constructs such as:
	//
	//	struct {} // empty struct
	//	[0]int    // array with zero size
	//
	// There is always only one trivial type.
	Trivial

	// Type for any integers (literals, constants and expressions)
	// excluding custom integer types.
	//
	// Flags specify whether this integer type has specified storage
	// size, is signed or unsigned, static or runtime type.
	//
	// Types with fixed bit size which can hold unsigned integers:
	//
	//	u8, u16, u32, u64, u128, uint
	//
	// Types with fixed bit size which can hold signed integers:
	//
	//	s8, s16, s32, s64, s128, sint
	//
	// Static integer type of arbitrary size is specified as:
	//
	//	<int>
	//
	// And this static type has size set to 0.
	Integer

	// Type for static floats (literals, constants and expressions).
	//
	// There is always only one Type of this category and it cannot
	// have Flavors
	StaticFloat

	// Builtin type which can hold fixed number of bytes. Inner structure
	// of this type is identical to chunk of bytes []u8. The difference
	// between these two types is logical, in most cases strings should
	// hold utf-8 encoded text opposed to arbitrary byte sequence in
	// chunk of bytes. However this is not a rule and not enforced by the
	// language in any way. Strings and chunks of bytes can be cast between
	// each other freely with no runtime overhead.
	//
	// Type for any strings (literals, constants and exressions).
	//
	// Static string type is specified as:
	//
	//	<str>
	//
	// And this static type has size set to 0.
	String

	// Builtin type which can hold one of two mutually exclusive values:
	//
	//	- true
	//	- false
	//
	// This type is not a flavor of any integer type. Furthermore boolean
	// flavors cannot be cast into integer flavors via cast operation.
	//
	// Type for any boolean values (literals, constants and exressions)
	//
	// Flags specify whether this boolean type is static or runtime type.
	//
	// Static boolean type is specified as:
	//
	//	<bool>
	//
	// And this static type has size set to 0.
	Boolean

	// Type for static nil value (literal, constants and exressions)
	//
	// There is always only one Type of this category and it cannot
	// have Flavors
	StaticNil

	AnyPointer

	// Types of symbols, created by importing other unit
	//
	//	import std {
	//		math => "math"
	//	}
	//
	// Symbol math in the example above has Type with Kind = Unit.
	// each distinct unit there will be each own Base type which
	// equals this one
	Unit

	// Types which are defined by giving a new name to another type.
	//
	// Created via language construct:
	//
	//	type Name <OtherType>
	Custom

	// Types which are defined as structs
	//
	// Specified via language construct:
	//
	//	struct { <Fields> }
	Struct

	// Types which are defined as tuples
	//
	// Specified via language construct:
	//
	//	(<type_1>, <type_2>, ...) - for unnamed tuples
	//	(<name_1>: <type_1>, <name_2>: <type_2>, ...) - for named tuples
	//
	// Tuples can only be used for specifying return values of functions and
	// cannot be named (or have Flavors through other means). Thus users of the
	// language cannot define their own tuple types with unique names
	//
	// Tuples are "behind-the-scenes" mechanism for implementing functions
	// with multiple return values and passing arguments from one function to
	// another in constructs like these:
	//
	//	foo(bar())
	//
	// or like these:
	//
	//	r := bar()
	//	foo(r)
	//
	// where function signatures of foo and bar are defined (for example) as:
	//
	//	fn bar() => (str, usz)
	//	fn foo(str, usz)
	//
	// Our reasoning for such restricted usage of tuples in the language are
	// mostly as follows:
	//
	// 1. For named tuples there are already structs present in the language,
	// no reason to multiply confusion by adding the same functionality in
	// different wrapping
	//
	// 2. Unnamed tuples should be treated as short-lived opaque containers and
	// unpacked into constants/variables with appropriate naming or handed over to
	// another function. There is no reason to introduce language constructs like these:
	//
	//	r := bar()
	//	first := r.0
	//
	// Accessing individual fields of unnamed tuples via "dot" syntax with numerical
	// "field" names is just confusing and should not be allowed. Instead programmers
	// are encouraged to use either named return values to clarify meaning of said values
	// or a struct to begin with, which is much easier for field naming and leaving
	// comments to those fields
	//
	// In short:
	//
	// Structs as named types are superior to tuples in all practical ways and thus
	// tuple serve a very niche purpose of passing around and unpacking returned values
	// from a function
	//
	// In case of named return values (named tuples) rules are not so strict. Individual
	// fields can be accessed via their names:
	//
	//	fn bar() => (my_string: str, count: usz)
	//
	//	r := bar()
	//	s := r.my_string
	//	c := r.count
	Tuple

	// Types which are defined as named tuples
	// TODO: consider renaming concept "named tuples" to "formations" across the project
	//
	// See comment to Tuple for more info
	Formation

	// Types which are defined as continuous region of memory holding
	// fixed number of elements of another Type
	//
	//	[5]u64 - an array of five u64 numbers
	Array

	// Pointer to another Type. An example of such a Type could be *str
	// (pointer to a string)
	Pointer

	// Pointer to continuous region of memory holding
	// variable number of elements of another type.
	ArrayPointer

	// Abstract enum Types. Values of this Types can only be obtained via
	// their enum respective symbols. Furthermore enum Flavors cannot be
	// cast into integer Flavors
	Enum

	// Types which are defined as function signatures
	//
	// Specified via language construct:
	//
	//	fn (<Args>) <Result>
	//
	// This type is a base type for functions, bound methods,
	// unbound methods, custom types (defined as signatures)
	// and closures.
	Signature

	Function

	// Created by usage of language construct:
	//
	//	T.add
	//
	// Where add is a method on type symbol a.
	UnboundMethod

	// Created by usage of language construct:
	//
	//	a.add
	//
	// Where add is a method on value symbol a.
	BoundMethod

	Closure

	// Types which are obtained via usage of type specifiers. In other
	// words this category contains Types of Types. For example:
	//
	//	- []u8: yields a type specifier for chunk of bytes
	//	- string: yields a type specifier for string
	//
	// Type of these specifiers is also a Type with Kind = Spec.
	//
	// Most common usage for Types from this category is narrowing type
	// parameters in Prototypes and Blueprints
	Spec

	// Parameterized Prototype (mechanism for generic types)
	Proto

	// Parameterized Blueprint (mechanism for generic functions and methods)
	Blue

	// Chunk of another Type. An example of such a Type could be []u8
	// (chunk of u8 integers), which is commonly used to represent
	// raw bytes of data
	Chunk

	// Types with fixed bit size which can hold floating point numbers:
	//
	//	f32, f64
	Float

	// Types which can have nil value at runtime
	Nillable
)

var text = [...]string{
	empty: "<nil>",

	Integer:     "integer",
	StaticFloat: "static.float",
	StaticNil:   "static.nil",

	AnyPointer: "pointer.rawmem",

	Trivial:      "trivial",
	Custom:       "custom",
	Struct:       "struct",
	Tuple:        "tuple",
	Formation:    "formation",
	Array:        "array",
	Pointer:      "pointer",
	ArrayPointer: "pointer.array",
	Enum:         "enum",
	Signature:    "function",
	Spec:         "spec",
	Proto:        "proto",
	Blue:         "blue",
	Chunk:        "chunk",
	Float:        "float",
	String:       "string",
	Boolean:      "boolean",
	Nillable:     "nillable",
}

func (k Kind) String() string {
	return text[k]
}
