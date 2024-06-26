package vm

import "fmt"

func Assemble(unit *ProgUnit) (*Prog, error) {
	a := Assembler{unit: unit}
	return a.Assemble()
}

type Assembler struct {
	text Buffer

	// local to currently processed function
	labels map[string]*ProgLabel

	funcs map[string]*ProgFunc
	defs  map[string]*ProgDef

	unit *ProgUnit

	// current function offset
	foff uint32
}

func (a *Assembler) Assemble() (*Prog, error) {
	err := a.index()
	if err != nil {
		return nil, err
	}
	err = a.assemble()
	if err != nil {
		return nil, err
	}
	return &Prog{
		Text: a.text.Bytes(),
	}, nil
}

func (a *Assembler) index() error {
	err := a.indexDefs()
	if err != nil {
		return err
	}
	return a.indexFuncs()
}

func (a *Assembler) indexDefs() error {
	if len(a.unit.Defs) == 0 {
		return nil
	}

	a.defs = make(map[string]*ProgDef, len(a.unit.Defs))
	for i := range len(a.unit.Defs) {
		d := &a.unit.Defs[i]
		_, ok := a.defs[d.Name]
		if ok {
			return fmt.Errorf("name \"%s\" has multiple definitions", d.Name)
		}

		a.defs[d.Name] = d
	}
	return nil
}

func (a *Assembler) indexFuncs() error {
	if len(a.unit.Funcs) == 0 {
		return nil
	}

	a.labels = make(map[string]*ProgLabel)

	// total size of program text
	var size uint32

	a.funcs = make(map[string]*ProgFunc, len(a.unit.Defs))
	for i := range len(a.unit.Funcs) {
		f := &a.unit.Funcs[i]
		_, ok := a.funcs[f.Name]
		if ok {
			return fmt.Errorf("function \"%s\" has multiple definitions", f.Name)
		}

		a.funcs[f.Name] = f
		size += f.TextSize
	}

	a.text.Init(int(size))
	return nil
}

func (a *Assembler) assemble() error {
	// TODO: add align padding between functions
	// and change opcode 0 to trap

	for i := range len(a.unit.Funcs) {
		f := &a.unit.Funcs[i]

		err := a.fn(f)
		if err != nil {
			return err
		}
		a.text.Align16()
	}

	return nil
}

func (a *Assembler) indexLabels(f *ProgFunc) error {
	a.foff = f.TextOffset

	clear(a.labels)
	for i := range len(f.Labels) {
		l := &f.Labels[i]
		name := l.Name

		_, ok := a.labels[name]
		if ok {
			return fmt.Errorf("label \"%s\" has multiple definitions inside \"%s\" function", name, f.Name)
		}

		a.labels[name] = l
	}
	return nil
}

func (a *Assembler) fn(f *ProgFunc) error {
	err := a.indexLabels(f)
	if err != nil {
		return err
	}

	for j := range len(f.Inst) {
		i := &f.Inst[j]

		var err error
		switch i.Op {
		case LoadValReg:
			if i.DefName == "" {
				a.loadLitValReg(uint8(i.Operand1), i.Operand2)
			} else {
				err = a.loadDefValReg(uint8(i.Operand1), i.DefName)
			}
		case LoadRegReg:
			a.loadRegReg(uint8(i.Operand1), uint8(i.Operand2))
		case JumpAddr:
			if i.DefName == "" {
				panic("empty label name")
			}
			err = a.jumpLabelAddr(i.DefName)
		case Halt:
			a.opcode(Halt)
		case Trap:
			a.opcode(Trap)
		case Nop:
			a.opcode(Nop)
		case SysCall:
			a.opcode(SysCall)
		default:
			panic(fmt.Sprintf("opcode 0x%x not implemented", i.Op))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Assembler) jumpLabelAddr(name string) error {
	label, ok := a.labels[name]
	if !ok {
		return fmt.Errorf("label \".%s\" is not defined", name)
	}

	addr := a.foff + label.Offset
	// TODO: add panic-check to validate address
	// is in range of program text

	a.jumpAddr(addr)
	return nil
}

func (a *Assembler) loadDefValReg(dr uint8, name string) error {
	def, ok := a.defs[name]
	if !ok {
		return fmt.Errorf("constant \"%s\" is not defined", name)
	}
	a.loadLitValReg(dr, def.Val)
	return nil
}

func (a *Assembler) jumpAddr(addr uint32) {
	a.opcode(JumpAddr)
	a.text.Val32(addr)
}

func (a *Assembler) loadLitValReg(dr uint8, v uint64) {
	a.opcode(LoadValReg)
	a.reg(dr)
	a.text.Val64(v)
}

func (a *Assembler) loadRegReg(dr, sr uint8) {
	a.opcode(LoadRegReg)
	a.reg(dr)
	a.reg(sr)
}

func (a *Assembler) reg(r uint8) {
	a.text.Val8(r)
}

func (a *Assembler) opcode(op Opcode) {
	a.text.Val8(uint8(op))
}
