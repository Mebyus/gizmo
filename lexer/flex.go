package lexer

import (
	"github.com/mebyus/gizmo/char"
	"github.com/mebyus/gizmo/token"
)

func (lx *Lexer) Flex() token.Token {
	tok := lx.flex()
	lx.num += 1
	return tok
}

func (lx *Lexer) flex() token.Token {
	if lx.EOF {
		return lx.create(token.EOF)
	}

	lx.SkipWhitespace()
	if lx.EOF {
		return lx.create(token.EOF)
	}

	if lx.C == '/' && lx.Next == '/' {
		return lx.lineComment()
	}
	if lx.C == '/' && lx.Next == '*' {
		return lx.blockComment()
	}

	return lx.codeToken()
}

func (lx *Lexer) lineComment() (tok token.Token) {
	tok.Pos = lx.pos()

	lx.Advance() // skip '/'
	lx.Advance() // skip '/'

	// skip until actual comment starts
	for !lx.EOF && lx.C != '\n' && char.IsWhitespace(lx.C) {
		lx.Advance()
	}

	if lx.EOF {
		tok.Kind = token.LineComment
		return
	}
	if lx.C == '\n' {
		lx.SkipWhitespace()
		tok.Kind = token.LineComment
		return
	}

	lx.Start()
	for !lx.EOF && lx.C != '\n' {
		lx.Advance()
	}

	lit, ok := lx.TrimWhitespaceSuffixTake()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return
	}
	tok.Lit = lit
	tok.Kind = token.LineComment

	if lx.EOF {
		return
	}

	lx.SkipWhitespace()
	return
}

func (lx *Lexer) blockComment() (tok token.Token) {
	tok.Pos = lx.pos()

	lx.Advance() // skip '/'
	lx.Advance() // skip '*'

	lx.Start()
	for !lx.EOF && !(lx.C == '*' && lx.Next == '/') {
		lx.Advance()
	}

	if lx.EOF {
		tok.SetIllegalError(token.MalformedBlockComment)
		return
	}

	tok.Kind = token.BlockComment
	tok.Lit = string(lx.View())

	lx.Advance() // skip '*'
	lx.Advance() // skip '/'

	return
}
