package gui

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/vm"
)

var Serve = &butler.Lackey{
	Name: "gui",

	Short: "serve web ui for vm debugger",
	Usage: "gizmo vm gui [options] <file>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}

	return serve(files[0])
}

func serve(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	prog, err := vm.Decode(f)
	if err != nil {
		return err
	}

	s := vm.NewDebugServer("localhost:4067", "vm/gui", prog)
	err = s.ListenAndServe()
	if err != nil {
		fmt.Printf("serve error: %v\n", err)
	}
	return nil
}
