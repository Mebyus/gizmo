package kbs

import (
	"fmt"

	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/token"
)

const procBufSize = 1 << 8

type MacroProcessor struct {
	buf [procBufSize]token.Token

	lx lexer.Stream

	err error

	// maps macro name to its definition (tokens which are produced upon macro placement)
	macros map[ /* macro name */ string][]token.Token

	// position inside buf
	pos int

	// number of tokens in buffer
	len int
}

func (p *MacroProcessor) SetInput(lx lexer.Stream) {
	if p.err != nil {
		panic("must not be called with saved error")
	}
	if p.len != 0 {
		panic("buffer was not exhausted")
	}

	if p.macros == nil {
		p.macros = make(map[string][]token.Token)
	}
	p.lx = lx
}

func (p *MacroProcessor) Lex() token.Token {
	if p.err != nil {
		return p.eofErrorToken()
	}

	p.process()
	if p.err != nil {
		return p.eofErrorToken()
	}

	tip := p.tip()
	p.advance() // skip tip
	return tip
}

// handle macro logic:
//   - defines
//   - usage
func (p *MacroProcessor) process() {
	for {
		tip := p.tip()

		var err error
		switch tip.Kind {
		case token.DirDefine:
			err = p.define()
		case token.Macro:
			err = p.expand()
		default:
			return
		}
		if err != nil {
			p.err = err
			return
		}
	}
}

func (p *MacroProcessor) eofErrorToken() token.Token {
	return token.Token{
		Kind: token.EOF,
		Lit:  p.err.Error(),
	}
}

func (p *MacroProcessor) expand() error {
	tip := p.tip()
	name := tip.Lit
	pos := tip.Pos
	tokens, ok := p.macros[name]
	if !ok {
		return fmt.Errorf("%s: usage of unknown macro #.%s", pos, name)
	}
	p.advance() // skip macro usage

	if p.len+len(tokens) > procBufSize {
		return fmt.Errorf("%s: expansion of macro #.%s does not fit into buffer", pos, name)
	}
	p.putback(tokens)
	return nil
}

func (p *MacroProcessor) putback(tokens []token.Token) {
	i := len(tokens)
	for i != 0 {
		i -= 1

		tok := tokens[i]
		p.buf[p.dec()] = tok
		p.len += 1
	}
}

// process macro definition and save it
func (p *MacroProcessor) define() error {
	p.advance() // skip "#define"

	tip := p.tip()
	if tip.Kind != token.Identifier {
		return fmt.Errorf("%s: unexpected %s token instead of identifier in macro definition", tip.Pos, tip.Kind)
	}

	name := tip.Lit
	_, ok := p.macros[name]
	if ok {
		return fmt.Errorf("%s: duplicate macro #.%s definition", tip.Pos, name)
	}
	p.advance() // skip macro name

	tip = p.tip()
	if tip.Kind != token.RightArrow {
		return fmt.Errorf("%s: unexpected %s token instead of \"=>\" before macro definition", tip.Pos, tip.Kind)
	}
	p.advance() // skip "=>"

	tip = p.tip()
	if tip.Kind != token.LeftCurly {
		return fmt.Errorf("%s: unexpected %s token instead of \"{\" in macro definition start", tip.Pos, tip.Kind)
	}
	p.advance() // skip "{"

	var tokens []token.Token
	for {
		// TODO: handle nested curly braces
		tip = p.tip()
		if tip.Kind != token.RightCurly {
			tokens = append(tokens, tip)
			p.advance()
			continue
		}

		if len(tokens) == 0 {
			return fmt.Errorf("%s: empty macro definition", tip.Pos)
		}
		p.advance() // skip "}"

		if len(tokens) > procBufSize {
			return fmt.Errorf("%s: macro definition is too big", tip.Pos)
		}
		p.macros[name] = tokens
		return nil
	}
}

func (p *MacroProcessor) load() {
	p.buf[p.tail()] = p.lx.Lex()
	p.len += 1
}

func (p *MacroProcessor) advance() {
	if p.len == 0 {
		panic("empty buffer")
	}

	p.len -= 1
	p.pos = p.nx()
}

func (p *MacroProcessor) tip() token.Token {
	if p.len == 0 {
		p.load()
	}
	return p.buf[p.pos]
}

func (p *MacroProcessor) nx() int {
	return (p.pos + 1) % procBufSize
}

func (p *MacroProcessor) dec() int {
	if p.pos == 0 {
		p.pos = procBufSize - 1
	} else {
		p.pos -= 1
	}
	return p.pos
}

func (p *MacroProcessor) tail() int {
	return (p.pos + p.len) % procBufSize
}
