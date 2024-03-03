package clean

import (
	"os"

	"github.com/mebyus/gizmo/butler"
)

var Clean = &butler.Lackey{
	Name:  "clean",
	Short: "clean project build output and local cache",
	Usage: "gizmo clean [options]",

	Exec: execute,
}

func execute(r *butler.Lackey, _ []string) error {
	return clean()
}

func clean() error {
	return os.RemoveAll("build")
}
