package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/genstf"
	"github.com/mebyus/gizmo/parser"
	"github.com/mebyus/gizmo/stg"
	"github.com/mebyus/gizmo/uwalk"
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

// returns path to executable
func buildTestExe(config *Config, program *uwalk.Program) (string, error) {
	if program.TestCount == 0 {
		return "", fmt.Errorf("no tests found")
	}

	err := os.MkdirAll("build/.test", 0o755)
	if err != nil {
		return "", err
	}

	base := filepath.Base(config.InitPath)
	outku := filepath.Join("build/.test", base+".gen.ku")
	err = genTestDriverCode(outku, program)
	if err != nil {
		return "", err
	}

	err = addTestMainUnitToProgram(outku, program)
	if err != nil {
		return "", err
	}

	outc := filepath.Join("build/.test", base+".gen.c")
	err = genCode(outc, program)
	if err != nil {
		return "", err
	}

	outobj := filepath.Join("build/.test", base+".o")
	err = compile(outobj, outc)
	if err != nil {
		return "", err
	}

	outexe := filepath.Join("build/.test", base+".exe")
	err = link(outexe, outobj)
	if err != nil {
		return "", err
	}
	return outexe, nil
}

func addTestMainUnitToProgram(src string, program *uwalk.Program) error {
	atom, err := parser.ParseFile(src)
	if err != nil {
		return err
	}

	u, err := stg.Merge(
		stg.UnitContext{
			Name:     "main",
			Resolver: stg.MapResolver(program.Map),
			Global:   program.Global,
		},
		[]*ast.Atom{atom},
	)
	if err != nil {
		return err
	}

	program.Main = u
	program.Units = append(program.Units, u)

	index := len(program.Graph.Nodes)
	program.Graph.Nodes = append(program.Graph.Nodes, uwalk.Node{
		Unit: u,
	})
	program.Graph.Cohorts = append(program.Graph.Cohorts, []int{index})

	return nil
}

func genTestDriverCode(out string, program *uwalk.Program) error {
	start := time.Now()

	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	err = genstf.Gen(f, program)
	if err != nil {
		return err
	}

	fmt.Printf("genstf: %s\n", time.Since(start))
	return nil
}
