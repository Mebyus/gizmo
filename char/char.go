package char

// utf-8 code point starts either from [0xxx xxxx] or [11xx xxxx].
// Non-starting byte of non-ascii code point has from [10xx xxxx].
// Thus we need to check that higher bits of a given byte are not 10.
func IsCodePointStart(c byte) bool {
	return (c & 0b1100_0000) != 0b1000_0000
}

func IsLetterOrUnderscore(c byte) bool {
	return IsLetter(c) || c == '_'
}

func IsAlphanum(c byte) bool {
	return IsLetterOrUnderscore(c) || IsDecDigit(c)
}

const capitalLetterMask = 0xDF

// ToUpperLetter transform latin letter to its upper (capital) form.
func ToUpperLetter(c byte) byte {
	return c & capitalLetterMask
}

func IsLetter(c byte) bool {
	c = ToUpperLetter(c)
	return 'A' <= c && c <= 'Z'
}

func IsDecDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func IsDecDigitOrPeriod(c byte) bool {
	return IsDecDigit(c) || c == '.'
}

func IsWhitespace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t' || c == '\r'
}

func ToString(b byte) string {
	return string([]byte{b})
}

func IsHexDigit(c byte) bool {
	if IsDecDigit(c) {
		return true
	}
	c = ToUpperLetter(c)
	return 'A' <= c && c <= 'F'
}

func IsOctDigit(c byte) bool {
	return '0' <= c && c <= '7'
}

func IsBinDigit(c byte) bool {
	return c == '0' || c == '1'
}

func DecDigitToNumber(d byte) uint8 {
	return d - '0'
}

// ParseBinDigits interprets ASCII digit characters as
// digits of binary number and returns the number.
//
// Does not validate the input.
// Input slice must satisfy 1 <= len(s) <= 64.
// As a special case returns zero for nil or empty input slice.
func ParseBinDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v <<= 1
		v += uint64(DecDigitToNumber(d))
	}
	return v
}

func ParseOctDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v <<= 3
		v += uint64(DecDigitToNumber(d))
	}
	return v
}

func ParseDecDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v *= 10
		v += uint64(DecDigitToNumber(d))
	}
	return v
}

func HexDigitToNumber(d byte) uint8 {
	if d <= '9' {
		return DecDigitToNumber(d)
	}
	return ToUpperLetter(d) - 'A' + 0x0A
}

func ParseHexDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v <<= 4
		v += uint64(HexDigitToNumber(d))
	}
	return v
}
