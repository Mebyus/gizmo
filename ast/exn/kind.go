package exn

// Kind indicates expression kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left inspecified
	empty Kind = iota

	// basic literal operand
	Basic

	// list literal operand
	List

	// direct substitution of identifier or scoped identifier operand
	Subs

	Unary
	Binary

	// parenthesized expression
	Paren

	// compound operands

	Select
	Index
	Call
	Start
)

var text = [...]string{
	empty: "<nil>",

	Basic:  "basic",
	List:   "list",
	Subs:   "subs",
	Unary:  "unary",
	Binary: "binary",
	Paren:  "paren",
	Select: "select",
	Index:  "index",
	Call:   "call",
	Start:  "start",
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
