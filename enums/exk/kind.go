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
	Member
	Index
	Call
	Indirect
	Address
	Slice
	Test
	IndirectIndex
	IndirectMember
	ChunkMember
	ChunkIndex
	ChunkSlice
	ArrayIndex
	ArraySlice

	Enum

	// Integer literal
	Integer

	// String literal
	String

	// True literal
	True

	// False literal
	False
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

	Chain:          "chain",
	Member:         "member",
	Indirect:       "indirect",
	Index:          "index",
	IndirectIndex:  "index.indirect",
	ChunkIndex:     "index.chunk",
	ChunkSlice:     "slice.chunk",
	ArrayIndex:     "index.array",
	ArraySlice:     "slice.array",
	IndirectMember: "member.indirect",
	ChunkMember:    "member.chunk",
	Address:        "address",
	Slice:          "slice",
	Call:           "call",

	IncompName: "name.incomp",

	Enum: "enum",

	Integer: "integer",
	String:  "string",
	True:    "true",
	False:   "false",
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
