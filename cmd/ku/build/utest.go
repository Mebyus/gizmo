package build

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/mebyus/gizmo/butler"
)

var Test = &butler.Lackey{
	Name: "test",

	Short: "test specified unit",
	Usage: "ku test [options] <unit>",

	Exec:   executeTest,
	Params: &Params{},
}

func executeTest(r *butler.Lackey, targets []string) error {
	if len(targets) == 0 {
		return fmt.Errorf("at least one unit must be specified")
	}

	start := time.Now()

	params := r.Params.(*Params)
	config, err := NewConfigFromParams(params, targets[0])
	if err != nil {
		return err
	}
	config.Test = true

	outexe, err := buildTarget(config)
	if err != nil {
		return err
	}
	if outexe == "" {
		panic("no test executable")
	}

	testCmd := exec.Command(outexe)
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	testCmd.Run()

	fmt.Printf("total: %s\n", time.Since(start))
	return nil
}
