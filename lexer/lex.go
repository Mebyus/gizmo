package lexer

import "github.com/mebyus/gizmo/token"

func (lx *Lexer) Lex() token.Token {
	if lx.isEOF() {
		return lx.create(token.EOF)
	}

	lx.skipWhitespaceAndComments()
	if lx.isEOF() {
		return lx.create(token.EOF)
	}

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
		return lx.lexCharacterLiteral()
	}

	return lx.lexOther()
}

func (lx *Lexer) create(k token.Kind) token.Token {
	return token.Token{
		Kind: k,
		Pos:  lx.pos,
	}
}

func (lx *Lexer) lexName() (tok token.Token) {
	tok.Pos = lx.pos

	overflow := lx.storeWord()
	if overflow {
		lx.drop()
		lx.skipWord()
		tok.Kind = token.Illegal
		tok.Val = token.IdentifierOverflow
		return
	}
	lit := lx.take()

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

	overflow := lx.storeBinaryDigits()
	if overflow {
		lx.drop()
		lx.skipBinaryDigits()
		tok.Kind = token.Illegal
		tok.Val = token.BinaryIntegerOverflow
		return
	}

	if isAlphanum(lx.c) {
		overflow := lx.storeWord()
		if overflow {
			lx.drop()
			lx.skipWord()
			tok.Val = token.MalformedIntergerOverflow
		} else {
			tok.Lit = lx.take()
		}
		tok.Kind = token.Illegal
		return
	}

	if lx.len == 0 {
		tok.Kind = token.Illegal
		tok.Lit = "0b"
		return
	}

	// TODO: handle integers which cannot be stored in 64 bits
	tok.Kind = token.BinaryInteger
	tok.Val = parseBinaryDigits(lx.view())
	lx.drop()
	return
}

func (lx *Lexer) lexOctalNumber() (tok token.Token) {
	tok.Pos = lx.pos

	lx.advance() // skip '0' byte
	lx.advance() // skip 'o' byte

	overflow := lx.storeOctalDigits()
	if overflow {
		lx.drop()
		lx.skipOctalDigits()
		tok.Kind = token.Illegal
		tok.Val = token.OctalIntegerOverflow
		return
	}

	if isAlphanum(lx.c) {
		overflow := lx.storeWord()
		if overflow {
			lx.drop()
			lx.skipWord()
			tok.Val = token.MalformedIntergerOverflow
		} else {
			tok.Lit = lx.take()
		}
		tok.Kind = token.Illegal
		return
	}

	if lx.len == 0 {
		tok.Kind = token.Illegal
		tok.Lit = "0o"
		return
	}

	// TODO: handle integers which cannot be stored in 64 bits
	tok.Kind = token.OctalInteger
	tok.Val = parseOctalDigits(lx.view())
	lx.drop()
	return
}

func (lx *Lexer) lexDecimalNumber() (tok token.Token) {
	tok.Pos = lx.pos

	scannedOnePeriod := false
	for !lx.isEOF() && isDecimalDigitOrPeriod(lx.c) {
		lx.store()
		if lx.c == '.' {
			if scannedOnePeriod {
				break
			} else {
				scannedOnePeriod = true
			}
		}
	}
	// TODO: handle overflow

	// TODO: handle trailing period

	if isAlphanum(lx.c) {
		overflow := lx.storeWord()
		if overflow {
			lx.drop()
			lx.skipWord()
			tok.Val = token.MalformedIntergerOverflow
		} else {
			tok.Lit = lx.take()
		}
		tok.Kind = token.Illegal
		return
	}

	if lx.prev == '.' {
		tok.Kind = token.Illegal
		tok.Lit = lx.take()
		return
	}

	if !scannedOnePeriod {
		// TODO: handle numbers which do not fit into 64 bits
		tok.Kind = token.DecimalInteger
		tok.Val = parseDecimalDigits(lx.view())
		lx.drop()
		return
	}

	tok.Kind = token.DecimalFloat
	tok.Lit = lx.take()
	return
}

