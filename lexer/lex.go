package lexer

import "github.com/mebyus/gizmo/token"

func (lx *Lexer) Lex() token.Token {
	tok := lx.lex()
	lx.pos.Num += 1
	return tok
}

func (lx *Lexer) lex() token.Token {
	if lx.eof {
		return lx.create(token.EOF)
	}

	lx.skipWhitespaceAndComments()
	if lx.eof {
		return lx.create(token.EOF)
	}

	return lx.codeToken()
}

func (lx *Lexer) codeToken() token.Token {
	if isLetterOrUnderscore(lx.c) {
		return lx.lexName()
	}

	if isDecimalDigit(lx.c) {
		return lx.lexNumber()
	}

	if lx.c == '"' {
		return lx.lexStringLiteral()
	}

	if lx.c == '\'' {
		return lx.runeLiteral()
	}

	if lx.c == '@' && lx.next == '.' {
		return lx.label()
	}

	return lx.lexOther()
}

// Create token (without literal) of specified kind at current lexer position
//
// Does not advance lexer scan position
func (lx *Lexer) create(k token.Kind) token.Token {
	return token.Token{
		Kind: k,
		Pos:  lx.pos,
	}
}

func (lx *Lexer) label() (tok token.Token) {
	tok.Pos = lx.pos

	lx.advance() // skip '@'
	lx.advance() // skip '.'

	lx.start()
	lx.skipWord()
	lit, ok := lx.take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}

	switch lit {
	case "next":
		tok.Kind = token.LabelNext
	case "end":
		tok.Kind = token.LabelEnd
	default:
		tok.Lit = lit
		panic("arbitrary labels not implemented: " + lit)
	}
	return
}

func (lx *Lexer) lexName() (tok token.Token) {
	tok.Pos = lx.pos

	lx.start()
	lx.skipWord()
	lit, ok := lx.take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}

	kind, ok := token.Lookup(lit)
	if ok {
		tok.Kind = kind
		return
	}

	tok.Kind = token.Identifier
	tok.Lit = lit
	return
}

func (lx *Lexer) lexBinaryNumber() (tok token.Token) {
	tok.Pos = lx.pos

	lx.advance() // skip '0'
	lx.advance() // skip 'b'

	lx.start()
	lx.skipBinaryDigits()

	if isAlphanum(lx.c) {
		lx.skipWord()
		lit, ok := lx.take()
		if ok {
			tok.SetIllegalError(token.MalformedBinaryInteger)
			tok.Lit = lit
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return
	}

	if lx.isLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}
	if lx.len() == 0 {
		tok.SetIllegalError(token.MalformedBinaryInteger)
		tok.Lit = "0b"
		return
	}

	tok.Kind = token.BinaryInteger
	if lx.len() > 64 {
		lit, ok := lx.take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Lit = lit
		return
	}

	tok.Val = parseBinaryDigits(lx.view())
	return
}

func (lx *Lexer) lexOctalNumber() (tok token.Token) {
	tok.Pos = lx.pos

	lx.advance() // skip '0' byte
	lx.advance() // skip 'o' byte

	lx.start()
	lx.skipOctalDigits()

	if isAlphanum(lx.c) {
		lx.skipWord()
		lit, ok := lx.take()
		if ok {
			tok.SetIllegalError(token.MalformedOctalInteger)
			tok.Lit = lit
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return
	}

	if lx.isLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}
	if lx.len() == 0 {
		tok.SetIllegalError(token.MalformedOctalInteger)
		tok.Lit = "0o"
		return
	}

	tok.Kind = token.OctalInteger
	if lx.len() > 21 {
		lit, ok := lx.take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Lit = lit
		return
	}

	tok.Val = parseOctalDigits(lx.view())
	return
}

