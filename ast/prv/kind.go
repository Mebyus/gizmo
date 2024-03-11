package prv

// Kind indicates prop value kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left unspecified
	empty Kind = iota

	Tag

	Integer
	String
	Bool
)

var text = [...]string{
	empty: "<nil>",

	Tag:     "tag",
	Integer: "integer",
	String:  "string",
	Bool:    "bool",
}

func (k Kind) String() string {
	return text[k]
}
