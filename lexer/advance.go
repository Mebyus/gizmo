package lexer

const (
	eof      = -1
	nonASCII = 1 << 7

	maxLiteralLength = 1 << 8
)

func (lx *Lexer) advance() {
	if lx.c == eof {
		return
	}

	if lx.c == '\n' {
		lx.pos.Line += 1
		lx.pos.Col = 0
	} else if lx.c < nonASCII {
		lx.pos.Col += 1
	}

	lx.prev = lx.c
	lx.c = lx.next
	lx.pos.Ofs += 1
	if int(lx.i) < len(lx.src) {
		lx.next = int(lx.src[lx.i])
		lx.i += 1
	} else {
		lx.next = eof
	}
}

func (lx *Lexer) isEOF() bool {
	return lx.c == eof
}

func (lx *Lexer) take() string {
	str := string(lx.buf[:lx.len])
	lx.drop()
	return str
}

func (lx *Lexer) view() []byte {
	return lx.buf[:lx.len]
}

func (lx *Lexer) drop() {
	lx.len = 0
}

func (lx *Lexer) store() {
	lx.add(byte(lx.c))
	lx.advance()
}

func (lx *Lexer) add(b byte) {
	lx.buf[lx.len] = b
	lx.len++
}

func (lx *Lexer) skipWhitespace() {
	for isWhitespace(lx.c) {
		lx.advance()
	}
}

func (lx *Lexer) skipLine() {
	for !lx.isEOF() && lx.c != '\n' {
		lx.advance()
	}
	if lx.c == '\n' {
		lx.advance()
	}
}

func (lx *Lexer) storeWord() (overflow bool) {
	for lx.len < maxLiteralLength && isAlphanum(lx.c) {
		lx.store()
	}
	return isAlphanum(lx.c)
}

func (lx *Lexer) storeBinaryDigits() (overflow bool) {
	for lx.len < maxLiteralLength && isBinaryDigit(lx.c) {
		lx.store()
	}
	return isBinaryDigit(lx.c)
}

func (lx *Lexer) storeOctalDigits() (overflow bool) {
	for lx.len < maxLiteralLength && isOctalDigit(lx.c) {
		lx.store()
	}
	return isOctalDigit(lx.c)
}

func (lx *Lexer) storeHexadecimalDigits() (overflow bool) {
	for lx.len < maxLiteralLength && isHexadecimalDigit(lx.c) {
		lx.store()
	}
	return isHexadecimalDigit(lx.c)
}

func (lx *Lexer) skipWord() {
	for isAlphanum(lx.c) {
		lx.advance()
	}
}

func (lx *Lexer) skipBinaryDigits() {
	for isBinaryDigit(lx.c) {
		lx.advance()
	}
}

func (lx *Lexer) skipOctalDigits() {
	for isOctalDigit(lx.c) {
		lx.advance()
	}
}

func (lx *Lexer) skipHexadecimalDigits() {
	for isHexadecimalDigit(lx.c) {
		lx.advance()
	}
}