func (lx *Lexer) lexDecimalNumber() (tok token.Token) {
	tok.Pos = lx.pos

	lx.start()
	scannedOnePeriod := false
	for !lx.eof && isDecimalDigitOrPeriod(lx.c) {
		lx.advance()
		if lx.c == '.' {
			if scannedOnePeriod {
				break
			} else {
				scannedOnePeriod = true
			}
		}
	}

	if lx.isLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}

	if lx.prev == '.' {
		tok.SetIllegalError(token.MalformedDecimalInteger)
		tok.Lit, _ = lx.take()
		return
	}

	if isAlphanum(lx.c) {
		lx.skipWord()
		lit, ok := lx.take()
		if ok {
			tok.SetIllegalError(token.MalformedDecimalInteger)
			tok.Lit = lit
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return
	}

	if !scannedOnePeriod {
		// decimal integer
		// TODO: handle numbers which do not fit into 64 bits
		tok.Kind = token.DecimalInteger
		tok.Val = parseDecimalDigits(lx.view())
		return
	}

	tok.Kind = token.DecimalFloat
	tok.Lit, _ = lx.take()
	return
}

func (lx *Lexer) lexHexadecimalNumber() (tok token.Token) {
	tok.Pos = lx.pos

	lx.advance() // skip '0' byte
	lx.advance() // skip 'x' byte

	lx.start()
	lx.skipHexadecimalDigits()

	if isAlphanum(lx.c) {
		lx.skipWord()
		lit, ok := lx.take()
		if ok {
			tok.SetIllegalError(token.MalformedHexadecimalInteger)
			tok.Lit = lit
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return
	}

	if lx.isLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}
	if lx.len() == 0 {
		tok.SetIllegalError(token.MalformedHexadecimalInteger)
		tok.Lit = "0x"
		return
	}

	tok.Kind = token.HexadecimalInteger
	if lx.len() > 16 {
		lit, ok := lx.take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Lit = lit
		return
	}

	tok.Val = parseHexadecimalDigits(lx.view())
	return
}

func (s *Lexer) lexNumber() (tok token.Token) {
	if s.c != '0' {
		return s.lexDecimalNumber()
	}

	if s.next == 'b' {
		return s.lexBinaryNumber()
	}

	if s.next == 'o' {
		return s.lexOctalNumber()
	}

	if s.next == 'x' {
		return s.lexHexadecimalNumber()
	}

	if s.next == '.' {
		return s.lexDecimalNumber()
	}

	if isAlphanum(s.next) {
		return s.illegalWord(token.MalformedDecimalInteger)
	}

	tok = token.Token{
		Kind: token.DecimalInteger,
		Pos:  s.pos,
		Val:  0,
	}
	s.advance()
	return
}

func (lx *Lexer) lexStringLiteral() (tok token.Token) {
	tok.Pos = lx.pos

	lx.advance() // skip '"'

	if lx.c == '"' {
		// common case of empty string literal
		lx.advance()
		tok.Kind = token.String
		return
	}

	lx.start()
	for !lx.eof && lx.c != '"' && lx.c != '\n' {
		if lx.c == '\\' && lx.next == '"' {
			// do not stop if we encounter escape sequence
			lx.advance() // skip "\"
			lx.advance() // skip quote
		} else {
			lx.advance()
		}
	}

	if lx.c != '"' {
		lit, ok := lx.take()
		if ok {
			tok.SetIllegalError(token.MalformedString)
			tok.Lit = lit
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return
	}

	lit, ok := lx.take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}

	lx.advance() // skip quote

	size, ok := token.ScanStringByteSize(lit)
	if !ok {
		tok.SetIllegalError(token.BadEscapeInString)
		return
	}

	tok.Lit = lit
	tok.Kind = token.String
	tok.Val = size
	return
}

func (lx *Lexer) skipWhitespaceAndComments() {
	for {
		lx.skipWhitespace()
		if lx.c == '/' && lx.next == '/' {
			lx.skipLineComment()
		} else if lx.c == '/' && lx.next == '*' {
			lx.skipMultilineComment()
		} else {
			return
		}
	}
}

