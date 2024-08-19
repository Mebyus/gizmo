package tps

// Kind indicates type specifier kind.
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	Name
	ImportName
	Struct
	Pointer
	RawMemoryPointer
	Array
	ArrayPointer
	Chunk
	Union
	Enum
	Bag
	Function
	Tuple
)

var text = [...]string{
	empty: "<nil>",

	Name:         "name",
	ImportName:   "name.import",
	Struct:       "struct",
	Pointer:      "pointer",
	Array:        "array",
	ArrayPointer: "pointer.array",
	Chunk:        "chunk",
	Union:        "union",
	Enum:         "enum",
	Bag:          "bag",
	Function:     "function",
	Tuple:        "tuple",

	RawMemoryPointer: "pointer.rawmem",
}

func (k Kind) String() string {
	return text[k]
}
