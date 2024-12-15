package token

const (
	LengthOverflow = iota + 1
	NonPrintableByte

	MalformedString
	MalformedRune
	MalformedBinaryInteger
	MalformedOctalInteger
	MalformedDecimalInteger
	MalformedHexadecimalInteger
	MalformedBlockComment
	MalformedMacro

	BadEscapeInString

	UnknownDirective
)

func (t *Token) SetIllegalError(code uint64) {
	t.Kind = Illegal
	t.Val = code
}
