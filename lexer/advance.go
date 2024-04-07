package lexer

const (
	nonASCII = 1 << 7

	maxTokenByteLength = 1 << 10
)

// advance forward lexer scan position (with respect of input read caching)
// if there is available input ahead, do nothing otherwise
func (lx *Lexer) advance() {
	if lx.eof {
		return
	}

	if lx.c == '\n' {
		lx.pos.Line += 1
		lx.pos.Col = 0
	} else if isCodePointStart(lx.c) {
		lx.pos.Col += 1
	}

	lx.prev = lx.c
	lx.c = lx.next
	lx.s += 1
	lx.pos.Ofs += 1

	if lx.i >= len(lx.src) {
		lx.next = 0
		lx.eof = lx.s >= len(lx.src)
		return
	}

	lx.next = lx.src[lx.i]
	lx.i += 1
}

func (lx *Lexer) init() {
	// prefill next and current bytes
	lx.advance()
	lx.advance()

	// adjust internal state to place scan position at input start
	lx.s = 0
	lx.pos.Reset()
	lx.eof = len(lx.src) == 0
}

// start recording new token literal
func (lx *Lexer) start() {
	lx.mark = lx.s
}

// returns recorded token literal string (if ok == true) and ok flag, if ok == false that
// means token overflowed max token length
func (lx *Lexer) take() (string, bool) {
	if lx.isLengthOverflow() {
		return "", false
	}
	return string(lx.view()), true
}

func (lx *Lexer) view() []byte {
	return lx.src[lx.mark:lx.s]
}

func (lx *Lexer) isLengthOverflow() bool {
	return lx.len() > maxTokenByteLength
}

// returns length of recorded token literal
func (lx *Lexer) len() int {
	return lx.s - lx.mark
}

func (lx *Lexer) skipWhitespace() {
	for isWhitespace(lx.c) {
		lx.advance()
	}
}

func (lx *Lexer) skipLine() {
	for !lx.eof && lx.c != '\n' {
		lx.advance()
	}
	if lx.c == '\n' {
		lx.advance()
	}
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
