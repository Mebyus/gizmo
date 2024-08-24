package clean

import (
	"os"
	"path/filepath"

	"github.com/mebyus/gizmo/butler"
)

var Clean = &butler.Lackey{
	Name: "clean",

	Short: "clean project build output and local cache",
	Usage: "ku clean [options]",

	Exec: execute,
}

func execute(r *butler.Lackey, _ []string) error {
	return clean()
}

func clean() error {
	dirs := []string{filepath.Join("build", "target"), filepath.Join("build", ".cache")}
	for _, dir := range dirs {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}
	return nil
}
