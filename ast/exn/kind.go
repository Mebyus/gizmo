package exn

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
	Receiver

	Unary
	Binary
	Cast
	BitCast

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
	IndirectIndex
	IndirectMember
	ChunkMember
	ChunkIndex

	// Template instance
	Instance

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

	Basic:    "basic",
	Object:   "object",
	List:     "list",
	Symbol:   "symbol",
	Unary:    "unary",
	Binary:   "binary",
	Cast:     "cast",
	BitCast:  "bitcast",
	Paren:    "paren",
	Receiver: "receiver",
	Instance: "instance",

	Chain:          "chain",
	Member:         "member",
	Indirect:       "indirect",
	Index:          "index",
	IndirectIndex:  "index.indirect",
	ChunkIndex:     "index.chunk",
	IndirectMember: "member.indirect",
	ChunkMember:    "member.chunk",
	Address:        "address",
	Slice:          "slice",
	Call:           "call",

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
