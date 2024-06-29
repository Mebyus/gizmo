package vm

import "fmt"

func Assemble(unit *ProgUnit) (*Prog, error) {
	a := Assembler{unit: unit}
	return a.Assemble()
}

type Assembler struct {
	prog Prog

	funcs map[string]*ProgFunc
	defs  map[string]*ProgDef

	unit *ProgUnit
}

func (a *Assembler) Assemble() (*Prog, error) {
	err := a.index()
	if err != nil {
		return nil, err
	}
	return &a.prog, nil
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

	a.funcs = make(map[string]*ProgFunc, len(a.unit.Defs))
	for i := range len(a.unit.Funcs) {
		f := &a.unit.Funcs[i]
		_, ok := a.funcs[f.Name]
		if ok {
			return fmt.Errorf("function \"%s\" has multiple definitions", f.Name)
		}

		a.funcs[f.Name] = f
	}
	return nil
}
