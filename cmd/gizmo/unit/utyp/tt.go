package utyp

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/er"
	"github.com/mebyus/gizmo/parser"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/tt"
)

var Utyp = &butler.Lackey{
	Name: "tt",

	Short: "construct type tree of a given unit",
	Usage: "gizmo unit tt <unit>",

	Exec: execute,
}

func execute(r *butler.Lackey, units []string) error {
	if len(units) == 0 {
		return fmt.Errorf("at least one unit must be specified")
	}
	err := utyp(units[0])
	if err != nil {
		e, ok := err.(er.Error)
		if !ok {
			return err
		}
		e.Format(os.Stdout, "text")
	}
	return err
}

func utyp(unit string) error {
	files, err := source.LoadUnitFiles(unit)
	if err != nil {
		return err
	}

	m := tt.New(tt.UnitContext{Global: tt.NewGlobalScope()})
	for _, file := range files {
		atom, err := parser.ParseSource(file)
		if err != nil {
			return err
		}
		err = m.Add(atom)
		if err != nil {
			return err
		}
	}
	_, err = m.Merge()
	if err != nil {
		return err
	}
	for _, warn := range m.Warns {
		fmt.Println(warn)
	}
	return nil
}
