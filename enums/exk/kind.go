package exk

// Kind indicates expression kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left unspecified
	empty Kind = iota

	// basic literal operand
	Basic

	// object literal operand
	Object

	// list literal operand
	List

	// direct substitution of identifier operand
	Symbol

	// incomplete name usage
	IncompName

	Unary
	Binary

	Cast
	Tint
	MemCast // TODO: maybe rename this to "punc" = "pun cast"

	// parenthesized expression
	Paren

	// compound operands

	Chain

	// Select chain on expression via construct:
	//
	//	exp.name
	//
	Select

	Index
	Call

	Indirect
	Address
	Slice
	Test
	IndirectIndex

	// Struct field access via select chain.
	// Analogous to "." operator in C.
	Field

	// Struct field access with poiter dereference via select chain.
	// Analogous to "->" operator in C.
	IndirectField

	BoundMethod

	ChunkMember
	ChunkIndex
	ChunkSlice
	ArrayIndex
	ArraySlice

	Enum

	// Integer literal.
	Integer

	// String literal.
	String

	// True literal.
	True

	// False literal.
	False

	// Dirty literal.
	Dirty

	// Desugared strings equality comparison.
	StringsEqual

	// Desugared strings inequality comparison.
	StringsNotEqual

	// Desugared empty string check.
	StringEmpty

	// Desugared not empty string check.
	StringNotEmpty
)

var text = [...]string{
	empty: "<nil>",

	Basic:   "basic",
	Object:  "object",
	List:    "list",
	Symbol:  "symbol",
	Unary:   "unary",
	Binary:  "binary",
	Cast:    "cast",
	Tint:    "tint",
	MemCast: "mcast",
	Paren:   "paren",

	Chain:         "chain",
	Field:         "field",
	Indirect:      "indirect",
	Index:         "index",
	Test:          "member.test",
	IndirectIndex: "index.indirect",
	ChunkIndex:    "index.chunk",
	ChunkSlice:    "slice.chunk",
	ArrayIndex:    "index.array",
	ArraySlice:    "slice.array",
	IndirectField: "field.indirect",
	BoundMethod:   "method.bound",
	ChunkMember:   "member.chunk",
	Address:       "address",
	Slice:         "slice",
	Call:          "call",

	IncompName: "name.incomp",

	Enum: "enum",

	Integer: "integer",
	String:  "string",
	True:    "true",
	False:   "false",
	Dirty:   "dirty",

	StringsEqual:    "eq.str",
	StringsNotEqual: "neq.str",

	StringEmpty:    "empty.str",
	StringNotEmpty: "notempty.str",
}

func (k Kind) String() string {
	return text[k]
}

func (k Kind) IsPrimary() bool {
	return k != Binary
}

func (k Kind) IsOperand() bool {
	return k != Unary && k != Binary
}
