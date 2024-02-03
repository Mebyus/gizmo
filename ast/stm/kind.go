package stm

// Kind indicates statement kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left inspecified
	empty Kind = iota

	Block
	Assign
	AddAssign
)

var text = [...]string{
	empty: "<nil>",

	Block:     "block",
	Assign:    "assign",
	AddAssign: "add_assign",
}

func (k Kind) String() string {
	return text[k]
}
