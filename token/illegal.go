package token

const (
	IdentifierOverflow = iota + 1
	MalformedIntergerOverflow
	BinaryIntegerOverflow
	OctalIntegerOverflow
	HexadecimalIntegerOverflow
	StringOverflow
	IllegalOverflow
)
