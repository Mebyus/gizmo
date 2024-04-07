package lexer

// utf-8 code point starts either from [0xxx xxxx] or [11xx xxxx].
// Non-starting byte of non-ascii code point has from [10xx xxxx].
// Thus we need to check that higher bits of a given byte are not 10 
func isCodePointStart(c byte) bool {
	return (c & 0b1100_0000) != 0b1000_0000
}

func isLetterOrUnderscore(c byte) bool {
	return isLetter(c) || c == '_'
}

func isAlphanum(c byte) bool {
	return isLetterOrUnderscore(c) || isDecimalDigit(c)
}

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func isDecimalDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isDecimalDigitOrPeriod(c byte) bool {
	return isDecimalDigit(c) || c == '.'
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t' || c == '\r'
}

func toString(b byte) string {
	return string([]byte{byte(b)})
}

func isHexadecimalDigit(c byte) bool {
	return isDecimalDigit(c) || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func isOctalDigit(c byte) bool {
	return '0' <= c && c <= '7'
}

func isBinaryDigit(c byte) bool {
	return c == '0' || c == '1'
}

func decDigitToNumber(d byte) uint8 {
	return d - '0'
}

func parseBinaryDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v <<= 1
		v += uint64(decDigitToNumber(d))
	}
	return v
}

func parseOctalDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v <<= 3
		v += uint64(decDigitToNumber(d))
	}
	return v
}

func parseDecimalDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v *= 10
		v += uint64(decDigitToNumber(d))
	}
	return v
}

func hexDigitToNumber(d byte) uint8 {
	if d <= '9' {
		return decDigitToNumber(d)
	}
	if d <= 'F' {
		return d - 'A' + 0x0A
	}
	return d - 'a' + 0x0A
}

func parseHexadecimalDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v <<= 4
		v += uint64(hexDigitToNumber(d))
	}
	return v
}
