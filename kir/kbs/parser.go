// package kbs implements Ku Build Script which can be used
// to setup unity builds for Ku IR language.
package kbs

import (
	"fmt"

	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/token"
)

type Parser struct {
	includes []string

	lx lexer.Stream

	// token at current parser position
	tok token.Token

	// next token
	next token.Token
}

func New(lx lexer.Stream) *Parser {
	p := &Parser{lx: lx}

	// init parser buffer
	p.advance()
	p.advance()
	return p
}

func FromFile(filename string) (*Parser, error) {
	lx, err := lexer.FromFile(filename)
	if err != nil {
		return nil, err
	}
	return New(lx), nil
}

func ParseFile(filename string) ([]string, error) {
	p, err := FromFile(filename)
	if err != nil {
		return nil, err
	}
	return p.parse()
}

func (p *Parser) advance() {
	p.tok = p.next
	p.next = p.lx.Lex()
}

func (p *Parser) isEOF() bool {
	return p.tok.Kind == token.EOF
}

// returns list of includes from parsed stream
func (p *Parser) parse() ([]string, error) {
	for {
		if p.isEOF() {
			return p.includes, nil
		}

		err := p.top()
		if err != nil {
			return nil, err
		}
	}
}

func (p *Parser) top() error {
	switch p.tok.Kind {
	case token.DirInclude:
		return p.directive()
	default:
		return fmt.Errorf("%s: unexpected %s token instead of directive", p.tok.Pos, p.tok.Kind)
	}
}

func (p *Parser) directive() error {
	p.advance() // skip "#include"

	if p.tok.Kind != token.String {
		return fmt.Errorf("%s: unexpected %s token instead of include string", p.tok.Pos, p.tok.Kind)
	}
	include := p.tok.Lit
	p.advance() // skip include string

	if p.tok.Kind != token.Semicolon {
		return fmt.Errorf("%s: missing semicolon after include", p.tok.Pos)
	}
	p.advance() // skip ";"

	p.includes = append(p.includes, include)
	return nil
}
