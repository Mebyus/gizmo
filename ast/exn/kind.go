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

	// direct substitution of identifier or scoped identifier operand
	Symbol

	Unary
	Binary
	Cast
	BitCast

	// parenthesized expression
	Paren

	// compound operands

	Select
	Index
	Call
	Start
	Receiver
	Indirect
	Address
	Slice

	// Indirect index expression
	Indirx

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
	Select:   "select",
	Index:    "index",
	Call:     "call",
	Start:    "start",
	Receiver: "receiver",
	Address:  "address",
	Slice:    "slice",
	Indirect: "indirect",
	Indirx:   "indirect_index",
	Instance: "instance",

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
