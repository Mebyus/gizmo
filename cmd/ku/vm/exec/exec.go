package exec

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/vm"
)

var Exec = &butler.Lackey{
	Name: "exec",

	Short: "execute kuvm bytecode program",
	Usage: "gizmo vm exec [options] <file>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}

	return exec(files[0])
}

func exec(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	prog, err := vm.Decode(f)
	if err != nil {
		return err
	}

	m := vm.Machine{}
	m.Init(&vm.Config{
		StackSize:    1 << 24,
		InitHeapSize: 0,
	})

	e := m.Exec(prog)
	return e.Render(os.Stderr)
}
