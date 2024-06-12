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

	// Simplified special case of general select expression,
	// when the target is an isolated symbol:
	//
	//	example_symbol.member
	//
	Member

	// Simplified special case of general address expression,
	// when the target is an isolated symbol:
	//
	//	example_symbol.&
	//
	SymbolAddress

	// Simplified special case of general indirect expression,
	// when the target is an isolated symbol:
	//
	//	example_symbol.@
	//
	SymbolIndirect

	// Simplified special case of general call expression,
	// when the callee is an isolated symbol:
	//
	//	example_symbol(arg1, arg2)
	//
	SymbolCall

	// Simplified special case of general call expression,
	// when the callee is an isolated symbol selector:
	//
	//	example_object.call(arg1, arg2)
	//
	MemberCall

	// Simplified special case of general call expression,
	// when the callee is an isolated receiver:
	//
	//	rv.call(arg1, arg2)
	//
	ReceiverMemberCall

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

	Member:             "member",
	SymbolCall:         "symbol_call",
	SymbolAddress:      "symbol_address",
	MemberCall:         "member_call",
	ReceiverMemberCall: "rvm_call",

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
