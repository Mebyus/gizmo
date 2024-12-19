package amd64

import (
	"fmt"
	"os"
)

func (p *Parser) Next(o *Operation) bool {
	if p.Error != nil {
		return false
	}

	switch p.tok.Kind {
	case TokEOF:
		return false
	case TokLabel:
		if p.next.Kind != TokColon {
			p.Error = fmt.Errorf("%s: unexpected %s token after label placement", p.next.Pos(), p.next.Kind)
			return false
		}

		o.Op = OpLabel
		o.Val = p.tok.Val // contains label index

		p.advance() // skip label
		p.advance() // skip ":"

		return true
	case TokMnemonic:
		o.Op = OpProto
		o.Val = p.tok.Val // contains mnemonic

		n := p.tok.Line
		p.advance() // skip mnemonic

		var k uint8
		for p.tok.Line == n {
			if k >= 2 {
				p.Error = fmt.Errorf("%s: too many line elements", p.tok.Pos())
				return false
			}

			switch p.tok.Kind {
			case TokRegister:
				o.Args[k] = OpArg{Kind: ArgRegister, Val: p.tok.Val}
			case TokHexInteger:
				o.Args[k] = OpArg{Kind: ArgInteger, Val: p.tok.Val}
			case TokLabel:
				o.Args[k] = OpArg{Kind: ArgLabel, Val: p.tok.Val}
			case TokFlag:
				o.Args[k] = OpArg{Kind: ArgFlag, Val: p.tok.Val}
			default:
				p.Error = fmt.Errorf("%s: unexpected %s token in mnemonic argument", p.tok.Pos(), p.tok.Kind)
				return false
			}
			k += 1

			p.advance() // skip arg

			if p.tok.Kind == TokComma {
				p.advance() // skip ","
			}
		}
		o.Num = k
		return true
	default:
		p.Error = fmt.Errorf("%s: unexpected %s token at line start", p.tok.Pos(), p.tok.Kind)
		return false
	}
}

type Parser struct {
	Error error

	Labels *LabelMap

	tok  *Token
	next *Token

	lx *Lexer
}

func NewParserFromFile(path string) (*Parser, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lx := NewLexer(text)
	return NewParser(lx), nil
}

func NewParser(lx *Lexer) *Parser {
	p := Parser{Labels: NewLabelMap()}
	p.init(lx)
	return &p
}

func (p *Parser) init(lx *Lexer) {
	p.lx = lx

	p.advance()
	p.advance()
}

func (p *Parser) advance() {
	tok := p.lx.Lex()
	if tok.Kind == TokLabel {
		tok.Val = p.Labels.Push(tok.Text)
	}

	p.tok = p.next
	p.next = tok
}
