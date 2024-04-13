package format

import (
	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/token"
)

// Holder decorates Lexer.Flex method to store all output tokens
// for later usage. It also removes comment tokens from the stream.
type Holder struct {
	// Stored tokens from lexer.
	tokens []token.Token

	lx *lexer.Lexer
}

func NewHolder(lx *lexer.Lexer) *Holder {
	return &Holder{lx: lx}
}

func (h *Holder) Lex() token.Token {
	for {
		tok := h.lx.Flex()
		h.tokens = append(h.tokens, tok)

		if tok.Kind.IsComment() {
			continue
		}

		return tok
	}
}

func (h *Holder) Tokens() []token.Token {
	return h.tokens
}
