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
	Val  uint64
}

type ProgFunc struct {
	Inst   []ProgInst
	Labels []ProgLabel

	Name string

	// Text (code) size of this function when encoded in binary.
	// Aligned by 16.
	TextSize uint32

	// Offset of text (code) of this function when encoded in binary
	// relative to start of the program binary text section.
	// Aligned by 16.
	TextOffset uint32
}

type ProgLabel struct {
	Name string

	// Offset of instruction under this label. Counted from
	// start of the function.
	Offset uint32

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

	// current text offset
	toff uint32
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
	val := p.tok.Val
	p.advance() // skip value token

	p.unit.Defs = append(p.unit.Defs, ProgDef{
		Name: name,
		Val:  val,
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

	f := ProgFunc{
		Name:       name,
		TextOffset: p.toff,
	}
	err := p.body(&f)
	if err != nil {
		return err
	}
	if len(f.Inst) == 0 {
		return fmt.Errorf("function \"%s\" has no body", name)
	}
	f.TextSize = AlignSizeBy16(f.TextSize)
	p.toff += f.TextSize

	p.unit.Funcs = append(p.unit.Funcs, f)
	return nil
}

func (p *Parser) body(f *ProgFunc) error {
	for {
		switch p.tok.Kind {
		case Mnemonic:
			var s ProgInst
			err := p.inst(&s)
			if err != nil {
				return err
			}
			f.Inst = append(f.Inst, s)
			f.TextSize += uint32(Size[s.Op])
		case Label:
			label, err := p.label()
			if err != nil {
				return err
			}
			label.Index = uint32(len(f.Inst))
			label.Offset = f.TextSize
			f.Labels = append(f.Labels, label)
		default:
			return nil
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
	case MeJump:
		err = p.jump(s)
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

func (p *Parser) jump(s *ProgInst) error {
	switch p.tok.Kind {
	case Register: // TODO: handle syntax [r3]
		panic("not implemented")
	case Label:
		p.jumpLabel(s)
	default:
		return fmt.Errorf("%s: unexpected %s token instead of destination operand in jump instruction",
			p.tok.Pos(), p.tok.Kind.String())
	}
	return nil
}

func (p *Parser) jumpLabel(s *ProgInst) {
	name := p.tok.Lit
	p.advance() // skip label name

	s.Op = JumpAddr
	s.DefName = name
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
