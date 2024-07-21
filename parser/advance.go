package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/token"
)

func (p *Parser) advance() {
	p.tok = p.next

	if p.back == 0 {
		p.next = p.lx.Lex()
	} else {
		p.next = p.buf.Pop()
		p.back -= 1
	}
}

func (p *Parser) skip(kind token.Kind) {
	if p.tok.Kind != kind {
		panic(fmt.Sprintf("%s: expected %s token, got %s",
			p.tok.Pos.String(), kind.String(), p.tok.Kind.String()))
	}
	p.advance()
}

// same as advance, but stores consumed tokens for optional backtrack later
func (p *Parser) advanceBackup() {
	p.buf.Push(p.tok)
	p.advance()
}

// backtrack parser by n tokens, this operation can only be used after
// at least n consecutive backups
func (p *Parser) backtrack(n uint) {
	if n == 0 {
		panic("bad argument")
	}
	if n > p.buf.len {
		panic("not enough stored tokens")
	}
	p.back = n
	p.buf.len = n

	tok := p.tok
	p.tok = p.buf.Pop()
	p.buf.Push(tok)

	next := p.next
	p.next = p.buf.Pop()
	p.buf.Push(next)
}

func (p *Parser) isEOF() bool {
	return p.tok.IsEOF()
}
