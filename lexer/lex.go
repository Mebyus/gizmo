package lexer

import (
	"github.com/mebyus/gizmo/char"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

func (lx *Lexer) Lex() token.Token {
	tok := lx.lex()
	lx.num += 1
	return tok
}

func (lx *Lexer) pos() source.Pos {
	return source.Pos{
		File: lx.file,
		Ofs:  uint32(lx.Pos),
		Num:  lx.num,
		Line: lx.Line,
		Col:  lx.Col,
	}
}

func (lx *Lexer) lex() token.Token {
	if lx.EOF {
		return lx.create(token.EOF)
	}

	lx.skipWhitespaceAndComments()
	if lx.EOF {
		return lx.create(token.EOF)
	}

	return lx.codeToken()
}

func (lx *Lexer) codeToken() token.Token {
	if char.IsLetterOrUnderscore(lx.C) {
		return lx.lexName()
	}

	if char.IsDecDigit(lx.C) {
		return lx.lexNumber()
	}

	if lx.C == '"' {
		return lx.lexStringLiteral()
	}

	if lx.C == '\'' {
		return lx.runeLiteral()
	}

	if lx.C == '@' && lx.Next == '.' {
		return lx.label()
	}

	return lx.lexOther()
}

// Create token (without literal) of specified kind at current lexer position
//
// Does not Advance lexer scan position
func (lx *Lexer) create(k token.Kind) token.Token {
	return token.Token{
		Kind: k,
		Pos:  lx.pos(),
	}
}

func (lx *Lexer) label() (tok token.Token) {
	tok.Pos = lx.pos()

	lx.Advance() // skip '@'
	lx.Advance() // skip '.'

	lx.Start()
	lx.SkipWord()
	lit, ok := lx.Take()
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
	tok.Pos = lx.pos()

	if !char.IsAlphanum(lx.Next) {
		// word is 1 character long
		c := lx.C
		lx.Advance() // skip character

		if c == 'g' {
			tok.Kind = token.Receiver
		} else if c == '_' {
			tok.Kind = token.Underscore
		} else {
			tok.Kind = token.Identifier
			tok.Lit = char.ToString(c)
		}
		return
	}

	// word is at least 2 characters long
	lx.Start()
	lx.SkipWord()
	lit, ok := lx.Take()
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
	tok.Pos = lx.pos()

	lx.Advance() // skip '0'
	lx.Advance() // skip 'b'

	lx.Start()
	lx.SkipBinDigits()

	if char.IsAlphanum(lx.C) {
		lx.SkipWord()
		lit, ok := lx.Take()
		if ok {
			tok.SetIllegalError(token.MalformedBinaryInteger)
			tok.Lit = lit
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}
	if lx.Len() == 0 {
		tok.SetIllegalError(token.MalformedBinaryInteger)
		tok.Lit = "0b"
		return
	}

	tok.Kind = token.BinaryInteger
	if lx.Len() > 64 {
		lit, ok := lx.Take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Lit = lit
		return
	}

	tok.Val = char.ParseBinDigits(lx.View())
	return
}

func (lx *Lexer) lexOctalNumber() (tok token.Token) {
	tok.Pos = lx.pos()

	lx.Advance() // skip '0' byte
	lx.Advance() // skip 'o' byte

	lx.Start()
	lx.SkipOctDigits()

	if char.IsAlphanum(lx.C) {
		lx.SkipWord()
		lit, ok := lx.Take()
		if ok {
			tok.SetIllegalError(token.MalformedOctalInteger)
			tok.Lit = lit
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}
	if lx.Len() == 0 {
		tok.SetIllegalError(token.MalformedOctalInteger)
		tok.Lit = "0o"
		return
	}

	tok.Kind = token.OctalInteger
	if lx.Len() > 21 {
		lit, ok := lx.Take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Lit = lit
		return
	}

	tok.Val = char.ParseOctDigits(lx.View())
	return
}

func (lx *Lexer) lexDecimalNumber() (tok token.Token) {
	tok.Pos = lx.pos()

	lx.Start()
	scannedOnePeriod := false
	for !lx.EOF && char.IsDecDigitOrPeriod(lx.C) {
		lx.Advance()
		if lx.C == '.' {
			if scannedOnePeriod {
				break
			} else {
				scannedOnePeriod = true
			}
		}
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}

	if lx.Prev == '.' {
		tok.SetIllegalError(token.MalformedDecimalInteger)
		tok.Lit, _ = lx.Take()
		return
	}

	if char.IsAlphanum(lx.C) {
		lx.SkipWord()
		lit, ok := lx.Take()
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
		tok.Val = char.ParseDecDigits(lx.View())
		return
	}

	tok.Kind = token.DecimalFloat
	tok.Lit, _ = lx.Take()
	return
}

func (lx *Lexer) lexHexadecimalNumber() (tok token.Token) {
	tok.Pos = lx.pos()

	lx.Advance() // skip '0' byte
	lx.Advance() // skip 'x' byte

	lx.Start()
	lx.SkipHexDigits()

	if char.IsAlphanum(lx.C) {
		lx.SkipWord()
		lit, ok := lx.Take()
		if ok {
			tok.SetIllegalError(token.MalformedHexadecimalInteger)
			tok.Lit = lit
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}
	if lx.Len() == 0 {
		tok.SetIllegalError(token.MalformedHexadecimalInteger)
		tok.Lit = "0x"
		return
	}

	tok.Kind = token.HexadecimalInteger
	if lx.Len() > 16 {
		lit, ok := lx.Take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Lit = lit
		return
	}

	tok.Val = char.ParseHexDigits(lx.View())
	return
}

func (lx *Lexer) lexNumber() (tok token.Token) {
	if lx.C != '0' {
		return lx.lexDecimalNumber()
	}

	if lx.Next == 'b' {
		return lx.lexBinaryNumber()
	}

	if lx.Next == 'o' {
		return lx.lexOctalNumber()
	}

	if lx.Next == 'x' {
		return lx.lexHexadecimalNumber()
	}

	if lx.Next == '.' {
		return lx.lexDecimalNumber()
	}

	if char.IsAlphanum(lx.Next) {
		return lx.illegalWord(token.MalformedDecimalInteger)
	}

	tok = token.Token{
		Kind: token.DecimalInteger,
		Pos:  lx.pos(),
		Val:  0,
	}
	lx.Advance()
	return
}

func (lx *Lexer) lexStringLiteral() (tok token.Token) {
	tok.Pos = lx.pos()

	lx.Advance() // skip '"'

	if lx.C == '"' {
		// common case of empty string literal
		lx.Advance()
		tok.Kind = token.String
		return
	}

	lx.Start()
	for !lx.EOF && lx.C != '"' && lx.C != '\n' {
		if lx.C == '\\' && lx.Next == '"' {
			// do not stop if we encounter escape sequence
			lx.Advance() // skip "\"
			lx.Advance() // skip quote
		} else {
			lx.Advance()
		}
	}

	if lx.C != '"' {
		lit, ok := lx.Take()
		if ok {
			tok.SetIllegalError(token.MalformedString)
			tok.Lit = lit
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return
	}

	lit, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}

	lx.Advance() // skip quote

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
		lx.SkipWhitespace()
		if lx.C == '/' && lx.Next == '/' {
			lx.skipLineComment()
		} else if lx.C == '/' && lx.Next == '*' {
			lx.skipMultilineComment()
		} else {
			return
		}
	}
}

func (lx *Lexer) skipLineComment() {
	lx.Advance() // skip '/'
	lx.Advance() // skip '/'
	lx.SkipLine()
}

func (lx *Lexer) skipMultilineComment() {
	lx.Advance() // skip '/'
	lx.Advance() // skip '*'

	for !lx.EOF && !(lx.C == '*' && lx.Next == '/') {
		lx.Advance()
	}

	if lx.EOF {
		return
	}

	lx.Advance() // skip '*'
	lx.Advance() // skip '/'
}

func (lx *Lexer) runeLiteral() (tok token.Token) {
	tok.Pos = lx.pos()

	lx.Advance() // skip "'"

	lx.Start()
	if lx.C == '\\' {
		// handle escape sequence
		var val uint64
		switch lx.Next {
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
			lx.Advance() // skip "\"
			lx.Advance() // skip unknown escape rune

			tok.SetIllegalError(token.MalformedRune)
			tok.Lit, _ = lx.Take()
			if lx.C == '\'' {
				lx.Advance()
			}
			return
		}

		lx.Advance() // skip "\"
		lx.Advance() // skip escape rune
		if lx.C != '\'' {
			tok.SetIllegalError(token.MalformedRune)
			tok.Lit, _ = lx.Take()
			return
		}
		lx.Advance() // skip "'"

		tok.Kind = token.Rune
		tok.Val = val
		return
	}

	if lx.Next == '\'' {
		// common case of ascii rune
		tok.Val = uint64(lx.C)
		tok.Kind = token.Rune
		lx.Advance()
		lx.Advance()
		return
	}

	// handle non-ascii runes
	for !lx.EOF && lx.C != '\'' && lx.C != '\n' {
		lx.Advance()
	}

	lit, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}

	if lx.C != '\'' {
		tok.SetIllegalError(token.MalformedRune)
		tok.Lit = lit
		return
	}
	lx.Advance() // skip "'"

	tok.Kind = token.Rune
	tok.Lit = lit // TODO: parse rune val
	return
}

func (s *Lexer) scanGreaterStart() (tok token.Token) {
	if s.Next == '=' {
		tok = s.create(token.GreaterOrEqual)
		s.Advance()
		s.Advance()
	} else {
		tok = s.create(token.RightAngle)
		s.Advance()
	}
	return
}

func (lx *Lexer) oneByteToken(k token.Kind) token.Token {
	tok := lx.create(k)
	lx.Advance()
	return tok
}

func (lx *Lexer) twoBytesToken(k token.Kind) token.Token {
	tok := lx.create(k)
	lx.Advance()
	lx.Advance()
	return tok
}

func (lx *Lexer) illegalWord(code uint64) (tok token.Token) {
	tok.Pos = lx.pos()

	lx.Start()
	lx.SkipWord()
	lit, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}

	tok.SetIllegalError(code)
	tok.Lit = lit
	return
}

func (lx *Lexer) illegalByteToken() (tok token.Token) {
	tok.Pos = lx.pos()
	tok.Kind = token.Illegal
	tok.Lit = char.ToString(byte(lx.C))
	lx.Advance()
	return
}

// scanOther scans next operator, punctucator or illegal byte Token
func (lx *Lexer) lexOther() token.Token {
	switch lx.C {
	case '(':
		return lx.oneByteToken(token.LeftParentheses)
	case ')':
		return lx.oneByteToken(token.RightParentheses)
	case '{':
		return lx.oneByteToken(token.LeftCurly)
	case '}':
		return lx.oneByteToken(token.RightCurly)
	case '[':
		if lx.Next == ']' {
			return lx.twoBytesToken(token.Chunk)
		}
		if lx.Next == '_' {
			pos := lx.pos()
			lx.Advance() // skip "["
			if lx.Next != ']' {
				lx.Advance() // skip "_"
				return token.Token{
					Pos:  pos,
					Kind: token.Illegal,
					Lit:  "[_",
				}
			}
			lx.Advance() // skip "_"
			lx.Advance() // skip "]"
			return token.Token{
				Pos:  pos,
				Kind: token.AutoLen,
			}
		}
		if lx.Next == '*' {
			pos := lx.pos()
			lx.Advance() // skip "["
			if lx.Next != ']' {
				lx.Advance() // skip "*"
				return token.Token{
					Pos:  pos,
					Kind: token.Illegal,
					Lit:  "[*",
				}
			}
			lx.Advance() // skip "*"
			lx.Advance() // skip "]"
			return token.Token{
				Pos:  pos,
				Kind: token.ArrayPointer,
			}
		}
		return lx.oneByteToken(token.LeftSquare)
	case ']':
		return lx.oneByteToken(token.RightSquare)
	case '<':
		switch lx.Next {
		case '=':
			return lx.twoBytesToken(token.LessOrEqual)
		case '<':
			return lx.twoBytesToken(token.LeftShift)
		case '-':
			return lx.twoBytesToken(token.LeftArrow)
		default:
			return lx.oneByteToken(token.LeftAngle)
		}
	case '>':
		return lx.scanGreaterStart()
	case '+':
		if lx.Next == '=' {
			return lx.twoBytesToken(token.AddAssign)
		}
		return lx.oneByteToken(token.Plus)
	case '-':
		if lx.Next == '=' {
			return lx.twoBytesToken(token.SubtractAssign)
		}
		return lx.oneByteToken(token.Minus)
	case ',':
		return lx.oneByteToken(token.Comma)
	case '=':
		switch lx.Next {
		case '=':
			return lx.twoBytesToken(token.Equal)
		case '>':
			return lx.twoBytesToken(token.RightArrow)
		default:
			return lx.oneByteToken(token.Assign)
		}
	case ':':
		if lx.Next == '=' {
			return lx.twoBytesToken(token.ShortAssign)
		}
		if lx.Next == ':' {
			return lx.twoBytesToken(token.DoubleColon)
		}
		return lx.oneByteToken(token.Colon)
	case ';':
		return lx.oneByteToken(token.Semicolon)
	case '.':
		switch lx.Next {
		case '&':
			return lx.twoBytesToken(token.Address)
		case '@':
			return lx.twoBytesToken(token.Indirect)
		case '{':
			return lx.twoBytesToken(token.Compound)
		case '[':
			return lx.twoBytesToken(token.IndirectIndex)
		case '!':
			return lx.twoBytesToken(token.Insist)
		case '?':
			return lx.twoBytesToken(token.Chain)
		default:
			return lx.oneByteToken(token.Period)
		}
	case '%':
		return lx.oneByteToken(token.Percent)
	case '*':
		return lx.oneByteToken(token.Asterisk)
	case '&':
		if lx.Next == '&' {
			return lx.twoBytesToken(token.LogicalAnd)
		}
		return lx.oneByteToken(token.Ampersand)
	case '/':
		return lx.oneByteToken(token.Slash)
	case '!':
		if lx.Next == '=' {
			return lx.twoBytesToken(token.NotEqual)
		}
		return lx.oneByteToken(token.Not)
	case '?':
		return lx.oneByteToken(token.Quest)
	case '^':
		return lx.oneByteToken(token.Caret)
	case '|':
		if lx.Next == '|' {
			return lx.twoBytesToken(token.LogicalOr)
		}
		return lx.oneByteToken(token.Pipe)
	case '#':
		if lx.Next == '[' {
			return lx.twoBytesToken(token.PropStart)
		}
		return lx.illegalByteToken()
	default:
		return lx.illegalByteToken()
	}
}
