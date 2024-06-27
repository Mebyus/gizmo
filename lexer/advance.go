package lexer

import "github.com/mebyus/gizmo/char"

const (
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
	} else if char.IsCodePointStart(lx.c) {
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

func lastByte(s []byte) byte {
	return s[len(s)-1]
}

// same as method take, but removes possible whitespace suffix from the resulting string
func (lx *Lexer) trimWhitespaceSuffixTake() (string, bool) {
	view := lx.view()
	if len(view) == 0 {
		return "", true
	}
	if !char.IsWhitespace(lastByte(view)) {
		if len(view) > maxTokenByteLength {
			return "", false
		}
		return string(view), true
	}

	i := len(view)
	for {
		if i == 0 {
			return "", true
		}
		i -= 1
		c := view[i]

		if !char.IsWhitespace(c) {
			length := i + 1
			if length > maxTokenByteLength {
				return "", false
			}
			return string(view[:length]), true
		}
	}
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
	for char.IsWhitespace(lx.c) {
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
	for char.IsAlphanum(lx.c) {
		lx.advance()
	}
}

func (lx *Lexer) skipBinaryDigits() {
	for char.IsBinDigit(lx.c) {
		lx.advance()
	}
}

func (lx *Lexer) skipOctalDigits() {
	for char.IsOctDigit(lx.c) {
		lx.advance()
	}
}

func (lx *Lexer) skipHexadecimalDigits() {
	for char.IsHexDigit(lx.c) {
		lx.advance()
	}
}
