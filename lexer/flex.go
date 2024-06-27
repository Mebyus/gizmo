package lexer

import (
	"github.com/mebyus/gizmo/char"
	"github.com/mebyus/gizmo/token"
)

func (lx *Lexer) Flex() token.Token {
	tok := lx.flex()
	lx.pos.Num += 1
	return tok
}

func (lx *Lexer) flex() token.Token {
	if lx.eof {
		return lx.create(token.EOF)
	}

	lx.skipWhitespace()
	if lx.eof {
		return lx.create(token.EOF)
	}

	if lx.c == '/' && lx.next == '/' {
		return lx.lineComment()
	}
	if lx.c == '/' && lx.next == '*' {
		return lx.blockComment()
	}

	return lx.codeToken()
}

func (lx *Lexer) lineComment() (tok token.Token) {
	tok.Pos = lx.pos

	lx.advance() // skip '/'
	lx.advance() // skip '/'

	// skip until actual comment starts
	for !lx.eof && lx.c != '\n' && char.IsWhitespace(lx.c) {
		lx.advance()
	}

	if lx.eof {
		tok.Kind = token.LineComment
		return
	}
	if lx.c == '\n' {
		lx.skipWhitespace()
		tok.Kind = token.LineComment
		return
	}

	lx.start()
	for !lx.eof && lx.c != '\n' {
		lx.advance()
	}

	lit, ok := lx.trimWhitespaceSuffixTake()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}
	tok.Lit = lit
	tok.Kind = token.LineComment

	if lx.eof {
		return
	}

	lx.skipWhitespace()
	return
}

func (lx *Lexer) blockComment() (tok token.Token) {
	tok.Pos = lx.pos

	lx.advance() // skip '/'
	lx.advance() // skip '*'

	lx.start()
	for !lx.eof && !(lx.c == '*' && lx.next == '/') {
		lx.advance()
	}

	if lx.eof {
		tok.SetIllegalError(token.MalformedBlockComment)
		return
	}

	tok.Kind = token.BlockComment
	tok.Lit = string(lx.view())

	lx.advance() // skip '*'
	lx.advance() // skip '/'

	return
}
