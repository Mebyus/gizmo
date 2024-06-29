package char

// Chopper is a text processing primitive providing low-level
// methods for splitting text into tokens in a sequential manner.
type Chopper struct {
	// Source text which is being scanned.
	text []byte

	// Mark index
	//
	// Mark is used to slice input text for token literals.
	mark int

	// next byte read index
	i int

	// Scanning index of source text (c = src[Pos])
	Pos int

	// Setting which should be configured before chopper usage.
	MaxTokenByteLength int

	// Zero-based line number of current position.
	//
	// This field is read-only outside of this package.
	// It should only be changed via Init() and Advance() methods.
	Line uint32

	// Zero-based column number (in utf-8 codepoints) of current position.
	//
	// This field is read-only outside of this package.
	// It should only be changed via Init() and Advance() methods.
	Col uint32

	// Byte that was previously at scan position.
	Prev byte

	// Byte at current scan position.
	//
	// This is a cached value. Look into advance() method
	// for details about how this caching algorithm works.
	C byte

	// Next byte that will be placed at scan position
	//
	// This is a cached value from previous read.
	Next byte

	// True if chopper reached end of input (by scan index).
	EOF bool
}

// Advance forward chopper scan position by one byte
// (with respect of input read caching)
// if there is available input ahead, do nothing otherwise.
func (c *Chopper) Advance() {
	if c.EOF {
		return
	}

	if c.C == '\n' {
		c.Line += 1
		c.Col = 0
	} else if IsCodePointStart(c.C) {
		c.Col += 1
	}

	c.Prev = c.C
	c.C = c.Next
	c.Pos += 1

	if c.i >= len(c.text) {
		c.Next = 0
		c.EOF = c.Pos >= len(c.text)
		return
	}

	c.Next = c.text[c.i]
	c.i += 1
}

const defaultMaxTokenByteLength = 1 << 10

// Init sets default settings and internal state.
// Must be used before any other method.
func (c *Chopper) Init(text []byte) {
	if c.MaxTokenByteLength < 0 {
		panic("negative value")
	}
	if c.MaxTokenByteLength == 0 {
		c.MaxTokenByteLength = defaultMaxTokenByteLength
	}

	c.text = text

	// prefill next and current bytes
	c.Advance()
	c.Advance()

	// adjust internal state to place scan position at input start
	c.Pos = 0
	c.Line = 0
	c.Col = 0
	c.EOF = len(c.text) == 0
}

// Start recording new token literal.
func (c *Chopper) Start() {
	c.mark = c.Pos
}

func (c *Chopper) View() []byte {
	return c.text[c.mark:c.Pos]
}

func (c *Chopper) IsLengthOverflow() bool {
	return c.Len() > c.MaxTokenByteLength
}

// Returns recorded token literal string (if ok == true) and ok flag, if ok == false that
// means token overflowed max token length.
func (c *Chopper) Take() (string, bool) {
	if c.IsLengthOverflow() {
		return "", false
	}
	return string(c.View()), true
}

// Returns length of recorded token literal.
func (c *Chopper) Len() int {
	return c.Pos - c.mark
}

func (c *Chopper) SkipWhitespace() {
	for IsWhitespace(c.C) {
		c.Advance()
	}
}

func (c *Chopper) SkipLine() {
	for !c.EOF && c.C != '\n' {
		c.Advance()
	}
	if c.C == '\n' {
		c.Advance()
	}
}

func (c *Chopper) SkipWord() {
	for IsAlphanum(c.C) {
		c.Advance()
	}
}

func (c *Chopper) SkipDecDigits() {
	for IsDecDigit(c.C) {
		c.Advance()
	}
}

func (c *Chopper) SkipBinDigits() {
	for IsBinDigit(c.C) {
		c.Advance()
	}
}

func (c *Chopper) SkipOctDigits() {
	for IsOctDigit(c.C) {
		c.Advance()
	}
}

func (c *Chopper) SkipHexDigits() {
	for IsHexDigit(c.C) {
		c.Advance()
	}
}

func LastByte(s []byte) byte {
	return s[len(s)-1]
}

// Same as method take, but removes possible whitespace suffix from the resulting string.
func (c *Chopper) TrimWhitespaceSuffixTake() (string, bool) {
	// TODO: transform implementation to pure function
	view := c.View()
	if len(view) == 0 {
		return "", true
	}
	if !IsWhitespace(LastByte(view)) {
		if len(view) > c.MaxTokenByteLength {
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
		b := view[i]

		if !IsWhitespace(b) {
			length := i + 1
			if length > c.MaxTokenByteLength {
				return "", false
			}
			return string(view[:length]), true
		}
	}
}

func (c *Chopper) SkipWhitespaceAndComments() {
	for {
		c.SkipWhitespace()
		if c.C == '/' && c.Next == '/' {
			c.SkipLineComment()
		} else if c.C == '/' && c.Next == '*' {
			c.SkipBlockComment()
		} else {
			return
		}
	}
}

func (c *Chopper) SkipLineComment() {
	c.Advance() // skip '/'
	c.Advance() // skip '/'
	c.SkipLine()
}

func (c *Chopper) SkipBlockComment() {
	c.Advance() // skip '/'
	c.Advance() // skip '*'

	for !c.EOF && !(c.C == '*' && c.Next == '/') {
		c.Advance()
	}

	if c.EOF {
		return
	}

	c.Advance() // skip '*'
	c.Advance() // skip '/'
}
