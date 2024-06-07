package tps

// Kind indicates type specifier kind.
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	Name
	Struct
	Pointer
	Array
	ArrayPointer
	Chunk
	Union
	Enum
	Bag
	Instance
	Function
	Tuple
)

var text = [...]string{
	empty: "<nil>",

	Name:         "name",
	Struct:       "struct",
	Pointer:      "pointer",
	Array:        "array",
	ArrayPointer: "array_pointer",
	Chunk:        "chunk",
	Union:        "union",
	Enum:         "enum",
	Bag:          "bag",
	Instance:     "instance",
	Function:     "function",
	Tuple:        "tuple",
}

func (k Kind) String() string {
	return text[k]
}
