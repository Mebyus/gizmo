// package kbs implements Ku Build Script which can be used
// to setup unity builds for Ku IR language.
package kbs

import (
	"fmt"
	"strings"

	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/token"
)

// ScriptOutput result of parsing a build script file(s).
type ScriptOutput struct {
	Includes []string
	Links    []string

	Name string
}

type Parser struct {
	out ScriptOutput

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

func ParseFile(filename string) (*ScriptOutput, error) {
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
func (p *Parser) parse() (*ScriptOutput, error) {
	for {
		if p.isEOF() {
			return &p.out, nil
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
		return p.include()
	case token.DirName:
		return p.name()
	case token.DirLink:
		return p.link()
	default:
		return fmt.Errorf("%s: unexpected %s token instead of directive", p.tok.Pos, p.tok.Kind)
	}
}

func (p *Parser) include() error {
	p.advance() // skip "#include"

	if p.tok.Kind != token.String {
		return fmt.Errorf("%s: unexpected %s token instead of include string", p.tok.Pos, p.tok.Kind)
	}
	include := p.tok.Lit
	if include == "" {
		return fmt.Errorf("%s: empty include string", p.tok.Pos)
	}
	p.advance() // skip include string

	if p.tok.Kind != token.Semicolon {
		return fmt.Errorf("%s: missing semicolon after include directive", p.tok.Pos)
	}
	p.advance() // skip ";"

	p.out.Includes = append(p.out.Includes, include)
	return nil
}

func (p *Parser) name() error {
	if p.out.Name != "" {
		return fmt.Errorf("%s: duplicate name directive", p.tok.Pos)
	}
	p.advance() // skip "#name"

	if p.tok.Kind != token.String {
		return fmt.Errorf("%s: unexpected %s token instead of name string", p.tok.Pos, p.tok.Kind)
	}
	name := strings.TrimSpace(p.tok.Lit)
	if name == "" {
		return fmt.Errorf("%s: empty name string", p.tok.Pos)
	}
	p.advance() // skip name string

	if p.tok.Kind != token.Semicolon {
		return fmt.Errorf("%s: missing semicolon after name directive", p.tok.Pos)
	}
	p.advance() // skip ";"

	p.out.Name = name
	return nil
}

func (p *Parser) link() error {
	p.advance() // skip "#link"

	if p.tok.Kind != token.String {
		return fmt.Errorf("%s: unexpected %s token instead of link string", p.tok.Pos, p.tok.Kind)
	}
	link := p.tok.Lit
	if link == "" {
		return fmt.Errorf("%s: empty link string", p.tok.Pos)
	}
	p.advance() // skip link string

	if p.tok.Kind != token.Semicolon {
		return fmt.Errorf("%s: missing semicolon after link directive", p.tok.Pos)
	}
	p.advance() // skip ";"

	p.out.Links = append(p.out.Links, link)
	return nil
}
