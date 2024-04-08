package typ

// Kind indicates type kind
type Kind uint8

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left unspecified
	empty Kind = iota

	// Type for static integers (literals, constants and expressions).
	// Includes positive, negative and zero integers.
	//
	// There is always only one type of this category and it is the base
	// of all static positive integer flavors.
	StaticInteger

	// Type for static floats (literals, constants and expressions).
	//
	// There is always only one Type of this category and it cannot
	// have Flavors
	StaticFloat

	// Type for static strings (literals, constants and exressions)
	//
	// There is always only one Type of this category and it cannot
	// have Flavors
	StaticString

	// Type for static boolean values (literals, constants and exressions)
	//
	// There is always only one Type of this category and it cannot
	// have Flavors
	StaticBoolean

	// Type for static nil value (literal, constants and exressions)
	//
	// There is always only one Type of this category and it cannot
	// have Flavors
	StaticNil

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
	// Most commonly created via language construct:
	//
	//	type Name <OtherType>
	Named

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

	// Abstract enum Types. Values of this Types can only be obtained via
	// their enum respective symbols. Furthermore enum Flavors cannot be
	// cast into integer Flavors
	Enum

	// Types which are defined as function signatures
	//
	// Specified via language construct:
	//
	//	fn (<Args>) <Result>
	Function

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

	// Types with fixed bit size which can hold unsigned integers:
	//
	//	u8, u16, u32, u64, u128, uint
	Unsigned

	// Types with fixed bit size which can hold signed integers:
	//
	//	i8, i16, i32, i64, i128, int
	Signed

	// Types with fixed bit size which can hold floating point numbers:
	//
	//	f32, f64
	Float

	// Builtin type which can hold fixed number of bytes. Inner structure
	// of this type is identical to chunk of bytes []u8. The difference
	// between these two types is logical, in most cases strings should
	// hold utf-8 encoded text opposed to arbitrary byte sequence in
	// chunk of bytes. However this is not a rule and not enforced by the
	// language in any way. Strings and chunks of bytes can be cast between
	// each other freely with no runtime overhead.
	String

	// Builtin type which can hold one of two mutually exclusive values:
	//
	//	- true
	//	- false
	//
	// There is always only one Type of this category and it is the Base
	// of all boolean Flavors
	//
	// This Type is not a Flavor of any Integer Type. Furthermore Boolean
	// Flavors cannot be cast into integer Flavors
	Boolean

	// Types which can have nil value at runtime
	Nillable
)

var text = [...]string{
	StaticInteger: "static_integer",
	StaticFloat:   "static_float",
	StaticString:  "static_string",
	StaticBoolean: "static_boolean",
	StaticNil:     "static_nil",
	Named:         "named",
	Struct:        "struct",
	Tuple:         "tuple",
	Formation:     "formation",
	Array:         "array",
	Pointer:       "pointer",
	Enum:          "enum",
	Function:      "function",
	Spec:          "spec",
	Proto:         "proto",
	Blue:          "blue",
	Chunk:         "chunk",
	Unsigned:      "unsigned",
	Signed:        "signed",
	Float:         "float",
	String:        "string",
	Boolean:       "boolean",
	Nillable:      "nillable",
}

func (k Kind) String() string {
	return text[k]
}

