package asm

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/vm"
)

var Asm = &butler.Lackey{
	Name: "asm",

	Short: "assemble bytecode program from a given kuvm asm file",
	Usage: "gizmo vm asm [options] <file>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}

	return asm(files[0])
}

func asm(path string) error {
	prog := vm.Prog{
		Text: []byte{
			byte(vm.Nop),
			byte(vm.LoadValReg), 0, 3, 0, 0, 0, 0, 0, 0, 0,
			byte(vm.Halt),
		},
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return vm.Encode(f, &prog)
}