func (lx *Lexer) skipLineComment() {
	lx.advance() // skip '/'
	lx.advance() // skip '/'
	lx.skipLine()
}

func (lx *Lexer) skipMultilineComment() {
	lx.advance() // skip '/'
	lx.advance() // skip '*'

	for !lx.eof && !(lx.c == '*' && lx.next == '/') {
		lx.advance()
	}

	if lx.eof {
		return
	}

	lx.advance() // skip '*'
	lx.advance() // skip '/'
}

func (lx *Lexer) runeLiteral() (tok token.Token) {
	tok.Pos = lx.pos

	lx.advance() // skip "'"

	lx.start()
	if lx.c == '\\' {
		// handle escape sequence
		var val uint64
		switch lx.next {
		case '\\':
			val = '\\'
		case 'n':
			val = '\n'
		case 't':
			val = '\t'
		case 'r':
			val = '\r'
		case '\'':
			val = '\''
		default:
			lx.advance() // skip "\"
			lx.advance() // skip unknown escape rune

			tok.SetIllegalError(token.MalformedRune)
			tok.Lit, _ = lx.take()
			if lx.c == '\'' {
				lx.advance()
			}
			return
		}

		lx.advance() // skip "\"
		lx.advance() // skip escape rune
		if lx.c != '\'' {
			tok.SetIllegalError(token.MalformedRune)
			tok.Lit, _ = lx.take()
			return
		}
		lx.advance() // skip "'"

		tok.Kind = token.Rune
		tok.Val = val
		return
	}

	if lx.next == '\'' {
		// common case of ascii rune
		tok.Val = uint64(lx.c)
		tok.Kind = token.Rune
		lx.advance()
		lx.advance()
		return
	}

	// handle non-ascii runes
	for !lx.eof && lx.c != '\'' && lx.c != '\n' {
		lx.advance()
	}

	lit, ok := lx.take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}

	if lx.c != '\'' {
		tok.SetIllegalError(token.MalformedRune)
		tok.Lit = lit
		return
	}
	lx.advance() // skip "'"

	tok.Kind = token.Rune
	tok.Lit = lit // TODO: parse rune val
	return
}

func (s *Lexer) scanGreaterStart() (tok token.Token) {
	if s.next == '=' {
		tok = s.create(token.GreaterOrEqual)
		s.advance()
		s.advance()
	} else {
		tok = s.create(token.RightAngle)
		s.advance()
	}
	return
}

func (lx *Lexer) lexByte(k token.Kind) token.Token {
	tok := lx.create(k)
	lx.advance()
	return tok
}

func (lx *Lexer) lexTwoBytes(k token.Kind) token.Token {
	tok := lx.create(k)
	lx.advance()
	lx.advance()
	return tok
}

func (lx *Lexer) illegalWord(code uint64) (tok token.Token) {
	tok.Pos = lx.pos

	lx.start()
	lx.skipWord()
	lit, ok := lx.take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}

	tok.SetIllegalError(code)
	tok.Lit = lit
	return
}

func (lx *Lexer) lexIllegalByte() (tok token.Token) {
	tok.Pos = lx.pos
	tok.Kind = token.Illegal
	tok.Lit = toString(byte(lx.c))
	lx.advance()
	return
}

