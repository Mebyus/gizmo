package rbs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mebyus/gizmo/butler"
)

var ReBuildSelf = &butler.Lackey{
	Name: "rbs",

	Short: "rebuild ku executable",
	Usage: "ku rbs",

	Exec: execute,
}

func execute(r *butler.Lackey, _ []string) error {
	return rebuild()
}

func rebuild() error {
	path, err := os.Executable()
	if err != nil {
		return err
	}
	rootDir := filepath.Join(path, "./../../..")
	fmt.Println("root dir:", rootDir)
	if rootDir == "." || rootDir == "/" || rootDir == "" {
		return fmt.Errorf("bad root dir \"%s\"", rootDir)
	}

	cmd := exec.Command("go", "build", "-o", path, "./cmd/ku")
	cmd.Dir = rootDir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("updated:", path)
	return nil
}
