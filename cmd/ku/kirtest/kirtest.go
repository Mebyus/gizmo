package kirtest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/compiler/build"
	"github.com/mebyus/gizmo/compiler/cc"
	"github.com/mebyus/gizmo/kir/kbs"
)

var KirTest = &butler.Lackey{
	Name: "kirtest",

	Short: "run kir tests starting from a specified unit",
	Usage: "ku kirtest [options] <file>",

	Exec: execute,
}

func execute(r *butler.Lackey, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}
	return kirtest(files[0])
}

func kirtest(filename string) error {
	script, err := kbs.WalkWithTests("src", filename)
	if err != nil {
		return err
	}

	outc := filepath.Join("build/.test", script.Name+".gen.c")
	err = kbs.GenWithTestsIntoFile(outc, script.Includes)
	if err != nil {
		return err
	}
	outobj := filepath.Join("build/.test", script.Name+".o")
	err = cc.CompileObj(build.Fast, outobj, outc)
	if err != nil {
		return err
	}
	outexe := filepath.Join("build/.test", script.Name)
	err = cc.LinkStatic(build.Fast, outexe, "_start", []string{outobj})
	if err != nil {
		return err
	}

	cmd := exec.Command(outexe)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if cmd.ProcessState == nil {
		return err
	}
	if !cmd.ProcessState.Exited() {
		return fmt.Errorf("(driver) abnormal exit in test executable")
	}
	if cmd.ProcessState.Success() {
		return nil
	}
	return fmt.Errorf("(driver) exit status %d from test executable", cmd.ProcessState.ExitCode())
}
