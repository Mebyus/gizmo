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
	dirs []Directive

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

func ParseFile(filename string) ([]Directive, error) {
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
func (p *Parser) parse() ([]Directive, error) {
	for {
		if p.isEOF() {
			return p.dirs, nil
		}

		dir, err := p.dir()
		if err != nil {
			return nil, err
		}
		p.dirs = append(p.dirs, dir)
	}
}

func (p *Parser) dir() (Directive, error) {
	switch p.tok.Kind {
	case token.DirInclude:
		return p.include()
	case token.DirName:
		return p.name()
	case token.DirLink:
		return p.link()
	case token.DirIf:
		return p.ifDir()
	default:
		return nil, fmt.Errorf("%s: unexpected %s token instead of directive", p.tok.Pos, p.tok.Kind)
	}
}

func (p *Parser) include() (Include, error) {
	p.advance() // skip "#include"

	if p.tok.Kind != token.String {
		return Include{}, fmt.Errorf("%s: unexpected %s token instead of include string", p.tok.Pos, p.tok.Kind)
	}
	include := p.tok.Lit
	pos := p.tok.Pos
	if include == "" {
		return Include{}, fmt.Errorf("%s: empty include string", p.tok.Pos)
	}
	p.advance() // skip include string

	if p.tok.Kind != token.Semicolon {
		return Include{}, fmt.Errorf("%s: missing semicolon after include directive", p.tok.Pos)
	}
	p.advance() // skip ";"

	return Include{
		Pos:    pos,
		String: include,
	}, nil
}

func (p *Parser) name() (Name, error) {
	p.advance() // skip "#name"

	if p.tok.Kind != token.String {
		return Name{}, fmt.Errorf("%s: unexpected %s token instead of name string", p.tok.Pos, p.tok.Kind)
	}
	name := strings.TrimSpace(p.tok.Lit)
	pos := p.tok.Pos
	if name == "" {
		return Name{}, fmt.Errorf("%s: empty name string", p.tok.Pos)
	}
	p.advance() // skip name string

	if p.tok.Kind != token.Semicolon {
		return Name{}, fmt.Errorf("%s: missing semicolon after name directive", p.tok.Pos)
	}
	p.advance() // skip ";"

	return Name{
		Pos:    pos,
		String: name,
	}, nil
}

func (p *Parser) link() (Link, error) {
	p.advance() // skip "#link"

	if p.tok.Kind != token.String {
		return Link{}, fmt.Errorf("%s: unexpected %s token instead of link string", p.tok.Pos, p.tok.Kind)
	}
	link := p.tok.Lit
	pos := p.tok.Pos
	if link == "" {
		return Link{}, fmt.Errorf("%s: empty link string", p.tok.Pos)
	}
	p.advance() // skip link string

	if p.tok.Kind != token.Semicolon {
		return Link{}, fmt.Errorf("%s: missing semicolon after link directive", p.tok.Pos)
	}
	p.advance() // skip ";"

	return Link{
		Pos:    pos,
		String: link,
	}, nil
}

func (p *Parser) ifDir() (If, error) {
	p.advance() // skip "#if"

	exp, err := p.exp()
	if err != nil {
		return If{}, nil
	}

	body, err := p.block()
	if err != nil {
		return If{}, nil
	}

	return If{
		Clause: IfClause{
			Exp:  exp,
			Body: body,
		},
	}, nil
}

func (p *Parser) block() (Block, error) {
	if p.tok.Kind != token.LeftCurly {
		return Block{}, fmt.Errorf("%s: unexpected token %s instead of block start", p.tok.Pos, p.tok.Kind)
	}
	pos := p.tok.Pos
	p.advance() // skip "{"

	var dirs []Directive
	for {
		if p.tok.Kind == token.RightCurly {
			p.advance() // skip "}"
			return Block{
				Pos:        pos,
				Directives: dirs,
			}, nil
		}

		dir, err := p.dir()
		if err != nil {
			return Block{}, err
		}
		dirs = append(dirs, dir)
	}
}

func (p *Parser) exp() (Exp, error) {
	return nil, nil
}
