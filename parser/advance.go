package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/token"
)

func (p *Parser) advance() {
	p.tok = p.next
	p.next = p.lx.Lex()
}

func (p *Parser) skip(kind token.Kind) {
	if p.tok.Kind != kind {
		panic(fmt.Sprintf("%s: expected %s token, got %s",
			p.tok.Pos.String(), kind.String(), p.tok.Kind.String()))
	}
	p.advance()
}

func (p *Parser) isEOF() bool {
	return p.tok.Kind == token.EOF
}
