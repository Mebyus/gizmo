package vm

import (
	"fmt"
	"os"
)

// ProgUnit represents parsed assembly program.
type ProgUnit struct {
	Defs  []ProgDef
	Funcs []ProgFunc
}

type ProgDef struct {
	Name string
}

type ProgFunc struct {
	Inst   []ProgInst
	Labels []ProgLabel

	Name string
}

type ProgLabel struct {
	Name string

	// Corresponds to instruction index inside the function.
	Index uint32
}

// ProgInst single instruction inside assembler program.
type ProgInst struct {
	// optional constant usage
	DefName string

	// usually it's destination operand
	Operand1 uint64

	// usually it's source operand
	Operand2 uint64

	Op Opcode
}

type Parser struct {
	unit ProgUnit

	tok  *Token
	next *Token

	lx *Lexer
}

func ParseFile(path string) (*ProgUnit, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lx := NewLexer(text)
	return Parse(lx)
}

func Parse(lx *Lexer) (*ProgUnit, error) {
	p := Parser{}
	p.init(lx)
	err := p.parse()
	if err != nil {
		return nil, err
	}
	return &p.unit, nil
}

func (p *Parser) parse() error {
	for {
		if p.tok.Kind == EOF {
			return nil
		}

		var err error
		switch p.tok.Kind {
		case Def:
			err = p.def()
		case Fn:
			err = p.fn()
		default:
			return fmt.Errorf("%s: unexpected top-level %s token", p.tok.Pos(), p.tok.Kind.String())
		}
		if err != nil {
			return err
		}
	}
}

func (p *Parser) init(lx *Lexer) {
	p.lx = lx

	p.advance()
	p.advance()
}

func (p *Parser) advance() {
	p.tok = p.next
	p.next = p.lx.Lex()
}

func (p *Parser) def() error {
	p.advance() // skip "def"

	if p.tok.Kind != Identifier {
		return fmt.Errorf("%s: unexpected %s token instead of name identifier in define construct",
			p.tok.Pos(), p.tok.Kind.String())
	}
	name := p.tok.Lit
	p.advance() // skip identifier

	if p.tok.Kind != HexInteger {
		return fmt.Errorf("%s: unexpected %s token instead of value in define construct",
			p.tok.Pos(), p.tok.Kind.String())
	}
	p.advance() // skip value token
	// TODO: store value in ProgDef struct

	p.unit.Defs = append(p.unit.Defs, ProgDef{
		Name: name,
	})
	return nil
}

func (p *Parser) fn() error {
	p.advance() // skip "fn"

	if p.tok.Kind != Identifier {
		return fmt.Errorf("%s: unexpected %s token instead of name identifier in function construct",
			p.tok.Pos(), p.tok.Kind.String())
	}
	name := p.tok.Lit
	p.advance() // skip identifier

	if p.tok.Kind != Colon {
		return fmt.Errorf("%s: unexpected %s token instead of \":\" in function construct",
			p.tok.Pos(), p.tok.Kind.String())
	}
	p.advance() // ":"

	insts, labels, err := p.body()
	if err != nil {
		return err
	}
	if len(insts) == 0 {
		return fmt.Errorf("function \"%s\" has no body", name)
	}

	p.unit.Funcs = append(p.unit.Funcs, ProgFunc{
		Name: name,

		Inst:   insts,
		Labels: labels,
	})
	return nil
}

func (p *Parser) body() ([]ProgInst, []ProgLabel, error) {
	var insts []ProgInst
	var labels []ProgLabel
	for {
		switch p.tok.Kind {
		case Mnemonic:
			var s ProgInst
			err := p.inst(&s)
			if err != nil {
				return nil, nil, err
			}
			insts = append(insts, s)
		case Label:
			label, err := p.label()
			if err != nil {
				return nil, nil, err
			}
			label.Index = uint32(len(insts))
			labels = append(labels, label)
		default:
			return insts, labels, nil
		}
	}
}

func (p *Parser) inst(s *ProgInst) error {
	me := p.tok.Val
	p.advance() // skip mnemonic

	var err error
	switch me {
	case MeLoad:
		err = p.load(s)
	case MeNop:
		s.Op = Nop
	case MeHalt:
		s.Op = Halt
	case MeSysCall:
		s.Op = SysCall
	case MeTrap:
		s.Op = Trap
	default:
		return fmt.Errorf("%s: instruction(s) with mnemonic \"%s\" is not implemented",
			p.tok.Pos(), meText[me])
	}
	return err
}

func (p *Parser) load(s *ProgInst) error {
	if p.tok.Kind != Register {
		return fmt.Errorf("%s: unexpected %s token instead of register in load instruction",
			p.tok.Pos(), p.tok.Kind.String())
	}
	s.Operand1 = p.tok.Val
	p.advance() // skip destination register

	if p.tok.Kind != Comma {
		return fmt.Errorf("%s: unexpected %s token instead of \",\" in load instruction",
			p.tok.Pos(), p.tok.Kind.String())
	}
	p.advance() // skip ","

	switch p.tok.Kind {
	case Register:
		s.Op = LoadRegReg
		s.Operand2 = p.tok.Val
	case Identifier:
		s.Op = LoadValReg
		s.DefName = p.tok.Lit
	case HexInteger:
		s.Op = LoadValReg
		s.Operand2 = p.tok.Val
	default:
		return fmt.Errorf("%s: unexpected %s token instead of source operand in load instruction",
			p.tok.Pos(), p.tok.Kind.String())
	}
	p.advance() // skip source operand token

	return nil
}

func (p *Parser) label() (ProgLabel, error) {
	name := p.tok.Lit
	p.advance() // skip label name

	if p.tok.Kind != Colon {
		return ProgLabel{}, fmt.Errorf("%s: unexpected %s token instead of \":\" in label construct",
			p.tok.Pos(), p.tok.Kind.String())
	}
	p.advance() // ":"

	return ProgLabel{Name: name}, nil
}
