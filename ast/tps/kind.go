package tps

// Kind indicates type specifier kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left inspecified
	empty Kind = iota

	Name
	Struct
	Pointer
	Array
	ArrayPointer
	Chunk
	Union
	Enum
	Instance
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
	Instance:     "instance",
}

func (k Kind) String() string {
	return text[k]
}