func (lx *Lexer) lexHexadecimalNumber() (tok token.Token) {
	tok.Pos = lx.pos

	lx.advance() // skip '0' byte
	lx.advance() // skip 'x' byte

	overflow := lx.storeHexadecimalDigits()
	if overflow {
		lx.drop()
		lx.skipHexadecimalDigits()
		tok.Kind = token.Illegal
		tok.Val = token.HexadecimalIntegerOverflow
		return
	}

	if isAlphanum(lx.c) {
		overflow := lx.storeWord()
		if overflow {
			lx.drop()
			lx.skipWord()
			tok.Val = token.MalformedIntergerOverflow
		} else {
			tok.Lit = lx.take()
		}
		tok.Kind = token.Illegal
		return
	}

	if lx.len == 0 {
		tok.Kind = token.Illegal
		tok.Lit = "0x"
		return
	}

	// TODO: handle numbers which do not fit into 64 bits
	tok.Kind = token.HexadecimalInteger
	tok.Val = parseHexadecimalDigits(lx.view())
	lx.drop()
	return
}

func (s *Lexer) lexNumber() (tok token.Token) {
	if s.c != '0' {
		return s.lexDecimalNumber()
	}

	if s.next == eof {
		tok = token.Token{
			Kind: token.DecimalInteger,
			Pos:  s.pos,
			Val:  0,
		}
		s.advance()
		return
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
		return s.lexIllegalWord()
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

	isIllegal := false // TODO: handle overflow
	for lx.len < maxLiteralLength && lx.c != eof && lx.c != '"' {
		if lx.c != '\\' || isIllegal {
			lx.store()
			continue
		}

		// handle escape sequence
		switch lx.next {
		case 'n':
			lx.add('\n')
		case 'r':
			lx.add('\r')
		case 't':
			lx.add('\t')
		case '"':
			lx.add('"')
		case '\\':
			lx.add('\\')
		case eof:
			isIllegal = true
			lx.store()
			continue
		default:
			isIllegal = true
			lx.store()
			lx.store()
			continue
		}
		lx.advance()
		lx.advance()
	}

	if lx.c == eof || isIllegal {
		tok.Kind = token.Illegal
		tok.Lit = lx.take()
		return
	}

	if lx.c != '"' {
		lx.drop()
		lx.skipLine()
		tok.Val = token.StringOverflow
		return
	}

	lx.advance() // consume "
	tok.Lit = lx.take()
	tok.Kind = token.String
	return
}

func (lx *Lexer) skipWhitespaceAndComments() {
	for {
		lx.skipWhitespace()
		if lx.c == '/' && lx.next == '/' {
			lx.skipLineComment()
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

func (lx *Lexer) lexCharacterLiteral() (tok token.Token) {
	tok.Pos = lx.pos

	lx.store() // consume "'"
	for !lx.isEOF() && lx.c != '\'' {
		lx.store()
	}
	// TODO: handle overflow or restrict size in bytes to 4

	if lx.c != '\'' {
		tok.Kind = token.Illegal
		tok.Lit = lx.take()
		return
	}

	lx.advance()         // skip "'"
	lit := lx.view()[1:] // this will contain only bytes between two ticks
	tok.Kind = token.Character
	tok.Lit = string(lit)
	lx.drop()
	return
}

func (s *Lexer) scanGreaterStart() (tok token.Token) {
	if s.next == '=' {
		tok = s.create(token.GreaterOrEqual)
		s.advance()
		s.advance()
	} else {
		tok = s.create(token.Greater)
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

func (lx *Lexer) lexIllegalWord() (tok token.Token) {
	tok.Pos = lx.pos
	overflow := lx.storeWord()
	if overflow {
		lx.drop()
		lx.skipWord()
		tok.Val = token.IllegalOverflow
	} else {
		tok.Lit = lx.take()
	}
	tok.Kind = token.Illegal
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
			return lx.lexByte(token.Less)
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
	default:
		return lx.lexIllegalByte()
	}
}
