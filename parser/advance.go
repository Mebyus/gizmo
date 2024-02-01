package parser

func (p *Parser) advance() {
	p.prev = p.tok
	p.tok = p.next
	p.next = p.lx.Lex()
}

func (p *Parser) isEOF() bool {
	return p.tok.IsEOF()
}