// scanOther scans next operator, punctucator or illegal byte Token
func (lx *Lexer) lexOther() token.Token {
	switch lx.c {
	case '(':
		return lx.lexByte(token.LeftParentheses)
	case ')':
		return lx.lexByte(token.RightParentheses)
	case '{':
		return lx.lexByte(token.LeftCurly)
	case '}':
		return lx.lexByte(token.RightCurly)
	case '[':
		if lx.next == ']' {
			return lx.lexTwoBytes(token.Chunk)
		}
		if lx.next == '[' {
			return lx.lexTwoBytes(token.LeftDoubleSquare)
		}
		if lx.next == '_' {
			pos := lx.pos
			lx.advance() // skip "["
			if lx.next != ']' {
				lx.advance() // skip "_"
				return token.Token{
					Pos:  pos,
					Kind: token.Illegal,
					Lit:  "[_",
				}
			}
			lx.advance() // skip "_"
			lx.advance() // skip "]"
			return token.Token{
				Pos:  pos,
				Kind: token.AutoLen,
			}
		}
		if lx.next == '*' {
			pos := lx.pos
			lx.advance() // skip "["
			if lx.next != ']' {
				lx.advance() // skip "*"
				return token.Token{
					Pos:  pos,
					Kind: token.Illegal,
					Lit:  "[*",
				}
			}
			lx.advance() // skip "*"
			lx.advance() // skip "]"
			return token.Token{
				Pos:  pos,
				Kind: token.ArrayPointer,
			}
		}
		return lx.lexByte(token.LeftSquare)
	case ']':
		if lx.next == ']' {
			return lx.lexTwoBytes(token.RightDoubleSquare)
		}
		return lx.lexByte(token.RightSquare)
	case '<':
		switch lx.next {
		case '=':
			return lx.lexTwoBytes(token.LessOrEqual)
		case '<':
			return lx.lexTwoBytes(token.LeftShift)
		case '-':
			return lx.lexTwoBytes(token.LeftArrow)
		default:
			return lx.lexByte(token.LeftAngle)
		}
	case '>':
		return lx.scanGreaterStart()
	case '+':
		if lx.next == '=' {
			return lx.lexTwoBytes(token.AddAssign)
		}
		return lx.lexByte(token.Plus)
	case '-':
		if lx.next == '=' {
			return lx.lexTwoBytes(token.SubtractAssign)
		}
		return lx.lexByte(token.Minus)
	case ',':
		return lx.lexByte(token.Comma)
	case '=':
		switch lx.next {
		case '=':
			return lx.lexTwoBytes(token.Equal)
		case '>':
			return lx.lexTwoBytes(token.RightArrow)
		default:
			return lx.lexByte(token.Assign)
		}
	case ':':
		if lx.next == '=' {
			return lx.lexTwoBytes(token.ShortAssign)
		}
		if lx.next == ':' {
			return lx.lexTwoBytes(token.DoubleColon)
		}
		return lx.lexByte(token.Colon)
	case ';':
		return lx.lexByte(token.Semicolon)
	case '.':
		switch lx.next {
		case '&':
			return lx.lexTwoBytes(token.Address)
		case '@':
			return lx.lexTwoBytes(token.Indirect)
		case '{':
			return lx.lexTwoBytes(token.Compound)
		case '[':
			return lx.lexTwoBytes(token.IndirectIndex)
		case '!':
			return lx.lexTwoBytes(token.Insist)
		case '?':
			return lx.lexTwoBytes(token.Chain)
		default:
			return lx.lexByte(token.Period)
		}
	case '%':
		return lx.lexByte(token.Percent)
	case '*':
		return lx.lexByte(token.Asterisk)
	case '&':
		if lx.next == '&' {
			return lx.lexTwoBytes(token.LogicalAnd)
		}
		return lx.lexByte(token.Ampersand)
	case '/':
		return lx.lexByte(token.Slash)
	case '!':
		if lx.next == '=' {
			return lx.lexTwoBytes(token.NotEqual)
		}
		return lx.lexByte(token.Not)
	case '?':
		return lx.lexByte(token.Quest)
	case '^':
		return lx.lexByte(token.Caret)
	case '|':
		if lx.next == '|' {
			return lx.lexTwoBytes(token.LogicalOr)
		}
		return lx.lexByte(token.Pipe)
	case '#':
		if lx.next == '[' {
			return lx.lexTwoBytes(token.PropStart)
		}
		return lx.lexIllegalByte()
	default:
		return lx.lexIllegalByte()
	}
}
